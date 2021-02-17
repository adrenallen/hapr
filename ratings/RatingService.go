package ratings

import (
	"time"

	"gitlab.com/garrettcoleallen/happy/factors"
)

type RatingService struct {
	UserID                 int
	Repository             *RatingRepository
	RatingFactorRepository *RatingFactorRepository
}

//Construct and return a new rating service for the user
func NewRatingService(userID int) *RatingService {
	s := new(RatingService)
	s.Repository = &RatingRepository{UserID: userID}
	s.RatingFactorRepository = &RatingFactorRepository{UserID: userID}
	return s
}

func (s *RatingService) GetAllRatings() ([]*Rating, error) {
	return s.Repository.GetAllRatings()
}

func (s *RatingService) GetRatingWithFactorsByID(ratingID int) (*RatingWithFactorsDTO, error) {

	rating, err := s.Repository.GetRatingByID(ratingID)

	if err != nil {
		return nil, err
	}
	ratingWithFactors := &RatingWithFactorsDTO{Rating: rating}

	ratingFactors, err := s.RatingFactorRepository.GetRatingFactorsForRating(rating.ID)
	if err != nil {
		return nil, err
	}

	ratingWithFactors.RatingFactors = ratingFactors

	return ratingWithFactors, nil
}

func (s *RatingService) GetAllRatingsWithFactors() ([]*RatingWithFactorsDTO, error) {

	ratings, err := s.Repository.GetAllRatings()
	if err != nil {
		return nil, err
	}

	ratingsWithFactors := []*RatingWithFactorsDTO{}

	//for each rating we want to add all the factors
	for _, rating := range ratings {
		newRatingWithFactors := &RatingWithFactorsDTO{Rating: rating}
		ratingFactors, err := s.RatingFactorRepository.GetRatingFactorsForRating(rating.ID)
		if err != nil {
			return nil, err
		}

		newRatingWithFactors.RatingFactors = ratingFactors
		ratingsWithFactors = append(ratingsWithFactors, newRatingWithFactors)
	}

	return ratingsWithFactors, nil
}

func (s *RatingService) GetRatingsWithFactorsByDateRange(startDate time.Time, endDate time.Time) ([]*RatingWithFactorsDTO, error) {

	ratings, err := s.Repository.GetRatingsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	ratingsWithFactors := []*RatingWithFactorsDTO{}

	//for each rating we want to add all the factors
	for _, rating := range ratings {
		newRatingWithFactors := &RatingWithFactorsDTO{Rating: rating}
		ratingFactors, err := s.RatingFactorRepository.GetRatingFactorsForRating(rating.ID)
		if err != nil {
			return nil, err
		}

		newRatingWithFactors.RatingFactors = ratingFactors
		ratingsWithFactors = append(ratingsWithFactors, newRatingWithFactors)
	}

	return ratingsWithFactors, nil
}

func (s *RatingService) SaveNewRating(rating int, createdDate time.Time, journalEntry string) (*Rating, error) {
	savedRating, err := s.Repository.SaveNewRating(&Rating{UserID: s.UserID, Rating: rating, CreatedDatetime: createdDate, JournalEntry: journalEntry})

	if err != nil {
		return nil, err
	}

	return savedRating, nil
}

func (s *RatingService) SaveNewRatingFactor(ratingID int, factorID int, factorTypeID int, rank int, factorAspectID *int) (*RatingFactor, error) {

	savedRating, err := s.RatingFactorRepository.SaveNewRatingFactor(&RatingFactor{RatingID: ratingID,
		Factor:       &factors.Factor{ID: factorID},
		FactorTypeID: factorTypeID,
		Rank:         rank,
		FactorAspect: &factors.FactorAspect{ID: factorAspectID}})

	if err != nil {
		return nil, err
	}

	return savedRating, nil
}

func (s *RatingService) UpdateRatingFactor(ratingFactorID int, rank int, factorAspectID *int) error {
	//TODO - check for factor aspect to be owned by user

	//get rating factor by id
	ratingFactor, err := s.RatingFactorRepository.GetRatingFactorByID(ratingFactorID)
	if err != nil {
		return err
	}

	//update rank
	ratingFactor.Rank = rank

	//TODO - make sure this is a valid aspect for this factor and user
	ratingFactor.FactorAspect.ID = factorAspectID

	//save update rating factor
	return s.RatingFactorRepository.UpdateRatingFactor(ratingFactor)
}

func (s *RatingService) DeleteRatingAndRatingFactors(ratingID int) error {
	err := s.RatingFactorRepository.DeleteRatingFactorsByRatingID(ratingID)
	if err != nil {
		return err
	}
	return s.Repository.DeleteRatingByID(ratingID)
}

func (s *RatingService) DeleteRatingFactor(ratingFactorID int) error {
	return s.RatingFactorRepository.DeleteRatingFactor(ratingFactorID)
}

func (s *RatingService) UpdateRating(ratingID int, rating int, journalEntry string) error {
	ratingRow, err := s.Repository.GetRatingByID(ratingID)
	if err != nil {
		return err
	}

	ratingRow.Rating = rating
	ratingRow.JournalEntry = journalEntry

	return s.Repository.UpdateRating(ratingRow)
}
