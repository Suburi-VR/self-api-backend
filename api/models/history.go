package models

// History 履歴モデル
type History struct {
	Callid string `json:"callid"`
	Calltime int `json:"calltime"`
	Duration int `json:"duration"`
	Caller map[string]interface{} `json:"caller,omitempty"`
	Receiver map[string]interface{} `json:"receiver,omitempty"`
}