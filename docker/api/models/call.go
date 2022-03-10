package models

type Call struct {
	CallID string `dynamodbav:"callid"`
	Password string `dynamodbav:"password"`
	Status int `dynamodbav:"status"`
	Supporter string `dynamodbav:"supporter"`
	Customer string `dynamodbav:"customer"`
}

type AnswerResponse struct {
		Caller string `json:"caller"`
		Nickename string `json:"nickname"`
		Callid string `json:"callid"`
		StartTime int `json:"starttime"`
}