package models

// Contact 通話情報
type Contact struct {
	Username string `dynamodbav:"username"`
	Nickname string `dynamodbav:"nickename"`
	Company string `dynamodbav:"company"`
	Department string `dynamodbav:"department"`
}