package user

type UserRepository struct {
	users map[string]User
}

func NewInMemoryRepository() *UserRepository {
	return &UserRepository{users: make(map[string]User)}
}

func (ur *UserRepository) GetByID(userID string) (User, error) {
	user, exists := ur.users[userID]

	if !exists {
		return User{}, ErrNotFound
	}

	return user, nil
}

func (ur *UserRepository) Create(user User) (User, error) {
	if emailAlreadyExists(user.Email, ur.users) {
		return User{}, ErrDuplicateEmail
	}
	ur.users[user.ID] = user
	return user, nil
}

func (ur *UserRepository) Update(user User) (User, error) {
	ur.users[user.ID] = user
	return user, nil
}

func (ur *UserRepository) Delete(userID string) error {
	if _, exists := ur.users[userID]; !exists {
		return ErrNotFound
	}
	delete(ur.users, userID)
	return nil
}

func emailAlreadyExists(email string, users map[string]User) bool {
	for _, u := range users {
		if u.Email == email {
			return true
		}
	}
	return false
}
