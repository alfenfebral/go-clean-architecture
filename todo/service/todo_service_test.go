package service_test

import (
	mockrepository "go-clean-architecture/todo/mocks/repository"
	"go-clean-architecture/todo/models"
	todoservice "go-clean-architecture/todo/service"
	errorsutil "go-clean-architecture/utils/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var DefaultID string = "1"

func TestTodoGetAll(t *testing.T) {
	t.Run("success when find all", func(t *testing.T) {
		mockList := make([]*models.Todo, 0)
		mockList = append(mockList, &models.Todo{})

		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("FindAll", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(mockList, nil)
		mockRepository.On("CountFindAll", mock.AnythingOfType("string")).Return(10, nil)

		results, count, err := service.GetAll("keyword", 10, 0)

		assert.NoError(t, err)
		assert.Equal(t, count, 10)
		assert.Equal(t, mockList, results)
	})

	t.Run("error when find all", func(t *testing.T) {
		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("FindAll", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, errorsutil.ErrDefault)
		mockRepository.On("CountFindAll", mock.AnythingOfType("string")).Return(10, nil)
		results, count, err := service.GetAll("keyword", 10, 0)

		assert.Nil(t, results)
		assert.Equal(t, 0, count)
		assert.Error(t, err)
	})

	t.Run("error when count find all", func(t *testing.T) {
		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("FindAll", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, nil)
		mockRepository.On("CountFindAll", mock.AnythingOfType("string")).Return(10, errorsutil.ErrDefault)

		results, count, err := service.GetAll("keyword", 10, 0)

		assert.Nil(t, results)
		assert.Equal(t, 0, count)
		assert.Error(t, err)
	})
}

func TestTodoGetByID(t *testing.T) {
	t.Run("success when find by id", func(t *testing.T) {
		var mockTodo = &models.Todo{}

		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("FindById", mock.AnythingOfType("string")).Return(mockTodo, nil)

		result, err := service.GetByID(DefaultID)

		assert.NoError(t, err)
		assert.Equal(t, mockTodo, result)
	})

	t.Run("error when find by id", func(t *testing.T) {
		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("FindById", mock.AnythingOfType("string")).Return(nil, errorsutil.ErrDefault)
		result, err := service.GetByID(DefaultID)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestTodoCreate(t *testing.T) {
	t.Run("success when create", func(t *testing.T) {
		var mockTodo = &models.Todo{}

		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("Store", mock.AnythingOfType("*models.Todo")).Return(mockTodo, nil)

		result, err := service.Create(&models.Todo{})

		assert.NoError(t, err)
		assert.Equal(t, mockTodo, result)
	})

	t.Run("error when create", func(t *testing.T) {
		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("Store", mock.AnythingOfType("*models.Todo")).Return(nil, errorsutil.ErrDefault)
		result, err := service.Create(&models.Todo{})

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestTodoUpdate(t *testing.T) {
	t.Run("success when update", func(t *testing.T) {
		var mockTodo = &models.Todo{}

		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("CountFindByID", mock.AnythingOfType("string")).Return(10, nil)
		mockRepository.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(mockTodo, nil)

		result, err := service.Update(DefaultID, &models.Todo{})

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("error when count find by id", func(t *testing.T) {
		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("CountFindByID", mock.AnythingOfType("string")).Return(0, errorsutil.ErrDefault)
		mockRepository.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(nil, nil)

		result, err := service.Update(DefaultID, &models.Todo{})

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("error when update", func(t *testing.T) {
		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("CountFindByID", mock.AnythingOfType("string")).Return(10, nil)
		mockRepository.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(nil, errorsutil.ErrDefault)

		result, err := service.Update(DefaultID, &models.Todo{})

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestTodoDelete(t *testing.T) {
	t.Run("success when delete", func(t *testing.T) {
		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("Delete", mock.AnythingOfType("string")).Return(nil)

		err := service.Delete(DefaultID)

		assert.NoError(t, err)
	})

	t.Run("error when delete", func(t *testing.T) {
		mockRepository := new(mockrepository.Repository)
		service := todoservice.New(mockRepository)

		mockRepository.On("Delete", mock.AnythingOfType("string")).Return(errorsutil.ErrDefault)

		err := service.Delete(DefaultID)

		assert.Error(t, err)
	})
}
