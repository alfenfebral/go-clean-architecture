package service

import (
	"go-clean-architecture/todo/models"
	todorepository "go-clean-architecture/todo/repository"
)

// Service represent the todo service
type Service interface {
	GetAll(keyword string, limit int, offset int) ([]*models.Todo, int, error)
	GetByID(id string) (*models.Todo, error)
	Create(value *models.Todo) (*models.Todo, error)
	Update(id string, value *models.Todo) (*models.Todo, error)
	Delete(id string) error
}

type ServiceImpl struct {
	todoRepo todorepository.Repository
}

// New will create new an ServiceImpl object representation of Service interface
func New(a todorepository.Repository) Service {
	return &ServiceImpl{
		todoRepo: a,
	}
}

// GetAll - get all todo service
func (a *ServiceImpl) GetAll(keyword string, limit int, offset int) ([]*models.Todo, int, error) {
	res, err := a.todoRepo.FindAll(keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Count total
	total, err := a.todoRepo.CountFindAll(keyword)
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

// GetByID - get todo by id service
func (a *ServiceImpl) GetByID(id string) (*models.Todo, error) {
	res, err := a.todoRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create - creating todo service
func (a *ServiceImpl) Create(value *models.Todo) (*models.Todo, error) {
	res, err := a.todoRepo.Store(&models.Todo{
		Title:       value.Title,
		Description: value.Description,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Update - update todo service
func (a *ServiceImpl) Update(id string, value *models.Todo) (*models.Todo, error) {
	_, err := a.todoRepo.CountFindByID(id)
	if err != nil {
		return nil, err
	}

	_, err = a.todoRepo.Update(id, &models.Todo{
		Title:       value.Title,
		Description: value.Description,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete - delete todo service
func (a *ServiceImpl) Delete(id string) error {
	err := a.todoRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
