package pagination

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		page           string
		perPage        string
		expectedPage   int
		expectedPerPg  int
		expectedOffset int
	}{
		{
			name:           "default values",
			page:           "",
			perPage:        "",
			expectedPage:   DefaultPage,
			expectedPerPg:  DefaultPerPage,
			expectedOffset: 0,
		},
		{
			name:           "valid page and per_page",
			page:           "2",
			perPage:        "20",
			expectedPage:   2,
			expectedPerPg:  20,
			expectedOffset: 20,
		},
		{
			name:           "invalid page defaults to 1",
			page:           "-1",
			perPage:        "10",
			expectedPage:   DefaultPage,
			expectedPerPg:  10,
			expectedOffset: 0,
		},
		{
			name:           "per_page exceeds max",
			page:           "1",
			perPage:        "200",
			expectedPage:   1,
			expectedPerPg:  MaxPerPage,
			expectedOffset: 0,
		},
		{
			name:           "invalid per_page defaults to 10",
			page:           "1",
			perPage:        "invalid",
			expectedPage:   1,
			expectedPerPg:  DefaultPerPage,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			c.Request = &gin.Request{URL: &gin.URL{}}

			query := c.Request.URL.Query()
			if tt.page != "" {
				query.Set("page", tt.page)
			}
			if tt.perPage != "" {
				query.Set("per_page", tt.perPage)
			}
			c.Request.URL.RawQuery = query.Encode()

			params := FromContext(c)

			assert.Equal(t, tt.expectedPage, params.Page)
			assert.Equal(t, tt.expectedPerPg, params.PerPage)
			assert.Equal(t, tt.expectedOffset, params.Offset)
		})
	}
}

func TestNewMeta(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		perPage       int
		totalData     int
		expectedPages int
	}{
		{
			name:          "basic pagination",
			page:          1,
			perPage:       10,
			totalData:     50,
			expectedPages: 5,
		},
		{
			name:          "last page not full",
			page:          3,
			perPage:       10,
			totalData:     25,
			expectedPages: 3,
		},
		{
			name:          "no data",
			page:          1,
			perPage:       10,
			totalData:     0,
			expectedPages: 1,
		},
		{
			name:          "single item",
			page:          1,
			perPage:       10,
			totalData:     1,
			expectedPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := NewMeta(tt.page, tt.perPage, tt.totalData)

			assert.Equal(t, tt.page, meta.Page)
			assert.Equal(t, tt.perPage, meta.PerPage)
			assert.Equal(t, tt.totalData, meta.TotalData)
			assert.Equal(t, tt.expectedPages, meta.TotalPages)
		})
	}
}
