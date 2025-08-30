package translations

type Locale string

const (
	LocaleID Locale = "id_ID"
	LocaleEN Locale = "en_EN" // bisa en_US kalau prefer
)

var SupportedLocales = []Locale{LocaleID, LocaleEN}

func IsSupportedLocale(l string) bool {
	for _, locale := range SupportedLocales {
		if string(locale) == l {
			return true
		}
	}
	return false
}
