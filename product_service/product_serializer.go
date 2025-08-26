package productservice

import (
	"gopkg.in/guregu/null.v3"
)

type AuthMiddlewareSerializer struct {
	Id   string `json:"id" binding:"required"`
	Role string `json:"role" binding:"required" valid:"-"`
}
type CategorySerializer struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" binding:"required" db:"name"`
}

type ImageSerializer struct {
	Name null.String `json:"name" db:"image_uuid"`
	Link null.String `json:"link"`
}

type ProductListSerailizer struct {
	Id         int         `json:"id" db:"id"`
	Name       string      `json:"name" db:"name"`
	FirstPrice int         `json:"first_price" db:"first_price"`
	Discount   int         `json:"discount" db:"discount"`
	Price      int         `json:"price" db:"price"`
	Category   null.String `json:"category" db:"category"`
	Image      null.String `json:"image_name" db:"image_uuid"`
	ImageLink  null.String `json:"image_link"`
}

type ProductDetailSerailizer struct {
	Id          int         `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	FirstPrice  int         `json:"first_price" db:"first_price"`
	Description null.String `json:"description" db:"description"`
	Discount    int         `json:"discount" db:"discount"`
	Price       int         `json:"price" db:"price"`
	Category    null.String `json:"category" db:"category"`
	Images      []string    `json:"images"`
}
