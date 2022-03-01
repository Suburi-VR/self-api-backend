package main

import (
	"crypto/rand"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"

	"fmt"
)

func main() {
    r := gin.Default()
    r.POST("/user/create", create)
    r.Run()
}


func create(c *gin.Context) {
    // emailIDPtr := flag.String("e", "", "The email address of the user")
    // userPoolIDPtr := flag.String("p", "", "The ID of the user pool")
    // userNamePtr := flag.String("n", "", "The name of the user")

    // flag.Parse()

    // fmt.Println("1")
    // if *emailIDPtr == "" || *userPoolIDPtr == "" || *userNamePtr == "" {
    //     fmt.Println("2")
    //     fmt.Println("You must supply an email address, user pool ID, and user name")
    //     fmt.Println("Usage: go run CreateUser.go -e EMAIL-ADDRESS -p USER-POOL-ID -n USER-NAME")
    //     os.Exit(1)
    // }
    random, _ := MakeRandomStr(10)
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        Config: aws.Config{
            Region: aws.String("ap-northeast-1"),
            CredentialsChainVerboseErrors: aws.Bool(true),
        },
    }))

    cognitoClient := cognitoidentityprovider.New(sess)

    cognitoClient.Client.Endpoint = "http://192.168.1.8:9229"
    newUserData := &cognitoidentityprovider.AdminCreateUserInput{
        DesiredDeliveryMediums: []*string{
            aws.String("EMAIL"),
        },
        UserAttributes: []*cognitoidentityprovider.AttributeType{
            {
                Name:  aws.String("email"),
                Value: aws.String(random),
            },
        },
    }
		
    newUserData.SetUserPoolId("local_5HgXw3xJ")
    newUserData.SetUsername(random)

    fmt.Println(newUserData)
    _, err := cognitoClient.AdminCreateUser(newUserData)
    fmt.Println(cognitoClient)
    if err != nil {
        fmt.Println("Got error creating user:", err)
    }

    c.JSON(200, gin.H{
        "user": "name",
        "pass": "pass",
    })

    return
}

func MakeRandomStr(digit uint32) (string, error) {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

    // 乱数を生成
    b := make([]byte, digit)
    if _, err := rand.Read(b); err != nil {
        return "", errors.New("unexpected error...")
    }

    // letters からランダムに取り出して文字列を生成
    var result string
    for _, v := range b {
        // index が letters の長さに収まるように調整
        result += string(letters[int(v)%len(letters)])
    }
    return result, nil
}