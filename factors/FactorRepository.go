package factors

import (
	"database/sql"
	"fmt"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type FactorRepository struct {
	UserID int
}

func (r *FactorRepository) GetFactors() ([]*Factor, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM factors WHERE user_id = $1 ORDER BY factor", helpers.GetSQLSelectForModel(Factor{})), r.UserID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return r.getFactorsListForAllRows(rows)
}

func (r *FactorRepository) GetFactorByID(factorID int) (*Factor, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM factors WHERE ID = $1 and user_id=$2", helpers.GetSQLSelectForModel(Factor{})), factorID, r.UserID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	rows.Next()

	return r.getFactorFromCurrentRow(rows)
}

func (r *FactorRepository) GetAllFactorTypes() []*FactorType {
	return []*FactorType{
		&FactorType{
			ID:         PositiveFactorType,
			FactorType: "Positive",
		},
		&FactorType{
			ID:         NegativeFactorType,
			FactorType: "Negative",
		},
	}
}

func (r *FactorRepository) SaveFactor(factor *Factor) error {
	db := helpers.NewDatabaseConnection()

	factor.Factor = helpers.CleanStringSpecials(factor.Factor)

	_, err := db.Query(fmt.Sprintf("UPDATE factors SET factor=%s, archived=$2 WHERE id = $3 and user_id=$4",
		helpers.GetEncryptSQLString(`($1)`)),
		factor.Factor, factor.Archived, factor.ID, r.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *FactorRepository) SaveNewFactor(newFactor *Factor) (*Factor, error) {
	db := helpers.NewDatabaseConnection()

	newFactor.Factor = helpers.CleanStringSpecials(newFactor.Factor)

	err := db.QueryRow(fmt.Sprintf("INSERT INTO factors (user_id, factor) VALUES ($1, %s) RETURNING id",
		helpers.GetEncryptSQLString(`($2)`)),
		r.UserID, newFactor.Factor).
		Scan(&newFactor.ID)
	if err != nil {
		return nil, err
	}

	return newFactor, nil
}

func (r *FactorRepository) getFactorsListForAllRows(rows *sql.Rows) ([]*Factor, error) {
	factorsList := []*Factor{} //so that we return empty not null if none

	for rows.Next() {
		currRow, err := r.getFactorFromCurrentRow(rows)
		if err != nil {
			return nil, err
		}
		factorsList = append(factorsList, currRow)
	}
	return factorsList, nil
}

func (r *FactorRepository) getFactorFromCurrentRow(row *sql.Rows) (*Factor, error) {
	factor := new(Factor)

	err := row.Scan(&factor.ID, &factor.UserID, &factor.Factor, &factor.Archived)

	if err != nil {
		return nil, err
	}

	return factor, nil
}
