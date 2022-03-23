package main

import (
	"api/config"
	"api/models"
	"api/utils"
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"

	"api/services"
)

var lastkey string

func start(c *gin.Context) {

	var body map[string]string
	c.BindJSON(&body)

	var supporter string
	var customer string
	call := models.Call{}
	if (body["calledpartyid"] != "") {
		supporter = body["calledpartyid"]
		customer = userName(c)
		call = models.Call{
			CallID: utils.CallId(),
			Password: utils.Password(),
			Supporter: supporter,
			Customer: customer,
			Status: 1,
			Caller: customer,
			Receiver: supporter,
		}
	} else {
		supporter = userName(c)
		call = models.Call{
			CallID: utils.CallId(),
			Password: utils.Password(),
			Supporter: supporter,
			Customer: "customer",
			Status: 0,
			Caller: "caller",
			Receiver: "receiver",
		}
	}

	/// CallTableにアイテム追加
	av, err := dynamodbattribute.MarshalMap(call)
	if err != nil {
		log.Fatalf("Got error marshalling map in start: %s", err)
	}

	input := &dynamodb.PutItemInput {
		Item: av,
		TableName: &config.CallTable,
	}

	_, err = config.Db.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem in call.go(start1): %s", err)
	}

	c.JSON(200, gin.H{
		"callid": call.CallID,
		"password": call.Password,
	})

	return
}

