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
	Description    string             `bson:"description" json:"description" validate:"required"`
	Source         string             `bson:"source" json:"source" validate:"required"`
	Group          primitive.ObjectID `bson:"group_id" json:"group" validate:"required"`
}

func SetNewExpenseId(expense *Expense) Expense {
	expense.Id = primitive.NewObjectID()
	expense.CreatedDate = time.Now()
	return *expense
}

func UpdateExpenseDateTimeToCurrent(expense *Expense) {
	currentTime := time.Now()
	date := expense.ExpenseDate
	updatedDate := time.Date(date.Year(), date.Month(), date.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond(), time.Local)
	expense.ExpenseDate = updatedDate
}
