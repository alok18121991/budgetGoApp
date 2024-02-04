package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FirstName    string             `bson:"firstName" json:"firstName,omitempty" validate:"required"`
	LastName     string             `bson:"lastName" json:"lastName,omitempty" validate:"required"`
	Email        string             `bson:"email" json:"email,omitempty" validate:"required,email"`
	Age          int32              `bson:"age" json:"age,omitempty" validate:"required"`
	Gender       string             `bson:"gender" json:"gender,omitempty" validate:"required"`
	Occupation   string             `bson:"occupation" json:"occupation,omitempty"`
	CreatedDate  time.Time          `bson:"createdDate" json:"createdDate,omitempty"`
	ModifiedDate time.Time          `bson:"modifiedDate" json:"modifiedDate,omitempty"`
}

func SetNewUserId(user *User) User {
	user.Id = primitive.NewObjectID()
	user.CreatedDate = time.Now()
	return *user
}

func CreateUser(user *User) User {
	return User{
		Id:          primitive.NewObjectID(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Age:         user.Age,
		Gender:      user.Gender,
		Occupation:  user.Occupation,
		CreatedDate: time.Now(),
	}
}
