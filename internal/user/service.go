package user

type Repository interface {
	GetByID(userID string) *User
}

type Service struct {
	repo Repository
}

func newService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(userID string) *User {
	user := s.repo.GetByID(userID)
	return user
}
