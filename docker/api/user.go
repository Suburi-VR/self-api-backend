package main

import (
	"crypto/rand"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
    "github.com/gin-gonic/gin"

	"fmt"
)

func main() {
    r := gin.Default()
    r.POST("/user/create", create)
    r.GET("/user/info", info)
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
        DesiredDeliveryMediums: []*string{
            aws.String("EMAIL"),
        },
        UserAttributes: []*cognitoidentityprovider.AttributeType{
            {
                Name:  aws.String("email"),
                Value: aws.String(email),
            },
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

    nawPass := MakeRandomStr(10) + "&&"
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

func info(c *gin.Context) {

    c.JSON(200, gin.H{
        "username": "testuser",
        "nickName": "testuser",
        "kana": "テスト ユーザ",
        "company": "test inc",
        "department": "develop",
      })

    return
}