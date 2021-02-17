package users

import (
	"time"
)

type SessionGuid struct {
	ID              int `table:"session_guids"`
	UserId          int
	Guid            string
	CreatedDatetime time.Time
	Active          bool
}
