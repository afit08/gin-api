package models

import (
	"mime/multipart"
	"time"
)

type Product struct {
	ID         string                `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string                `json:"name,omitempty" bson:"name,omitempty"`
	Image      *multipart.FileHeader `json:"image,omitempty" bson:"image,omitempty"`
	Price      string                `json:"price,omitempty" bson:"price,omitempty"`
	Desc       string                `json:"desc,omitempty" bson:"desc,omitempty"`
	Stock      string                `json:"stock,omitempty" bson:"stock,omitempty"`
	Weight     string                `json:"weight,omitempty" bson:"weight,omitempty"`
	Created_at time.Time             `json:"created_at"`
	Updated_at time.Time             `json:"updated_at"`
}

type CreateProductRequest struct {
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

type UpdateProductRequest struct {
	ID        string                `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string                `form:"name" binding:"required"`
	Image     *multipart.FileHeader `form:"image" binding:"-"`
	Price     string                `form:"price" binding:"required"`
	Desc      string                `form:"desc"`
	Stock     string                `form:"stock" binding:"required"`
	Weight    string                `form:"weight" binding:"required"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}
