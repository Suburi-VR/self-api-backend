package models

import "time"

type user struct {
	ID uint
	username string
	secret string
	orgid int
	nickname string
	kana string
	company string
	department string
	CreatedAt time.Time
  UpdatedAt time.Time
}