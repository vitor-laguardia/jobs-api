package user

import "regexp"

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

const (
	keyName             = "name"
	MsgNameRequired     = "name is required"
	keyEmail            = "email"
	MsgEmailRequired    = "email is required"
	MsgWrongEmailFormat = "wrong email format"
)

func isEmailValidRegex(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (cr CreateUserRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	if cr.Name == "" {
		problems[keyName] = MsgNameRequired
	}

	switch {
	case cr.Email == "":
		problems[keyEmail] = MsgEmailRequired
	case !isEmailValidRegex(cr.Email):
		problems[keyEmail] = MsgWrongEmailFormat
	}

	return
}

func (ur UpdateUserRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	if ur.Name == "" {
		problems[keyName] = MsgNameRequired
	}

	return
}
