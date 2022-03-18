package models

// History 履歴モデル
type History struct {
	Callid string `json:"callid"`
	Calltime int `json:"calltime"`
	Duration int `json:"duration"`
	Caller map[string]interface{} `json:"caller,omitempty"`
	Receiver map[string]interface{} `json:"receiver,omitempty"`
}

// Caller Callerモデル
type Caller struct {
	Username string `json:"name"`
	Nickname string `json:"nickname"`
	Kana string `json:"kana"`
	Company string `json:"company"`
	Department string `json:"department"`
}

// Receiver Receiverモデル
type Receiver struct {
	Username string `json:"name"`
	Nickname string `json:"nickname"`
	Kana string `json:"kana"`
	Company string `json:"company"`
	Department string `json:"department"`
}