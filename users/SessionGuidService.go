package users

import (
	"errors"
	"fmt"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type SessionGuidService struct {
	Repository *SessionGuidRepository
}

func NewSessionGuidService() *SessionGuidService {
	s := new(SessionGuidService)
	s.Repository = new(SessionGuidRepository)
	return s
}

func (s *SessionGuidService) VerifySessionGuid(guid string) bool {
	row, err := s.Repository.GetSessionGuidByGuid(guid)

	if err != nil {
		return false
	}

	return row.Active
}

func (s *SessionGuidService) AuthenticateUserAndCreateSessionGuid(user string, password string) (string, error) {
	userService := NewUserService()
	authenticated, userObj, err := userService.AuthenticateUser(user, password)
	if err != nil {
		return "", err
	}

	if !authenticated {
		return "", errors.New("authentication: Failed to authenticate user")
	}

	newGuid := new(SessionGuid)
	newGuid.Active = true
	newGuid.Guid = helpers.GenerateAuthenticationToken()
	newGuid.UserId = userObj.ID

	//Deactivate all the old guids, if error then return and dont insert new
	// if err := repositories.DeactivateAllUserSessionGuids(newGuid.UserId); err != nil{
	// 	return "", err
	// }

	newGuid, err = s.Repository.SaveOrUpdateSessionGuid(newGuid)

	//return the new guid generated
	return newGuid.Guid, err
}

func (s *SessionGuidService) DiableSessionGuid(guid string) error {
	guidObj, err := s.Repository.GetSessionGuidByGuid(guid)
	if err != nil {
		fmt.Printf("Got a bad guid %s, ignoring error since we dont care if the guid doesnt exist to deactivate.\nerr: %v", guid, err)
		return nil
	}
	guidObj.Active = false
	_, err = s.Repository.SaveOrUpdateSessionGuid(guidObj)
	return err

}
