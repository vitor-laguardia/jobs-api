package user

type Repository interface {
	GetByID(userID string) (User, error)
	Create(user User) User
}

type Service struct {
	repo Repository
}

func newService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(userID string) (User, error) {
	user, err := s.repo.GetByID(userID)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Service) Create(reqInput CreateUserRequest) User {
	user := NewUser(reqInput.Name, reqInput.Email)

	newUser := s.repo.Create(user)

	return newUser
}
