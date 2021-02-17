package notifications

import (
	"database/sql"
	"errors"
	"fmt"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type UserReminderSettingRepository struct {
	UserID       int
	Unrestricted bool
}

func (r *UserReminderSettingRepository) SaveNew(urs *UserReminderSetting) (*UserReminderSetting, error) {
	db := helpers.NewDatabaseConnection()
	
	err := db.QueryRow(`INSERT INTO user_reminder_settings (user_id, reminder_time, sunday, monday, tuesday, wednesday, thursday, friday, saturday, email, push, sms)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
		r.UserID, urs.ReminderTime, urs.Sunday, urs.Monday, urs.Tuesday, urs.Wednesday, urs.Thursday, urs.Friday, urs.Saturday, urs.Email, urs.Push, urs.SMS).
		Scan(&urs.ID)
	if err != nil {
		return nil, err
	}

	return urs, nil
}

func (r *UserReminderSettingRepository) Update(urs *UserReminderSetting) (*UserReminderSetting, error) {
	db := helpers.NewDatabaseConnection()
	
	_, err := db.Query(`UPDATE user_reminder_settings 
		SET 
			reminder_time=$3,
			sunday=$4,
			monday=$5,
			tuesday=$6,
			wednesday=$7,
			thursday=$8,
			friday=$9,
			saturday=$10,
			email=$11,
			push=$12,
			sms=$13
		WHERE id=$1 AND user_id=$2`,
		urs.ID, r.UserID, urs.ReminderTime, urs.Sunday, urs.Monday, urs.Tuesday, urs.Wednesday, urs.Thursday, urs.Friday, urs.Saturday, urs.Email, urs.Push, urs.SMS)
	if err != nil {
		return nil, err
	}

	return urs, nil
}

func (r *UserReminderSettingRepository) GetByID(id int) (*UserReminderSetting, error) {
	db := helpers.NewDatabaseConnection()
	
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM user_reminder_settings d WHERE d.id = $1 AND d.user_id = $2",
		helpers.GetSQLSelectForModelWithTableAlias(UserReminderSetting{}, "d")), id, r.UserID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	rows.Next()

	return r.getFromCurrentRow(rows)
}

func (r *UserReminderSettingRepository) DeleteByID(id int) error {
	db := helpers.NewDatabaseConnection()
	
	_, err := db.Query("DELETE FROM user_reminder_settings where id = $1 and user_id=$2", id, r.UserID)
	return err
}

func (r *UserReminderSettingRepository) GetAll() ([]*UserReminderSetting, error) {
	db := helpers.NewDatabaseConnection()
	
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM user_reminder_settings WHERE user_id = $1 ORDER BY ID",
		helpers.GetSQLSelectForModel(UserReminderSetting{})), r.UserID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return r.getListForAllRows(rows)
}

func (r *UserReminderSettingRepository) GetAllValidForNotification(dayOfWeek int, reminderTime int, email bool, push bool, sms bool) ([]*UserReminderSetting, error) {
	err := r.isUnrestrictedMode()
	if err != nil {
		return nil, err
	}

	whereClause := ``
	switch dayOfWeek {
	case 0:
		whereClause = `sunday=true`
		break
	case 1:
		whereClause = `monday=true`
		break
	case 2:
		whereClause = `tuesday=true`
		break
	case 3:
		whereClause = `wednesday=true`
		break
	case 4:
		whereClause = `thursday=true`
		break
	case 5:
		whereClause = `friday=true`
		break
	case 6:
		whereClause = `saturday=true`
		break
	}

	whereClause = fmt.Sprintf(`%s AND reminder_time=%v`, whereClause, reminderTime)

	if email {
		whereClause = fmt.Sprintf(`%s AND email=true`, whereClause)
	}

	if sms {
		whereClause = fmt.Sprintf(`%s AND sms=true`, whereClause)
	}

	if push {
		whereClause = fmt.Sprintf(`%s AND push=true`, whereClause)
	}

	db := helpers.NewDatabaseConnection()
	
	rows, err := db.Query(
		fmt.Sprintf(`SELECT %s FROM user_reminder_settings 
				WHERE 
					%s
				ORDER BY ID`,
			helpers.GetSQLSelectForModel(UserReminderSetting{}),
			whereClause),
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return r.getListForAllRows(rows)
}

func (r *UserReminderSettingRepository) isUnrestrictedMode() error {
	if !r.Unrestricted {
		return errors.New(`repository is in restricted access mode and does not have access to this method`)
	}
	return nil
}

func (r *UserReminderSettingRepository) getListForAllRows(rows *sql.Rows) ([]*UserReminderSetting, error) {
	list := []*UserReminderSetting{} //so that we return empty not null if none

	for rows.Next() {
		currRow, err := r.getFromCurrentRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, currRow)
	}
	return list, nil
}

func (r *UserReminderSettingRepository) getFromCurrentRow(row *sql.Rows) (*UserReminderSetting, error) {
	urs := new(UserReminderSetting)

	err := row.Scan(&urs.ID, &urs.UserID, &urs.ReminderTime, &urs.Sunday, &urs.Monday, &urs.Tuesday, &urs.Wednesday, &urs.Thursday, &urs.Friday, &urs.Saturday, &urs.Email, &urs.Push, &urs.SMS)

	if err != nil {
		return nil, err
	}

	return urs, nil
}
