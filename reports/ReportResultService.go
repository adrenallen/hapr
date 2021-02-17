package reports

type ReportResultService struct {
	UserID     int
	Repository *ReportResultRepository
}

func NewReportResultService(userID int) *ReportResultService {
	s := &ReportResultService{UserID: userID,
		Repository: &ReportResultRepository{UserID: userID}}

	return s
}

func (r *ReportResultService) GetResultsForReportType(reportTypeID int) ([]*ReportResult, error) {
	return r.Repository.GetResultsForReportType(reportTypeID)
}

func (r *ReportResultService) NewResult(reportTypeID int, reportRequestID int, result string) (*ReportResult, error) {
	return r.Repository.SaveResult(&ReportResult{
		UserID:          r.UserID,
		ReportTypeID:    reportTypeID,
		ReportRequestID: reportRequestID,
		Result:          result,
	})
}
