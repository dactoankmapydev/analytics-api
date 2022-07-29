package security

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func ExtractAccessToken(r *http.Request) string {
	atCookie, err := r.Cookie("access_token")
	if err != nil {
		return ""
	}
	return atCookie.Value
}

func VerifyAccessToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractAccessToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func AccessTokenValid(r *http.Request) error {
	token, err := VerifyAccessToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractAccessTokenMetadata(r *http.Request) (*TokenDetails, error) {
	token, err := VerifyAccessToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}
		return &TokenDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}
