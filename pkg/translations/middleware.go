package translations

import (
	"context"
	"strings"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/constants"
	"github.com/gin-gonic/gin"
)

var localeAlias = map[string]Locale{
	"id": LocaleID,
	"en": LocaleEN,
}

func normalizeLocale(lang string, defaultLocale Locale) Locale {
	if lang == "" {
		return defaultLocale
	}

	lang = strings.ReplaceAll(lang, "-", "_")
	parts := strings.Split(lang, "_")

	var normalized string
	if len(parts) == 2 {
		normalized = strings.ToLower(parts[0]) + "_" + strings.ToUpper(parts[1])
	} else {
		normalized = strings.ToLower(lang)
	}

	if alias, ok := localeAlias[normalized]; ok {
		return alias
	}

	if IsSupportedLocale(normalized) {
		return Locale(normalized)
	}

	return defaultLocale
}

func Middleware(defaultLocale Locale) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query(constants.XLocale)
		if lang == "" {
			lang = c.GetHeader(constants.XLocale)
		}
		locale := normalizeLocale(lang, defaultLocale)
		c.Set(constants.XLocale, locale)

		// juga bisa simpan ke context.Context, misal untuk pass ke usecase
		ctx := SetContext(c.Request.Context(), locale)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// SetContext adds locale info into context
func SetContext(ctx context.Context, locale Locale) context.Context {
	return context.WithValue(ctx, constants.XLocale, locale)
}

func GetLocaleFromGin(c *gin.Context, defaultLocale Locale) Locale {
	if l, exists := c.Get(constants.XLocale); exists {
		if lang, ok := l.(Locale); ok {
			return lang
		}
	}
	return defaultLocale
}

// GetLocaleFromContext retrieves locale from context
func GetLocaleFromContext(ctx context.Context, defaultLocale Locale) Locale {
	if l, ok := ctx.Value(constants.XLocale).(Locale); ok {
		return l
	}
	return defaultLocale
}
