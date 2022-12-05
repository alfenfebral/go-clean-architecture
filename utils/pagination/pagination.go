package paginationutil

import (
	"math"
)

// PerPage - get per_page, the default value is 10
func PerPage(value int) int {
	if value <= 0 {
		return 10
	}

	return value
}

// CurrentPage - get current pages, the default value is 1
func CurrentPage(value int) int {
	if value < 1 {
		return 1
	}

	return value
}

// TotalPage - get total pages, based on ceil total/per page
func TotalPage(total int, perPage int) int {
	totalFloat := float64(total)
	perPageFloat := float64(perPage)
	resultFloat := math.Ceil(totalFloat / perPageFloat)
	result := int(resultFloat)

	return result
}

// Offset - offset of pages
func Offset(currentPage int, perPage int) int {
	result := (currentPage - 1) * perPage
	if result < 0 {
		return 0
	}

	return result
}
