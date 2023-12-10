package controllers

import (
	"context"
	"gin-api/configs"
	"gin-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var roleCollection *mongo.Collection = configs.GetCollection(configs.DB, "roles")

func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	role.ID = uuid.New().String()
	role.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	role.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err := roleCollection.InsertOne(context.Background(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created",
		"data":    role,
	})
}

func GetAllRoles(c *gin.Context) {
	cur, err := roleCollection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}
	defer cur.Close(context.Background())

	var roles []models.Role
	if err := cur.All(context.Background(), &roles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding roles"})
		return
	}

	if len(roles) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No Data Roles",
			"data":    []models.Role{}, // Empty slice of models.User
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get Roles",
		"data":    roles,
	})
}
