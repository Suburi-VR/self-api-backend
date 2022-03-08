package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"

	"api/config"
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

    _, err := config.CognitoClient.AdminCreateUser(newUserData)
    fmt.Println(config.CognitoClient.Endpoint)
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

    _, e := config.CognitoClient.AdminSetUserPassword(newPassword)
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

        
    av, dbErr := dynamodbattribute.MarshalMap(dbData)
    if dbErr != nil {
        log.Fatalf("Got error marshalling map: %s", dbErr)
    }

    input := &dynamodb.PutItemInput{
        Item:      av,
        TableName: aws.String(config.TableName),
    }
    
    _, err = config.Db.PutItem(input)
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

func userName(c *gin.Context) string {
    // idToken取得
    idToken := c.Request.Header.Get("Authorization")
    sprited := strings.Split(idToken, ".")

    //　ユーザー情報を取り出す([]byte)
    userInfo, err := base64.RawStdEncoding.DecodeString(sprited[1])
    if err != nil {
        log.Fatalf("Got error calling PutItem: %s", err)
    }

    var mapData map[string]string
    json.Unmarshal(userInfo, &mapData)
    username := mapData["cognito:username"]

    return username
}

func getInfo(c *gin.Context) {

    username := userName(c)

    input := &dynamodb.GetItemInput{
        TableName: aws.String(config.TableName),
        Key: map[string]*dynamodb.AttributeValue{
            "username": {
                S: &username,
            },
        },
    }

    result, err := config.Db.GetItem(input)
    if err != nil {
        fmt.Println("[GetItem Error]", err)
        return
    }

    user := result.Item

    fmt.Println("---------user------")
    fmt.Println(user)
    fmt.Println("---------user------")

    c.JSON(200, gin.H{
        "username": user["username"].S,
        "orgid": user["orgid"].N,
        "nickName": user["nickname"].S,
        "kana": user["kana"].S,
        "company": user["company"].S,
        "department": user["department"].S,
      })

    return
}

func updateInfo(c *gin.Context) {

    var body map[string]string
    c.BindJSON(&body)

    username := userName(c)

    params := &dynamodb.UpdateItemInput {
        TableName: aws.String(config.TableName),
        Key: map[string]*dynamodb.AttributeValue{
            "username": {
                S: &username,
            },
        },
        UpdateExpression: aws.String("set #platform = :platform, #deviceToken = :deviceToken, #nickname = :nickname, #company = :company, #department = :department"),
        ExpressionAttributeNames: map[string]*string {
            "#platform": aws.String("platform"),
            "#deviceToken": aws.String("deviceToken"),
            "#nickname": aws.String("nickname"),
            "#company": aws.String("company"),
            "#department": aws.String("department"),
        },
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
            ":platform": {
                S: aws.String(body["platform"]),
            },
            ":deviceToken": {
                S: aws.String(body["deviceToken"]),
            },
            ":nickname": {
                S: aws.String(body["nickname"]),
            },
            ":company": {
                S: aws.String(body["company"]),
            },
            ":department": {
                S: aws.String(body["department"]),
            },
        },
        ReturnValues: aws.String("ALL_NEW"),
        ReturnConsumedCapacity: aws.String("TOTAL"),
        ReturnItemCollectionMetrics: aws.String("SIZE"),
    }

    _, err := config.Db.UpdateItem(params)
    if err != nil {
        log.Fatalf("Got error calling UpdateItem: %s", err)
    }

    return
}