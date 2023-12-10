package models

import (
	"time"
)

type Product struct {
	ID         string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Image      string    `json:"image,omitempty" bson:"image,omitempty"`
	Price      string    `json:"price,omitempty" bson:"price,omitempty"`
	Desc       string    `json:"desc,omitempty" bson:"desc,omitempty"`
	Stock      string    `json:"stock,omitempty" bson:"stock,omitempty"`
	Weight     string    `json:"weight,omitempty" bson:"weight,omitempty"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
