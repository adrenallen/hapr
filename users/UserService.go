package users

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"gitlab.com/garrettcoleallen/happy/helpers"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repository *UserRepository
}

func NewUserService() *UserService {
	s := new(UserService)
	s.Repository = &UserRepository{}
	return s
}

func (s *UserService) AuthenticateUser(user string, password string) (bool, *User, error) {
	userObj := s.GetUserByUsername(user)
	if userObj == nil {
		return false, nil, errors.New("invalid user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userObj.Password), []byte(password)); err != nil {
		return false, nil, errors.New("invalid password")
	}

	return true, userObj, nil
}

func (s *UserService) CreateNewUserWithCode(username string, password string, email string, code string) (*User, error) {
	user, err := s.CreateNewUser(username, password, email)
	if err != nil {
		return nil, err
	}
	eacService := NewEarlyAccessCodeService()
	err = eacService.ClaimCodeForUser(code, user.ID)
	return user, err
}

func (s *UserService) CreateNewUser(username string, password string, email string) (*User, error) {
	newUser := User{}

	userRegex := regexp.MustCompile(`[0-9a-zA-Z]{5,32}`)

	if !userRegex.MatchString(username) {
		return nil, fmt.Errorf(`Username does not meet requirements`)
	}

	err := helpers.CheckPasswordRequirements(password)
	if err != nil {
		return nil, err
	}

	newUser.Username = helpers.CleanStringSpecials(username)
	newUser.Email = email

	passwordHash, err := helpers.HashPassword(password)

	if err != nil {
		return nil, err
	}
	newUser.Password = passwordHash
	newUserObj, err := s.Repository.SaveNewUser(&newUser)

	if err != nil {
		return nil, err
	}

	url, err := helpers.GetConfigValue("web.url")
	if err != nil {
		log.Panicf("failed to find url from config %v", err)
	}

	subject := `Welcome to Hapr.io!`
	body := fmt.Sprintf(`Thank you for joining <a href="%v">Hapr.io</a>, your new account <b>%v</b> is ready!<br><br>If you have any questions, just respond to this email.<br><br>- The Hapr Team`, url, newUserObj.Username)

	m := helpers.NewEmailMessage(subject, body, []string{newUserObj.Email})
	err = helpers.SendEmailMessage(m)
	if err != nil {
		log.Panicf("Error sending email %v", err)
	}

	return newUserObj, nil
}

func (s *UserService) CheckIfUsernameExists(username string) (bool, error) {
	return s.Repository.CheckIfUsernameExists(username)
}

func (s *UserService) CheckIfEmailExists(email string) (bool, error) {
	return s.Repository.CheckIfEmailExists(email)
}

func (s *UserService) GetUserByUsername(user string) *User {
	user = strings.ToLower(user)
	userObj, _ := s.Repository.GetUserByUsername(user)

	return userObj
}

func (s *UserService) GetUserBySessionGuid(guid string) (*User, error) {
	return s.Repository.GetUserByActiveSessionGuid(guid)
}

func (s *UserService) SetUserMobileNumber(userID int, mobileNumber string) error {
	userObj, err := s.Repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	userObj.MobileNumber = &mobileNumber
	return s.Repository.UpdateUserMobileNumber(userObj)
}

//This will hash the password as well
func (s *UserService) SetUserPassword(userID int, password string) error {

	err := helpers.CheckPasswordRequirements(password)
	if err != nil {
		return err
	}

	userObj, err := s.Repository.GetUserByID(userID)
	if err != nil {
		return err
	}
	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return err
	}
	userObj.Password = passwordHash
	err = s.Repository.UpdateUserPassword(userObj)
	if err == nil {
		//no error so notify user that password changed
		s.sendPasswordChangedNotificationEmail(userObj.Email)
	}
	return err
}

func (s *UserService) sendPasswordChangedNotificationEmail(email string) {
	m := helpers.NewEmailMessage("Your password was changed!", `Your password at <b>Hapr.io</b> has been changed.  
	<br><br>
	If this was not something you requested, please contact support immediately (<a href='mailto:support@hapr.io'>support@hapr.io</a>)
	<br><br>
	- The Hapr Team`, []string{email})

	err := helpers.SendEmailMessage(m)
	if err != nil {
		log.Panicf("failed to send password update email %v", err)
	}
}

func (s *UserService) GetUserByEmail(email string) (*User, error) {
	return s.Repository.GetUserByEmail(email)
}
