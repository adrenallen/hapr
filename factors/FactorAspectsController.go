package factors

import (
	"encoding/json"
	"net/http"

	"gitlab.com/garrettcoleallen/happy/helpers"
	"gitlab.com/garrettcoleallen/happy/users"
)

func GetAllFactorAspects(w http.ResponseWriter, r *http.Request) {
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

	aspectService := NewFactorAspectService(userObj.ID)
	fas, err := aspectService.GetFactorAspects()
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	} else {
		fasJSON, _ := json.Marshal(fas)

		w.WriteHeader(200)
		w.Write(fasJSON)
	}
}

func SetFactorAspectArchive(w http.ResponseWriter, r *http.Request) {
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

	userService := users.NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	archiveSetting, err := helpers.GetBoolValueFromRequestPayload(requestJSON, "archive")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	faID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "factorAspectID")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}
	aspectService := NewFactorAspectService(userObj.ID)

	err = aspectService.SetArchiveFactorAspect(faID, archiveSetting)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}

func RenameFactorAspect(w http.ResponseWriter, r *http.Request) {
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

	userService := users.NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	faName, err := helpers.GetValueFromRequestPayload(requestJSON, "factorAspect")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	faID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "factorAspectID")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	aspectService := NewFactorAspectService(userObj.ID)

	err = aspectService.RenameFactorAspect(faID, faName)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}

func NewFactorAspect(w http.ResponseWriter, r *http.Request) {
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

	userService := users.NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	faName, err := helpers.GetValueFromRequestPayload(requestJSON, "factorAspect")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	factorID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "factorID")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	aspectService := NewFactorAspectService(userObj.ID)
	newFA, err := aspectService.SaveNewFactorAspect(faName, factorID)

	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	newFAJSON, _ := json.Marshal(newFA)

	w.WriteHeader(200)
	w.Write(newFAJSON)

}
