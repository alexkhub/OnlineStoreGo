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
	Discount    int         `json:"discount" binding:"required" valid:"-"`
	Description null.String `json:"description" valid:"-" `
	Category    null.Int    `json:"category" valid:"-"`
}
