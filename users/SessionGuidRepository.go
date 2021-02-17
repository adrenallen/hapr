package users

import (
	"database/sql"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type SessionGuidRepository struct {
}

func (r *SessionGuidRepository) GetSessionGuidByGuid(guid string) (*SessionGuid, error) {
	db := helpers.NewDatabaseConnection()

	row, err := db.Query("SELECT * FROM session_guids WHERE guid  = $1", guid)
	defer row.Close()
	if err != nil {
		return nil, err
	}

	row.Next()
	return r.getSessionGuidFromCurrentRow(row)
}

func (r *SessionGuidRepository) DeactivateAllUserSessionGuids(userId int) error {
	db := helpers.NewDatabaseConnection()

	_, err := db.Query("UPDATE session_guids SET active=false WHERE user_id = $1", userId)
	return err
}

func (r *SessionGuidRepository) SaveOrUpdateSessionGuid(guid *SessionGuid) (*SessionGuid, error) {
	if guid.ID <= 0 {
		return r.saveNewSessionGuid(guid)
	}

	return r.updateSessionGuid(guid)
}

func (r *SessionGuidRepository) updateSessionGuid(guid *SessionGuid) (*SessionGuid, error) {
	db := helpers.NewDatabaseConnection()

	_, err := db.Exec("UPDATE session_guids SET user_id=$2, active=$3 WHERE id=$1", guid.ID, guid.UserId, guid.Active)

	if err != nil {
		return nil, err
	}

	return guid, nil
}

func (r *SessionGuidRepository) saveNewSessionGuid(guid *SessionGuid) (*SessionGuid, error) {
	db := helpers.NewDatabaseConnection()

	err := db.QueryRow("INSERT INTO session_guids (user_id, guid, active) VALUES ($1, $2, $3) RETURNING id, created_datetime", guid.UserId, guid.Guid, guid.Active).
		Scan(&guid.ID, &guid.CreatedDatetime)
	if err != nil {
		return nil, err
	}

	return guid, nil
}

func (r *SessionGuidRepository) getSessionGuidFromCurrentRow(row *sql.Rows) (*SessionGuid, error) {
	guid := new(SessionGuid)

	err := row.Scan(&guid.ID, &guid.UserId, &guid.Guid, &guid.CreatedDatetime, &guid.Active)

	if err != nil {
		return nil, err
	}

	return guid, nil
}
