package controllers

import (
	"context"
	"fmt"
	"gin-api/configs"
	"gin-api/helpers"
	"gin-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoriCollection *mongo.Collection = configs.GetCollection(configs.DB, "categories")

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Show an account
// @Description get string by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Success 200 {object} model.Account
// @Failure 400 {object} model.HTTPError
// @Router /accounts/{id} [get]
func CreateCategori(c *gin.Context) {
	var categori models.CreateCategoriRequest

	// Use ShouldBind instead of ShouldBindJSON for form data
	if err := c.ShouldBind(&categori); err != nil {
		// Use StatusJSON to set both HTTP status and JSON response
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid form data: %s", err.Error())})
		return
	}

	// Assign other fields and generate ID
	categori.ID = uuid.New().String()
	categori.CreatedAt = time.Now()
	categori.UpdatedAt = time.Now()
	// Upload gambar ke MinIO
	err := helpers.UploadImageToMinio(categori.Image, categori.Image.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading image to MinIO"})
		return
	}
	// Insert the categori into the database
	_, err = categoriCollection.InsertOne(context.Background(), categori)
	if err != nil {
		// Use StatusJSON for consistent response format
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		return
	}

	result := gin.H{
		"id":         categori.ID,
		"name":       categori.Name,
		"image":      categori.Image.Filename,
		"created_at": categori.CreatedAt,
		"update_at":  categori.UpdatedAt,
	}

	// Use StatusJSON for consistent response format
	c.JSON(http.StatusCreated, gin.H{
		"message": "Categories Created",
		"data":    result,
	})
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Show an account
// @Description get string by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Success 200 {object} model.Account
// @Failure 400 {object} model.HTTPError
// @Router /accounts/{id} [get]
func AllCategories(c *gin.Context) {
	var categories []models.Categori

	cur, err := categoriCollection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching categories"})
		return
	}
	defer cur.Close(context.Background())

	if err := cur.All(context.Background(), &categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding categories"})
		return
	}

	var result []gin.H

	for _, v := range categories {
		data := gin.H{
			"id":    v.ID,
			"name":  v.Name,
			"image": v.Image.Filename,
		}
		result = append(result, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get All Categories",
		"data":    result,
	})
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Show an account
// @Description get string by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Success 200 {object} model.Account
// @Failure 400 {object} model.HTTPError
// @Router /accounts/{id} [get]
func OneCategori(c *gin.Context) {
	categoriID := c.Param("id")

	filter := bson.M{"_id": categoriID}
	var categories models.Categori
	err := categoriCollection.FindOne(context.Background(), filter).Decode(&categories)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching categori"})
		return
	}

	result := gin.H{
		"id":    categories.ID,
		"name":  categories.Name,
		"image": categories.Image.Filename,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get One Categories",
		"data":    result,
	})
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Show an account
// @Description get string by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Success 200 {object} model.Account
// @Failure 400 {object} model.HTTPError
// @Router /accounts/{id} [get]
func UpdateCategori(c *gin.Context) {

}
