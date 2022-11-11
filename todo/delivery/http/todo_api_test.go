package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	todo_http "go-clean-architecture/todo/delivery/http"
	"go-clean-architecture/utils"

	mockService "go-clean-architecture/todo/mocks/service"

	"go-clean-architecture/todo/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ErrDefault error = errors.New("error")
var ErrNotFound error = errors.New("not found")
var WhenError400EOF string = "when return 400 bad request (error EOF)"
var WhenError500Service string = "when return 500 internal error (error service)"
var WhenError500Query string = "when return 500 internal error (error query)"
var WhenError400Validation string = "when return 400 bad request (error validation)"
var WhenError404NotFound string = "when return 404 not found (resouce not found)"
var WhenSuccess201Created string = "when return 201 created"
var WhenSuccess200OK string = "when return 200 ok"

func TestNewTodoHTTPHandler(t *testing.T) {
	utils.InitializeValidator()

	router := chi.NewRouter()
	mockService := new(mockService.TodoService)

	mockService.On("Create", mock.AnythingOfType("*models.Todo")).Return(&models.Todo{}, nil)

	todo_http.NewTodoHTTPHandler(router, mockService)
}

// TestTodoGetAll - testing GetAll [200]
func TestTodoGetAll(t *testing.T) {
	t.Run(WhenError400Validation, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?page=-1&per_page=-1", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetAll)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?page=1&per_page=10", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("GetAll", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, 1, ErrDefault)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetAll)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		utils.InitializeValidator()

		mockListTodo := make([]*models.Todo, 0)
		mockListTodo = append(mockListTodo, &models.Todo{})

		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?page=1&per_page=10", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("GetAll", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(mockListTodo, 1, nil)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetAll)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
}

// TestTodoCreate - testing create [201]
func TestTodoCreate(t *testing.T) {
	t.Run(WhenError400EOF, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader([]byte("")))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("when return 400 bad request (error validation) ", func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		mockPostBody := map[string]interface{}{
			"title":       "",
			"description": "",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("when error 500 internal error (error service)", func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		mockPostBody := map[string]interface{}{
			"title":       "lorem ipsum",
			"description": "desc",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("Create", mock.AnythingOfType("*models.Todo")).Return(nil, errors.New(""))

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run("when return 201 created", func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		mockPostBody := map[string]interface{}{
			"title":       "lorem ipsum",
			"description": "desc",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("Create", mock.AnythingOfType("*models.Todo")).Return(&models.Todo{}, nil)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusCreated, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
}

// TestTodoGetByID - testing GetByID [200]
func TestTodoGetByID(t *testing.T) {
	t.Run(WhenError404NotFound, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("GetByID", mock.AnythingOfType("string")).Return(nil, ErrNotFound)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetByID)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusNotFound, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("GetByID", mock.AnythingOfType("string")).Return(nil, ErrDefault)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetByID)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("GetByID", mock.AnythingOfType("string")).Return(&models.Todo{}, nil)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetByID)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
}

// TestTodoUpdate - testing update [200]
func TestTodoUpdate(t *testing.T) {
	t.Run(WhenError400EOF, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/product?id=1", bytes.NewReader([]byte("")))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(&models.Todo{}, nil)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run(WhenError400Validation, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		mockPostBody := map[string]interface{}{
			"title":       "",
			"description": "",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenError404NotFound, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		mockPostBody := map[string]interface{}{
			"title":       "a",
			"description": "a",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(nil, ErrNotFound)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusNotFound, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		mockPostBody := map[string]interface{}{
			"title":       "a",
			"description": "a",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(nil, ErrDefault)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		mockPostBody := map[string]interface{}{
			"title":       "a",
			"description": "a",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(&models.Todo{}, nil)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
}

// TestDeleteSuccess - testing delete [200]
func TestTodoDelete(t *testing.T) {
	t.Run(WhenError404NotFound, func(t *testing.T) {
		utils.InitializeValidator()
		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/product?id=", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService.On("Delete", mock.AnythingOfType("string")).Return(ErrNotFound)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Delete)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		utils.InitializeValidator()
		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService.On("Delete", mock.AnythingOfType("string")).Return(ErrDefault)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Delete)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		utils.InitializeValidator()

		mockService := new(mockService.TodoService)

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService.On("Delete", mock.AnythingOfType("string")).Return(nil)

		todoHandler := todo_http.NewTodoHTTPHandler(chi.NewRouter(), mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Delete)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
}
