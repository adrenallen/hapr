package reports

import (
	"time"

	"gitlab.com/garrettcoleallen/happy/helpers"
)

type ReportRepository struct {
	UserID int
}

func (r *ReportRepository) GetFactorToRatingOccurrences() ([]*FactorToRatingOccurrenceDTO, error) {
	db := helpers.NewDatabaseConnection()

	rows, err := db.Query(`
	with r as (
		select
			id,
			rating
		from
			ratings
		where
			user_id = $1),
		rf as (
		select
			rating_id,
			factor_id,
			factor_type_id
		from
			rating_factors
		where
			rating_id in (
			select
				id
			from
				r)),
		rf_group as (
		select
			count(*) occurrences,
			factor_id,
			factor_type_id,
			rating
		from
			rf
		left join r on
			r.id = rf.rating_id
		group by
			factor_id,
			factor_type_id,
			rating ) select
			occurrences, factor_id, factor_type_id, rating
		from
			rf_group
		order by
			occurrences desc
	`, r.UserID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	reportList := []*FactorToRatingOccurrenceDTO{}

	for rows.Next() {
		result := new(FactorToRatingOccurrenceDTO)

		err := rows.Scan(&result.Occurrences, &result.FactorID, &result.FactorTypeID, &result.Rating)
		if err != nil {
			return nil, err
		}

		reportList = append(reportList, result)
	}

	return reportList, nil
}

//Get a rating's normalized happiness rating
func (r *ReportRepository) GetRatingNHR(ratingID int) (float64, error) {
	db := helpers.NewDatabaseConnection()

	row := db.QueryRow(`
	with main_rating as (
		select *
	from
		ratings
	where
		id = $1),
	prior_ratings as (
		select *
	from
		ratings
	where
		user_id = $2
		and created_datetime < (
			select created_datetime
		from
			main_rating)),
	avg_prior_rating as (
		select avg(rating) average
	from
		prior_ratings) select
		(main_rating.rating - avg_prior_rating.average) nhr
	from
		main_rating,
		avg_prior_rating
	`, ratingID, r.UserID)

	var nhr float64
	err := row.Scan(&nhr)
	return nhr, err

}

//Get a rating's normalized happiness rating
func (r *ReportRepository) GetAverageRatingAsOfDate(date time.Time) (float64, error) {
	db := helpers.NewDatabaseConnection()

	row := db.QueryRow(`
	select avg(rating) from ratings where user_id=$1 and created_datetime <= $2
	`, r.UserID, date)

	var avg float64
	err := row.Scan(&avg)
	return avg, err

}
