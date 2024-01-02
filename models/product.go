package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `json:"_id,omitempty"`
	Name        string             `json:"name" validate:"required"`
	Type        string             `json:"type" validate:"required"`
	Description string             `json:"description" validate:"required"`
	Sku         string             `json:"sku" validate:"required"`
	ImgUrl      string             `json:"imgUrl"`
	AddedBy     string             `json:"addedBy"`
}
