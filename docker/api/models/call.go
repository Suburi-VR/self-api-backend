package models

// Call 通話情報
type Call struct {
	CallID string `dynamodbav:"callid"`
	Password string `dynamodbav:"password"`
	Supporter string `dynamodbav:"supporter"`
	Customer string `dynamodbav:"customer"`
	Status int `dynamodbav:"status"`
	Caller string `dynamodbav:"caller"`
	Receiver string `dynamodbav:"receiver"`
}

type AnswerResponse struct {
	Caller string `json:"caller"`
	Nickename string `json:"nickname"`
	Callid string `json:"callid"`
	StartTime int `json:"starttime"`
}