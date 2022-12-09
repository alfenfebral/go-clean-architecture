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
	repository todorepository.Repository
}

// New will create new an ServiceImpl object representation of Service interface
func New(repository todorepository.Repository) Service {
	return &ServiceImpl{
		repository: repository,
	}
}

// GetAll - get all todo service
func (s *ServiceImpl) GetAll(keyword string, limit int, offset int) ([]*models.Todo, int, error) {
	res, err := s.repository.FindAll(keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Count total
	total, err := s.repository.CountFindAll(keyword)
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

// GetByID - get todo by id service
func (s *ServiceImpl) GetByID(id string) (*models.Todo, error) {
	res, err := s.repository.FindById(id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create - creating todo service
func (r *ServiceImpl) Create(value *models.Todo) (*models.Todo, error) {
	res, err := r.repository.Store(&models.Todo{
		Title:       value.Title,
		Description: value.Description,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Update - update todo service
func (r *ServiceImpl) Update(id string, value *models.Todo) (*models.Todo, error) {
	_, err := r.repository.CountFindByID(id)
	if err != nil {
		return nil, err
	}

	_, err = r.repository.Update(id, &models.Todo{
		Title:       value.Title,
		Description: value.Description,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete - delete todo service
func (r *ServiceImpl) Delete(id string) error {
	err := r.repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
