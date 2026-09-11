package user

type Repository interface {
	GetByID(userID string) (User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
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

func (s *Service) Create(reqInput CreateUserRequest) (User, error) {
	user := NewUser(reqInput.Name, reqInput.Email)

	newUser, err := s.repo.Create(user)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (s *Service) Update(userID string, reqInput UpdateUserRequest) (User, error) {
	user, err := s.repo.GetByID(userID)

	if err != nil {
		return User{}, err
	}

	if reqInput.Name != "" {
		user.Name = reqInput.Name
	}

	updatedUser, updateErr := s.repo.Update(user)
	if updateErr != nil {
		return User{}, updateErr
	}
	return updatedUser, nil
}

func (s *Service) Delete(userID string) error {
	return s.repo.Delete(userID)
}
