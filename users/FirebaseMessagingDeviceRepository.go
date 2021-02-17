package users

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type FirebaseMessagingDeviceRepository struct {
}

func (r *FirebaseMessagingDeviceRepository) SaveNewFMD(fmd *FirebaseMessagingDevice) (*FirebaseMessagingDevice, error) {
	db := helpers.NewDatabaseConnection()

	err := db.QueryRow(`INSERT INTO firebase_messaging_devices (device_id, user_id) VALUES ($1, $2) RETURNING id`,
		fmd.DeviceID, fmd.UserID).
		Scan(&fmd.ID)
	if err != nil {
		return nil, err
	}

	return fmd, nil
}

func (r *FirebaseMessagingDeviceRepository) UpdateFMD(fmd *FirebaseMessagingDevice) (*FirebaseMessagingDevice, error) {
	db := helpers.NewDatabaseConnection()

	_, err := db.Query(`UPDATE firebase_messaging_devices SET device_id=$1, user_id=$2 WHERE id=$3`,
		fmd.DeviceID, fmd.UserID, fmd.ID)
	if err != nil {
		return nil, err
	}

	return fmd, nil
}

func (r *FirebaseMessagingDeviceRepository) GetFMDByDeviceID(deviceID string) (*FirebaseMessagingDevice, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM firebase_messaging_devices d WHERE d.device_id = $1",
		helpers.GetSQLSelectForModelWithTableAlias(FirebaseMessagingDevice{}, "d")), deviceID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	rows.Next()

	return r.getFromCurrentRow(rows)
}

func (r *FirebaseMessagingDeviceRepository) GetFMDsForUserID(userID int) ([]*FirebaseMessagingDevice, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM firebase_messaging_devices d WHERE d.user_id = $1",
		helpers.GetSQLSelectForModelWithTableAlias(FirebaseMessagingDevice{}, "d")), userID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return r.getListForAllRows(rows)
}

func (r *FirebaseMessagingDeviceRepository) GetFMDsForUserIDs(userIDs []int) ([]*FirebaseMessagingDevice, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM firebase_messaging_devices d WHERE d.user_id = any($1)",
		helpers.GetSQLSelectForModelWithTableAlias(FirebaseMessagingDevice{}, "d")), pq.Array(userIDs))
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return r.getListForAllRows(rows)
}

func (r *FirebaseMessagingDeviceRepository) getListForAllRows(rows *sql.Rows) ([]*FirebaseMessagingDevice, error) {
	list := []*FirebaseMessagingDevice{} //so that we return empty not null if none

	for rows.Next() {
		currRow, err := r.getFromCurrentRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, currRow)
	}
	return list, nil
}

func (r *FirebaseMessagingDeviceRepository) getFromCurrentRow(row *sql.Rows) (*FirebaseMessagingDevice, error) {
	fmd := new(FirebaseMessagingDevice)

	err := row.Scan(&fmd.ID, &fmd.DeviceID, &fmd.UserID)

	if err != nil {
		return nil, err
	}

	return fmd, nil
}
