package users

import (
	"database/sql"
	"fmt"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type PasswordResetRepository struct {
	passwordResetMaxAgeDays int
}

func NewPasswordResetRepository() *PasswordResetRepository {
	var passwordResetMaxAgeDays = 1 //how many days old we look at
	return &PasswordResetRepository{passwordResetMaxAgeDays: passwordResetMaxAgeDays}
}

func (r *PasswordResetRepository) SaveNewPasswordReset(pwr *PasswordReset) (*PasswordReset, error) {
	db := helpers.NewDatabaseConnection()

	err := db.QueryRow(`INSERT INTO password_resets (user_id, token) VALUES ($1, $2) RETURNING id, created_datetime`,
		pwr.UserID, pwr.Token).
		Scan(&pwr.ID, &pwr.CreatedDatetime)
	if err != nil {
		return nil, err
	}

	return pwr, nil
}

func (r *PasswordResetRepository) DeletePasswordResetByID(pwrID int) error {
	db := helpers.NewDatabaseConnection()

	_, err := db.Query("DELETE FROM password_resets where id = $1", pwrID)
	return err
}

func (r *PasswordResetRepository) GetPasswordResetByToken(token string) (*PasswordReset, error) {
	db := helpers.NewDatabaseConnection()

	row, err := db.Query(fmt.Sprintf("SELECT * FROM password_resets WHERE token=$1 and created_datetime >= NOW() - INTERVAL '%v DAY'", r.passwordResetMaxAgeDays),
		token)
	defer row.Close()
	row.Next()
	if err != nil {
		return nil, err
	}

	return r.getPasswordResetFromCurrentRow(row)
}

func (r *PasswordResetRepository) getPasswordResetFromCurrentRow(row *sql.Rows) (*PasswordReset, error) {
	pwr := new(PasswordReset)

	err := row.Scan(&pwr.ID, &pwr.UserID, &pwr.CreatedDatetime, &pwr.Token)

	if err != nil {
		return nil, err
	}

	return pwr, nil
}
