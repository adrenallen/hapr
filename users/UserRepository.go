package users

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type UserRepository struct {
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	db := helpers.NewDatabaseConnection()

	row, err := db.Query("SELECT * FROM users WHERE Email ILIKE $1", email)
	defer row.Close()
	row.Next()
	if err != nil {
		return nil, err
	}

	return r.getUserFromCurrentRow(row)
}

func (r *UserRepository) GetUserByUsername(username string) (*User, error) {
	db := helpers.NewDatabaseConnection()

	row, err := db.Query("SELECT * FROM users WHERE  Username ILIKE $1", username)
	defer row.Close()
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		return nil, fmt.Errorf("the user %v was not found", username)
	}
	return r.getUserFromCurrentRow(row)
}

func (r *UserRepository) GetUsersByIDs(userIDs []int) ([]*User, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query("SELECT * FROM users WHERE id=any($1)", pq.Array(userIDs))
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return r.getListForAllRows(rows)
}

func (r *UserRepository) GetUserByID(userID int) (*User, error) {
	db := helpers.NewDatabaseConnection()

	row, err := db.Query("SELECT * FROM users WHERE id=$1", userID)
	defer row.Close()
	row.Next()
	if err != nil {
		return nil, err
	}

	return r.getUserFromCurrentRow(row)
}

func (r *UserRepository) CheckIfUsernameExists(username string) (bool, error) {
	db := helpers.NewDatabaseConnection()

	row, err := db.Query("SELECT * FROM users WHERE  Username ILIKE $1", username)
	defer row.Close()
	if err != nil {
		return false, err
	}
	return row.Next(), nil
}

func (r *UserRepository) CheckIfEmailExists(email string) (bool, error) {
	db := helpers.NewDatabaseConnection()

	row, err := db.Query("SELECT * FROM users WHERE email ILIKE $1", email)
	defer row.Close()
	if err != nil {
		return false, err
	}
	return row.Next(), nil
}

func (r *UserRepository) SaveNewUser(newUser *User) (*User, error) {
	db := helpers.NewDatabaseConnection()

	err := db.QueryRow("INSERT INTO users (username, password, email, mobile_number) VALUES ($1, $2, $3, $4) RETURNING id", newUser.Username, newUser.Password, newUser.Email, newUser.MobileNumber).
		Scan(&newUser.ID)
	if err != nil {
		return nil, err
	}

	newUser.Password = "" //clear it out we dont need to return this

	return newUser, nil
}

func (r *UserRepository) UpdateUserPassword(user *User) error {
	db := helpers.NewDatabaseConnection()

	_, err := db.Query("UPDATE users set password=$1 where id = $2", user.Password, user.ID)
	return err
}

func (r *UserRepository) UpdateUserMobileNumber(user *User) error {
	db := helpers.NewDatabaseConnection()

	_, err := db.Query("UPDATE users set mobile_number=$1 where id = $2", user.MobileNumber, user.ID)
	return err
}

func (r *UserRepository) GetUserByActiveSessionGuid(guid string) (*User, error) {
	db := helpers.NewDatabaseConnection()

	row, err := db.Query("SELECT users.* FROM users LEFT JOIN session_guids ON session_guids.user_id = users.id and session_guids.active=true WHERE guid=$1", guid)
	defer row.Close()
	row.Next()
	if err != nil {
		return nil, err
	}

	return r.getUserFromCurrentRow(row)
}

func (r *UserRepository) getListForAllRows(rows *sql.Rows) ([]*User, error) {
	list := []*User{} //so that we return empty not null if none

	for rows.Next() {
		currRow, err := r.getUserFromCurrentRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, currRow)
	}
	return list, nil
}

func (r *UserRepository) getUserFromCurrentRow(row *sql.Rows) (*User, error) {
	user := new(User)

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.MobileNumber)

	if err != nil {
		return nil, err
	}

	return user, nil
}
