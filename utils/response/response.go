package response

import (
	"go-clean-architecture/pkg/logger"
	pkgvalidator "go-clean-architecture/pkg/validator"
	"net/http"

	"github.com/go-chi/render"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

type ResponseSuccessList struct {
	Data interface{} `json:"data"`
	Meta *Meta       `json:"meta"`
}

type Meta struct {
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"page"`
	TotalPage   int `json:"page_count"`
	TotalData   int `json:"total_count"`
}

type ResponseSuccess struct {
	Data interface{} `json:"data"`
}

func ResponseErrorValidation(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, H{
		"success": false,
		"code":    http.StatusBadRequest,
		"message": "Validation errors in your request",
		"errors":  pkgvalidator.ValidatonError(err).Errors,
	})
}

func ResponseBodyError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, H{
		"success": false,
		"code":    http.StatusBadRequest,
		"message": "Validation errors in your request",
		"error":   "Check your body request",
	})
}

// ResponseError - send response error (500)
func ResponseError(w http.ResponseWriter, r *http.Request, err error) {
	logger.Error(err)

	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, H{
		"success": false,
		"code":    http.StatusInternalServerError,
		"message": "There is something error",
	})
}

// ResponseNotFound - send response not found (404)
func ResponseNotFound(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, H{
		"success": false,
		"code":    http.StatusNotFound,
		"message": message,
	})
}

func ResponseCreated(w http.ResponseWriter, r *http.Request, data *ResponseSuccess) {
	render.Status(r, http.StatusCreated)

	render.JSON(w, r, H{
		"success": true,
		"code":    http.StatusCreated,
		"data":    data.Data,
	})
}

func ResponseOK(w http.ResponseWriter, r *http.Request, data *ResponseSuccess) {
	render.Status(r, http.StatusOK)

	render.JSON(w, r, H{
		"success": true,
		"code":    http.StatusOK,
		"data":    data.Data,
	})
}

func ResponseOKList(w http.ResponseWriter, r *http.Request, data *ResponseSuccessList) {
	render.Status(r, http.StatusOK)

	render.JSON(w, r, H{
		"success": true,
		"code":    http.StatusOK,
		"data":    data.Data,
		"meta":    data.Meta,
	})
}

func ResponseInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusOK)

	logger.Error(err)
	render.JSON(w, r, H{
		"success": false,
		"code":    http.StatusInternalServerError,
		"message": "Internal server error",
	})
}
