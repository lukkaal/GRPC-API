package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

// viper already inited in ./config/config.go
var jwtSecret = viper.GetString("server.jwtSecret")

// base model + userid -> jwt token
type Claims struct {
	// information for jwt token
	jwt.RegisteredClaims
	UserId int64 `json:"user_id"`
}

// generate token
func GenerateToken(userid int64) (string, error) {
	now := time.Now()
	expiretime := now.Add(time.Hour * 24)

	// token validation
	claim := &Claims{
		UserId: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				expiretime),
			Issuer: "38384-SearchEngine",
		},
	}
	// assign the algorithm HS256
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// encrypt to string using jwtsecret
	token, err := newToken.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseToken(token string) (*Claims, error) {
	// verify and fill the Claims struct
	// get *jwt.token
	jwttoken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// verify jwttoken
	if jwttoken != nil {
		// assert and validate the jwttoken
		if claims, ok := jwttoken.Claims.(*Claims); ok && jwttoken.Valid {
			return claims, nil
		}
	}
	return nil, err
}
