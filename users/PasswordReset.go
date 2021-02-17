package users

import "time"

type PasswordReset struct {
	ID              int `table:"pasword_resets" json:"id"`
	UserID          int `json:"userID"`
	CreatedDatetime time.Time
	Token           string
}
