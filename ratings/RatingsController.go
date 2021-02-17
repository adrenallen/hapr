package ratings

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gitlab.com/garrettcoleallen/happy/users"

	"github.com/Jeffail/gabs"
	"gitlab.com/garrettcoleallen/happy/dtos"
	"gitlab.com/garrettcoleallen/happy/helpers"
)

func GetAllRatingsForUser(w http.ResponseWriter, r *http.Request) {
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

	ratingService := NewRatingService(userObj.ID)

	allRatings, err := ratingService.GetAllRatings()
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	} else {
		allRatingsJSON, _ := json.Marshal(AllRatingsDTO{Ratings: allRatings})

		w.WriteHeader(200)
		w.Write(allRatingsJSON)
	}
}

func GetRatingByID(w http.ResponseWriter, r *http.Request) {
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

	ratingID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "ratingID")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	ratingService := NewRatingService(userObj.ID)

	rating, err := ratingService.GetRatingWithFactorsByID(ratingID)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	} else {
		ratingJSON, _ := json.Marshal(rating)

		w.WriteHeader(200)
		w.Write(ratingJSON)
	}
}

func GetAllRatingsWithFactorsForUser(w http.ResponseWriter, r *http.Request) {
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

	ratingService := NewRatingService(userObj.ID)
	ratingsWithFactors, err := ratingService.GetAllRatingsWithFactors()
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	ratingsWithFactorsJSON, err := json.Marshal(ratingsWithFactors)

	//if we had an error serializing return it
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(ratingsWithFactorsJSON)
}

func GetRatingsWithFactorsForUserByDateRange(w http.ResponseWriter, r *http.Request) {
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

	startDateJSON, err := helpers.GetContainerFromRequestPayload(requestJSON, "startdate")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}
	startDate := &dtos.DateDTO{}
	err = startDate.LoadDateFromMap(startDateJSON)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	endDateJSON, err := helpers.GetContainerFromRequestPayload(requestJSON, "enddate")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}
	endDate := &dtos.DateDTO{}
	err = endDate.LoadDateFromMap(endDateJSON)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	ratingService := NewRatingService(userObj.ID)
	ratingsWithFactors, err := ratingService.GetRatingsWithFactorsByDateRange(startDate.ConvertToTime(), endDate.ConvertToTime())

	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	ratingsWithFactorsJSON, err := json.Marshal(ratingsWithFactors)

	//if we had an error serializing return it
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(ratingsWithFactorsJSON)
}

func UpdateRating(w http.ResponseWriter, r *http.Request) {
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

	rating, err := helpers.GetIntValueFromRequestPayload(requestJSON, "rating")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	ratingJournalEntry, err := helpers.GetValueFromRequestPayload(requestJSON, "journalEntry")
	if err != nil {
		ratingJournalEntry = ``
	}

	ratingID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "id")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	ratingService := NewRatingService(userObj.ID)
	if err := ratingService.UpdateRating(ratingID, rating, ratingJournalEntry); err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}

func NewRating(w http.ResponseWriter, r *http.Request) {
	requestJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}
	ratingDate := time.Now()

	dateJSON, err := helpers.GetContainerFromRequestPayload(requestJSON, "date")
	if err == nil && dateJSON != nil {
		date := &dtos.DateDTO{}
		err = date.LoadDateFromMap(dateJSON)
		if err != nil {
			helpers.SerializeAndWriteError(w, err, 403)
			return
		}
		ratingDate = date.ConvertToTime()
	}

	ratingString, err := helpers.GetValueFromRequestPayload(requestJSON, "rating")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}
	rating, err := strconv.Atoi(ratingString)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 400)
		return
	}

	ratingJournalEntry, err := helpers.GetValueFromRequestPayload(requestJSON, "journalEntry")
	if err != nil {
		ratingJournalEntry = ``
	}

	//get user id from session
	guid, err := helpers.GetValueFromRequestPayload(requestJSON, "guid")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}
	userService := users.NewUserService()
	userObj, err := userService.GetUserBySessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	ratingService := NewRatingService(userObj.ID)
	newRating, err := ratingService.SaveNewRating(rating, ratingDate, ratingJournalEntry)

	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	ratingJSON, err := json.Marshal(newRating)

	//if we had an error serializing return it
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(ratingJSON)
}

func DeleteRating(w http.ResponseWriter, r *http.Request) {
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

	ratingID, err := helpers.GetIntValueFromRequestPayload(requestJSON, "ratingID")
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	ratingService := NewRatingService(userObj.ID)
	err = ratingService.DeleteRatingAndRatingFactors(ratingID)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	return

}

