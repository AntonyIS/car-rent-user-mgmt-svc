package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AntonyIS/notelify-users-service/internal/core/domain"
	"github.com/AntonyIS/notelify-users-service/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type middleware struct {
	svc       ports.UserService
	logger    ports.LoggingService
	secretKey string
}

func NewMiddleware(svc ports.UserService, logger ports.LoggingService, secretKey string) *middleware {

	return &middleware{
		svc:       svc,
		logger:    logger,
		secretKey: secretKey,
	}
}

func (m middleware) GenerateToken(user_id string) (string, error) {
	key := []byte(m.secretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	user, err := m.svc.ReadUserWithId(user_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "users",
			Message:  err.Error(),
		}
		m.logger.LogError(logEntry)
		return "", err
	}

	claims["user_id"] = user.UserId
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(key)

	if err != nil {
		if err != nil {
			logEntry := domain.LogMessage{
				LogLevel: "ERROR",
				Service:  "users",
				Message:  err.Error(),
			}
			m.logger.LogError(logEntry)
			return "", err
		}
		return "", err
	}
	return tokenString, nil
}

func (m middleware) Authorize(c *gin.Context) {
	tokenString := c.GetHeader("token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logEntry := domain.LogMessage{
				LogLevel: "ERROR",
				Service:  "users",
				Message:  fmt.Sprintf("unexpected signing method: %v", token.Header["sub"]),
			}
			m.logger.LogError(logEntry)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		if err != nil {
			logEntry := domain.LogMessage{
				LogLevel: "ERROR",
				Service:  "users",
				Message:  err.Error(),
			}
			m.logger.LogError(logEntry)

		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	} else {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "users",
			Message:  "request not authorized",
		}
		m.logger.LogError(logEntry)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New("request not authorized"),
		})
		return
	}
}
