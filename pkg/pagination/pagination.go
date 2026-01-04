package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Params represents pagination parameters
type Params struct {
	Page    int
	PerPage int
	Offset  int
}

// Meta represents pagination metadata
type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalData  int `json:"total_data"`
	TotalPages int `json:"total_pages"`
}

// DefaultPage is the default page number
const DefaultPage = 1

// DefaultPerPage is the default items per page
const DefaultPerPage = 10

// MaxPerPage is the maximum items per page
const MaxPerPage = 100

// FromContext extracts pagination parameters from Gin context
func FromContext(c *gin.Context) Params {
	page := getIntParam(c, "page", DefaultPage)
	perPage := getIntParam(c, "per_page", DefaultPerPage)

	// Validate page
	if page < 1 {
		page = DefaultPage
	}

	// Validate per_page
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	offset := (page - 1) * perPage

	return Params{
		Page:    page,
		PerPage: perPage,
		Offset:  offset,
	}
}

// NewMeta creates pagination metadata
func NewMeta(page, perPage, totalData int) Meta {
	totalPages := int(math.Ceil(float64(totalData) / float64(perPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	return Meta{
		Page:       page,
		PerPage:    perPage,
		TotalData:  totalData,
		TotalPages: totalPages,
	}
}

// getIntParam safely extracts integer parameter from query string
func getIntParam(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.DefaultQuery(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// Response wraps data with pagination metadata
type Response struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

// NewResponse creates a new paginated response
func NewResponse(data interface{}, meta Meta) Response {
	return Response{
		Data: data,
		Meta: meta,
	}
}
