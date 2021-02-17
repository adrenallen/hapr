package reports

import "time"

type ReportResult struct {
	ID                int       `table:"report_requests" json:"id" column:"id"`
	UserID            int       `json:"userID" column:"user_id"`
	ReportTypeID      int       `json:"reportTypeID" column:"report_type_id"`
	ReportRequestID   int       `json:"reportRequestID" column:"report_request_id"`
	CompletedDatetime time.Time `json:"completedDatetime" column:"completed_datetime"`
	Result            string    `json:"result" column:"result"`
}
