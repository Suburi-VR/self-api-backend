package main

import (
	"api/config"
	"api/models"
	"api/utils"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
)

var callData models.Call
var callDatas []models.Call

func start(c *gin.Context) {

	username := userName(c)
	callid := utils.CallId()
	password := utils.Password()
	supporter := username

	callData.CallID = callid
	callData.Password = password
	callData.Supporter = supporter
	callData.Customer = "customer"

	callDatas = append(callDatas, callData)

	/// CallTableにアイテム追加
	callDataAv, err := dynamodbattribute.MarshalMap(callData)
	if err != nil {
		log.Fatalf("Got error marshalling map: %s", err)
	}

	callDataParams := &dynamodb.PutItemInput {
		Item: callDataAv,
		TableName: &config.CallTable,
	}

	_, err = config.Db.PutItem(callDataParams)
	if err != nil {
		log.Fatalf("Got error calling PutItem in call.go(start1): %s", err)
	}

	/// UserTableのcallDatas(リスト)更新
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
	
	testDatas := result.Item

	var av []*dynamodb.AttributeValue

	if (testDatas["callDatas"] == nil) {
		av, _ = dynamodbattribute.MarshalList(callDatas)
	} else {
		additionalCallData, _ := dynamodbattribute.Marshal(callData)
		av = append(testDatas["callDatas"].L, additionalCallData)
	}

	if err != nil {
		log.Fatalf("Got error marshalling map: %s", err)
	}

	callDatasParams := &dynamodb.UpdateItemInput {
		TableName: &config.UserTable,
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: &username,
			},
		},
		UpdateExpression: aws.String("set #callDatas = :callDatas"),
		ExpressionAttributeNames: map[string]*string {
			"#callDatas": aws.String("callDatas"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
			":callDatas": {
				L: av,
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
		ReturnConsumedCapacity: aws.String("TOTAL"),
		ReturnItemCollectionMetrics: aws.String("SIZE"),
	}

	_, err = config.Db.UpdateItem(callDatasParams)
	if err != nil {
		log.Fatalf("Got error calling PutItem in call.go(start2): %s", err)
	}

	c.JSON(200, gin.H{
		"callid": callid,
		"password": password,
	})

	return
}

func answer(c *gin.Context) {

	customer := userName(c)

	/// callidとpassword受け取る
	var body map[string]string
	c.BindJSON(&body)

	callid := body["callid"]
	password := body["password"]

	/// 受け取ったcallidとpasswordでCallTableを検索
	getCallItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(config.CallTable),
		Key: map[string]*dynamodb.AttributeValue{
			"callid": {
				S: &callid,
			},
			"password": {
				S: &password,
			},
		},
	}

	callItem, err := config.Db.GetItem(getCallItemInput)
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
		return
	}

	supporter := callItem.Item["supporter"].S

	/// 取得したsupporter(username)でUserTable内のアイテム特定

	getUserItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(config.UserTable),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: supporter,
			},
		},
	}

	userItem, err := config.Db.GetItem(getUserItemInput)
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
		return
	}

	previousCallDatas := userItem.Item["callDatas"].L
	
	for i, v := range previousCallDatas {
		if (*v.M["callid"].S == callid) {
			*previousCallDatas[i].M["status"].N = *aws.String("1")
			*previousCallDatas[i].M["customer"].S = *aws.String(customer)
		}
	}

	callDatasParams := &dynamodb.UpdateItemInput {
		TableName: &config.UserTable,
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: supporter,
			},
		},
		UpdateExpression: aws.String("set #callDatas = :callDatas"),
		ExpressionAttributeNames: map[string]*string {
			"#callDatas": aws.String("callDatas"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
			":callDatas": {
					L: previousCallDatas,
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
		ReturnConsumedCapacity: aws.String("TOTAL"),
		ReturnItemCollectionMetrics: aws.String("SIZE"),
	}

	_, err = config.Db.UpdateItem(callDatasParams)
	if err != nil {
		log.Fatalf("Got error calling PutItem in call.go(answer): %s", err)
	}
}

func get(c *gin.Context) {
	/// UserTableから、Autherizationで読み取ったusernameに合致するデータ取得
	username := userName(c)
	getUserItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(config.UserTable),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: &username,
			},
		},
	}

	userItem, err := config.Db.GetItem(getUserItemInput)
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
		return
	}
	
	getCallData := userItem.Item["callDatas"].L

	var res models.AnswerResponse
	var resList []models.AnswerResponse
	
	for i := range getCallData {
		if (*getCallData[i].M["status"].N == *aws.String("1")) {
			res.Caller = *getCallData[i].M["customer"].S
			res.Nickename = *userItem.Item["nickname"].S
			res.Callid = *getCallData[i].M["callid"].S
			res.StartTime = int(time.Now().Unix())
			resList = append(resList, res)
		}
	}

	_, err = json.Marshal(resList)

	c.JSON(200, gin.H{
		"calls": resList,
	})

	return

	// /// UserTableはcallidとpasswordのmapのリストを持っており、そこからすべてのcallidとpasswordを取得
	// /// (はじめは空。startするごとにappend?していく)
	// /// すべてのcallidとpasswordで、CallTableのデータを取得し、statusが１だったらレスポンスを返す


}