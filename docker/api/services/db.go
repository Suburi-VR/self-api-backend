package services

import (
	"api/config"
	"api/models"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

// PutUserItem UserTableのPutItem
func PutUserItem(user models.User) {
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Fatalf("Got error marshalling map in PutUserItem: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(config.UserTable),
	}
	
	_, err = config.Db.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	return
}

// GetCallItem CallTableのGetItem
