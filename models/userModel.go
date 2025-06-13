package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Username      *string            `bson:"username" validate:"required,min=2,max=100"`
	Email         *string            `bson:"email" validate:"email,required"`
	Password      *string            `bson:"password" validate:"required,min=6"`
	Token         *string            `bson:"token,omitempty"`
	Refresh_Token *string            `bson:"refresh_token,omitempty"`
	Created_at    time.Time          `bson:"created_at"`
	Updated_at    time.Time          `bson:"updated_at"`
	User_id       string             `bson:"user_id"`
}
