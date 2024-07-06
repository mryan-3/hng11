package utils

import (
    "fmt"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// Signs a JWT Token
func SignJwtToken(text string) (jwtToken string, err error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": text,
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
    })

    // Sign and get the complete encoded token as a string using the secret
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

    return tokenString, err
}

// Verifys a JWT Token
func VerifyJwtToken(tokenString string) (jwt.MapClaims, bool, error) {

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    claims, ok := token.Claims.(jwt.MapClaims)

    isValid := ok && token.Valid


    return claims, isValid, err
}

