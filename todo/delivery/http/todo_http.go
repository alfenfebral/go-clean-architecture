package httpdelivery

import (
	"net/http"
	"strconv"

	pkgvalidator "go-clean-architecture/pkg/validator"
	"go-clean-architecture/todo/models"
	todoservice "go-clean-architecture/todo/service"
	paginationutil "go-clean-architecture/utils/pagination"
	responseutil "go-clean-architecture/utils/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type HTTPHandler interface {
	RegisterRoutes(router *chi.Mux)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type HTTPHandlerImpl struct {
	todoService todoservice.Service
}

// New - make http handler
func New(service todoservice.Service) HTTPHandler {
	return &HTTPHandlerImpl{
		todoService: service,
	}
}

func (h *HTTPHandlerImpl) RegisterRoutes(router *chi.Mux) {
	router.Get("/todo", h.GetAll)
	router.Get("/todo/{id}", h.GetByID)
	router.Post("/todo", h.Create)
	router.Put("/todo/{id}", h.Update)
	router.Delete("/todo/{id}", h.Delete)
}

// GetAll - get all todo http handler
func (h *HTTPHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	qQuery := r.URL.Query().Get("q")
	pageQueryStr := r.URL.Query().Get("page")
	perPageQueryStr := r.URL.Query().Get("per_page")

	err := pkgvalidator.ValidateStruct(&models.TodoListRequest{
		Keywords: &models.SearchForm{
			Keywords: qQuery,
		},
		Page:    pageQueryStr,
		PerPage: perPageQueryStr,
	})
	if err != nil {
		responseutil.ResponseErrorValidation(w, r, err)
		return
	}

	pageQuery, _ := strconv.Atoi(pageQueryStr)
	perPageQuery, _ := strconv.Atoi(perPageQueryStr)

	currentPage := paginationutil.CurrentPage(pageQuery)
	perPage := paginationutil.PerPage(perPageQuery)
	offset := paginationutil.Offset(currentPage, perPage)

	results, totalData, err := h.todoService.GetAll(qQuery, perPage, offset)
	if err != nil {
		responseutil.ResponseError(w, r, err)
		return
	}
	totalPages := paginationutil.TotalPage(totalData, perPage)

	responseutil.ResponseOKList(w, r, &responseutil.ResponseSuccessList{
		Data: results,
		Meta: &responseutil.Meta{
			PerPage:     perPage,
			CurrentPage: currentPage,
			TotalPage:   totalPages,
			TotalData:   totalData,
		},
	})
}

// GetByID - get todo by id http handler
func (h *HTTPHandlerImpl) GetByID(w http.ResponseWriter, r *http.Request) {
	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Get detail
	result, err := h.todoService.GetByID(id)
	if err != nil {
		if err.Error() == "not found" {
			responseutil.ResponseNotFound(w, r, "Item not found")
			return
		}

		responseutil.ResponseError(w, r, err)
		return
	}

	responseutil.ResponseOK(w, r, &responseutil.ResponseSuccess{
		Data: result,
	})

}

// Create - create todo http handler
func (h *HTTPHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		if err.Error() == "EOF" {
			responseutil.ResponseBodyError(w, r, err)
			return
		}

		responseutil.ResponseErrorValidation(w, r, err)
		return
	}

	result, err := h.todoService.Create(&models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})
	if err != nil {
		responseutil.ResponseError(w, r, err)
		return
	}

	responseutil.ResponseCreated(w, r, &responseutil.ResponseSuccess{
		Data: result,
	})
}

// Update - update todo by id http handler
func (h *HTTPHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	// Get and filter id param
	id := chi.URLParam(r, "id")

	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		if err.Error() == "EOF" {
			responseutil.ResponseBodyError(w, r, err)
			return
		}

		responseutil.ResponseErrorValidation(w, r, err)
		return
	}

	// Edit data
	_, err := h.todoService.Update(id, &models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})

	if err != nil {
		if err.Error() == "not found" {
			responseutil.ResponseNotFound(w, r, "Item not found")
			return
		}

		responseutil.ResponseError(w, r, err)
		return
	}

	responseutil.ResponseOK(w, r, &responseutil.ResponseSuccess{
		Data: responseutil.H{
			"id": id,
		},
	})
}

// Delete - delete todo by id http handler
func (h *HTTPHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Delete record
	err := h.todoService.Delete(id)
	if err != nil {
		if err.Error() == "not found" {
			responseutil.ResponseNotFound(w, r, "Item not found")
			return
		}

		responseutil.ResponseError(w, r, err)
		return
	}

	responseutil.ResponseOK(w, r, &responseutil.ResponseSuccess{
		Data: responseutil.H{
			"id": id,
		},
	})
}
