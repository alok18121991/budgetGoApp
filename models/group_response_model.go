package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupResponse struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Owners      []UserResponse     `json:"owners,omitempty"`
	CreatedDate time.Time          `json:"createdDate,omitempty"`
	UpdatedOn   time.Time          `json:"updatedOn,omitempty"`
	IsActive    bool               `json:"isActive,omitempty"`
}

// UserResponse represents the response structure for user details
type UserResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	FirstName string             `json:"firstName,omitempty"`
	LastName  string             `json:"lastName,omitempty"`
}
