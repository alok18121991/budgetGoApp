package configs

import "go.mongodb.org/mongo-driver/mongo"

var ExpenseCollection *mongo.Collection = GetCollection(DB, "expenses")
var UserCollection *mongo.Collection = GetCollection(DB, "user")
var GroupCollection *mongo.Collection = GetCollection(DB, "group")
