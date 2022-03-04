package models

type DbData struct {
	Username string `dynamodbav:"username"`
	Secret string `dynamodbav:"secret"`
	Orgid int `dynamodbav:"orgid"`
	Nickname string `dynamodbav:"nickname"`
	Kana string `dynamodbav:"kana"`
	Company string `dynamodbav:"company"`
	Department string `dynamodbav:"department"`
}
