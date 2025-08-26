package authservice

type LoginUser struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required"`
}

type JWTToken struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type RefreshToken struct {
	Refresh string `json:"refresh" binding:"required"`
}

type LoginPostgresData struct {
	Id       int    `json:"id" db:"id"`
	Role     int    `json:"role" db:"role_id"`
	Password string `json:"password" db:"hash_password"`
	Activate bool   `json:"activate" db:"activate"`
	Block    bool   `json:"block" db:"block"`
}

type RefreshCheckUser struct {
	Role     int  `json:"role" db:"role_id"`
	Activate bool `json:"activate" db:"activate"`
	Block    bool `json:"block" db:"block"`
}
