package translations

type Translator interface {
	T(locale Locale, key string, data map[string]any, fallback Locale) string
}
