package controllers

import (
	"context"
	"fmt"
	"gin-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		results, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}
		defer results.Close(ctx)

		var users []models.User
		for results.Next(ctx) {
			var user models.User
			if err = results.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "success", "data": users})
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		userId, _ := primitive.ObjectIDFromHex(c.Param("id"))

		var user bson.M
		fmt.Printf("userId %v \n", userId)
		err := userCollection.FindOne(ctx, bson.M{"id": userId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "success", "data": user})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userId, _ := primitive.ObjectIDFromHex(c.Param("id"))

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err.Error()})
			return
		}

		fmt.Printf("The request body is %v \n", user)

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err.Error()})
			return
		}
		fmt.Printf("The request body is %v", user)
		result, err := userCollection.UpdateOne(ctx, bson.M{"id": userId}, bson.M{"$set": user})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "error", "data": result})
	}
}
