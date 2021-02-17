package helpers

import (
	"fmt"
	"log"
	"regexp"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAuthenticationToken() string {
	u1 := uuid.Must(uuid.NewV4())

	return u1.String()
}

func GenerateForgotPasswordToken(username string, userID string) string {
	uu, err := uuid.NewV1()
	if err != nil {
		log.Fatalf("failed to generate forgot password token %v", err)
	}

	u1 := uuid.Must(uuid.NewV3(uu, username+userID), nil)
	return u1.String()
}

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func CheckPasswordRequirements(password string) error {
	passwordRegex := regexp.MustCompile(`.{6,64}`)

	if !passwordRegex.MatchString(password) {
		return fmt.Errorf(`Password does not meet requirements`)
	}
	return nil
}
