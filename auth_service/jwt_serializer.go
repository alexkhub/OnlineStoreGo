package authservice

type LoginUser struct {
	Username string `json:"username"  valid:"-"`
	Email string `json:"email" valid:"email"`
	Password string `json:"password" binding:"required" valid:"-"`
}

type JWTToken struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type RefreshToken struct {
	Refresh string `json:"refresh" valid:"-"`
}

type LoginPostgresData struct{
	Id int `json:"id" db:"id"`
	Role int `json:"role" db:"role_id"`
	Password string `json:"password" db:"hash_password"`
	Activate bool `json:"activate" db:"activate"`
	Block bool `json:"block" db:"block"`

}

type RefreshCheckUser struct{
	Role int `json:"role" db:"role_id"`
	Activate bool `json:"activate" db:"activate"`
	Block bool `json:"block" db:"block"`
}
