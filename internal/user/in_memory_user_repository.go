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
