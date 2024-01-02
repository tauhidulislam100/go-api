package controllers

import (
	"context"
	"gin-api/configs"
	"gin-api/helpers"
	"gin-api/models"
	"gin-api/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginRequset struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

var userCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "users")
var validate *validator.Validate = validator.New()

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()
		//validate the request body
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err.Error()})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err.Error()})
			log.Panic(err)
			return
		}

		emailCounts, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}
		if emailCounts > 0 {
			c.JSON(http.StatusConflict, gin.H{"status": http.StatusBadRequest, "message": "Email already exist."})
			return
		}

		userNameExist, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}

		if userNameExist > 0 {
			c.JSON(http.StatusConflict, gin.H{"status": http.StatusBadRequest, "message": "User name already exist."})
			return
		}

		newUser := models.User{
			Id:          primitive.NewObjectID(),
			Username:    user.Username,
			Email:       user.Email,
			Password:    utils.HashPassword(user.Password),
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Birthday:    user.Birthday,
			Gender:      user.Gender,
			Address:     user.Address,
			PhoneNumber: user.PhoneNumber,
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "User created successfully.", "data": result})
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user LoginRequset
		var foundUser models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err})
			return
		}

		if err := validate.Struct(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Email or password is incorrect"})
			return
		}

		passwordIsValid, msg := utils.VerifyPassword(foundUser.Password, user.Password)
		defer cancel()

		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": msg})
			return
		}

		token, _ := helpers.GenerateToken(foundUser.Id, foundUser.Email, foundUser.Username)
		defer cancel()

		helpers.UpdateToken(token, foundUser.Id)
		foundError := userCollection.FindOne(ctx, bson.M{"id": foundUser.Id}).Decode(&foundUser)
		defer cancel()

		if foundError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "success", "data": foundUser})

	}
}
