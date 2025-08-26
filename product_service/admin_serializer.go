package productservice

import "gopkg.in/guregu/null.v3"

type FileUploadSerializer struct {
	FileName string
	Size     int64
	Data     []byte
}

type AdminCreateProductSerializer struct {
	Name        string      `json:"name" binding:"required"`
	Price       int         `json:"price" binding:"required"`
	Discount    null.Int     `json:"discount" binding:"required"`
	Description null.String `json:"description"`
	Category    null.Int    `json:"category"`
}

type AdminProductDetailSerailizer struct {
	Id          int               `json:"id" db:"id"`
	Name        string            `json:"name" db:"name"`
	FirstPrice  int               `json:"first_price" db:"first_price"`
	Description null.String       `json:"description" db:"description"`
	Discount    int               `json:"discount" db:"discount"`
	Price       int               `json:"price" db:"price"`
	Category    null.String       `json:"category" db:"category"`
	Images      []ImageSerializer `json:"images"`
}

type AdminUpdateProductSerializer struct {
	Name        string      `json:"name"`
	Price       int         `json:"price"`
	Discount    null.Int    `json:"discount"`
	Description null.String `json:"description"`
	Category    null.Int    `json:"category"`
}


type AdminKafkaUpdateProductSerializer struct {
	Id int `json:"name"`
	Price int  `json:"price"`
}