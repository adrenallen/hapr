package notifications

import (
	"encoding/json"
	"net/http"

	"github.com/Jeffail/gabs"
	"gitlab.com/garrettcoleallen/happy/helpers"
	"gitlab.com/garrettcoleallen/happy/users"
)

func UpdateUserReminderSettings(w http.ResponseWriter, r *http.Request) {
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

	settingsToUpdate := []UserReminderSetting{}

	settings, _ := helpers.GetContainerFromRequestPayload(requestJSON, "settings")
	settingsChildren, _ := settings.Children()
	for _, child := range settingsChildren {
		settingsToUpdate = append(settingsToUpdate, convertContainerToUserReminderSetting(child))
	}

	reminderService := NewUserReminderSettingService(userObj.ID)
	for _, setting := range settingsToUpdate {
		if *setting.ID < 1 {
			_, err := reminderService.SaveNew(setting)
			if err != nil {
				helpers.SerializeAndWriteError(w, err, 500)
				return
			}
		} else {
			_, err := reminderService.Update(setting)
			if err != nil {
				helpers.SerializeAndWriteError(w, err, 500)
				return
			}
		}
	}

	w.WriteHeader(200)

}

//TODO - remove this monstrosity and use reflection properly
//Since we have an array this isn't super ideal, we do  this a lot
//throughout the project and it needs handled better
func convertContainerToUserReminderSetting(c *gabs.Container) UserReminderSetting {
	urs := UserReminderSetting{}

	//This has to happen due to weirdness with float64 instead of int coming from json parses
	id, _ := helpers.GetIntValueFromRequestPayload(c, `id`)
	urs.ID = &id

	rt, _ := helpers.GetIntValueFromRequestPayload(c, `reminderTime`)
	urs.ReminderTime = rt

	urs.Sunday = c.Path(`sunday`).Data().(bool)
	urs.Monday = c.Path(`monday`).Data().(bool)
	urs.Tuesday = c.Path(`tuesday`).Data().(bool)
	urs.Wednesday = c.Path(`wednesday`).Data().(bool)
	urs.Thursday = c.Path(`thursday`).Data().(bool)
	urs.Friday = c.Path(`friday`).Data().(bool)
	urs.Saturday = c.Path(`saturday`).Data().(bool)

	urs.Push = c.Path(`push`).Data().(bool)
	urs.SMS = c.Path(`sms`).Data().(bool)
	urs.Email = c.Path(`email`).Data().(bool)

	return urs
}

func DeleteUserReminderSettings(w http.ResponseWriter, r *http.Request) {
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

	settingsToDelete := []UserReminderSetting{}

	settings, _ := helpers.GetContainerFromRequestPayload(requestJSON, "settings")
	settingsChildren, _ := settings.Children()
	for _, child := range settingsChildren {
		settingsToDelete = append(settingsToDelete, convertContainerToUserReminderSetting(child))
	}

	reminderService := NewUserReminderSettingService(userObj.ID)
	for _, setting := range settingsToDelete {
		err := reminderService.DeleteByID(*setting.ID)
		if err != nil {
			helpers.SerializeAndWriteError(w, err, 500)
			return
		}
	}

	w.WriteHeader(200)

}

func GetAllUserReminderSettings(w http.ResponseWriter, r *http.Request) {
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

	reminderService := NewUserReminderSettingService(userObj.ID)
	reminders, err := reminderService.GetAll()

	//if we had an error serializing return it
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	remindersJSON, _ := json.Marshal(reminders)

	w.Write([]byte(remindersJSON))
	w.WriteHeader(200)
}
