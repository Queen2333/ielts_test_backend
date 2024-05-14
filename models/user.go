package models
type UserQuery struct {
	ID       string    `json:"id"`
	Email    string `json:"email"`
	RoleID   int `json:"role_id"`
}

type UserCreate struct {
	UserQuery

	Password string `json:"password"`
}
