package users

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type EarlyAccessCodeService struct {
	Repository *EarlyAccessCodeRepository
}

func NewEarlyAccessCodeService() *EarlyAccessCodeService {
	return &EarlyAccessCodeService{
		Repository: &EarlyAccessCodeRepository{},
	}
}

func (s *EarlyAccessCodeService) ClaimCodeForUser(code string, userID int) error {
	eac, err := s.Repository.GetByCodeUnclaimed(code)
	if err != nil {
		return err
	}
	claimedTime := time.Now()
	eac.ClaimedDatetime = &claimedTime
	eac.ClaimedByUserID = &userID

	return s.Repository.Update(eac)
}

//TODO - maybe refactor this code?  It's only for beta so we can nuke it later
//but it's making my eyes hurt :(
//Generates codes and sends to the provided emails
func (s *EarlyAccessCodeService) GenerateAndSendCodes(emails []string) error {
	existingEACs, err := s.Repository.GetAvailableCodes()
	if err != nil {
		return err
	}

	existingCodes := []string{}
	for _, eac := range existingEACs {
		existingCodes = append(existingCodes, eac.Code)
	}

	newCodes := []string{}
	for idx := 0; idx < len(emails); idx++ {

		maxGenerations := 25
		generations := 0

		newCode := s.generateCode()

		for helpers.StringArrayContainsItem(existingCodes, newCode) {

			newCode = s.generateCode()

			//Sanity check for inf loops so we don't crash prod
			generations = generations + 1
			if generations > maxGenerations {
				return fmt.Errorf(`infinite loop detected while generating early access codes, aborting generation`)
			}
		}

		newCodes = append(newCodes, newCode)
		existingCodes = append(existingCodes, newCode)
	}

	for _, code := range newCodes {
		err = s.Repository.SaveNewCode(code)
		if err != nil {
			return err
		}
	}

	//Get our new EAC rows
	availableCodes, err := s.Repository.GetCodesToEmail()
	if err != nil {
		return err
	}

	for idx, email := range emails {

		eac := availableCodes[idx]
		//If this is nil then we d ont have enough codes so error
		if eac == nil {
			return fmt.Errorf(`failed to generate enough early access codes, aborting`)
		}

		err = s.emailCode(eac.Code, email)
		if err != nil {
			return err
		}

		//sent the code so mark it sent
		eac.EmailedTo = &email
		err = s.Repository.Update(eac)
		if err != nil {
			return err
		}
	}

	return nil
}

//Email the code to the email address
func (s *EarlyAccessCodeService) emailCode(code string, email string) error {
	newMessage := helpers.NewEmailMessage(`Welcome to Early Access!`,
		fmt.Sprintf(`Congratulations, you have been selected to participate in Hapr's Early Access program!
		<br><br>
		We have been working hard the last few months to make Hapr the perfect happiness tracking service, but now we need feedback from actual users.  
		We need <b>YOU</b> to tell us what you love, and what you hate about the current interface and functionality.
		<br><br>
		Your early access code is <b>%s</b>
		<br><br>
		You can visit the website <a href='app.hapr.io/signup'>HERE</a> and sign-up using your code!
		<br><br>
		An introductory video can be found <a href='hapr.io/tutorial'>HERE</a>.
		<br><br>
		We look forward to hearing about your experience, and thank you for helping us build something great!
		<br><br>
		- The Hapr Team`, code),
		[]string{email})

	return helpers.SendEmailMessage(newMessage)
}

const EAC_CODE_CHARS = "ABCDEFGHIJKLMNPQRSTUVWXYZ23456789"
const EAC_CODE_LENGTH = 6

func (s *EarlyAccessCodeService) generateCode() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = EAC_CODE_CHARS[rand.Intn(len(EAC_CODE_CHARS))]
	}
	return string(b)
}
