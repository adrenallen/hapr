package users

import (
	"fmt"
	"strconv"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type PasswordResetService struct {
	Repository *PasswordResetRepository
}

func NewPasswordResetService() *PasswordResetService {
	s := new(PasswordResetService)
	s.Repository = NewPasswordResetRepository()
	return s
}

func (s *PasswordResetService) ResetPasswordForUserByToken(token string, password string) error {
	pwr, err := s.Repository.GetPasswordResetByToken(token)
	if err != nil {
		return err
	}

	userService := NewUserService()
	err = userService.SetUserPassword(pwr.UserID, password)
	if err != nil {
		return err
	}

	return s.Repository.DeletePasswordResetByID(pwr.ID)

}

func (s *PasswordResetService) GenerateResetPasswordForEmail(email string) error {
	userService := NewUserService()
	userObj, err := userService.GetUserByEmail(email)
	if err != nil {
		return err
	}

	newToken := helpers.GenerateForgotPasswordToken(userObj.Username, strconv.Itoa(userObj.ID))

	newPwr := &PasswordReset{
		Token:  newToken,
		UserID: userObj.ID,
	}

	_, err = s.Repository.SaveNewPasswordReset(newPwr)
	if err != nil {
		return err
	}

	url, err := helpers.GetConfigValue("web.url")
	if err != nil {
		return err
	}

	newMessage := helpers.NewEmailMessage(`Forgot your password?`,
		fmt.Sprintf(`Hey %v!<br><br> Looks like you may have forgot your password. No worries, you can <b><a href='%v/password-reset?token=%v'>click here to set a new password!</a></b>
			<br><br>
			If you did not request a password reset, please contact <a href='mailto:support@hapr.io'>support@hapr.io</a>.
			<br><br>
			- The Hapr Team`, userObj.Username, url, newToken),
		[]string{userObj.Email})

	err = helpers.SendEmailMessage(newMessage)
	if err != nil {
		return err
	}

	return nil
}
