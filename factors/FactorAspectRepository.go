package factors

import (
	"database/sql"
	"fmt"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type FactorAspectRepository struct {
	UserID int
}

func (r *FactorAspectRepository) GetFactorAspects() ([]*FactorAspect, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM factor_aspects fa left join factors f on f.id=fa.factor_id WHERE f.user_id = $1 ORDER BY factor_aspect",
		helpers.GetSQLSelectForModelWithTableAlias(FactorAspect{}, "fa")), r.UserID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return r.getFactorAspectListForAllRows(rows)
}

func (r *FactorAspectRepository) GetFactorAspectByID(faID int) (*FactorAspect, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM factor_aspects fa left join factors f on f.id=fa.factor_id WHERE f.user_id = $1 AND fa.ID = $2 ORDER BY factor_aspect",
		helpers.GetSQLSelectForModelWithTableAlias(FactorAspect{}, "fa")), r.UserID, faID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	rows.Next()

	return r.getFactorAspectFromCurrentRow(rows)
}

func (r *FactorAspectRepository) SaveFactorAspect(fa *FactorAspect) error {
	db := helpers.NewDatabaseConnection()

	fa.FactorAspect = helpers.CleanStringSpecials(fa.FactorAspect)

	_, err := db.Query(fmt.Sprintf(`UPDATE factor_aspects SET factor_aspect=%s, archived=$2
		FROM factors
			WHERE factors.id=factor_aspects.factor_id
			AND factor_aspects.id = $3
				AND factors.user_id=$4`,
		helpers.GetEncryptSQLString(`($1)`)),
		fa.FactorAspect, fa.Archived, *fa.ID, r.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *FactorAspectRepository) SaveNewFactorAspect(fa *FactorAspect) (*FactorAspect, error) {
	db := helpers.NewDatabaseConnection()

	fa.FactorAspect = helpers.CleanStringSpecials(fa.FactorAspect)

	err := db.QueryRow(fmt.Sprintf("INSERT INTO factor_aspects (factor_id, factor_aspect) VALUES ($1, %s) RETURNING id",
		helpers.GetEncryptSQLString(`($2)`)),
		fa.FactorID, fa.FactorAspect).
		Scan(&fa.ID)
	if err != nil {
		return nil, err
	}

	return fa, nil
}

func (r *FactorAspectRepository) getFactorAspectListForAllRows(rows *sql.Rows) ([]*FactorAspect, error) {
	factorAspects := []*FactorAspect{} //so that we return empty not null if none

	for rows.Next() {
		currRow, err := r.getFactorAspectFromCurrentRow(rows)
		if err != nil {
			return nil, err
		}
		factorAspects = append(factorAspects, currRow)
	}
	return factorAspects, nil
}

func (r *FactorAspectRepository) getFactorAspectFromCurrentRow(row *sql.Rows) (*FactorAspect, error) {
	fa := new(FactorAspect)

	err := row.Scan(&fa.ID, &fa.FactorID, &fa.FactorAspect, &fa.Archived)

	if err != nil {
		return nil, err
	}

	return fa, nil
}
