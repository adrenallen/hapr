package users

import (
	"encoding/json"
	"net/http"

	"gitlab.com/garrettcoleallen/happy/dtos"
	"gitlab.com/garrettcoleallen/happy/helpers"
)

func SetUserChangelogSeenVersion(w http.ResponseWriter, r *http.Request) {
	requestJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	guid, err := helpers.GetValueFromRequestPayload(requestJSON, "guid")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	userService := NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	version, err := helpers.GetValueFromRequestPayload(requestJSON, "version")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	supportService := NewSupportService()
	err = supportService.SetUserChangelogSeenVersion(userObj.ID, version)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	w.WriteHeader(200)
}

func GetUserChangelogSeenVersion(w http.ResponseWriter, r *http.Request) {
	requestJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	guid, err := helpers.GetValueFromRequestPayload(requestJSON, "guid")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	userService := NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	supportService := NewSupportService()
	version, err := supportService.GetUserChangelogSeenVersionByUserID(userObj.ID)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}
	w.WriteHeader(200)
	versionDTO := UserChangelogSeen{}
	versionDTO.VersionString = version
	versionJSON, err := json.Marshal(versionDTO)

	//if we had an error serializing return it
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}
	w.Write(versionJSON)
}

func SendFeedbackEmail(w http.ResponseWriter, r *http.Request) {
	requestJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		w.WriteHeader(400)
		errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "Invalid format for request"})
		w.Write([]byte(errorResponse))
	}

	guid, err := helpers.GetValueFromRequestPayload(requestJSON, "guid")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	userService := NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	feedback, err := helpers.GetValueFromRequestPayload(requestJSON, "feedback")

	supportService := NewSupportService()
	err = supportService.SendFeedbackEmail(feedback, userObj)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}

func AddInterestedPeopleEmail(w http.ResponseWriter, r *http.Request) {
	userJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		w.WriteHeader(400)
		errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "Invalid format for request"})
		w.Write([]byte(errorResponse))
	}

	email, err := helpers.GetValueFromRequestPayload(userJSON, "email")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	supportService := NewSupportService()
	supportService.AddInterestedPeopleEmail(email)
	w.WriteHeader(200)
}
