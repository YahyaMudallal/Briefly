package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// GenerateToken generates a JWT token for the given user ID.
func GenerateToken(userID bson.ObjectID) (string, error) {

	secret := os.Getenv("JWT_SERCRET")
	if secret == "" {
		secret = "default_secret_phrase" // remove in production
	}

	// define the token content
	claims := jwt.MapClaims{
		"sub" : userID.Hex(),
		"exp" : time.Now().Add(24 * 7 * time.Hour).Unix(), // token expires in 7 days
		"iat" : time.Now().Unix(),
	}

	// create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)


	// sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	
	return signedToken, nil
}