package authservice

type AuthMiddlewareSerializer struct {
	Id   string `json:"id" binding:"required"`
	Role string `json:"role" binding:"required"`
}

type AuthRegistrationSerializer struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repet_password" binding:"required"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
}

type AuthRegistrationResponseSerializer struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type ConfirmUserSerializer struct {
	Id int `json:"id" db:"id"`
}
