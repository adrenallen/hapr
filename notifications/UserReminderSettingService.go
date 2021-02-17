package notifications

import (
	"errors"
	"log"
)

type UserReminderSettingService struct {
	Repository *UserReminderSettingRepository
}

func NewUserReminderSettingService(userID int) *UserReminderSettingService {
	s := new(UserReminderSettingService)
	s.Repository = &UserReminderSettingRepository{
		UserID:       userID,
		Unrestricted: false,
	}
	return s
}

//Initializes and unrestricted UserReminderSettingService and repository.
//Only for use in non-front facing projects
//TODO -reevaluate this approach
func NewUserReminderSettingServiceUnrestricted(userID int) *UserReminderSettingService {
	if userID > 0 {
		log.Fatal(errors.New(`unrestricted repository cannot be intialized with a user id over 0`))
	}
	s := new(UserReminderSettingService)
	s.Repository = &UserReminderSettingRepository{
		UserID:       userID,
		Unrestricted: true,
	}
	return s
}

func (s *UserReminderSettingService) DeleteByID(id int) error {
	return s.Repository.DeleteByID(id)
}

func (s *UserReminderSettingService) SaveNew(urs UserReminderSetting) (*UserReminderSetting, error) {
	urs.ID = nil
	return s.Repository.SaveNew(&urs)
}

func (s *UserReminderSettingService) Update(urs UserReminderSetting) (*UserReminderSetting, error) {
	return s.Repository.Update(&urs)
}

func (s *UserReminderSettingService) GetAll() ([]*UserReminderSetting, error) {
	return s.Repository.GetAll()
}
