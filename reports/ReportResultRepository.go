package reports

import (
	"database/sql"
	"time"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type ReportResultRepository struct {
	UserID int
}

func (r *ReportResultRepository) SaveResult(result *ReportResult) (*ReportResult, error) {
	db := helpers.NewDatabaseConnection()

	err := db.QueryRow("INSERT INTO report_results (user_id, report_type_id, report_request_id, result) VALUES ($1, $2, $3, $4) RETURNING id",
		r.UserID, result.ReportTypeID, result.ReportRequestID, result.Result).
		Scan(&result.ID)
	if err != nil {
		return nil, err
	}

	result.CompletedDatetime = time.Now()

	return result, nil
}

func (r *ReportResultRepository) GetResultsForReportType(reportTypeID int) ([]*ReportResult, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query("SELECT * FROM report_results WHERE user_id = $1 and report_type_id = $2 order by completed_datetime desc", r.UserID, reportTypeID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return r.getListForAllRows(rows)
}

func (r *ReportResultRepository) getListForAllRows(rows *sql.Rows) ([]*ReportResult, error) {
	list := []*ReportResult{} //so that we return empty not null if none

	for rows.Next() {
		currRow, err := r.getFromCurrentRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, currRow)
	}
	return list, nil
}

func (r *ReportResultRepository) getFromCurrentRow(row *sql.Rows) (*ReportResult, error) {
	item := new(ReportResult)

	err := row.Scan(&item.ID, &item.UserID, &item.ReportTypeID, &item.ReportRequestID, &item.CompletedDatetime, &item.Result)

	if err != nil {
		return nil, err
	}

	return item, nil
}
