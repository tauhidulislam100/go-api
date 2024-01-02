package helpers

import (
	"context"
	"fmt"
	"gin-api/configs"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// for testing
type SignedDetails struct {
	Id       primitive.ObjectID `json:"id"`
	Email    string             `json:"email"`
	Username string             `json:"username"`
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "users")

func GenerateToken(id primitive.ObjectID, email string, username string) (token string, err error) {

	nowTime := time.Now().Local()
	expireTime := nowTime.Add(24 * time.Hour)

	mySigningKey := []byte(configs.EnvSecretKey())

	// Create claims while leaving out some of the optional fields
	claims := SignedDetails{
		id,
		email,
		username,
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "test",
		},
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := generateToken.SignedString(mySigningKey)
	fmt.Printf("%v %v", ss, err)
	return ss, err

}

func UpdateToken(signedToken string, userId primitive.ObjectID) {

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "id", Value: userId}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "token", Value: signedToken}}}}

	result, err := userCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}
}
