package models

import "time"

type dbData struct {
	ID uint
	username map[string]string
	secret map[string]string
	orgid map[string]int
	nickname map[string]string
	kana map[string]string
	company map[string]string
	department map[string]string
	CreatedAt time.Time
  UpdatedAt time.Time
}