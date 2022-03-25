package models

// History 履歴モデル
type History struct {
	Callid string `json:"callid"`
	Calltime int `json:"calltime"`
	Duration int `json:"duration"`
	Caller Caller `json:"caller,omitempty"`
	Receiver Receiver `json:"receiver,omitempty"`
}