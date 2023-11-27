package services

import (
	"github.com/AntonyIS/notelify-users-service/internal/core/domain"
	"github.com/AntonyIS/notelify-users-service/internal/core/ports"
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
	user.UserId = uuid.New().String()
	// hash user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	return svc.repo.CreateUser(user)
}

func (svc *UserManagementService) ReadUserWithId(user_id string) (*domain.User, error) {
	return svc.repo.ReadUserWithId(user_id)
}

func (svc *UserManagementService) ReadUserWithEmail(email string) (*domain.User, error) {
	return svc.repo.ReadUserWithEmail(email)
}

func (svc *UserManagementService) ReadUsers() ([]domain.User, error) {
	return svc.repo.ReadUsers()
}

func (svc *UserManagementService) UpdateUser(user *domain.User) (*domain.User, error) {
	return svc.repo.UpdateUser(user)
}

func (svc *UserManagementService) DeleteUser(user_id string) (string, error) {
	return svc.repo.DeleteUser(user_id)
}

func (svc *UserManagementService) DeleteAllUsers() (string, error) {
	return svc.repo.DeleteAllUsers()
}
