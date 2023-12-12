// models/user_model.go
package models

import (
	"mime/multipart"
	"time"
)

type User struct {
	ID         string                `json:"id,omitempty" bson:"_id,omitempty"`
	Username   string                `form:"username" binding:"required"`
	Password   string                `form:"password" binding:"required"`
	Name       string                `form:"name"`
	Image      *multipart.FileHeader `form:"image" binding:"-"`
	Role_id    string                `form:"role_id"`
	Created_at time.Time             `json:"created_at"`
	Updated_at time.Time             `json:"updated_at"`
}
