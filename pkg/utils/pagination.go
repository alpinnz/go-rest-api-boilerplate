package utils

import (
	"math"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
)

func NewPagination(page, limit, total int) response.Pagination {
	if limit <= 0 {
		limit = 10 // default limit
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return response.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
