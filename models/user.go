package models

type UserQuery struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
}

type UserCreate struct {
	UserQuery

	Password string `json:"password"`
}
