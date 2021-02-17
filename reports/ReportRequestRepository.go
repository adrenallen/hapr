package reports

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type ReportRequestRepository struct {
	UserID       int
	Unrestricted bool
}

func (r *ReportRequestRepository) isUnrestrictedMode() error {
	if !r.Unrestricted {
		return errors.New(`repository is in restricted access mode and does not have access to this method`)
	}
	return nil
}

func (r *ReportRequestRepository) GetPendingRequestsAllUsers() ([]*ReportRequest, error) {
	err := r.isUnrestrictedMode()
	if err != nil {
		return nil, err
	}

	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM report_requests WHERE completed_datetime is null ORDER BY ID",
		helpers.GetSQLSelectForModel(ReportRequest{})))
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return r.getListForAllRows(rows)
}

func (r *ReportRequestRepository) GetPendingRequestsForReportType(reportTypeID int) ([]*ReportRequest, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM report_requests WHERE user_id = $1 and completed_datetime is null and report_type_id=$2 ORDER BY ID",
		helpers.GetSQLSelectForModel(ReportRequest{})), r.UserID, reportTypeID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return r.getListForAllRows(rows)
}

func (r *ReportRequestRepository) GetPendingRequests() ([]*ReportRequest, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM report_requests WHERE user_id = $1 and completed_datetime is null ORDER BY ID",
		helpers.GetSQLSelectForModel(ReportRequest{})), r.UserID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return r.getListForAllRows(rows)
}

func (r *ReportRequestRepository) MarkRequestComplete(requestID int, completionDate time.Time) error {
	db := helpers.NewDatabaseConnection()

	_, err := db.Exec(`UPDATE report_requests set completed_datetime=$1 where id = $2`, completionDate, requestID)

	return err
}

func (r *ReportRequestRepository) SaveNewRequest(request *ReportRequest) (*ReportRequest, error) {
	db := helpers.NewDatabaseConnection()

	request.UserID = r.UserID

	err := db.QueryRow("INSERT INTO report_requests (user_id, report_type_id, additional_parameters) VALUES ($1, $2, $3) RETURNING id",
		request.UserID, request.ReportTypeID, request.AdditionalParameters).
		Scan(&request.ID)
	if err != nil {
		return nil, err
	}

	request.RequestedDatetime = time.Now()

	return request, nil
}

func (r *ReportRequestRepository) getListForAllRows(rows *sql.Rows) ([]*ReportRequest, error) {
	list := []*ReportRequest{} //so that we return empty not null if none

	for rows.Next() {
		currRow, err := r.getFromCurrentRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, currRow)
	}
	return list, nil
}

func (r *ReportRequestRepository) getFromCurrentRow(row *sql.Rows) (*ReportRequest, error) {
	item := new(ReportRequest)

	err := row.Scan(&item.ID, &item.UserID, &item.ReportTypeID, &item.RequestedDatetime, &item.AdditionalParameters, &item.CompletedDatetime)

	if err != nil {
		return nil, err
	}

	return item, nil
}
