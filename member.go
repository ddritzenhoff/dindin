package dinny

// Member represents a member of dinner rotation.
type Member struct {
	ID          int64  `json:"id"`
	SlackUID    string `json:"slackUID"`
	FullName    string `json:"fullName"`
	MealsEaten  int64  `json:"mealsEaten"`
	MealsCooked int64  `json:"mealsCooked"`
	Leader      bool   `json:"leader"`
}

// MemberService represents a service for managing members.
type MemberService interface {
	// FindMemberByID retrieves a member by ID.
	// Returns ErrNotFound if meal does not exist.
	FindMemberByID(id int64) (*Member, error)

	// FindMemberBySlackUID retrieves a member by SlackID.
	// Returns ErrNotFound if meal does not exist.
	FindMemberBySlackUID(slackUID string) (*Member, error)

	// ListMembers retrieves a list of members.
	ListMembers() ([]*Member, error)

	// CreateMember creates a new member.
	CreateMember(m *Member) error

	// UpdateMember updates a member object.
	UpdateMember(id int64, upd MemberUpdate) error

	// DeleteMember permanently deletes a member.
	DeleteMember(id int64) error
}

// MemberUpdate represents a set of fields to be updated via UpdateMember().
type MemberUpdate struct {
	MealsEaten  *int64
	MealsCooked *int64
	Leader      *bool
}
