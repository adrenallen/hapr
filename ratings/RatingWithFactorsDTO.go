package ratings

type RatingWithFactorsDTO struct {
	Rating        *Rating         `json:"rating"`
	RatingFactors []*RatingFactor `json:"ratingFactors"`
}
