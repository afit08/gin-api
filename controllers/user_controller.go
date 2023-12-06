// controllers/user_controller.go
package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"gin-api/configs"
	"gin-api/models"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	user.ID = uuid.New().String()

	_, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "data": user})
}

// GetUsers returns all users
func GetUsers(c *gin.Context) {
	cur, err := userCollection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}
	defer cur.Close(context.Background())

	var users []models.User
	if err := cur.All(context.Background(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding users"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No Data Users",
			"data":    []models.User{}, // Empty slice of models.User
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get Users",
		"data":    users,
	})
}

// GetUserByID returns a user by ID
func GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	filter := bson.M{"_id": userID}
	var user models.User
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		// Check if the user is not found
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get Users By id",
		"data":    user,
	})
}

// UpdateUser updates a user by ID
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var updateUser models.User
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": updateUser}
	result, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated",
		"data":    updateUser,
	})
}

// DeleteUser deletes a user by ID
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	filter := bson.M{"_id": userID}
	result, err := userCollection.DeleteOne(context.Background(), filter)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User deleted. %d documents deleted", result.DeletedCount)})
}
