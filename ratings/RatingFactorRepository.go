package ratings

import (
	"database/sql"

	"gitlab.com/garrettcoleallen/happy/factors"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

//Rating factor repository
type RatingFactorRepository struct {
	UserID int
}

func (repo *RatingFactorRepository) UpdateRatingFactor(ratingFactor *RatingFactor) error {
	db := helpers.NewDatabaseConnection()
	
	_, err := db.Query(`UPDATE rating_factors set rank=$1, factor_aspect_id=$2 
		from ratings 
			where ratings.id=rating_factors.rating_id 
				and rating_factors.id = $3
				and ratings.user_id=$4`, ratingFactor.Rank, ratingFactor.FactorAspect.ID, ratingFactor.ID, repo.UserID)
	return err
}

func (repo *RatingFactorRepository) GetRatingFactorByID(ratingFactorID int) (*RatingFactor, error) {
	db := helpers.NewDatabaseConnection()
	
	rows, err := db.Query(`SELECT rating_factors.*
		FROM rating_factors 
			left join ratings on ratings.id=rating_factors.rating_id
			WHERE rating_factors.id = $1
			and ratings.user_id = $2`, ratingFactorID, repo.UserID)
	if err != nil {
		return nil, err
	}
	rows.Next()
	return repo.getRatingFactorFromCurrentRow(rows)
}

func (repo *RatingFactorRepository) DeleteRatingFactor(ratingFactorID int) error {
	db := helpers.NewDatabaseConnection()
	
	_, err := db.Query(`DELETE FROM rating_factors 
		USING ratings
		WHERE ratings.id=rating_factors.rating_id and rating_factors.id = $1 and ratings.user_id=$2`, ratingFactorID, repo.UserID)
	return err
}

func (repo *RatingFactorRepository) DeleteRatingFactorsByRatingID(ratingID int) error {
	db := helpers.NewDatabaseConnection()
	
	_, err := db.Query(`DELETE FROM rating_factors 
		USING ratings
		WHERE ratings.id=rating_factors.rating_id and ratings.id = $1 and ratings.user_id=$2`, ratingID, repo.UserID)
	return err
}

func (repo *RatingFactorRepository) SaveNewRatingFactor(newRatingFactor *RatingFactor) (*RatingFactor, error) {
	db := helpers.NewDatabaseConnection()
	

	//Handle via database constraint on userid matching rating, factor, and factor aspect
	err := db.QueryRow("INSERT INTO rating_factors (rating_id, factor_id, factor_type_id, rank, factor_aspect_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		newRatingFactor.RatingID, newRatingFactor.Factor.ID, newRatingFactor.FactorTypeID, newRatingFactor.Rank, newRatingFactor.FactorAspect.ID).
		Scan(&newRatingFactor.ID)
	if err != nil {
		return nil, err
	}

	return newRatingFactor, nil
}

func (repo *RatingFactorRepository) GetRatingFactorsForRating(ratingID int) ([]*RatingFactor, error) {
	db := helpers.NewDatabaseConnection()
	
	rows, err := db.Query(`SELECT rating_factors.*
		FROM rating_factors 
			left join ratings on ratings.id=rating_factors.rating_id
			WHERE ratings.id = $1
			and ratings.user_id = $2`, ratingID, repo.UserID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return repo.getRatingFactorsListForAllRows(rows)
}

func (repo *RatingFactorRepository) GetRatingFactorsByFactorID(factorID int) ([]*RatingFactor, error) {
	db := helpers.NewDatabaseConnection()
	
	rows, err := db.Query(`SELECT rating_factors.*
		FROM rating_factors 
			left join ratings on ratings.id=rating_factors.rating_id
			WHERE rating_factors.factor_id = $1
			and ratings.user_id = $2`, factorID, repo.UserID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return repo.getRatingFactorsListForAllRows(rows)
}

func (repo *RatingFactorRepository) getRatingFactorsListForAllRows(rows *sql.Rows) ([]*RatingFactor, error) {
	ratingsList := []*RatingFactor{} //so that we return empty not null if none

	for rows.Next() {
		currRow, err := repo.getRatingFactorFromCurrentRow(rows)
		if err != nil {
			return nil, err
		}
		ratingsList = append(ratingsList, currRow)
	}
	return ratingsList, nil
}

func (repo *RatingFactorRepository) getRatingFactorFromCurrentRow(row *sql.Rows) (*RatingFactor, error) {
	rating := new(RatingFactor)
	rating.Factor = &factors.Factor{}
	rating.FactorAspect = &factors.FactorAspect{}

	err := row.Scan(&rating.ID, &rating.RatingID, &rating.Factor.ID, &rating.FactorTypeID, &rating.Rank, &rating.FactorAspect.ID)

	if err != nil {
		return nil, err
	}

	//TODO - make this more performant, we dont want a lookup everytime bruh
	//fill the factor
	factorService := factors.NewFactorService(repo.UserID)
	rating.Factor, err = factorService.GetFactorByID(rating.Factor.ID)

	if err != nil {
		return nil, err
	}

	aspectService := factors.NewFactorAspectService(repo.UserID)
	//TODO - dont do a lookup everytime!
	if rating.FactorAspect.ID != nil {
		rating.FactorAspect, err = aspectService.GetFactorAspectByID(*rating.FactorAspect.ID)
		if err != nil {
			return nil, err
		}
	}

	return rating, nil
}
