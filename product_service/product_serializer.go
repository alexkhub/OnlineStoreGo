package productservice

type AuthMiddlewareSerializer struct {
	Id   string `json:"id" binding:"required" valid:"-"`
	Role string `json:"role" binding:"required" valid:"-"`
}
type CategorySerializer struct {
	Id   int    `json:"id" valid:"-" db:"id"`
	Name string `json:"name" binding:"required" valid:"-" db:"name"`
}


type ImageSerializer struct {
	Name string `json:""`
}