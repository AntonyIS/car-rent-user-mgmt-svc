package service

import (
	"errors"

	"github.com/AntonyIS/car-rent-user-mgmt-svc/internal/core/domain"
	"github.com/AntonyIS/car-rent-user-mgmt-svc/internal/core/ports"
	"github.com/google/uuid"
)

type UserManagementService struct {
	repo ports.UserManagementRepository
}

func NewUserManagementService(repo *ports.UserManagementRepository) *UserManagementService {
	return &UserManagementService{
		repo: *repo,
	}
}

func (svc *UserManagementService) CreateUser(user *domain.User) (*domain.User, error) {
	// Check if user exists using user email
	users, err := svc.repo.ReadUsers()
	if err != nil {
		return nil, err
	}
	// Loop through existing users and get the user with email
	for _ , item := range users {
		if item.Email == user.Email {
			return nil, errors.New("user with email exists")
		}
	}
	// Create new user ID
	user.ID = uuid.New().String()

	return svc.CreateUser(user)
}

func (svc *UserManagementService) ReadUser(id string) (*domain.User, error) {
	return svc.ReadUser(id)
}

func (svc *UserManagementService) ReadUsers() ([]*domain.User, error) {
	return svc.ReadUsers()
}

func (svc *UserManagementService) UpdateUsers(user *domain.User) (*domain.User, error) {
	return svc.UpdateUsers(user)
}

func (svc *UserManagementService) DeleteUser(id string) error {
	return svc.DeleteUser(id)
}