func answer(c *gin.Context) {

	/// callidとpassword受け取る
	var body map[string]string
	c.BindJSON(&body)

	callid := body["callid"]
	password := body["password"]
	username := userName(c)

	item := services.GetCallItem(callid)

	supporter := item["supporter"].S

	if (password != "") {
		updateCallItemInput := &dynamodb.UpdateItemInput {
			TableName: &config.CallTable,
			Key: map[string]*dynamodb.AttributeValue{
				"callid": {
					S: &callid,
				},
			},
			UpdateExpression: aws.String("set #customer = :customer, #status = :status, #caller = :caller, #receiver = :receiver"),
			ExpressionAttributeNames: map[string]*string {
				"#customer": aws.String("customer"),
				"#status": aws.String("status"),
				"#caller": aws.String("caller"),
				"#receiver": aws.String("receiver"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
				":customer": {
					S: aws.String(username),
				},
				":status": {
					N: aws.String("1"),
				},
				":caller": {
					S: aws.String(username),
				},
				":receiver": {
					S: aws.String(*supporter),
				},

			},
			ReturnValues: aws.String("ALL_NEW"),
			ReturnConsumedCapacity: aws.String("TOTAL"),
			ReturnItemCollectionMetrics: aws.String("SIZE"),
		}
	
		_, err := config.Db.UpdateItem(updateCallItemInput)
		if err != nil {
			log.Fatalf("Got error calling PutItem in call.go(answer): %s", err)
		}
	} else {
		updateCallItemInput := &dynamodb.UpdateItemInput {
			TableName: &config.CallTable,
			Key: map[string]*dynamodb.AttributeValue{
				"callid": {
					S: &callid,
				},
			},
			UpdateExpression: aws.String("set #status = :status, #receiver = :receiver, #calltime = :calltime"),
			ExpressionAttributeNames: map[string]*string {
				"#status": aws.String("status"),
				"#receiver": aws.String("receiver"),
				"#calltime": aws.String("calltime"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
				":status": {
					N: aws.String("2"),
				},
				":receiver": {
					S: aws.String(username),
				},
				":calltime": {
					N: aws.String(strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)),
				},
			},
			ReturnValues: aws.String("ALL_NEW"),
			ReturnConsumedCapacity: aws.String("TOTAL"),
			ReturnItemCollectionMetrics: aws.String("SIZE"),
		}
	
		_, err := config.Db.UpdateItem(updateCallItemInput)
		if err != nil {
			log.Fatalf("Got error calling PutItem in call.go(answer): %s", err)
		}
	}

	c.JSON(200, gin.H{
		"callid": callid,
	})

	return
}

func get(c *gin.Context) {
	/// UserTableから、Autherizationで読み取ったusernameに合致するデータ取得
	username := userName(c)

	queryCallItemInput := &dynamodb.QueryInput{
		TableName: aws.String(config.CallTable),
		IndexName: aws.String("SupporterIndex"),
		KeyConditionExpression: aws.String("#supporter = :supporter"),
		ExpressionAttributeNames: map[string]*string{
			"#supporter": aws.String("supporter"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":supporter": {
				S: aws.String(username),
			},
		},
	}

	callItems, err := config.Db.Query(queryCallItemInput)
	if err != nil {
		log.Fatalf("Got error calling GetItem in get(call.go): %s", err)
	}

	calls := callItems.Items
	var res models.AnswerResponse
	resList := []models.AnswerResponse{}
	for _, v := range calls {
		if (*v["status"].N == *aws.String("1")) {
			/// UserTableを、customer取得したで検索してnicknameを取得する必要あり。
			res.Caller = *v["customer"].S
			res.Nickename = ""
			res.Callid = *v["callid"].S
			res.StartTime = int(time.Now().Unix())
			resList = append(resList, res)
		}
	}

	_, err = json.Marshal(resList)

	c.JSON(200, gin.H{
		"calls": resList,
	})

	return
}

func status(c *gin.Context) {
	var body map[string]string
	c.BindJSON(&body)

	callid := body["callid"]

	/// 受け取ったcallidでCallTableを検索
	item := services.GetCallItem(callid)

	status := item["status"].N
	response, _ := strconv.Atoi(*status)

	c.JSON(200, gin.H{
		"status": response,
	})
}

func end(c *gin.Context) {
	/// callidを受け取る
	var body map[string]string
	c.BindJSON(&body)

	callid := body["callid"]
	callItem := services.GetCallItem(callid)

	var duration int
	var calltime int
	if (callItem["calltime"] == nil) {
		calltime = int(time.Now().UnixNano() / 1000000)
		duration = 0
	} else {
		calltime, _ = strconv.Atoi(*callItem["calltime"].N)
		duration = int(time.Now().UnixNano() / 1000000) - calltime
	}

	input := &dynamodb.UpdateItemInput {
		TableName: aws.String(config.CallTable),
		Key: map[string]*dynamodb.AttributeValue{
				"callid": {
						S: &callid,
				},
		},
		UpdateExpression: aws.String("set #status = :status, #duration = :duration, #calltime = :calltime"),
		ExpressionAttributeNames: map[string]*string {
			"#status": aws.String("status"),
			"#duration": aws.String("duration"),
			"#calltime": aws.String("calltime"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
			":status": {
				N: aws.String("4"),
			},
			":duration": {
				N: aws.String(strconv.Itoa(duration)),
			},
			":calltime": {
				N: aws.String(strconv.Itoa(calltime)),
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
		ReturnConsumedCapacity: aws.String("TOTAL"),
		ReturnItemCollectionMetrics: aws.String("SIZE"),
	}

	_, err := config.Db.UpdateItem(input)
	if err != nil {
			log.Fatalf("Got error calling UpdateItem: %s", err)
	}
}

func history(c *gin.Context) {
	username := userName(c)

	item := services.GetUserItem(username)

	var user models.User
	user.Username = *item["username"].S
	user.Nickname = *item["nickname"].S
	user.Company = *item["company"].S
	user.Department = *item["department"].S
	user.Anonflg = *item["anonflg"].BOOL

	// caller
	queryCallerItemInput := &dynamodb.QueryInput{
		TableName: aws.String(config.CallTable),
		IndexName: aws.String("CallerIndex"),
		KeyConditionExpression: aws.String("#caller = :caller"),
		ExpressionAttributeNames: map[string]*string{
			"#caller": aws.String("caller"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":caller": {
				S: aws.String(username),
			},
		},
	}

	callerItems, err := config.Db.Query(queryCallerItemInput)
	if err != nil {
		log.Fatalf("Got error calling GetItem in QueryInput(caller): %s", err)
	}

	callers := callerItems.Items

	histories := []models.History{}
	for _, v := range callers {
		var historyItem models.History
		historyItem.Callid = *v["callid"].S
		historyItem.Calltime, _ = strconv.Atoi(*v["calltime"].N)
		historyItem.Duration, _ = strconv.Atoi(*v["duration"].N)
		historyItem.Caller = map[string]interface{}{}
		historyItem.Caller["name"] = user.Username
		historyItem.Caller["nickname"] = user.Nickname
		historyItem.Caller["kana"] = user.Kana
		historyItem.Caller["company"] = user.Company
		historyItem.Caller["department"] = user.Department
		historyItem.Caller["anonflg"] = user.Anonflg

		item := services.GetUserItem(*v["receiver"].S)
		var receiver models.User
		receiver.Username = *item["username"].S
		receiver.Nickname = *item["nickname"].S
		receiver.Kana = *item["kana"].S
		receiver.Company = *item["company"].S
		receiver.Department = *item["department"].S
		receiver.Anonflg = *item["anonflg"].BOOL

		historyItem.Receiver = map[string]interface{}{}
		historyItem.Receiver["name"] = receiver.Username
		historyItem.Receiver["nickname"] = receiver.Nickname
		historyItem.Receiver["kana"] = receiver.Kana
		historyItem.Receiver["company"] = receiver.Company
		historyItem.Receiver["department"] = receiver.Department
		historyItem.Receiver["anonflg"] = receiver.Anonflg

		histories = append(histories, historyItem)
	}

	// receiver
	queryReceiverItemInput := &dynamodb.QueryInput{
		TableName: aws.String(config.CallTable),
		IndexName: aws.String("ReceiverIndex"),
		KeyConditionExpression: aws.String("#receiver = :receiver"),
		ExpressionAttributeNames: map[string]*string{
			"#receiver": aws.String("receiver"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":receiver": {
				S: aws.String(username),
			},
		},
	}

	receiverItems, err := config.Db.Query(queryReceiverItemInput)
	if err != nil {
		log.Fatalf("Got error calling GetItem in QueryInput(caller): %s", err)
	}

	receivers := receiverItems.Items

	for _, v := range receivers {
		var historyItem models.History
		historyItem.Callid = *v["callid"].S
		historyItem.Calltime, _ = strconv.Atoi(*v["calltime"].N)
		historyItem.Duration, _ = strconv.Atoi(*v["duration"].N)
		historyItem.Receiver = map[string]interface{}{}
		historyItem.Receiver["name"] = user.Username
		historyItem.Receiver["nickname"] = user.Nickname
		historyItem.Receiver["kana"] = user.Kana
		historyItem.Receiver["company"] = user.Company
		historyItem.Receiver["department"] = user.Department
		historyItem.Receiver["anonflg"] = user.Anonflg

		item := services.GetUserItem(*v["caller"].S)
		var caller models.User
		caller.Username = *item["username"].S
		caller.Nickname = *item["nickname"].S
		caller.Kana = *item["kana"].S
		caller.Company = *item["company"].S
		caller.Department = *item["department"].S
		caller.Anonflg = *item["anonflg"].BOOL

		historyItem.Caller = map[string]interface{}{}
		historyItem.Caller["name"] = caller.Username
		historyItem.Caller["nickname"] = caller.Nickname
		historyItem.Caller["kana"] = caller.Kana
		historyItem.Caller["company"] = caller.Company
		historyItem.Caller["department"] = caller.Department
		historyItem.Caller["anonflg"] = caller.Anonflg

		histories = append(histories, historyItem)
	}

	sort.SliceStable(histories, func(i, j int) bool { return histories[i].Calltime > histories[j].Calltime })

	lastkey = c.Query("lastkey")
	var endIndex int
	var response []models.History

	if (lastkey == "") {
		length := len(histories)
		switch {
		case length == 0:
			return
		case length < 10 && 0 < length:
			endIndex = len(histories)
			response = histories[:endIndex]
			lastkey = response[len(response)-1].Callid
			c.JSON(200, gin.H{
				"histories": response,
			})
		case length >= 10:
			endIndex = 10
			response = histories[:endIndex]
			lastkey = response[len(response)-1].Callid
			c.JSON(200, gin.H{
				"histories": response,
				"lastkey": lastkey,
			})
		}
	} else if (lastkey == histories[len(histories)-1].Callid) {
		return
	} else {
		var index int
		for i, v := range histories {
			if (v.Callid == lastkey) {
				index = i
			}
		}

		if (histories[index].Callid == lastkey && len(histories[index + 1:]) < 10) {
			response := histories[index + 1:]
			lastkey = response[len(response)-1].Callid
			c.JSON(200, gin.H{
				"histories": response,
			})
		} else if (histories[index].Callid == lastkey && len(histories[index + 1:]) >= 10) {
			response := histories[index + 1:index + 10]
			lastkey = response[len(response)-1].Callid
			c.JSON(200, gin.H{
				"histories": response,
				"lastkey": lastkey,
			})
		}
	}
}