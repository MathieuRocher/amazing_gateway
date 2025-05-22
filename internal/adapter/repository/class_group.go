package repository

import (
	"amazing_gateway/internal/infrastructure/database"
	domain "github.com/MathieuRocher/amazing_domain"
	"gorm.io/gorm"
)

type ClassGroup struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Trainees []User
}

type ClassGroupRepository struct {
	db *gorm.DB
}

func NewClassGroupRepository() *ClassGroupRepository {
	return &ClassGroupRepository{database.DB}
}

func (r *ClassGroupRepository) FindAll() ([]domain.ClassGroup, error) {
	var classGroups []ClassGroup
	if err := r.db.Find(&classGroups).Error; err != nil {
		return nil, err
	}

	var domainClassGroups []domain.ClassGroup
	for _, repoClassGroup := range classGroups {
		domainClassGroup := repoClassGroup.ToDomain()
		domainClassGroups = append(domainClassGroups, *domainClassGroup)
	}
	return domainClassGroups, nil
}

func (r *ClassGroupRepository) FindByID(id uint) (*domain.ClassGroup, error) {
	var classGroup ClassGroup
	if err := r.db.First(&classGroup, id).Error; err != nil {
		return nil, err
	}
	return classGroup.ToDomain(), nil
}

func (r *ClassGroupRepository) Create(classGroup *domain.ClassGroup) error {
	return r.db.Create(ClassGroupFromDomain(classGroup)).Error
}

func (r *ClassGroupRepository) Update(classGroup *domain.ClassGroup) error {
	return r.db.Save(ClassGroupFromDomain(classGroup)).Error
}

func (r *ClassGroupRepository) Delete(id uint) error {
	return r.db.Delete(&ClassGroup{}, id).Error
}

// ToDomain converts a repository ClassGroup to a domain ClassGroup
func (cg *ClassGroup) ToDomain() *domain.ClassGroup {
	return &domain.ClassGroup{
		ID:   cg.ID,
		Name: cg.Name,
	}
}

// ClassGroupFromDomain converts a domain ClassGroup to a repository ClassGroup
func ClassGroupFromDomain(cg *domain.ClassGroup) *ClassGroup {
	return &ClassGroup{
		ID:   cg.ID,
		Name: cg.Name,
	}
}
