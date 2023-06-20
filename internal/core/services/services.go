package service

import (
	"github.com/AntonyIS/car-rent-user-mgmt-svc/internal/core/domain"
	"github.com/AntonyIS/car-rent-user-mgmt-svc/internal/core/ports"
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
