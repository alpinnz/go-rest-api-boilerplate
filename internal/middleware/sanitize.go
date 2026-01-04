package middleware

import (
	"html"
	"strings"

	"github.com/gin-gonic/gin"
)

// Sanitize middleware sanitizes all string inputs to prevent XSS attacks
func Sanitize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		query := c.Request.URL.Query()
		for key, values := range query {
			for i, value := range values {
				query[key][i] = html.EscapeString(strings.TrimSpace(value))
			}
		}
		c.Request.URL.RawQuery = query.Encode()

		// Sanitize path parameters
		for _, param := range c.Params {
			c.Params = append(c.Params[:0], gin.Param{
				Key:   param.Key,
				Value: html.EscapeString(strings.TrimSpace(param.Value)),
			})
		}

		c.Next()
	}
}

// SanitizeString sanitizes a single string
func SanitizeString(s string) string {
	return html.EscapeString(strings.TrimSpace(s))
}

// SanitizeStrings sanitizes a slice of strings
func SanitizeStrings(strs []string) []string {
	result := make([]string, len(strs))
	for i, s := range strs {
		result[i] = SanitizeString(s)
	}
	return result
}
