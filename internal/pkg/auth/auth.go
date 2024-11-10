package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// user authorizations settings
const (
	// TokenExpired token time life
	TokenExpired = time.Hour * 24
	// SecretKey salt
	SecretKey = "7y8i24^&&(G*WEFOo23euh2o3gnyutFUDEFnopo2efnTD##@k"
)

var (
	// ErrUUIDNotExists raized if UUID not exists in claims
	ErrUUIDNotExists = errors.New("user UUID not exists in claims")
)

// Claims jwt struct
type Claims struct {
	jwt.RegisteredClaims
	UserUUID string
}

// BuildJWT create jwt token
func BuildJWT(userUUID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpired)),
		},
		UserUUID: userUUID,
	})
	tk, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tk, nil
}

// GetUserUUIDFromJWT return user ID from token
func GetUserUUIDFromJWT(tk string) (string, error) {
	claims, err := getClaims(tk)
	if err != nil {
		return "", err
	}
	if claims.UserUUID == "" {
		return "", ErrUUIDNotExists
	}
	logger.Log.Infof("app user ID: %s", claims.UserUUID)
	return claims.UserUUID, nil
}

func getClaims(tk string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tk, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil || !token.Valid {
		logger.Log.Error(err.Error())
		return nil, err
	}
	return claims, nil
}

// CreateUUID return UUID v5
func CreateUUID(name string) string {
	return uuid.NewV5(uuid.NewV4(), name).String()
}

// CreateHashPassword generating hash for password
func CreateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CompareHashPassword comparing has with password
func CompareHashPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
