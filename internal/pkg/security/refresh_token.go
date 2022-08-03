package security

// import (
// 	"fmt"
// 	"net/http"
// 	"os"

// 	"github.com/golang-jwt/jwt"
// )

// func ExtractRefreshToken(r *http.Request) string {
// 	rtCookie, err := r.Cookie("refresh_token")
// 	if err != nil {
// 		return ""
// 	}
// 	return rtCookie.Value
// }

// func VerifyRefreshToken(r *http.Request) (*jwt.Token, error) {
// 	tokenString := ExtractRefreshToken(r)
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("REFRESH_SECRET")), nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return token, nil
// }

// func RefreshTokenValid(r *http.Request) error {
// 	token, err := VerifyRefreshToken(r)
// 	if err != nil {
// 		return err
// 	}
// 	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
// 		return err
// 	}
// 	return nil
// }

// func ExtractRefreshTokenMetadata(r *http.Request) (*TokenDetails, error) {
// 	token, err := VerifyRefreshToken(r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if ok && token.Valid {
// 		refreshUUID, ok := claims["refresh_uuid"].(string)
// 		if !ok {
// 			return nil, err
// 		}
// 		userID, ok := claims["user_id"].(string)
// 		if !ok {
// 			return nil, err
// 		}
// 		return &TokenDetails{
// 			RefreshUUID: refreshUUID,
// 			UserID:      userID,
// 		}, nil
// 	}
// 	return nil, err
// }
