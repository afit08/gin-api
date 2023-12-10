package models

import "time"

type Role struct {
	ID         string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
