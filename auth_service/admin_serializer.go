package authservice

type UserBlockResponseSerializer struct {
	Email string `json:"email" binding:"required" valid:"email" db:"email"`
	Block bool   `json:"block" binding:"required" valid:"email" db:"block"`
}
