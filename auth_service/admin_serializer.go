package authservice

type AuthMiddlewareSerializer struct{
	Id string `json:"id" binding:"required" valid:"-"`
	Role string `json:"role" binding:"required" valid:"-"`

}

type UserBlockResponseSerializer struct{
	Email string `json:"email" binding:"required" valid:"email" db:"email"`
	Block bool `json:"block" binding:"required" valid:"email" db:"block"`
}