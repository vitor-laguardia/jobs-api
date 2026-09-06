package user

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

const (
	keyName          = "name"
	MsgNameRequired  = "name is required"
	keyEmail         = "email"
	MsgEmailRequired = "email is required"
)

func (cr CreateUserRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	if cr.Name == "" {
		problems[keyName] = MsgNameRequired
	}

	if cr.Email == "" {
		problems[keyEmail] = MsgEmailRequired
	}

	return
}
