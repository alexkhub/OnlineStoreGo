package authservice

type UserBlockResponseSerializer struct {
	Email string `json:"email" binding:"required,email" db:"email"`
	Block bool   `json:"block" binding:"required"  db:"block"`
}
