package ratings

import (
	"database/sql"
	"time"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

//Rating repository
type RatingRepository struct {
	UserID int
}

func (repo *RatingRepository) UpdateRating(rating *Rating) error {
	db := helpers.NewDatabaseConnection()

	_, err := db.Query("UPDATE ratings set rating=$1, journal_entry=$2 where id = $3 and user_id=$4", rating.Rating, rating.JournalEntry, rating.ID, repo.UserID)
	return err
}

func (repo *RatingRepository) DeleteRatingByID(ratingID int) error {
	db := helpers.NewDatabaseConnection()

	_, err := db.Query("DELETE FROM ratings where id = $1 and user_id=$2", ratingID, repo.UserID)
	return err
}

func (repo *RatingRepository) GetRatingByID(ratingID int) (*Rating, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query("SELECT * FROM ratings WHERE id = $1 and user_id=$2", ratingID, repo.UserID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	rows.Next()
	return repo.getRatingFromCurrentRow(rows)

}

func (repo *RatingRepository) GetAllRatings() ([]*Rating, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query("SELECT * FROM ratings WHERE user_id = $1 order by created_datetime", repo.UserID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return repo.getRatingsListForAllRows(rows)
}

//inclusive
func (repo *RatingRepository) GetRatingsByDateRange(startDate time.Time, endDate time.Time) ([]*Rating, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query("SELECT * FROM ratings WHERE user_id = $1 and created_datetime >= $2 and created_datetime <= $3", repo.UserID, startDate, endDate)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return repo.getRatingsListForAllRows(rows)
}

func (repo *RatingRepository) GetFirstRatingBeforeDate(date time.Time) (*Rating, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query("SELECT * FROM ratings WHERE created_datetime < $1 AND user_id=$2", date, repo.UserID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	rows.Next()
	return repo.getRatingFromCurrentRow(rows)
}

func (repo *RatingRepository) SaveNewRating(newRating *Rating) (*Rating, error) {
	db := helpers.NewDatabaseConnection()

	err := db.QueryRow("INSERT INTO ratings (user_id, rating, created_datetime, journal_entry) VALUES ($1, $2, $3, $4) RETURNING id, created_datetime", repo.UserID, newRating.Rating, newRating.CreatedDatetime, newRating.JournalEntry).
		Scan(&newRating.ID, &newRating.CreatedDatetime)
	if err != nil {
		return nil, err
	}

	return newRating, nil
}

func (repo *RatingRepository) getRatingsListForAllRows(rows *sql.Rows) ([]*Rating, error) {
	ratingsList := []*Rating{} //so that we return empty not null if none

	for rows.Next() {
		currRow, err := repo.getRatingFromCurrentRow(rows)
		if err != nil {
			return nil, err
		}
		ratingsList = append(ratingsList, currRow)
	}
	return ratingsList, nil
}

func (repo *RatingRepository) getRatingFromCurrentRow(row *sql.Rows) (*Rating, error) {
	rating := new(Rating)

	err := row.Scan(&rating.ID, &rating.UserID, &rating.Rating, &rating.CreatedDatetime, &rating.JournalEntry)

	if err != nil {
		return nil, err
	}

	return rating, nil
}
