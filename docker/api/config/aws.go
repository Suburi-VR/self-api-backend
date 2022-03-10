package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var UserTable = "UserTable"
var CallTable = "CallTable"
var Sess = session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}))
	var Db = dynamodb.New(Sess, &aws.Config{Endpoint: aws.String("http://192.168.1.8:8000")})
var CognitoClient = cognitoidentityprovider.New(Sess)