package validator

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type Validator struct{}

var validate *validator.Validate

// CommonError - error response format
type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

func New() *Validator {
	validate = validator.New()

	return &Validator{}
}

func ValidatonError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)

	for _, v := range errs {
		field := strcase.ToSnake(v.Field())

		switch v.Tag() {
		case "required":
			res.Errors[field] = fmt.Sprintf("%v is %v", field, v.Tag())
		case "sinteger":
			res.Errors[field] = fmt.Sprintf("%v is number only", field)
		case "sgte":
			res.Errors[field] = fmt.Sprintf("%v must higher than equal %v", field, v.Param())
		case "slte":
			res.Errors[field] = fmt.Sprintf("%v must less than equal %v", field, v.Param())
		case "max":
			res.Errors[field] = fmt.Sprintf("%v must less than %v character", field, v.Param())
		case "min":
			res.Errors[field] = fmt.Sprintf("%v must higher than %v character", field, v.Param())
		case "email":
			res.Errors[field] = fmt.Sprintf("%v is not a valid email address", v.Value())
		case "username":
			res.Errors[field] = fmt.Sprintf("%v is not a valid username", v.Value())
		}
	}

	return res
}

func ValidateStruct(i interface{}) error {
	validate = validator.New()
	validate.RegisterValidation("sinteger", Integer)
	validate.RegisterValidation("sgte", GreaterThanEqual)
	validate.RegisterValidation("slte", LessThanEqual)
	validate.RegisterValidation("username", Username)

	err := validate.Struct(i)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		return err
	}

	return err
}

// Integer - integer only validation
func Integer(fl validator.FieldLevel) bool {
	// If empty skip
	if fl.Field().String() == "" {
		return true
	}

	_, err := strconv.Atoi(fl.Field().String())
	if err != nil {
		return false
	}

	return true
}

// GreaterThanEqual - greater than equal
func GreaterThanEqual(fl validator.FieldLevel) bool {
	// If empty skip
	if fl.Field().String() == "" {
		return true
	}

	i, err := strconv.Atoi(fl.Field().String())
	if err != nil {
		return false
	}

	param, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}

	if i < param {
		return false
	}

	return true
}

// LessThanEqual - less than equal
func LessThanEqual(fl validator.FieldLevel) bool {
	// If empty skip
	if fl.Field().String() == "" {
		return true
	}

	i, err := strconv.Atoi(fl.Field().String())
	if err != nil {
		return false
	}

	param, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}

	if i > param {
		return false
	}

	return true
}

// Username - username regex only alphanumeric
func Username(fl validator.FieldLevel) bool {
	// If empty skip
	if fl.Field().String() == "" {
		return true
	}

	var regex = regexp.MustCompile(`^[A-Za-z0-9]+(?:[_-][A-Za-z0-9]+)*$`)
	return regex.MatchString(fl.Field().String())
}
