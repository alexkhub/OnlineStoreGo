package productservice


type ProductGRPCSerializer struct{
	Id int64 `db:"id"`
	Price int64 `db:"price"`
	Name string `db:"name"`
}