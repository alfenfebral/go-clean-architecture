package errorsutil

import "errors"

var ErrDefault error = errors.New("error")
var ErrNotFound error = errors.New("not found")