func DeleteRatingFactors(w http.ResponseWriter, r *http.Request) {
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

	ratingFactorIDs, err := requestJSON.Path("idList").Children()
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	ratingService := NewRatingService(userObj.ID)

	for _, rfID := range ratingFactorIDs {
		err = ratingService.DeleteRatingFactor(int(rfID.Data().(float64)))
		if err != nil {
			helpers.SerializeAndWriteError(w, err, 500)
			return
		}
	}

	w.WriteHeader(200)
}

func UpdateRatingFactors(w http.ResponseWriter, r *http.Request) {
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

	ratingService := NewRatingService(userObj.ID)

	ratingFactorsContainerArray, err := requestJSON.Path("ratingFactors").Children()
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	for _, ratingFactorContainer := range ratingFactorsContainerArray {
		err := updateRatingFactor(ratingFactorContainer, ratingService)
		if err != nil {
			helpers.SerializeAndWriteError(w, err, 500)
			return
		}
	}

	w.WriteHeader(200)
}

func NewRatingFactors(w http.ResponseWriter, r *http.Request) {
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

	ratingService := NewRatingService(userObj.ID)

	ratingFactorsContainerArray, err := requestJSON.Path("ratingFactors").Children()
	if err != nil {
		helpers.WriteInvalidRequestPayloadResponse(w)
		return
	}

	newRatingFactors := []*RatingFactor{}

	for _, ratingFactorContainer := range ratingFactorsContainerArray {
		newRF, err := saveNewRatingFactor(ratingFactorContainer, ratingService)
		if err != nil {
			helpers.SerializeAndWriteError(w, err, 500)
			return
		}
		newRatingFactors = append(newRatingFactors, newRF)
	}

	newRatingFactorsJSON, err := json.Marshal(newRatingFactors)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(newRatingFactorsJSON)
}

func NewRatingFactor(w http.ResponseWriter, r *http.Request) {
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

	ratingService := NewRatingService(userObj.ID)

	newRatingFactor, err := saveNewRatingFactor(requestJSON, ratingService)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	newRatingFactorJSON, err := json.Marshal(newRatingFactor)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(newRatingFactorJSON)
}

func saveNewRatingFactor(ratingJSON *gabs.Container, ratingService *RatingService) (*RatingFactor, error) {
	ratingIDJSON, err := helpers.GetValueFromRequestPayload(ratingJSON, "ratingID")
	if err != nil {
		return nil, err
	}

	ratingID, err := strconv.Atoi(ratingIDJSON)
	if err != nil {
		return nil, err
	}

	factorIDJSON, err := helpers.GetValueFromRequestPayload(ratingJSON, "factorID")
	if err != nil {
		return nil, err
	}
	factorID, err := strconv.Atoi(factorIDJSON)
	if err != nil {
		return nil, err
	}

	var factorAspectID *int
	factorAspectIDJSON, err := helpers.GetRawValueFromRequestPayload(ratingJSON, "factorAspectID")
	if err == nil && factorAspectIDJSON != nil {
		factorAspectIDConv := int(factorAspectIDJSON.(float64))
		factorAspectID = &factorAspectIDConv
	}

	rankJSON, err := helpers.GetValueFromRequestPayload(ratingJSON, "rank")
	if err != nil {
		return nil, err
	}
	rank, err := strconv.Atoi(rankJSON)
	if err != nil {
		return nil, err
	}

	factorTypeIDJSON, err := helpers.GetValueFromRequestPayload(ratingJSON, "factorTypeID")
	if err != nil {
		return nil, err
	}

	factorTypeID, err := strconv.Atoi(factorTypeIDJSON)
	if err != nil {
		return nil, err
	}

	//TODO - this is inefficient and the whole method should be in the service
	newRatingFactor, err := ratingService.SaveNewRatingFactor(ratingID, factorID, factorTypeID, rank, factorAspectID)
	if err != nil {
		return nil, err
	}

	return newRatingFactor, nil
}

func updateRatingFactor(ratingJSON *gabs.Container, ratingService *RatingService) error {

	ratingFactorID, err := helpers.GetIntValueFromRequestPayload(ratingJSON, "id")
	if err != nil {
		return err
	}

	rank, err := helpers.GetIntValueFromRequestPayload(ratingJSON, "rank")
	if err != nil {
		return err
	}

	var factorAspectID *int
	factorAspectIDJSON, err := helpers.GetRawValueFromRequestPayload(ratingJSON, "factorAspectID")
	if err == nil && factorAspectIDJSON != nil {
		factorAspectIDConv := int(factorAspectIDJSON.(float64))
		factorAspectID = &factorAspectIDConv
	}

	//TODO - this is bad, whole method should be in the service
	return ratingService.UpdateRatingFactor(ratingFactorID, rank, factorAspectID)
}
