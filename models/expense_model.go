package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	UserId         string             `bson:"user_id" json:"userId" validate:"required"`
	Amount         float32            `bson:"amount" json:"amount" validate:"required"`
	ExpenseType    string             `bson:"type" json:"type" validate:"required"`
	ExpenseSubType string             `bson:"subType" json:"subType" validate:"required"`
	CreatedDate    time.Time          `bson:"createdDate" json:"createdDate"`
	ExpenseDate    time.Time          `bson:"expenseDate" json:"expenseDate"`
	Description    string             `bson:"description" json:"description"`
	Source         string             `bson:"source" json:"source" validate:"required"`
}

func SetNewExpenseId(expense *Expense) Expense {
	expense.Id = primitive.NewObjectID()
	expense.CreatedDate = time.Now()
	return *expense
}
