package users

// create request
type UsersCreateRequest struct {
	Name     string `json:"name" `
	Username string `json:"username" `
	Password string `json:"password" `
	Email    string `json:"email" `
}

// update request
type UsersUpdateRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}
