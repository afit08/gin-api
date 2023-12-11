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

	var product models.CreateProductRequest

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

	result := gin.H{
		"id":         product.ID,
		"name":       product.Name,
		"image":      product.Image.Filename,
		"price":      product.Price,
		"desc":       product.Desc,
		"stock":      product.Stock,
		"weight":     product.Weight,
		"created_at": product.CreatedAt,
		"update_at":  product.UpdatedAt,
	}

	// Use StatusJSON for consistent response format
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created",
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
func AllProduct(c *gin.Context) {
	var products []models.Product

	cur, err := productCollection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
		return
	}
	defer cur.Close(context.Background())

	if err := cur.All(context.Background(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding products"})
		return
	}

	var result []gin.H

	for _, v := range products {
		data := gin.H{
			"id":     v.ID,
			"name":   v.Name,
			"image":  v.Image.Filename,
			"price":  v.Price,
			"desc":   v.Desc,
			"stock":  v.Stock,
			"weight": v.Weight,
		}
		result = append(result, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get All Products",
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
func OneProduct(c *gin.Context) {
	productID := c.Param("id")

	filter := bson.M{"_id": productID}
	var products models.Product
	err := productCollection.FindOne(context.Background(), filter).Decode(&products)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching product"})
		return
	}

	result := gin.H{
		"id":     products.ID,
		"name":   products.Name,
		"image":  products.Image.Filename,
		"price":  products.Price,
		"desc":   products.Desc,
		"stock":  products.Stock,
		"weight": products.Weight,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get One Product",
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
func UpdateProduct(c *gin.Context) {
	productID := c.Param("id")

	// Convert productID to UUID
	uuid, err := uuid.Parse(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.UpdateProductRequest
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the ID field with the parsed UUID
	product.ID = uuid.String()

	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": product}

	// Upload image to MinIO
	if product.Image != nil {
		err := helpers.UploadImageToMinio(product.Image, product.Image.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading image to MinIO"})
			return
		}
	}

	_, err = productCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating product"})
		return
	}

	result := gin.H{
		"id":     product.ID,
		"name":   product.Name,
		"image":  product.Image.Filename,
		"price":  product.Price,
		"desc":   product.Desc,
		"stock":  product.Stock,
		"weight": product.Weight,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated",
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
func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	filter := bson.M{"_id": productID}
	result, err := productCollection.DeleteOne(context.Background(), filter)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User deleted. %d documents deleted", result.DeletedCount),
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
// func Transaction(c *gin.Context) {
// 	var product productsDTO.ProductRequest

// 	// Use ShouldBind instead of ShouldBindJSON for form data
// 	if err := c.ShouldBind(&product); err != nil {
// 		// Use StatusJSON to set both HTTP status and JSON response
// 		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid form data: %s", err.Error())})
// 		return
// 	}

// 	// Start a MongoDB transaction
// 	session, err := configs.ConnectDB().StartSession()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting session"})
// 		return
// 	}
// 	defer session.EndSession(context.Background())

// 	// Use WithTransaction to execute the operations within a transaction
// 	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
// 		// Assign other fields and generate ID
// 		product.ID = uuid.New().String()
// 		product.CreatedAt = time.Now()
// 		product.UpdatedAt = time.Now()

// 		// Upload gambar ke MinIO
// 		err := helpers.UploadImageToMinio(product.Image, product.Image.Filename)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Insert the product into the database
// 		_, err = productCollection.InsertOne(sessCtx, product)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Return a dummy value (can be nil) to satisfy the (interface{}, error) signature
// 		return nil, nil
// 	})

// 	if err != nil {
// 		// Handle transaction error
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product in transaction"})
// 		return
// 	}

// 	// Use StatusJSON for consistent response format
// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "Product created",
// 		"data":    product,
// 	})
// }
