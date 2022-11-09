package models

import (
	"net/http"
	"time"

	"go-clean-architecture/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo - todo model
type Todo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updatedAt"`
}

// TodoRequest - todo request
type TodoRequest struct {
	Title       string `form:"title" json:"title" validate:"required"`
	Description string `form:"description" json:"description" validate:"required"`
}

func (tr *TodoRequest) Bind(r *http.Request) error {
	return utils.ValidateStruct(tr)
}

// TodoListRequest - form for list validation
type TodoListRequest struct {
	Keywords *SearchForm
	Page     string `form:"page" json:"page" validate:"sgte=1"`
	PerPage  string `form:"per_page" json:"per_page" validate:"sgte=1,slte=100"`
}

// SearchForm - search list struct
type SearchForm struct {
	Keywords string `form:"q" json:"q" validate:"max=255"`
}
