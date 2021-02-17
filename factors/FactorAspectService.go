package factors

type FactorAspectService struct {
	UserID     int
	Repository *FactorAspectRepository
}

func NewFactorAspectService(userID int) *FactorAspectService {
	s := new(FactorAspectService)
	s.UserID = userID
	s.Repository = &FactorAspectRepository{UserID: userID}
	return s
}

func (s *FactorAspectService) GetFactorAspects() ([]*FactorAspect, error) {
	return s.Repository.GetFactorAspects()
}

func (s *FactorAspectService) SetArchiveFactorAspect(faID int, archive bool) error {
	fa, err := s.Repository.GetFactorAspectByID(faID)
	if err != nil {
		return err
	}

	fa.Archived = archive
	return s.Repository.SaveFactorAspect(fa)
}

func (s *FactorAspectService) RenameFactorAspect(faID int, faName string) error {
	fa, err := s.Repository.GetFactorAspectByID(faID)
	if err != nil {
		return err
	}

	fa.FactorAspect = faName

	return s.Repository.SaveFactorAspect(fa)
}

func (s *FactorAspectService) SaveNewFactorAspect(faName string, factorID int) (*FactorAspect, error) {
	factorService := NewFactorService(s.UserID)
	factor, err := factorService.GetFactorByID(factorID)
	if err != nil || factor == nil {
		return nil, err
	}

	fa, err := s.Repository.SaveNewFactorAspect(&FactorAspect{FactorAspect: faName, FactorID: factorID})

	if err != nil {
		return nil, err
	}

	return fa, nil

}

func (s *FactorAspectService) GetFactorAspectByID(faID int) (*FactorAspect, error) {
	return s.Repository.GetFactorAspectByID(faID)
}
