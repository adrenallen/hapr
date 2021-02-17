package ratings

import "gitlab.com/garrettcoleallen/happy/factors"

type RatingFactor struct {
	ID           int                   `table:"rating_factors" json:"id"`
	RatingID     int                   `json:"ratingID"`
	Factor       *factors.Factor       `json:"factor"`
	FactorTypeID int                   `json:"factorTypeID"`
	Rank         int                   `json:"rank"`
	FactorAspect *factors.FactorAspect `json:"factorAspect"`
}
