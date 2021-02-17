package reports

import "time"

type ReportRequest struct {
	ID                   int        `table:"report_requests" json:"id" column:"id"`
	UserID               int        `json:"userID" column:"user_id"`
	ReportTypeID         int        `json:"reportTypeID" column:"report_type_id"`
	RequestedDatetime    time.Time  `json:"requestedDatetime" column:"requested_datetime"`
	AdditionalParameters *string    `json:"additionalParameters" column:"additional_parameters"`
	CompletedDatetime    *time.Time `json:"completedDatetime" column:"completed_datetime"`
}
