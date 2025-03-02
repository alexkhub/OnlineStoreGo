package authservice


type AuthRegistrationSerializer struct{
	Username string `json:"username" binding:"required" valid:"-"`
	Email    string `json:"email" binding:"required" valid:"email"`
	Password string `json:"password" binding:"required" valid:"-"`
	FirstName string `json:"first_name" valid:"-"`
	LastName string `json:"last_name" valid:"-"`
}

type AuthRegistrationResponseSerializer  struct{
	Id int `json:"id" binding:"required" valid:"-"`
	Email string `json:"email" binding:"required" valid:"email"`

}

type ConfirmUserSerializer struct{
	Id int `json:"id" db:"id"`
}

