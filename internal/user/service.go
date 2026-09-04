package user

type Repository interface {
	GetByID(userID string) (*User, error)
}

type Service struct {
	repo Repository
}

func newService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(userID string) (*User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
