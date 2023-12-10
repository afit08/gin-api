// controllers/user_controller.go
package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"gin-api/configs"
	"gin-api/helpers"
	"gin-api/models"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	return err == nil
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	password := HashPassword(user.Password)

	user.ID = uuid.New().String()
	user.Password = password
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created", "data": user})
}

// Login is the api used to tget a single user
func Login(c *gin.Context) {
	request := new(LoginRequest)
	if err := c.Bind(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	user := models.User{
		Username: request.Username,
		Password: request.Password,
	}

	err := userCollection.FindOne(context.Background(), bson.M{"username": request.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password is incorrect"})
		return
	}

	var role models.Role
	err = roleCollection.FindOne(context.Background(), bson.M{"_id": user.Role_id}).Decode(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Role not found"})
		return
	}

	passwordIsValid := VerifyPassword(request.Password, user.Password)
	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Invalid password",
		})
		return
	}

	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["roleType"] = role.Name
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix() // expired in 3 hours

	// jsonClaims, _ := json.Marshal(claims)
	// fmt.Println(string(jsonClaims))

	token, errGenerateToken := helpers.GenerateAllTokens(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		c.Status(http.StatusUnauthorized)
		return
	}

	c.SetCookie("jwt", token, 86400, "/", "localhost", false, true)

	loginResponse := loginResponse{
		Username: user.Username,
		Token:    token,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Successfully!!!",
		"status":  http.StatusOK,
		"data":    loginResponse,
	})
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

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User deleted. %d documents deleted", result.DeletedCount),
	})
}

func OneUsersHandler(c *gin.Context) {
	id := c.Param("id")

	query := []bson.M{
		{
			"$match": bson.M{"_id": id},
		},
		{
			"$lookup": bson.M{
				"from":         "roles",
				"localField":   "role_id",
				"foreignField": "_id",
				"as":           "user_roles",
			},
		},
		{
			"$unwind": "$user_roles",
		},
		{
			"$project": bson.M{
				"_id":           0,
				"user_id":       "$_id",
				"user_name":     "$name",
				"user_password": "$password",
				"role_id":       "$role_id",
				"role_name":     "$user_roles.name",
			},
		},
	}

	cur, err := userCollection.Aggregate(context.TODO(), query)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cur.Close(context.TODO())

	var result []bson.M
	if err := cur.All(context.TODO(), &result); err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Find One Users",
		"data":    result,
	})
}
