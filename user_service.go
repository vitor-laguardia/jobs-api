package main

type UserService struct {
}

func (us *UserService) GetByID(userID string) *User {
	user := &User{ID: userID, Name: "alfred", Email: "alfred@gmail.com"}
	return user
}
