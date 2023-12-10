package productsDTO

import (
	"mime/multipart"
	"time"
)

type ProductRequest struct {
	ID        string                `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string                `form:"name" binding:"required"`
	Image     *multipart.FileHeader `form:"image" binding:"required"`
	Price     string                `form:"price" binding:"required"`
	Desc      string                `form:"desc"`
	Stock     string                `form:"stock"`
	Weight    string                `form:"weight"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}
