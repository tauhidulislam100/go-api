package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(hashedPassword string, userPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		log.Panic(err)
		msg = "email of password is incorrect"
		check = false
	}

	return check, msg
}

func HashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	println(string(bytes))
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
