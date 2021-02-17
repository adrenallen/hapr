package factors

import (
	"encoding/json"
	"net/http"

	"gitlab.com/garrettcoleallen/happy/helpers"
	"gitlab.com/garrettcoleallen/happy/users"
)

func GetAllFactors(w http.ResponseWriter, r *http.Request) {
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

	factorService := NewFactorService(userObj.ID)
	allFactors, err := factorService.GetFactors()
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	} else {
		allFactorsJSON, _ := json.Marshal(AllFactorsDTO{Factors: allFactors})

		w.WriteHeader(200)
		w.Write(allFactorsJSON)
	}
}

func GetAllFactorTypes(w http.ResponseWriter, r *http.Request) {
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

	factorService := NewFactorService(userObj.ID)
	allFactorTypesJSON, _ := json.Marshal(factorService.GetFactorTypes())
	w.WriteHeader(200)
	w.Write(allFactorTypesJSON)
}

func SetFactorArchive(w http.ResponseWriter, r *http.Request) {
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

	factorID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "factorID")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	factorService := NewFactorService(userObj.ID)

	err = factorService.SetArchiveFactor(factorID, archiveSetting)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}

func RenameFactor(w http.ResponseWriter, r *http.Request) {
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

	factorName, err := helpers.GetValueFromRequestPayload(requestJSON, "factor")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	factorID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "factorID")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	factorService := NewFactorService(userObj.ID)
	err = factorService.RenameFactor(factorID, factorName)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}

func NewFactor(w http.ResponseWriter, r *http.Request) {
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

	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	factor, err := helpers.GetValueFromRequestPayload(requestJSON, "factor")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	factorService := NewFactorService(userObj.ID)
	newFactor, err := factorService.SaveNewFactor(factor)

	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	newFactorJSON, _ := json.Marshal(newFactor)

	w.WriteHeader(200)
	w.Write(newFactorJSON)

}
