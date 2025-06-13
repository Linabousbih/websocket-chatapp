package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
)
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.StandardClaims
}
var jwtKey []byte

func SetJWTKey(key string) {
	jwtKey = []byte(key)
}

func GetJWTKey() []byte {
	return []byte(jwtKey)
}

func GenerateTokens(email, userID string) (string, string, error) {
    tokenExpiry := time.Now().Add(24 * time.Hour).Unix()
    refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour).Unix()

    claims := &Claims{
        Email:  email,
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: tokenExpiry,
        },
    }

    refreshClaims := &Claims{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: refreshTokenExpiry,
        },
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedAccessToken, err := accessToken.SignedString(jwtKey)
    if err != nil {
        return "", "", err
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    signedRefreshToken, err := refreshToken.SignedString(jwtKey)
    if err != nil {
        return "", "", err
    }

    return signedAccessToken, signedRefreshToken, nil
}

func HashPassword(password *string) (*string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    hashedPwd := string(bytes)
    return &hashedPwd, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	// Use the dynamically set JWT key here
	secretKey := GetJWTKey() // This retrieves the key set in SetJWTKey

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid and return the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func UpdateAllTokens(signedToken, signedRefreshToken, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	userCollection := OpenCollection("users")
	//Crate an update object

	updateObj := bson.D{
		{"$set", bson.D{
			{"token", signedToken},
			{"refresh_token", signedRefreshToken},
			{"updated_at", time.Now()},
		}},
	}

	//Create a filter
	filter := bson.M{"user_id": userID}

	// upsert := true
	// opt := options.UpdateOptions{
	// 	Upsert: &upsert,
	// }

	_, err := userCollection.UpdateOne(ctx, filter, updateObj)

	return err
}

func VerifyPassword(foundPwd, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(foundPwd), []byte(pwd))

	return err == nil, err
}