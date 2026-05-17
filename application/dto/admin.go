package dto

// LoginInput matches the swagger LoginRequest schema.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AdminProfileData matches the swagger AdminProfile schema.
type AdminProfileData struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LoginData matches the swagger LoginResponse.data schema.
type LoginData struct {
	Token     string           `json:"token"`
	ExpiresAt string           `json:"expires_at"`
	Admin     AdminProfileData `json:"admin"`
}

// LoginOutput matches the swagger LoginResponse schema.
type LoginOutput struct {
	Data LoginData `json:"data"`
}

// AdminMeOutput wraps the admin profile for GET /admin/auth/me.
type AdminMeOutput struct {
	Data AdminProfileData `json:"data"`
}
