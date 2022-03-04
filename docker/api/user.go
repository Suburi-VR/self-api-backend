package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
	// "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"

	"fmt"
)

type dbData struct {
	ID uint
	username map[string]string
	secret map[string]string
	orgid map[string]int
	nickname map[string]string
	kana map[string]string
	company map[string]string
	department map[string]string
	CreatedAt time.Time
    UpdatedAt time.Time
}

func (u *dbData) UnmarshalJSON(b []byte) error {
	// 自分で新しく定義した構造体
    u2 := &struct {
        Username map[string]string
        Secret map[string]string
        Orgid map[string]int
        Nickname map[string]string
        Kana map[string]string
        Company map[string]string
        Department map[string]string
    }{}
    err := json.Unmarshal(b, u2)
	if err != nil {
		panic(err)
	}
	// 新しく定義した構造体の結果をもとのpに詰める
	u.username = u2.Username
    u.secret = u2.Secret
    u.orgid = u2.Orgid
    u.kana = u2.Kana
    u.company = u2.Company
    u.department = u2.Department
 
	return err
}

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

//     jsonData := 
//         `{
//     "username": {"S": "` + *userName + `"},
//     "secret": {"S": "?"},
//     "orgid": {"N": "0"},
//     "nickname": {"S": "testNickname"},
//     "kana": {"S": "Kana"},
//     "company": {"S": "kaisha"},
//     "department": {"S": "busho"}
// }`
//     data := []byte(jsonData)
//     f, createFileError := os.Create("item.json")
//     f.Write(data)
//     if createFileError != nil {
//         fmt.Println(createFileError)
//         fmt.Println("fail to write file")
//     }

//     db := dynamodb.New(sess)
//     item := getItems()
//     tableName := "UserTable"
        
//     av, dbErr := dynamodbattribute.MarshalMap(item)
//     if dbErr != nil {
//         log.Fatalf("Got error marshalling map: %s", dbErr)
//     }
//     input := &dynamodb.PutItemInput{
//         Item:      av,
//         TableName: aws.String(tableName),
//     }
//     _, err = db.PutItem(input)
//     if err != nil {
//         log.Fatalf("Got error calling PutItem: %s", err)
//     }

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

func getItems() []dbData {
    raw, err := os.ReadFile("item.json")
    if err != nil {
        log.Fatalf("Got error reading file: %s", err)
    }

    var item []dbData

    json.Unmarshal(raw, &item)
    fmt.Println("------items------")
    fmt.Println(string(raw))
    fmt.Println(item)
    fmt.Println("------------")
    return item
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