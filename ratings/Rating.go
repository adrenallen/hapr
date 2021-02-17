package ratings

import "time"

type Rating struct {
	ID              int       `table:"ratings" json:"id"`
	UserID          int       `json:"userID"`
	Rating          int       `json:"rating"`
	CreatedDatetime time.Time `json:"createdDatetime"`
	JournalEntry    string    `json:"journalEntry"`
}
