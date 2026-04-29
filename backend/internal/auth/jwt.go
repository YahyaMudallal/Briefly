package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// GenerateToken generates a JWT token for the given user ID.
func GenerateToken(userID bson.ObjectID) (string, error) {

	// get the secret key
	secret := os.Getenv("JWT_SECRET")
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

// VerifyToken verifies the JWT token returns an error if the token is invalid.
func VerifyToken(tokenString string) (bson.ObjectID, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" { secret = "default_secret_phrase" }

	// parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signature method wanted : %v", token.Header["alg"])
		}
		return []byte(secret), nil
    })

	if err != nil || !token.Valid {
		return bson.NilObjectID, fmt.Errorf("token invalid: %v", err)
	}

	// extract userID
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        sub, _ := claims["sub"].(string)
        return bson.ObjectIDFromHex(sub)
    }

    return bson.NilObjectID, fmt.Errorf("invalid claims")
}