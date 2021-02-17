package users

import (
	"database/sql"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type UserChangelogSeenRepository struct{}

func (r *UserChangelogSeenRepository) SaveOrUpdateUserChangelogSeen(ucs *UserChangelogSeen) (*UserChangelogSeen, error) {
	if ucs.ID <= 0 {
		return r.saveNewUserChangelogSeen(ucs)
	}

	return r.updateUserChangelogSeen(ucs)
}

func (r *UserChangelogSeenRepository) updateUserChangelogSeen(ucs *UserChangelogSeen) (*UserChangelogSeen, error) {
	db := helpers.NewDatabaseConnection()

	_, err := db.Exec("UPDATE user_changelog_seen SET version_string=$2 WHERE id=$1", ucs.ID, ucs.VersionString)

	if err != nil {
		return nil, err
	}

	return ucs, nil
}

func (r *UserChangelogSeenRepository) saveNewUserChangelogSeen(ucs *UserChangelogSeen) (*UserChangelogSeen, error) {
	db := helpers.NewDatabaseConnection()

	err := db.QueryRow("INSERT INTO user_changelog_seen (user_id, version_string) VALUES ($1, $2) RETURNING id", ucs.UserID, ucs.VersionString).
		Scan(&ucs.ID)
	if err != nil {
		return nil, err
	}

	return ucs, nil
}

func (r *UserChangelogSeenRepository) GetUserChangelogSeenByUserID(userID int) (*UserChangelogSeen, error) {
	db := helpers.NewDatabaseConnection()

	row, err := db.Query("SELECT * FROM user_changelog_seen WHERE user_id=$1", userID)
	defer row.Close()
	if !row.Next() {
		return nil, nil //doesn't have an entry yet
	}
	if err != nil {
		return nil, err
	}

	return r.getUserChangelogSeenFromCurrentRow(row)
}

func (r *UserChangelogSeenRepository) getUserChangelogSeenFromCurrentRow(row *sql.Rows) (*UserChangelogSeen, error) {
	ucs := new(UserChangelogSeen)

	err := row.Scan(&ucs.ID, &ucs.UserID, &ucs.VersionString)

	if err != nil {
		return nil, err
	}

	return ucs, nil
}
