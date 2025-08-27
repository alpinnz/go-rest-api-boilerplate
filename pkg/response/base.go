package response

// BaseResponse defines a unified response structure
type BaseResponse struct {
	TraceID string      `json:"trace_id,omitempty"` // Untuk debugging / tracing
	Code    string      `json:"code,omitempty"`     // e.g. USER_NOT_FOUND, INTERNAL_ERROR
	Message string      `json:"message"`            // Human-readable message
	Data    interface{} `json:"data,omitempty"`     // Untuk response sukses
	Errors  interface{} `json:"errors,omitempty"`   // Untuk response gagal
}

type FieldError struct {
	Field   string `json:"field"`   // e.g. "password"
	Code    string `json:"code"`    // e.g. "required", "email"
	Type    string `json:"type"`    // e.g. "validation"
	Message string `json:"message"` // Human-readable message
}

type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type ListWithPagination[T any] struct {
	Items      []T        `json:"-"`
	Key        string     `json:"-"`
	Pagination Pagination `json:"pagination"`
}

func (r ListWithPagination[T]) ToMap() map[string]interface{} {
	return map[string]interface{}{
		r.Key:        r.Items,
		"pagination": r.Pagination,
	}
}
