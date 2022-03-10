package models

type User struct {
	Username string `dynamodbav:"username"`
	Secret string `dynamodbav:"secret"`
	Orgid int `dynamodbav:"orgid"`
	Nickname string `dynamodbav:"nickname"`
	Kana string `dynamodbav:"kana"`
	DeviceToken string `dynamodbav:"deviceToken"`
	Platform string `dynamodbav:"platform"`
	Company string `dynamodbav:"company"`
	Department string `dynamodbav:"department"`
}