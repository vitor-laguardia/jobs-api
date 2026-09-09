package user

type Repository interface {
	GetByID(userID string) (User, error)
	Create(user User) User
	Update(user User) User
	Delete(userID string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
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

func (s *Service) Update(userID string, reqInput UpdateUserRequest) (User, error) {
	user, err := s.repo.GetByID(userID)

	if err != nil {
		return User{}, err
	}

	if reqInput.Name != "" {
		user.Name = reqInput.Name
	}

	updatedUser := s.repo.Update(user)
	return updatedUser, nil
}

func (s *Service) Delete(userID string) error {
	return s.repo.Delete(userID)
}
