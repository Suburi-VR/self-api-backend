package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Sess aws session
var Sess = session.Must(session.NewSessionWithOptions(session.Options{
	Config: aws.Config{
		Region: aws.String("ap-northeast-1"),
		CredentialsChainVerboseErrors: aws.Bool(true),
	},
}))

// Db dynamodb
var Db = dynamodb.New(Sess, &aws.Config{Endpoint: aws.String("http://192.168.1.8:8000")})

// UserTable dynamodb UserTable
var UserTable = "UserTable"

// CallTable dynamodb CallTable
var CallTable = "CallTable"

// CognitoClient CognitoIdentityProvider
var CognitoClient = cognitoidentityprovider.New(Sess)

// CognitoUserPoolID CognitoUserPoolID
var CognitoUserPoolID = "ap-northeast-1_Kjb4vUZPh"