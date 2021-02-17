package reports

type FactorToRatingOccurrenceDTO struct {
	Occurrences  int `json:"occurrences"`
	FactorID     int `json:"factorID"`
	FactorTypeID int `json:"factorTypeID"`
	Rating       int `json:"rating"`
}
