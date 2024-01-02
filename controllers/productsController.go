package controllers

import (
	"context"
	"gin-api/configs"
	"gin-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "products")

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var product models.Product

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err.Error()})
			return
		}

		if err := validate.Struct(product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err.Error()})
			return
		}

		newProduct := models.Product{
			Id:          primitive.NewObjectID(),
			Name:        product.Name,
			Description: product.Description,
			Type:        product.Type,
			Sku:         product.Sku,
			ImgUrl:      product.Sku,
			AddedBy:     product.AddedBy,
		}

		result, err := productCollection.InsertOne(ctx, newProduct)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "success", "data": result})
	}
}

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := productCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}
		var products []bson.M
		if err = cursor.All(ctx, &products); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "error", "data": products})
	}
}

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		productId, _ := primitive.ObjectIDFromHex(c.Param("id"))

		var product bson.M
		err := productCollection.FindOne(ctx, bson.M{"id": productId}).Decode(&product)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error", "data": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "success", "data": product})
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		productId, _ := primitive.ObjectIDFromHex(c.Param("id"))

		var product models.User

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err.Error()})
			return
		}

		if err := validate.Struct(product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err.Error()})
			return
		}
		result, err := userCollection.UpdateOne(ctx, bson.M{"id": productId}, bson.M{"$set": product})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusBadRequest, "message": "error", "data": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "error", "data": result})
	}
}
