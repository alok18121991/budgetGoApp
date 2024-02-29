package controllers

import (
	"alok/web-service-budget/configs"
	"alok/web-service-budget/models"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var groupCollection *mongo.Collection = configs.GetCollection(configs.DB, "group")

func CreateGroupHandler(c echo.Context) error {
	// Parse request body into Group struct
	group := new(models.Group)
	if err := c.Bind(group); err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	// Validate the group data
	if err := validate.Struct(group); err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "invalid group data")
	}

	group.CreatedDate = time.Now()
	group.IsActive = true
	group.ID = primitive.NewObjectID()

	// Start MongoDB session
	session, err := configs.GetSession(configs.DB)
	if err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to start session")
	}
	defer session.EndSession(context.Background())

	// Start MongoDB transaction
	err = session.StartTransaction()
	if err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to start transaction")
	}

	// Insert group document
	_, err = groupCollection.InsertOne(context.Background(), group)
	if err != nil {
		c.Echo().Logger.Error(err.Error())
		session.AbortTransaction(context.Background())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create group")
	}

	for _, userID := range group.Owners {
		filter := bson.M{"_id": userID}
		update := bson.M{"$addToSet": bson.M{"groups": group.ID}}
		_, err := userCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.Echo().Logger.Error(err.Error())
			session.AbortTransaction(context.Background())
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user")
		}
	}

	// Commit transaction
	err = session.CommitTransaction(context.Background())
	if err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit transaction")
	}

	return c.JSON(http.StatusCreated, group)
}

func AddOwnersToGroupHandler(c echo.Context) error {
	// Parse request body into GroupOwner struct
	request := new(struct {
		GroupID primitive.ObjectID   `json:"groupId" validate:"required"`
		Owners  []primitive.ObjectID `json:"owners" validate:"required"`
	})
	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	// Validate the request data
	if err := validate.Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request data")
	}

	session, err := configs.GetSession(configs.DB)
	if err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to start session")
	}
	defer session.EndSession(context.Background())

	err = session.StartTransaction()
	if err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to start transaction")
	}

	// Define filter to find the group by ID
	filter := bson.M{"_id": request.GroupID}

	// Define update to add owner to the group
	update := bson.M{"$addToSet": bson.M{"owners": bson.M{"$each": request.Owners}},
		"$set": bson.M{"updatedOn": time.Now()},
	}

	// Perform update operation
	_, err = groupCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		session.AbortTransaction(context.Background())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update group")
	}

	for _, userID := range request.Owners {
		filter := bson.M{"_id": userID}
		update := bson.M{"$addToSet": bson.M{"groups": request.GroupID}}
		_, err := userCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.Echo().Logger.Error(err.Error())
			session.AbortTransaction(context.Background())
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user")
		}
	}

	err = session.CommitTransaction(context.Background())
	if err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit transaction")
	}
	return c.NoContent(http.StatusOK)
}

func MarkGroupInactiveHandler(c echo.Context) error {
	request := new(struct {
		GroupID primitive.ObjectID `json:"groupId" validate:"required"`
	})
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	if err := validate.Struct(request); err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request data")
	}

	filter := bson.M{"_id": request.GroupID}

	update := bson.M{
		"$set": bson.M{"isActive": false, "updatedOn": time.Now()},
	}

	_, err := groupCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to mark group as inactive")
	}

	return c.NoContent(http.StatusOK)
}

func GetGroupDetailsHandler(c echo.Context) error {
	// Parse group ID from request parameters
	groupID := c.Param("id")
	if groupID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "group ID is required")
	}

	// Convert group ID string to ObjectID
	objID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "invalid group ID format")
	}

	var group models.Group
	err = groupCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Echo().Logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusNotFound, "group not found")
		}
		c.Echo().Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch group details")
	}

	return c.JSON(http.StatusOK, group)
}
