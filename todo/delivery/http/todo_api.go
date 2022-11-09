package handlers

import (
	"net/http"

	"go-clean-architecture/todo/models"
	"go-clean-architecture/todo/services"
	"go-clean-architecture/utils"
	response "go-clean-architecture/utils/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// todoHandler represent the http handler
type todoHandler struct {
	router      *chi.Mux
	todoService services.TodoService
}

// NewTodoHTTPHandler - make http handler
func NewTodoHTTPHandler(router *chi.Mux, service services.TodoService) *todoHandler {
	return &todoHandler{
		router:      router,
		todoService: service,
	}
}

func (handler *todoHandler) RegisterRoutes() {
	handler.router.Get("/todo", handler.GetAll)
	handler.router.Get("/todo/{id}", handler.GetByID)
	handler.router.Post("/todo", handler.Create)
	handler.router.Put("/todo/{id}", handler.Update)
	handler.router.Delete("/todo/{id}", handler.Delete)
}

// GetAll - get all todo http handler
func (handler *todoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	qQuery := r.URL.Query().Get("q")
	pageQuery := r.URL.Query().Get("page")
	perPageQuery := r.URL.Query().Get("per_page")

	err := utils.ValidateStruct(&models.TodoListRequest{
		Keywords: &models.SearchForm{
			Keywords: qQuery,
		},
		Page:    pageQuery,
		PerPage: perPageQuery,
	})
	if err != nil {
		response.ResponseErrorValidation(w, r, err)
		return
	}

	currentPage := utils.CurrentPage(pageQuery)
	perPage := utils.PerPage(perPageQuery)
	offset := utils.Offset(currentPage, perPage)

	results, totalData, err := handler.todoService.GetAll(qQuery, perPage, offset)
	if err != nil {
		response.ResponseError(w, r, err)
		return
	}
	totalPages := utils.TotalPage(totalData, perPage)

	response.ResponseOKList(w, r, &response.ResponseSuccessList{
		Data: results,
		Meta: &response.Meta{
			PerPage:     perPage,
			CurrentPage: currentPage,
			TotalPage:   totalPages,
			TotalData:   totalData,
		},
	})
}

// GetByID - get todo by id http handler
func (handler *todoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Get detail
	result, err := handler.todoService.GetByID(id)
	if err != nil {
		if err.Error() == "not found" {
			response.ResponseNotFound(w, r, "Item not found")
			return
		}

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseOK(w, r, &response.ResponseSuccess{
		Data: result,
	})

}

// Create - create todo http handler
func (handler *todoHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		if err.Error() == "EOF" {
			response.ResponseBodyError(w, r, err)
			return
		}

		response.ResponseErrorValidation(w, r, err)
		return
	}

	result, err := handler.todoService.Create(&models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})
	if err != nil {
		response.ResponseError(w, r, err)
		return
	}

	response.ResponseCreated(w, r, &response.ResponseSuccess{
		Data: result,
	})
}

// Update - update todo by id http handler
func (handler *todoHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Get and filter id param
	id := chi.URLParam(r, "id")

	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		if err.Error() == "EOF" {
			response.ResponseBodyError(w, r, err)
			return
		}

		response.ResponseErrorValidation(w, r, err)
		return
	}

	// Edit data
	_, err := handler.todoService.Update(id, &models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})

	if err != nil {
		if err.Error() == "not found" {
			response.ResponseNotFound(w, r, "Item not found")
			return
		}

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseOK(w, r, &response.ResponseSuccess{
		Data: response.H{
			"id": id,
		},
	})
}

// Delete - delete todo by id http handler
func (handler *todoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Delete record
	err := handler.todoService.Delete(id)
	if err != nil {
		if err.Error() == "not found" {
			response.ResponseNotFound(w, r, "Item not found")
			return
		}

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseOK(w, r, &response.ResponseSuccess{
		Data: response.H{
			"id": id,
		},
	})
}
