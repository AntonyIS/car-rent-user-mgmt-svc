package services

import (
	"errors"

	"github.com/AntonyIS/notlify-user-svc/internal/core/domain"
	"github.com/AntonyIS/notlify-user-svc/internal/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserManagementService struct {
	repo ports.UserRepository
}

func NewUserManagementService(repo ports.UserRepository) *UserManagementService {
	svc := UserManagementService{
		repo: repo,
	}

	return &svc
}

func (svc *UserManagementService) CreateUser(user *domain.User) (*domain.User, error) {
	// Assign new user with a unique id
	user.Id = uuid.New().String()
	// hash user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	return svc.repo.CreateUser(user)
}

func (svc *UserManagementService) ReadUser(id string) (*domain.User, error) {
	return svc.repo.ReadUser(id)
}

func (svc *UserManagementService) ReadUserWithEmail(email string) (*domain.User, error) {
	return svc.repo.ReadUser(email)
}

func (svc *UserManagementService) ReadUsers() ([]domain.User, error) {
	return svc.repo.ReadUsers()
}

func (svc *UserManagementService) UpdateUser(user *domain.User) (*domain.User, error) {
	return svc.repo.UpdateUser(user)
}

func (svc *UserManagementService) DeleteUser(id string) (string, error) {
	// Check if user exists
	_, err := svc.ReadUser(id)
	if err != nil {
		return " ", errors.New("Error, item not found!")
	}

	return svc.repo.DeleteUser(id)
}
