package middleware

import (
	"context"
	"strings"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/localization"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	// LangKey is context key for language dictionary
	LangKey ctxKey = "lang_dict"
	// LangCodeKey is context key for language code
	LangCodeKey ctxKey = "lang_code"
)

// pickLang extracts language code from Accept-Language header.
// Supports formats: "id", "id-ID", "en-US,en;q=0.9"
// Returns first language tag in lowercase, stripped of region code.
func pickLang(header string) string {
	if header == "" {
		return ""
	}

	// Split by comma (multiple languages)
	parts := strings.Split(header, ",")
	if len(parts) == 0 {
		return ""
	}

	// Take first language
	tag := strings.TrimSpace(parts[0])

	// Remove quality value (;q=0.9)
	tag = strings.Split(tag, ";")[0]

	// Lowercase
	tag = strings.ToLower(tag)

	// Strip region code (en-US -> en)
	if strings.Contains(tag, "-") {
		tag = strings.Split(tag, "-")[0]
	}

	return tag
}

// Language middleware extracts Accept-Language header and injects dictionary into context.
// Falls back to bundle's default language if header is missing or invalid.
func Language(bundle *localization.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := pickLang(c.GetHeader("Accept-Language"))
		if code == "" {
			code = bundle.DefaultLang()
		}

		dict := bundle.Get(code)

		ctx := context.WithValue(c.Request.Context(), LangKey, dict)
		ctx = context.WithValue(ctx, LangCodeKey, code)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetDict retrieves dictionary from Gin context.
func GetDict(c *gin.Context) localization.Dictionary {
	dict, ok := c.Request.Context().Value(LangKey).(localization.Dictionary)
	if !ok {
		return localization.Dictionary{}
	}
	return dict
}

// GetLangCode retrieves language code from Gin context.
func GetLangCode(c *gin.Context) string {
	code, ok := c.Request.Context().Value(LangCodeKey).(string)
	if !ok {
		return "en"
	}
	return code
}
