// models/user_model.go
package models

import "time"

type User struct {
	ID         string    `json:"id,omitempty" bson:"_id,omitempty"`
	Username   string    `json:"username,omitempty" bson:"username,omitempty"`
	Password   string    `json:"password,omitempty" bson:"password,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Role_id    string    `json:"role_id,omitempty" bson:"role_id,omitempty"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
