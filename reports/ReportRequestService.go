package reports

import (
	"errors"
	"log"
	"time"
)

type ReportRequestService struct {
	UserID     int
	Repository *ReportRequestRepository
}

func NewReportRequestService(userID int) *ReportRequestService {
	s := &ReportRequestService{UserID: userID,
		Repository: &ReportRequestRepository{UserID: userID, Unrestricted: false}}

	return s
}

func NewReportRequestServiceUnrestricted(userID int) *ReportRequestService {
	if userID > 0 {
		log.Fatal(errors.New(`unrestricted repository cannot be intialized with a user id over 0`))
	}
	s := &ReportRequestService{UserID: userID,
		Repository: &ReportRequestRepository{UserID: userID, Unrestricted: true}}

	return s
}

func (r *ReportRequestService) GetPendingRequests() ([]*ReportRequest, error) {
	return r.Repository.GetPendingRequests()
}

func (r *ReportRequestService) GetPendingRequestsForReportType(reportTypeID int) ([]*ReportRequest, error) {
	return r.Repository.GetPendingRequestsForReportType(reportTypeID)
}

func (r *ReportRequestService) MarkRequestComplete(requestID int, completionDate time.Time) error {
	return r.Repository.MarkRequestComplete(requestID, completionDate)
}

func (r *ReportRequestService) NewRequest(reportTypeID int, additionalParameters *string) (*ReportRequest, error) {

	newRequest := &ReportRequest{
		UserID:               r.UserID,
		ReportTypeID:         reportTypeID,
		AdditionalParameters: additionalParameters,
	}

	return r.Repository.SaveNewRequest(newRequest)

}
