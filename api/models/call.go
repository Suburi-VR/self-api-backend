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

// Caller 発信者モデル
type Caller struct {
	Name string `json:"name"`
	Nickname string `json:"nickname"`
	Kana string `json:"kana"`
	Company string `json:"company"`
	Department string `json:"department"`
	Anonflg bool `json:"anonflg"`
}

// Receiver 着信者モデル
type Receiver struct {
	Name string `json:"name"`
	Nickname string `json:"nickname"`
	Kana string `json:"kana"`
	Company string `json:"company"`
	Department string `json:"department"`
	Anonflg bool `json:"anonflg"`
}

// AnswerResponse /call/answerしたときのresponse
type AnswerResponse struct {
	Caller string `json:"caller"`
	Nickename string `json:"nickname"`
	Callid string `json:"callid"`
	StartTime int `json:"starttime"`
}