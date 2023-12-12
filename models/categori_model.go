package models

import (
	"mime/multipart"
	"time"
)

type Categori struct {
	ID        string                `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string                `json:"name,omitempty" bson:"name,omitempty"`
	Image     *multipart.FileHeader `json:"image,omitempty" bson:"image,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// required to form data
type CreateCategoriRequest struct {
	ID        string                `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string                `form:"name" binding:"required"`
	Image     *multipart.FileHeader `form:"image" binding:"required"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// required to update image form data
type UpdateCategoriRequest struct {
	ID        string                `form:"id,omitempty" bson:"_id,omitempty"`
	Name      string                `form:"name"`
	Image     *multipart.FileHeader `form:"image" binding:"-"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// update not required to upload image form data
type UpdateCategoriRequestNoImage struct {
	ID        string                `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string                `form:"name" json:"name,omitempty" bson:"name,omitempty"`
	Image     *multipart.FileHeader `json:"image,omitempty" bson:"-"` // Remove "form" and "binding" tags
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}
