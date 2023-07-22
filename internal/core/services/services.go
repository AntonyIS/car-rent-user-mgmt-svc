package services

import (
	"github.com/AntonyIS/notlify-user-svc/internal/core/domain"
	"github.com/AntonyIS/notlify-user-svc/internal/core/ports"
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
	return svc.repo.CreateUser(user)
}

func (svc *UserManagementService) ReadUser(id string) (*domain.User, error) {
	return svc.repo.ReadUser(id)
}

func (svc *UserManagementService) ReadUsers() ([]*domain.User, error) {
	return svc.repo.ReadUsers()
}

func (svc *UserManagementService) UpdateUser(user *domain.User) (*domain.User, error) {
	return svc.repo.UpdateUser(user)
}

func (svc *UserManagementService) DeleteUser(id string) (string, error) {
	return svc.repo.DeleteUser(id)
}
