package productservice

import "gopkg.in/guregu/null.v3"

type FileUploadSerializer struct {
	FileName string
	Size     int64
	Data     []byte
}

type AdminCreateProductSerializer struct {
	Name        string      `json:"name" binding:"required" valid:"-"`
	Price       int         `json:"price" binding:"required" valid:"-"`
	Discount    null.Int         `json:"discount" binding:"required" valid:"-"`
	Description null.String `json:"description" valid:"-" `
	Category    null.Int    `json:"category" valid:"-"`
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
	Name        string      `json:"name"  valid:"-"`
	Price       int         `json:"price" valid:"-"`
	Discount    null.Int    `json:"discount" valid:"-"`
	Description null.String `json:"description" valid:"-" `
	Category    null.Int    `json:"category" valid:"-"`
}


type AdminKafkaUpdateProductSerializer struct {
	Id int `json:"name"`
	Price int  `json:"price" valid:"-"`
}