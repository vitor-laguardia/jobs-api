package job

const (
	MsgTitleRequired      = "title is required"
	KeyTitle              = "title"
	MsgUserIDRequired     = "userId is required"
	KeyUserID             = "userId"
	MsgPriorityOutOfRange = "priority must be in range [1-3]"
	KeyPriority           = "priority"
	KeyBody               = "body"
	MsgEmptyJSONBody      = "request body must contain at least one field to update"
	KeyStatus             = "status"
	MsgInvalidStatus      = "invalid status value"
)

type CreateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	UserID      string `json:"userId"`
}

type UpdateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Status      string `json:status`
}

func (r CreateJobRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	if r.Title == "" {
		problems[KeyTitle] = MsgTitleRequired
	}

	if r.UserID == "" {
		problems[KeyUserID] = MsgUserIDRequired
	}

	if r.Priority < 0 || r.Priority > 3 {
		problems[KeyPriority] = MsgPriorityOutOfRange
	}
	return
}

func (u UpdateJobRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	if u.Title == "" && u.Description == "" && u.Priority == 0 && u.Status == "" {
		problems[KeyBody] = MsgEmptyJSONBody
	}

	if u.Priority < 0 || u.Priority > 3 {
		problems[KeyPriority] = MsgPriorityOutOfRange
	}

	if u.Status != "" && !JobStatus(u.Status).IsValid() {
		problems[KeyStatus] = MsgInvalidStatus
	}

	return
}
