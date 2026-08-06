package main

const (
	MsgTitleRequired      = "title is required"
	KeyTitle              = "title"
	MsgUserIDRequired     = "userId is required"
	KeyUserID             = "userId"
	MsgPriorityOutOfRange = "priority must be in range [1-3]"
	KeyPriority           = "priority"
)

type CreateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	UserID      string `json:"userId"`
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
