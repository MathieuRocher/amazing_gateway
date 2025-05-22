package user

import (
	domain "github.com/MathieuRocher/amazing_domain"
)

type CreateUserInput struct {
	Name     string `form:"name" validate:"required,min=2"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=8"`
	Role     string `form:"role" validate:"required,oneof=Trainee Trainer Admin"`
}

func (i *CreateUserInput) ToDomain() (*domain.User, error) {
	return &domain.User{
		Name:     i.Name,
		Email:    i.Email,
		Password: i.Password,
	}, nil
}
