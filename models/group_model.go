package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id"`
	Name        string               `bson:"name" json:"name" validate:"required"`
	Owners      []primitive.ObjectID `bson:"owners" json:"owners" validate:"required"`
	CreatedDate time.Time            `bson:"createdDate" json:"createdDate"`
	UpdatedOn   time.Time            `bson:"updatedOn" json:"updatedOn"`
	IsActive    bool                 `bson:"isActive" json:"isActive"`
}
