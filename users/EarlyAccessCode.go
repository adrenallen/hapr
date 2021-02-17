package users

import (
	"time"
)

type EarlyAccessCode struct {
	ID              int        `table:"early_access_codes" json:"id" column:"id"`
	CreatedDatetime time.Time  `column:"created_datetime" json:"createdDatetime"`
	Code            string     `column:"code" json:"code"`
	ClaimedByUserID *int       `column:"claimed_by_user_id" json:"claimedByUserID"`
	ClaimedDatetime *time.Time `column:"claimed_datetime" json:"claimedDatetime"`
	EmailedTo       *string    `column:"emailed_to" json:"emailedTo"`
}
