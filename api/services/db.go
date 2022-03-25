package services

import (
	"api/config"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// GetUserItem UserTableのGetItem
func GetUserItem(username string) (map[string]*dynamodb.AttributeValue) {
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
		log.Fatalf("Got error calling GetUserItem: %s", err)
	}

	userItem := result.Item

	return userItem
}

// GetCallItem CallTableのGetItem
func GetCallItem(callid string) (map[string]*dynamodb.AttributeValue) {
	Input := &dynamodb.GetItemInput{
		TableName: aws.String(config.CallTable),
		Key: map[string]*dynamodb.AttributeValue{
			"callid": {
				S: &callid,
			},
		},
	}

	result, err := config.Db.GetItem(Input)
	if err != nil {
		log.Fatalf("Got error calling GetCallItem: %s", err)
	}

	callItem := result.Item

	return callItem
}

