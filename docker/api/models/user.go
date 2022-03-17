package models

// User 現在ログイン中の法人または匿名アカウントユーザー
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