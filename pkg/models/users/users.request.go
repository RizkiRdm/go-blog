package users

// create request
type UsersCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// update request
type UsersUpdateRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

// login request
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
