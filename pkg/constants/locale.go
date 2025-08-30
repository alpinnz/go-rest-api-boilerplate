package constants

type Locale string

const (
	LocaleID Locale = "id_ID"
	LocaleEN Locale = "en_EN" // atau en_US kalau kamu lebih prefer American English
)

// Daftar semua locale yang didukung
var SupportedLocales = []Locale{
	LocaleID,
	LocaleEN,
}

// Helper untuk validasi
func IsSupportedLocale(l string) bool {
	for _, locale := range SupportedLocales {
		if string(locale) == l {
			return true
		}
	}
	return false
}
