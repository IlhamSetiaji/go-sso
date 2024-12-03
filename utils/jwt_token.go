package utils

import (
	"app/go-sso/internal/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GenerateToken(user *entity.User) (string, error) {
	viper := viper.New()
	logger := logrus.New()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()

	if err != nil {
		logger.Fatalf("Fatal error config file: %v", err)
	}

	// Prepare roles and permissions
	roles := make([]map[string]interface{}, len(user.Roles))
	for i, role := range user.Roles {
		permissions := make([]string, len(role.Permissions))
		for j, permission := range role.Permissions {
			permissions[j] = permission.Name
		}
		roles[i] = map[string]interface{}{
			"name":        role.Name,
			"permissions": permissions,
		}
	}

	// prepare token claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"name":     user.Name,
		"username": user.Username,
		"email":    user.Email,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateTokenForOAuth2(data *map[string]interface{}) (string, error) {
	viper := viper.New()
	logger := logrus.New()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()

	if err != nil {
		logger.Fatalf("Fatal error config file: %v", err)
	}

	// prepare token claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
