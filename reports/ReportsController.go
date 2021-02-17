package reports

import (
	"encoding/json"
	"net/http"

	"gitlab.com/garrettcoleallen/happy/helpers"
	"gitlab.com/garrettcoleallen/happy/users"
)

func GetPendingRequestsForReport(w http.ResponseWriter, r *http.Request) {
	requestJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	guid, err := helpers.GetValueFromRequestPayload(requestJSON, "guid")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	userService := users.NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	reportTypeID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "reportTypeID")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	service := NewReportRequestService(userObj.ID)
	pending, err := service.GetPendingRequestsForReportType(reportTypeID)
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	resultJSON, err := json.Marshal(pending)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(resultJSON)
}

func RequestReport(w http.ResponseWriter, r *http.Request) {
	requestJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	guid, err := helpers.GetValueFromRequestPayload(requestJSON, "guid")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	userService := users.NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	reportTypeID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "reportTypeID")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	var additionalParams *string
	//ignore error here, we may have reports with no additional params so this will be null
	reportParamsFromRequest, err := helpers.GetValueFromRequestPayload(requestJSON, "additionalParameters")

	//if we have no error then we found a value so assign it
	if err == nil {
		additionalParams = &reportParamsFromRequest
	}

	requestService := NewReportRequestService(userObj.ID)
	_, err = requestService.NewRequest(reportTypeID, additionalParams)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)

}

func GetReportResultsForReportType(w http.ResponseWriter, r *http.Request) {
	requestJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	guid, err := helpers.GetValueFromRequestPayload(requestJSON, "guid")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	userService := users.NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	reportTypeID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "reportTypeID")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	rs := NewReportResultService(userObj.ID)
	results, err := rs.GetResultsForReportType(reportTypeID)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	resultJSON, err := json.Marshal(results)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(resultJSON)
}

func GetReportFactorToRatingOccurrences(w http.ResponseWriter, r *http.Request) {
	requestJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	guid, err := helpers.GetValueFromRequestPayload(requestJSON, "guid")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	userService := users.NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	reportService := NewReportService(userObj.ID)
	result, err := reportService.GetFactorToRatingOccurrences()
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(resultJSON)
}

func GetFullDataExport(w http.ResponseWriter, r *http.Request) {
	requestJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	guid, err := helpers.GetValueFromRequestPayload(requestJSON, "guid")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	userService := users.NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	s := NewReportService(userObj.ID)
	exData, err := s.GetFullDataExport()
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	resultJSON, err := json.Marshal(exData)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=export.json")
	w.WriteHeader(200)
	w.Write(resultJSON)
}
