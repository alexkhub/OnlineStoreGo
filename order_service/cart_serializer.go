package orderservice


type AuthMiddlewareSerializer struct {
	Id   string `json:"id" binding:"required" valid:"-"`
	Role string `json:"role" binding:"required" valid:"-"`
}

type CartPostgresSerializer struct {
	Id int64 `json:"id" db:"id"`
	Product int64 `json:"product" db:"product_id"`
	Amount  int64 `json:"amount" db:"amount"`

}


type CartProductSerializer struct {
	Id int64 `json:"id"`
	Price int64 `json:"price"`
	Name string `json:"name"`
}

type CartSerializer struct {
	Id int64 `json:"id" db:"id"`
	Product CartProductSerializer `json:"product"`
	Amount  int64 `json:"amount" db:"amount"`
}


type CreateCartSerializer struct {
	Product int64 `json:"product"`
}

type UpdateCartSerializer struct{
	Amount int `json:"amount" binding:"required" valid:"-"`
}