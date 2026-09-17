package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arunkumar-1311/mongo-db/auth"
	"github.com/arunkumar-1311/mongo-db/models"
	"github.com/arunkumar-1311/mongo-db/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Service interface {
	CreateUsers(ctx context.Context, req models.User) (models.User, error)
	GetUser(ctx context.Context, objID primitive.ObjectID) (models.User, error)
	Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error)
}

type userService struct {
	repo repository.Repository
}

func NewUserService() Service {
	return userService{repo: repository.NewUserRepository()}
}

// CreateUsers creates a Department, a Role and (optionally) one or more
// Orders, then creates the User document referencing all of them by
// ObjectID. If any step after the first succeeds and a later step fails,
// the already-created related documents are rolled back so no orphaned
// records are left behind.
func (s userService) CreateUsers(ctx context.Context, req models.User) (models.User, error) {
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("unable to hash password: %w", err)
	}
	req.PasswordHash = passwordHash
	req.Password = ""

	departmentId, err := s.repo.CreateDepartment(ctx, req.Department)
	if err != nil {
		return models.User{}, fmt.Errorf("unable to create department: %w", err)
	}

	roleId, err := s.repo.CreateRole(ctx, req.Role)
	if err != nil {
		_ = s.repo.DeleteDepartment(ctx, departmentId)
		return models.User{}, fmt.Errorf("unable to create role: %w", err)
	}

	orderIds, err := s.repo.CreateOrders(ctx, req.Orders)
	if err != nil {
		_ = s.repo.DeleteDepartment(ctx, departmentId)
		_ = s.repo.DeleteRole(ctx, roleId)
		return models.User{}, fmt.Errorf("unable to create orders: %w", err)
	}

	req.DepartmentId = departmentId
	req.RoleId = roleId
	req.OrderIds = orderIds

	userId, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		_ = s.repo.DeleteDepartment(ctx, departmentId)
		_ = s.repo.DeleteRole(ctx, roleId)
		_ = s.repo.DeleteOrders(ctx, orderIds)
		return models.User{}, fmt.Errorf("unable to create user: %w", err)
	}

	req.Id = userId
	req.DepartmentId = departmentId
	req.RoleId = roleId
	req.OrderIds = orderIds
	return req, nil
}

func (s userService) GetUser(ctx context.Context, objID primitive.ObjectID) (models.User, error) {
	result, err := s.repo.GetUser(ctx, objID)
	if err != nil {
		return models.User{}, err
	}
	return result, nil
}

func (s userService) Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return models.LoginResponse{}, ErrInvalidCredentials
	}

	if err := auth.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return models.LoginResponse{}, ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.Id, user.Email, user.Role.Level)
	if err != nil {
		return models.LoginResponse{}, fmt.Errorf("unable to generate token: %w", err)
	}

	return models.LoginResponse{Token: token}, nil
}
