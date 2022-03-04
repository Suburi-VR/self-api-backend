package main

import (
	"crypto/rand"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"

	"api/models"
	"fmt"
)

func main() {
    r := gin.Default()
    r.POST("/user/create", create)
    r.GET("/user/info", getInfo)
    r.POST("/user/info", updateInfo)
    r.Run()
}

func create(c *gin.Context) {
    email := MakeRandomStr(10) + "@example.com"
    pass := MakeRandomStr(10) + "&1"
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        Config: aws.Config{
            Region: aws.String("ap-northeast-1"),
            CredentialsChainVerboseErrors: aws.Bool(true),
        },
    }))

    cognitoClient := cognitoidentityprovider.New(sess)

    newUserData := &cognitoidentityprovider.AdminCreateUserInput{
        UserAttributes: []*cognitoidentityprovider.AttributeType{
            {
                Name:  aws.String("custom:supporter"),
                Value: aws.String("false"),
            },
        },
    }
		
    newUserData.SetUserPoolId("ap-northeast-1_Kjb4vUZPh")
    newUserData.SetUsername(email)
    newUserData.SetTemporaryPassword(pass)

    fmt.Println(newUserData)
    _, err := cognitoClient.AdminCreateUser(newUserData)
    fmt.Println(cognitoClient.Endpoint)
    if err != nil {
        fmt.Println("Got error creating user:", err)
    }

    nawPass := MakeRandomStr(10) + "&&1"
    userName := newUserData.Username

    newPassword := &cognitoidentityprovider.AdminSetUserPasswordInput{
        Password: aws.String(nawPass),
        Permanent: aws.Bool(true),
        UserPoolId: aws.String("ap-northeast-1_Kjb4vUZPh"),
        Username: aws.String(*userName),
    }

    _, e := cognitoClient.AdminSetUserPassword(newPassword)
    if e != nil {
        fmt.Println("Got error creating new password:", e)
    }

    c.JSON(200, gin.H{
        "user": email,
        "pass": nawPass,
    })

    dbData := models.DbData{
        Username: *userName,
        Secret: "?",
        Orgid: 0,
        Nickname: "testNickname",
        Kana: "Kana",
        Company: "kaisha",
        Department: "busho",
    }

    db := dynamodb.New(sess, &aws.Config{Endpoint: aws.String("http://192.168.1.8:8000")})
    tableName := "UserTable"
        
    av, dbErr := dynamodbattribute.MarshalMap(dbData)
    if dbErr != nil {
        log.Fatalf("Got error marshalling map: %s", dbErr)
    }

    input := &dynamodb.PutItemInput{
        Item:      av,
        TableName: aws.String(tableName),
    }
    
    _, err = db.PutItem(input)
    if err != nil {
        log.Fatalf("Got error calling PutItem: %s", err)
    }

    return
}

func MakeRandomStr(digit uint32) (string) {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

    // 乱数を生成
    b := make([]byte, digit)
    if _, err := rand.Read(b); err != nil {
        return ""
    }

    // letters からランダムに取り出して文字列を生成
    var result string
    for _, v := range b {
        // index が letters の長さに収まるように調整
        result += string(letters[int(v)%len(letters)])
    }
    return result
}

func getInfo(c *gin.Context) {

    c.JSON(200, gin.H{
        "username": "testuser",
        "orgid": 1,
        "nickName": "testuser",
        "kana": "テスト ユーザ",
        "company": "test inc",
        "department": "develop",
      })

    return
}

func updateInfo(c *gin.Context) {

    return
}