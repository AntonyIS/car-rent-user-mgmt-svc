package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AntonyIS/notelify-users-service/internal/core/domain"
	"github.com/AntonyIS/notelify-users-service/internal/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserManagementService struct {
	repo   ports.UserRepository
	logger ports.LoggingService
}

type loggingManagementService struct {
	loggerURL string
}

func NewUserManagementService(repo ports.UserRepository, logger ports.LoggingService) *UserManagementService {
	svc := UserManagementService{
		repo:   repo,
		logger: logger,
	}
	return &svc
}

func (svc *UserManagementService) CreateUser(user *domain.User) (*domain.User, error) {
	// Assign new user with a unique id
	if user.UserId == "" {
		user.UserId = uuid.New().String()
	}

	// hash user password
	if user.Password == "" {
		user.Password = "pass"
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		if err != nil {
			logEntry := domain.LogMessage{
				LogLevel: "ERROR",
				Service:  "users",
				Message:  err.Error(),
			}
			svc.logger.LogError(logEntry)
			return nil, err
		}
		return nil, err
	}
	user.Password = string(hashedPassword)
	user.Handle = fmt.Sprintf(`%s@notelify`, user.Firstname)
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "users",
		Message:  fmt.Sprintf("User with ID [%s] created successfuly", user.UserId),
	}
	svc.logger.LogInfo(logEntry)
	return svc.repo.CreateUser(user)
}

func (svc *UserManagementService) ReadUserWithId(user_id string) (*domain.User, error) {
	user, err := svc.repo.ReadUserWithId(user_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return nil, err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "users",
		Message:  fmt.Sprintf("User with ID [%s] created successfuly", user.UserId),
	}
	svc.logger.LogInfo(logEntry)
	return user, nil
}

func (svc *UserManagementService) ReadUserWithEmail(email string) (*domain.User, error) {
	user, err := svc.repo.ReadUserWithEmail(email)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return nil, err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "users",
		Message:  fmt.Sprintf("User with EMAIL [%s] created successfuly", user.Email),
	}
	svc.logger.LogInfo(logEntry)
	return user, nil
}

func (svc *UserManagementService) ReadUsers() ([]domain.User, error) {

	users, err := svc.repo.ReadUsers()
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return nil, err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "users",
		Message:  "Users found successfuly",
	}
	fmt.Println("Users", users)
	svc.logger.LogInfo(logEntry)
	return users, nil
}

func (svc *UserManagementService) UpdateUser(user *domain.User) (*domain.User, error) {
	user, err := svc.repo.UpdateUser(user)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return nil, err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "users",
		Message:  fmt.Sprintf("User with ID [%s] updated successfuly", user.UserId),
	}
	svc.logger.LogInfo(logEntry)
	return user, nil
}

func (svc *UserManagementService) DeleteUser(user_id string) (string, error) {
	message, err := svc.repo.DeleteUser(user_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return "", err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "users",
		Message:  fmt.Sprintf("User with ID [%s] deleted successfuly", user_id),
	}
	svc.logger.LogInfo(logEntry)
	return message, nil
}

func (svc *UserManagementService) DeleteAllUsers() (string, error) {
	message, err := svc.repo.DeleteAllUsers()
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return "", err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "users",
		Message:  "Users deleted successfuly",
	}
	svc.logger.LogInfo(logEntry)
	return message, nil
}

func NewLoggingManagementService(loggerURL string) *loggingManagementService {
	svc := loggingManagementService{
		loggerURL: loggerURL,
	}
	return &svc
}

func (svc *loggingManagementService) SendLog(logEntry domain.LogMessage) {
	// Marshal the struct into JSON
	payloadBytes, err := json.Marshal(logEntry)
	if err != nil {
		fmt.Println("Error encoding JSON payload:", err)
		return
	}
	// Create a new POST request with the JSON payload
	resp, err := http.Post(svc.loggerURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()
}

func (svc *loggingManagementService) LogDebug(logEntry domain.LogMessage) {
	message := fmt.Sprintf("[%s] [DEBUG] %s %s", logEntry.Service, getCurrentDateTime(), logEntry.Message)
	logEntry.Message = message
	svc.SendLog(logEntry)
}

func (svc *loggingManagementService) LogInfo(logEntry domain.LogMessage) {
	message := fmt.Sprintf("[%s] [INFO] %s %s", logEntry.Service, getCurrentDateTime(), logEntry.Message)
	logEntry.Message = message
	svc.SendLog(logEntry)
}

func (svc *loggingManagementService) LogWarning(logEntry domain.LogMessage) {
	message := fmt.Sprintf("[%s] [WARNING] %s %s", logEntry.Service, getCurrentDateTime(), logEntry.Message)
	logEntry.Message = message
	svc.SendLog(logEntry)
}

func (svc *loggingManagementService) LogError(logEntry domain.LogMessage) {
	message := fmt.Sprintf("[%s] [ERROR] %s %s", logEntry.Service, getCurrentDateTime(), logEntry.Message)
	logEntry.Message = message
	svc.SendLog(logEntry)
}

func getCurrentDateTime() string {
	currentTime := time.Now()
	return currentTime.Format("2006/01/02 15:04:05")
}
