package domain

const RoleAdmin = "admin"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

type UserRepository interface {
	FindByUsername(username string) (*User, error)
	Create(user *User) error
}

type JWTService interface {
	GenerateToken(user *User) (string, int64, error)
	ValidateToken(tokenString string) (*User, error)
}
