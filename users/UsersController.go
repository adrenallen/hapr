package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/garrettcoleallen/happy/dtos"
	"gitlab.com/garrettcoleallen/happy/helpers"
)

func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	userPassData, err := UserPassDTO{}.ImportPayloadData(r)

	var authenticated bool

	userService := NewUserService()
	if err == nil {
		authenticated, _, err = userService.AuthenticateUser(userPassData.User, userPassData.Password)
	}

	if err != nil || !authenticated {
		w.WriteHeader(401)
		errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "Invalid credentials"})
		w.Write([]byte(errorResponse))
	} else {
		w.WriteHeader(200)
	}
}

func AuthorizeRequest(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//We copy the body reader stream so that we can forward it to the next function
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close() //  must close
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		guidJSON, err := helpers.GetGabsContainerFromBytes(bodyBytes)
		if err != nil {
			invalidGUIDResponse(w)
			return
		}

		guid, ok := guidJSON.Path("guid").Data().(string)
		if !ok {
			invalidGUIDResponse(w)
			return
		}

		guidService := NewSessionGuidService()
		if !guidService.VerifySessionGuid(guid) {
			invalidGUIDResponse(w)
			return
		} else {
			//We are g2g so send to the next route
			f(w, r)
		}
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
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

	userJSON, _ := json.Marshal(userObj)
	w.WriteHeader(200)
	w.Write(userJSON)
}

func invalidGUIDResponse(w http.ResponseWriter) {
	errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "GUID is invalid"})

	//guid is invalid
	w.WriteHeader(403)
	w.Write([]byte(errorResponse))
}

func UserForgotPassword(w http.ResponseWriter, r *http.Request) {
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

	passResetService := NewPasswordResetService()
	passResetService.GenerateResetPasswordForEmail(email)

	// We dont want them to know if it was successful so they can't mine emails
	// if err != nil {
	// 	helpers.SerializeAndWriteError(w, err, 403)
	// 	return
	// }

	w.WriteHeader(200)
}

func UserResetPassword(w http.ResponseWriter, r *http.Request) {
	userJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		w.WriteHeader(400)
		errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "Invalid format for request"})
		w.Write([]byte(errorResponse))
	}

	token, err := helpers.GetValueFromRequestPayload(userJSON, "token")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}
	password, err := helpers.GetValueFromRequestPayload(userJSON, "password")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	passService := NewPasswordResetService()
	err = passService.ResetPasswordForUserByToken(token, password)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	w.WriteHeader(200)

}

func UserSetPhoneNumber(w http.ResponseWriter, r *http.Request) {
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

	number, err := helpers.GetValueFromRequestPayload(requestJSON, "mobileNumber")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	err = userService.SetUserMobileNumber(userObj.ID, number)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	return

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	userJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		w.WriteHeader(400)
		errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "Invalid format for request"})
		w.Write([]byte(errorResponse))
	}

	username, err := helpers.GetValueFromRequestPayload(userJSON, "username")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	email, err := helpers.GetValueFromRequestPayload(userJSON, "email")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	password, err := helpers.GetValueFromRequestPayload(userJSON, "password")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	userService := NewUserService()
	exists, err := userService.CheckIfUsernameExists(username)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	if exists {
		helpers.SerializeAndWriteError(w, fmt.Errorf("username exists already"), 406)
		return
	}

	exists, err = userService.CheckIfEmailExists(email)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	if exists {
		helpers.SerializeAndWriteError(w, fmt.Errorf("email exists already"), 406)
		return
	}

	_, err = userService.CreateNewUser(username, password, email)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
	return
}

func CheckIfUsernameExists(w http.ResponseWriter, r *http.Request) {
	userJSON, err := helpers.GetGabsContainerFromRequest(r)
	if err != nil {
		w.WriteHeader(400)
		errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "Invalid format for request"})
		w.Write([]byte(errorResponse))
	}

	username, err := helpers.GetValueFromRequestPayload(userJSON, "username")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	userService := NewUserService()

	exists, err := userService.CheckIfUsernameExists(username)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	if exists {
		helpers.SerializeAndWriteError(w, fmt.Errorf("username exists already"), 406)
		return
	}

	w.WriteHeader(200)
	return
}

func CheckIfEmailExists(w http.ResponseWriter, r *http.Request) {
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

	userService := NewUserService()
	exists, err := userService.CheckIfEmailExists(email)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	if exists {
		helpers.SerializeAndWriteError(w, fmt.Errorf("email exists already"), 406)
		return
	}

	w.WriteHeader(200)
	return
}

func CreateNewSessionForUser(w http.ResponseWriter, r *http.Request) {
	userPassJSON, err := helpers.GetGabsContainerFromRequest(r)

	var guid string

	if err != nil {
		w.WriteHeader(400)
		errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "Invalid format for request"})
		w.Write([]byte(errorResponse))
	}

	user, hasUser := userPassJSON.Path("user").Data().(string)
	pass, hasPass := userPassJSON.Path("password").Data().(string)

	if !hasUser || !hasPass {
		w.WriteHeader(401)
		errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "Invalid credentials, did not receive a username or password"})
		w.Write([]byte(errorResponse))
		return
	}

	//we didn't error out on pass/user so auth
	guidService := NewSessionGuidService()
	guid, err = guidService.AuthenticateUserAndCreateSessionGuid(user, pass)

	if err != nil {
		w.WriteHeader(401)
		errorResponse, _ := json.Marshal(dtos.ErrorDTO{ErrorMessage: "Invalid credentials", ErrorDetails: err.Error()})
		w.Write([]byte(errorResponse))
	} else {
		guidResponse, _ := json.Marshal(dtos.NewSessionGuidDTO{Guid: guid})

		w.WriteHeader(201)
		w.Write([]byte(guidResponse))
	}
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
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
	guidService := NewSessionGuidService()
	err = guidService.DiableSessionGuid(guid)
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}

func LogDeviceToUser(w http.ResponseWriter, r *http.Request) {
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

	deviceID, err := helpers.GetValueFromRequestPayload(requestJSON, "deviceID")
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 403)
		return
	}

	s := NewFirebaseMessagingDeviceService()
	_, err = s.RecordDeviceToUser(deviceID, userObj.ID)

	//if we had an error serializing return it
	if err != nil {
		helpers.SerializeAndWriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}
