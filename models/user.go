package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Username    string             `json:"username" validate:"required"`
	Email       string             `json:"email" validate:"required,email"`
	Password    string             `json:"password" validate:"required,min=6,max=20"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	Birthday    time.Time          `json:"birthday"`
	Gender      string             `json:"gender"`
	Address     string             `json:"address"`
	PhoneNumber string             `json:"phoneNumber"`
	Token       string             `json:"token"`
}
