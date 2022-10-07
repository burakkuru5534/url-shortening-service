package auth

type AuthResponse struct {
	UserID       int64
	UserName     string
	UserEmail    string
	UserFullName string
	AccessToken  string
}
