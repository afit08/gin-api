package controllers

import (
	"context"
	"fmt"
	"gin-api/configs"
	productsDTO "gin-api/dto/Products"
	"gin-api/helpers"
	"gin-api/models"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection = configs.GetCollection(configs.DB, "products")

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
func CreateProduct(c *gin.Context) {

	var product productsDTO.ProductRequest

	// Use ShouldBind instead of ShouldBindJSON for form data
	if err := c.ShouldBind(&product); err != nil {
		// Use StatusJSON to set both HTTP status and JSON response
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid form data: %s", err.Error())})
		return
	}

	// Assign other fields and generate ID
	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.Image = product.Image * multipart.FileHeader.Filename
	// Upload gambar ke MinIO
	err := helpers.UploadImageToMinio(product.Image, product.Image.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading image to MinIO"})
		return
	}
	// Insert the product into the database
	_, err = productCollection.InsertOne(context.Background(), product)
	if err != nil {
		// Use StatusJSON for consistent response format
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		return
	}

	// Use StatusJSON for consistent response format
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created",
		"data":    product,
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
func AllProduct(c *gin.Context) {
	cur, err := productCollection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
		return
	}
	defer cur.Close(context.Background())

	var products []models.Product
	if err := cur.All(context.Background(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding products"})
		return
	}
	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No Data Product",
			"data":    []models.Product{}, // Empty slice of models.User
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get Roles",
		"data":    products,
	})
}
