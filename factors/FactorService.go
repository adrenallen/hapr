package factors

type FactorService struct {
	UserID     int
	Repository *FactorRepository
}

func NewFactorService(userID int) *FactorService {
	s := new(FactorService)
	s.UserID = userID
	s.Repository = &FactorRepository{UserID: userID}
	return s
}

func (s *FactorService) GetFactors() ([]*Factor, error) {
	return s.Repository.GetFactors()
}

func (s *FactorService) GetFactorTypes() AllFactorTypesDTO {
	return AllFactorTypesDTO{
		FactorTypes: s.Repository.GetAllFactorTypes(),
	}
}

func (s *FactorService) SetArchiveFactor(factorID int, archive bool) error {
	factor, err := s.Repository.GetFactorByID(factorID)
	if err != nil {
		return err
	}

	factor.Archived = archive
	return s.Repository.SaveFactor(factor)
}

func (s *FactorService) RenameFactor(factorID int, factorName string) error {
	factor, err := s.Repository.GetFactorByID(factorID)
	if err != nil {
		return err
	}

	factor.Factor = factorName

	return s.Repository.SaveFactor(factor)
}

func (s *FactorService) SaveNewFactor(factor string) (*Factor, error) {
	savedFactor, err := s.Repository.SaveNewFactor(&Factor{UserID: s.UserID, Factor: factor})

	if err != nil {
		return nil, err
	}

	return savedFactor, nil

}

func (s *FactorService) GetFactorByID(factorID int) (*Factor, error) {
	return s.Repository.GetFactorByID(factorID)
}
