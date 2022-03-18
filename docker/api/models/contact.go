package models

// Contact 通話情報
type Contact struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Company string `json:"company"`
	Department string `json:"department"`
}