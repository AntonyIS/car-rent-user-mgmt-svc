package services

import (
	"fmt"

	"github.com/AntonyIS/notelify-users-service/internal/core/domain"
	"github.com/AntonyIS/notelify-users-service/internal/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserManagementService struct {
	repo   ports.UserRepository
	logger ports.LoggingService
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
	user.UserId = uuid.New().String()
	// hash user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		if err != nil {
			logEntry := domain.LogMessage{
				LogLevel: "critical",
				Service:  "users",
				Message:  err.Error(),
			}
			svc.logger.CreateLog(logEntry)
			return nil, err
		}
		return nil, err
	}
	user.Password = string(hashedPassword)
	user.Handle = fmt.Sprintf(`%s@notelify`, user.Firstname)
	return svc.repo.CreateUser(user)

}

func (svc *UserManagementService) ReadUserWithId(user_id string) (*domain.User, error) {
	user, err := svc.repo.ReadUserWithId(user_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return nil, err
	}
	return user, nil
}

func (svc *UserManagementService) ReadUserWithEmail(email string) (*domain.User, error) {
	user, err := svc.repo.ReadUserWithEmail(email)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return nil, err
	}
	return user, nil
}

func (svc *UserManagementService) ReadUsers() ([]domain.User, error) {

	users, err := svc.repo.ReadUsers()
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return nil, err
	}
	return users, nil
}

func (svc *UserManagementService) UpdateUser(user *domain.User) (*domain.User, error) {
	user, err := svc.repo.UpdateUser(user)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return nil, err
	}
	return user, nil
}

func (svc *UserManagementService) DeleteUser(user_id string) (string, error) {
	message, err := svc.repo.DeleteUser(user_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return "", err
	}
	return message, nil
}

func (svc *UserManagementService) DeleteAllUsers() (string, error) {
	message, err := svc.repo.DeleteAllUsers()
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "users",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return "", err
	}
	return message, nil
}
