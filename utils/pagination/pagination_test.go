package paginationutil_test

import (
	"testing"

	paginationutil "go-clean-architecture/utils/pagination"

	"github.com/stretchr/testify/assert"
)

func TestPerPage(t *testing.T) {
	value := paginationutil.PerPage(0)
	assert.Equal(t, value, 10)

	value = paginationutil.PerPage(10)
	assert.Equal(t, value, 10)
}

func TestCurrentPage(t *testing.T) {
	value := paginationutil.CurrentPage(0)
	assert.Equal(t, value, 1)

	value = paginationutil.CurrentPage(10)
	assert.Equal(t, value, 10)
}

func TestTotalPage(t *testing.T) {
	value := paginationutil.TotalPage(20, 10)
	assert.Equal(t, value, 2)
}

func TestOffset(t *testing.T) {
	value := paginationutil.Offset(1, 10)
	assert.Equal(t, value, 0)

	value = paginationutil.Offset(-1, 10)
	assert.Equal(t, value, 0)
}
