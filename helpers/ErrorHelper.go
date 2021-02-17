package helpers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/garrettcoleallen/happy/dtos"
)

func SerializeAndWriteError(w http.ResponseWriter, err error, statusCode int) bool {
	errorResponse, parseError := json.Marshal(dtos.ErrorDTO{ErrorMessage: err.Error()})

	if statusCode == 0 {
		statusCode = 500
	}

	if parseError != nil {
		w.WriteHeader(500)
		panic(err)
	}
	w.WriteHeader(statusCode)
	w.Write(errorResponse)
	return true
}

func WriteInvalidRequestPayloadResponse(w http.ResponseWriter) {
	errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "Request payload was invalid"})

	w.WriteHeader(403)
	w.Write([]byte(errorResponse))
}
