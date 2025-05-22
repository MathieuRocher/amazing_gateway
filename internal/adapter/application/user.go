package application

import (
	"amazing_gateway/internal/adapter/handler/dto/user"
	"amazing_gateway/internal/auth"
	"errors"
	"fmt"
	domain "github.com/MathieuRocher/amazing_domain"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseInterface interface {
	FindAll() ([]domain.User, error)
	FindByID(id uint) (*domain.User, error)
	Create(user *domain.User, roleStr string) error
	Update(user *domain.User, input *user.UpdateUserInput) error
	Delete(id uint) error
	Authenticate(email string, password string) (string, error)
}

type UserRepositoryInterface interface {
	FindAll() ([]domain.User, error)
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id uint) error
}

type UserUsecase struct {
	userRepository UserRepositoryInterface
}

func NewUserUsecase(userRepository UserRepositoryInterface) UserUsecaseInterface {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (u *UserUsecase) FindAll() ([]domain.User, error) {
	return u.userRepository.FindAll()
}

func (u *UserUsecase) FindByID(id uint) (*domain.User, error) {
	return u.userRepository.FindByID(id)
}

func (u *UserUsecase) Create(user *domain.User, roleStr string) error {
	existing, err := u.userRepository.FindByEmail(user.Email)
	if err == nil && existing != nil {
		return fmt.Errorf("email already used")
	}

	role := domain.ParseRole(roleStr)
	if role == -1 {
		return fmt.Errorf("invalid role")
	}
	user.Role = role

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashed)
	return u.userRepository.Create(user)
}

func (uc *UserUsecase) Update(user *domain.User, input *user.UpdateUserInput) error {
	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	if input.Role != nil {
		role := domain.ParseRole(*input.Role)
		if role == -1 {
			return fmt.Errorf("invalid role: %s", *input.Role)
		}
		user.Role = role
	}

	return uc.userRepository.Update(user)
}

func (u *UserUsecase) Delete(id uint) error {
	return u.userRepository.Delete(id)

}

func (u *UserUsecase) Authenticate(email, password string) (string, error) {
	user, err := u.userRepository.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	return auth.GenerateJWT(*user)
}
