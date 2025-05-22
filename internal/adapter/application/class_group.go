package application

import (
	domain "github.com/MathieuRocher/amazing_domain"
)

type ClassGroupUsecaseInterface interface {
	FindAll() ([]domain.ClassGroup, error)
	FindByID(id uint) (*domain.ClassGroup, error)
	Create(user *domain.ClassGroup) error
	Update(user *domain.ClassGroup) error
	Delete(id uint) error
}

type ClassGroupRepositoryInterface interface {
	FindAll() ([]domain.ClassGroup, error)
	FindByID(id uint) (*domain.ClassGroup, error)
	Create(classGroup *domain.ClassGroup) error
	Update(classGroup *domain.ClassGroup) error
	Delete(id uint) error
}

type ClassGroupUsecase struct {
	classGroupRepository ClassGroupRepositoryInterface
}

func NewClassGroupUseCase(classGroupRepository ClassGroupRepositoryInterface) ClassGroupUsecaseInterface {
	return &ClassGroupUsecase{
		classGroupRepository: classGroupRepository,
	}
}

func (u *ClassGroupUsecase) FindAll() ([]domain.ClassGroup, error) {
	return u.classGroupRepository.FindAll()
}

func (u *ClassGroupUsecase) FindByID(id uint) (*domain.ClassGroup, error) {
	return u.classGroupRepository.FindByID(id)
}

func (u *ClassGroupUsecase) Create(user *domain.ClassGroup) error {
	return u.classGroupRepository.Create(user)
}

func (u *ClassGroupUsecase) Update(user *domain.ClassGroup) error {
	return u.classGroupRepository.Update(user)

}

func (u *ClassGroupUsecase) Delete(id uint) error {
	return u.classGroupRepository.Delete(id)

}
