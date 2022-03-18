package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"

	"api/config"
	"api/models"
	"api/utils"
)

var User models.User

func create(c *gin.Context) {

    email := utils.MakeRandomStr(10) + "@example.com"
    pass := utils.MakeRandomStr(10) + "&1"

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
    if err != nil {
        log.Fatalf("Got error creating user: %s", err)
    }

    newPass := utils.MakeRandomStr(10) + "&&1"
    userName := newUserData.Username

    newPassword := &cognitoidentityprovider.AdminSetUserPasswordInput{
        Password: aws.String(newPass),
        Permanent: aws.Bool(true),
        UserPoolId: aws.String("ap-northeast-1_Kjb4vUZPh"),
        Username: aws.String(*userName),
    }

    _, e := config.CognitoClient.AdminSetUserPassword(newPassword)
    if e != nil {
        log.Fatalf("Got error create new password: %s", e)
    }

    c.JSON(200, gin.H{
        "user": email,
        "pass": newPass,
    })

    User = models.User{
        Username: *userName,
        Secret: newPass,
        Orgid: 0,
        Nickname: "testNickname",
        Kana: "Kana",
        Company: "kaisha",
        Department: "busho",
    }

        

    av, dbErr := dynamodbattribute.MarshalMap(User)
    if dbErr != nil {
        log.Fatalf("Got error marshalling map: %s", dbErr)
    }

    input := &dynamodb.PutItemInput{
        Item:      av,
        TableName: aws.String(config.UserTable),
    }
    
    _, err = config.Db.PutItem(input)
    if err != nil {
        log.Fatalf("Got error calling PutItem: %s", err)
    }

    return
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
        TableName: aws.String(config.UserTable),
        Key: map[string]*dynamodb.AttributeValue{
            "username": {
                S: &username,
            },
        },
    }

    result, err := config.Db.GetItem(input)
    if err != nil {
        log.Fatalf("Got error calling GetItem: %s", err)
        return
    }

    
    user := result.Item

    User.Username = *user["username"].S
    User.Orgid, _ = strconv.Atoi(*user["orgid"].N)
    User.Nickname = *user["nickname"].S
    User.Kana = *user["kana"].S
    User.Company = *user["company"].S
    User.Department = *user["department"].S

    c.JSON(200, gin.H{
        "username": User.Username,
        "orgid": User.Orgid,
        "nickname": User.Nickname,
        "kana": User.Kana,
        "company": User.Company,
        "department": User.Department,
      })
    return
}

func updateInfo(c *gin.Context) {

    var body map[string]string
    c.BindJSON(&body)

    User.Platform = body["platform"]
    User.DeviceToken = body["deviceToken"]
    User.Nickname = body["nickname"]
    User.Company = body["company"]
    User.Department = body["department"]

    username := User.Username

    params := &dynamodb.UpdateItemInput {
        TableName: aws.String(config.UserTable),
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
                S: aws.String(User.Platform),
            },
            ":deviceToken": {
                S: aws.String(User.DeviceToken),
            },
            ":nickname": {
                S: aws.String(User.Nickname),
            },
            ":company": {
                S: aws.String(User.Company),
            },
            ":department": {
                S: aws.String(User.Department),
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

func contact(c *gin.Context) {
    username := userName(c)

    input := &dynamodb.GetItemInput{
        TableName: aws.String(config.UserTable),
        Key: map[string]*dynamodb.AttributeValue{
            "username": {
                S: &username,
            },
        },
    }

    result, err := config.Db.GetItem(input)
    if err != nil {
        log.Fatalf("Got error calling GetItem in contact: %s", err)
        return
    }

    user := result.Item
    orgid := user["orgid"]

    queryInput := &dynamodb.QueryInput{
		TableName: aws.String(config.UserTable),
		IndexName: aws.String("OrgIdIndex"),
		KeyConditionExpression: aws.String("#orgid = :orgid"),
		ExpressionAttributeNames: map[string]*string{
			"#orgid": aws.String("orgid"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":orgid": {
				N: aws.String(*orgid.N),
			},
		},
	}

	items, err := config.Db.Query(queryInput)
	if err != nil {
		log.Fatalf("Got error calling Query in contact: %s", err)
	}

    contactsItems := items.Items

    var contacts []models.Contact
    for _, v := range contactsItems {
        var contact models.Contact
        contact.Username = *v["username"].S
        contact.Nickname = *v["nickname"].S
        contact.Company = *v["company"].S
        contact.Department = *v["department"].S
        contacts = append(contacts, contact)
    }
    
    c.JSON(200, gin.H{
		"contact": contacts,
	})
}