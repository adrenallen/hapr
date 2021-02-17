package users

import (
	"fmt"
	"time"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

var feedbackMessageActionTags = `
<br>
/label ~"User Feedback"
<br>
/assign me
<br>
/todo
<br>
/confidential
<br>
`

type SupportService struct {
	changelogSeenRepository *UserChangelogSeenRepository
}

func NewSupportService() *SupportService {
	s := new(SupportService)
	s.changelogSeenRepository = new(UserChangelogSeenRepository)
	return s
}

func (s *SupportService) SendFeedbackEmail(feedback string, user *User) error {
	supportEmail, err := helpers.GetConfigValue("support.feedback_email")
	if err != nil {
		return err
	}

	feedback = helpers.CleanStringSpecials(feedback)

	m := helpers.NewEmailMessage(s.createFeedbackSubject(feedback), fmt.Sprintf(`%s - submitted by %s (%s) on %v %s`,
		feedback,
		user.Username,
		user.Email,
		time.Now().Format("2006-01-02 15:04:05"),
		feedbackMessageActionTags), []string{supportEmail})

	return helpers.SendEmailMessage(m)
}

func (s *SupportService) SetUserChangelogSeenVersion(userID int, version string) error {
	ucs, err := s.changelogSeenRepository.GetUserChangelogSeenByUserID(userID)
	if err != nil {
		return err
	}
	if ucs == nil {
		ucs = &UserChangelogSeen{}
		ucs.UserID = userID
		ucs.VersionString = version
	} else {
		ucs.VersionString = version
	}

	_, err = s.changelogSeenRepository.SaveOrUpdateUserChangelogSeen(ucs)
	return err
}

func (s *SupportService) GetUserChangelogSeenVersionByUserID(userID int) (string, error) {
	ucs, err := s.changelogSeenRepository.GetUserChangelogSeenByUserID(userID)
	if err != nil {
		return "", err
	}
	if ucs == nil {
		return "", nil
	}

	return ucs.VersionString, nil
}

func (s *SupportService) AddInterestedPeopleEmail(email string) {
	db := helpers.NewDatabaseConnection()
	
	db.QueryRow("INSERT INTO interested_people_emails (email) VALUES ($1)", email)
}

func (s *SupportService) createFeedbackSubject(feedback string) string {
	const MAX_SUBJECT_LEN = 64
	if len(feedback) > MAX_SUBJECT_LEN {
		feedback = feedback[:MAX_SUBJECT_LEN]
	}
	return feedback
}
