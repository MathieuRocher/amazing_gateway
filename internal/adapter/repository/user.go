package repository

import (
	"amazing_gateway/internal/infrastructure/database"
	domain "github.com/MathieuRocher/amazing_domain"
	"gorm.io/gorm"
)

type Role int

const (
	Trainee Role = iota
	Trainer
	Administrator
)

var RoleName = map[Role]string{
	Trainee:       "Trainee",
	Trainer:       "Trainer",
	Administrator: "Admin",
}

func (r Role) String() string {
	return RoleName[r]
}

type User struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Email        string `gorm:"unique"`
	Role         Role
	Password     string
	ClassGroupID *uint `gorm:"column:class_group_id"`
	ClassGroup   *ClassGroup
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{database.DB}
}

func (r *UserRepository) FindAll() ([]domain.User, error) {
	var repoUsers []User
	if err := r.db.Find(&repoUsers).Error; err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, repoUser := range repoUsers {
		domainReview := repoUser.ToDomain()
		domainUsers = append(domainUsers, *domainReview)
	}

	return domainUsers, nil
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(UserFromDomain(user)).Error
}

func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(UserFromDomain(user)).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

// ToDomain converts a repository User to a domain User
func (u *User) ToDomain() *domain.User {
	var classGroup *domain.ClassGroup
	if u.ClassGroup != nil {
		classGroup = u.ClassGroup.ToDomain()
	}

	return &domain.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Role:         domain.Role(u.Role),
		Password:     u.Password,
		ClassGroupID: u.ClassGroupID,
		ClassGroup:   classGroup,
	}
}

// UserFromDomain converts a domain User to a repository User
func UserFromDomain(u *domain.User) *User {
	var classGroup *ClassGroup
	if u.ClassGroup != nil {
		classGroup = ClassGroupFromDomain(u.ClassGroup)
	}

	return &User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Role:         Role(u.Role),
		Password:     u.Password,
		ClassGroupID: u.ClassGroupID,
		ClassGroup:   classGroup,
	}
}
