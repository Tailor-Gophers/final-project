package utils

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"

	"golang.org/x/crypto/bcrypt"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// HashPassword hashes the given password and returns the hashed password and error
func HashPassword(pass string) (string, error) {
	if len(pass) == 0 {
		return "", errors.New("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hash), err
}

func HashToken(pass string) (string, error) {
	// if len(pass) == 0 {
	// 	return "", errors.New("password cannot be empty")
	// }

	// h := sha256.New()
	// h.Write([]byte(pass))
	// // Calculate and print the hash
	// s := fmt.Sprintf("%x", h.Sum(nil))
	return pass, nil
}

// ValidatePassword compares given password with hashed password
func ValidatePassword(givenPass, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(givenPass), []byte(pass)) == nil
}

func ENV(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
