package apis

import (
	"log"
	"net/http"

	"../models"
	"../services"
	"../utils"
	response "../utils/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gopkg.in/mgo.v2/bson"
)

// Todohandler represent the httphandler for file
type Todohandler struct {
	TodoService services.TodoService
}

// NewTodoHTTPHandler - make http handler
func NewTodoHTTPHandler(router *chi.Mux, service services.TodoService) {
	handler := &Todohandler{
		TodoService: service,
	}

	router.Get("/todo", handler.GetAll)
	router.Get("/todo/{id}", handler.GetByID)
	router.Post("/todo", handler.Create)
}

// Create - create todo http handler
func (handler *Todohandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		if err.Error() == "EOF" {
			utils.ResponseBodyError(w, r, err)
			return
		}

		utils.ResponseErrorValidation(w, r, err)
		return
	}
	timeNow := utils.GetTimeNow()

	result, err := handler.TodoService.Create(bson.M{
		"title":       data.Title,
		"description": data.Description,
		"createdAt":   timeNow,
		"updatedAt":   timeNow,
		"deletedAt":   timeNow,
	})
	if err != nil {
		utils.ResponseError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": "Create Todo",
		"data":    result,
	})
}

func (handler *Todohandler) GetAll(w http.ResponseWriter, r *http.Request) {
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
		log.Printf(err.Error())
		utils.ResponseErrorValidation(w, r, err)
		return
	}

	currentPage := utils.CurrentPage(pageQuery)
	perPage := utils.PerPage(perPageQuery)
	offset := utils.Offset(currentPage, perPage)

	results, totalData, err := handler.TodoService.GetAll(qQuery, perPage, offset)
	if err != nil {
		utils.ResponseError(w, r, err)
		return
	}
	totalPages := utils.TotalPage(totalData, perPage)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.H{
		"success": true,
		"code":    http.StatusOK,
		"message": "Get All Todo",
		"data":    results,
		"meta": response.H{
			"per_page":   perPage,
			"page":       currentPage,
			"pageCount":  totalPages,
			"totalCount": totalData,
		},
	})
}

// GetByID - get TODO by id http handler
func (handler *Todohandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Get and filter id param
	id := chi.URLParam(r, "id")
	log.Print(id)

	// Get detail
	result, err := handler.TodoService.GetByID(id)
	if err != nil {
		if err.Error() == "not found" {
			utils.ResponseNotFound(w, r, "Item not found")
			return
		}

		utils.ResponseError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.H{
		"success": true,
		"code":    http.StatusOK,
		"message": "Get Todo",
		"data":    result,
	})

}
