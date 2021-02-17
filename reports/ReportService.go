package reports

import (
	"log"
	"sort"

	"gitlab.com/garrettcoleallen/happy/factors"
	"gitlab.com/garrettcoleallen/happy/ratings"
)

type ReportService struct {
	UserID     int
	Repository *ReportRepository
}

func NewReportService(userID int) *ReportService {
	s := &ReportService{UserID: userID,
		Repository: &ReportRepository{UserID: userID}}

	return s
}

func (s *ReportService) GetFactorToRatingOccurrences() ([]*FactorToRatingOccurrenceDTO, error) {
	return s.Repository.GetFactorToRatingOccurrences()
}

func (s *ReportService) GetRatingNHR(ratingID int) (float64, error) {
	return s.Repository.GetRatingNHR(ratingID)
}

//Returns current factor impacts for a user
func (s *ReportService) CalculateFactorImpactsForUser(userID int) (map[int]float64, error) {
	impactCalcs := map[int]float64{}

	ratingService := ratings.NewRatingService(userID)

	ratings, err := ratingService.GetAllRatingsWithFactors()
	if err != nil {
		return nil, err
	}

	impactTotals := map[int]float64{}
	impactCounts := map[int]int{}
	for _, rating := range ratings {
		calcs, err := s.CalculateFactorImpactsForRating(rating)
		if err != nil {
			return nil, err
		}
		for factorID, impact := range calcs {
			impactCounts[factorID] = impactCounts[factorID] + 1
			impactTotals[factorID] = impactTotals[factorID] + impact
		}
	}

	//average out the total with the counts
	for factorID, impact := range impactTotals {
		impactCalcs[factorID] = impact / float64(impactCounts[factorID])
	}

	return impactCalcs, nil
}

//Calculate the impact of each factor in a rating and return a map of factor.id => the weight of the factor in the rating
func (s *ReportService) CalculateFactorImpactsForRating(rating *ratings.RatingWithFactorsDTO) (map[int]float64, error) {
	impactCalcs := map[int]float64{}

	nhr, err := s.Repository.GetRatingNHR(rating.Rating.ID)
	if err != nil {
		log.Printf("Skipping NHR for rating %v due to an error: %v", rating.Rating.ID, err)
		return impactCalcs, nil
	}

	positiveFactors := []*ratings.RatingFactor{}
	negativeFactors := []*ratings.RatingFactor{}

	//sort positive and negative into their arrays
	for _, factor := range rating.RatingFactors {
		if factor.FactorTypeID == factors.PositiveFactorType {
			positiveFactors = append(positiveFactors, factor)
		} else if factor.FactorTypeID == factors.NegativeFactorType {
			negativeFactors = append(negativeFactors, factor)
		}
	}

	//sort the factors by their rank
	sort.Slice(positiveFactors, func(i, j int) bool {
		return positiveFactors[i].Rank < positiveFactors[j].Rank
	})

	sort.Slice(negativeFactors, func(i, j int) bool {
		return negativeFactors[i].Rank < negativeFactors[j].Rank
	})

	allFactors := []*ratings.RatingFactor{}

	//Positive nhr means positive get more weight and vice-versa
	if nhr > 0 {
		allFactors = append(allFactors, positiveFactors...)
		allFactors = append(allFactors, negativeFactors...)
	} else {
		allFactors = append(allFactors, negativeFactors...)
		allFactors = append(allFactors, positiveFactors...)
	}

	totalFactorCount := len(allFactors)

	//Build the total weight value for distibuting
	totalWeight := 0
	for i := 1; i <= totalFactorCount; i++ {
		totalWeight = totalWeight + 1
	}

	//Calculate impact by weight and record
	for idx, factor := range allFactors {
		//calculate the impact weight of the factor, multiply by the nhr to get the weight total
		impact := float64(totalFactorCount-idx) / float64(totalWeight)
		impactCalcs[factor.Factor.ID] = impact * float64(nhr)
	}

	return impactCalcs, nil
}

//Get the full data export for this user
func (s *ReportService) GetFullDataExport() (*FullDataExport, error) {
	ex := &FullDataExport{}

	ratingService := ratings.NewRatingService(s.UserID)
	rs, err := ratingService.GetAllRatingsWithFactors()
	if err != nil {
		return nil, err
	}

	ex.RatingData = rs

	return ex, nil
}

type FullDataExport struct {
	RatingData []*ratings.RatingWithFactorsDTO `json:"ratingData"`
}
