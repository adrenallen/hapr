package users

import (
	"database/sql"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type EarlyAccessCodeRepository struct {
}

func (r *EarlyAccessCodeRepository) Update(eac *EarlyAccessCode) error {
	db := helpers.NewDatabaseConnection()
	_, err := db.Exec(`update early_access_codes set claimed_by_user_id=$1, claimed_datetime=$2, emailed_to=$3 where id = $4`,
		eac.ClaimedByUserID,
		eac.ClaimedDatetime,
		eac.EmailedTo,
		eac.ID)
	return err
}

func (r *EarlyAccessCodeRepository) GetByID(id int) (*EarlyAccessCode, error) {
	db := helpers.NewDatabaseConnection()
	row, err := db.Query(`select * from early_access_codes where id=$1`, id)
	if err != nil {
		return nil, err
	}
	row.Next()
	return r.getFromCurrentRow(row)
}

func (r *EarlyAccessCodeRepository) GetByCodeUnclaimed(code string) (*EarlyAccessCode, error) {
	db := helpers.NewDatabaseConnection()
	row, err := db.Query(`select * from early_access_codes where code ilike $1 and claimed_by_user_id is null`, code)
	if err != nil {
		return nil, err
	}
	row.Next()
	return r.getFromCurrentRow(row)
}

func (r *EarlyAccessCodeRepository) SaveNewCode(code string) error {
	db := helpers.NewDatabaseConnection()
	_, err := db.Exec(`insert into early_access_codes (code) values ($1)`, code)
	return err
}

func (r *EarlyAccessCodeRepository) GetCodesToEmail() ([]*EarlyAccessCode, error) {
	db := helpers.NewDatabaseConnection()
	rows, err := db.Query(`select * from early_access_codes where claimed_by_user_id is null and emailed_to is null`)
	if err != nil {
		return nil, err
	}
	return r.getListForAllRows(rows)
}

func (r *EarlyAccessCodeRepository) GetAvailableCodes() ([]*EarlyAccessCode, error) {
	db := helpers.NewDatabaseConnection()
	rows, err := db.Query(`select * from early_access_codes where claimed_by_user_id is null`)
	if err != nil {
		return nil, err
	}
	return r.getListForAllRows(rows)
}

func (r *EarlyAccessCodeRepository) getListForAllRows(rows *sql.Rows) ([]*EarlyAccessCode, error) {
	list := []*EarlyAccessCode{} //so that we return empty not null if none

	for rows.Next() {
		currRow, err := r.getFromCurrentRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, currRow)
	}
	return list, nil
}

func (r *EarlyAccessCodeRepository) getFromCurrentRow(row *sql.Rows) (*EarlyAccessCode, error) {
	item := new(EarlyAccessCode)

	err := row.Scan(&item.ID, &item.CreatedDatetime, &item.Code, &item.ClaimedByUserID, &item.ClaimedDatetime, &item.EmailedTo)

	if err != nil {
		return nil, err
	}

	return item, nil
}
