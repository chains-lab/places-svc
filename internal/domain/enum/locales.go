package enum

import "fmt"

const (
	LocaleEN = "en"
	LocaleRU = "ru"
	LocaleUK = "uk"
)

const DefaultLocale = LocaleEN

var supportedLocales = []string{
	LocaleEN,
	LocaleRU,
	LocaleUK,
}

var ErrorUnsupportedLocale = fmt.Errorf("unsupported locale, supported locales are: %v", supportedLocales)

// CheckLocale checks if the given locale is supported
func CheckLocale(locale string) error {
	for _, l := range supportedLocales {
		if l == locale {
			return nil
		}
	}
	return fmt.Errorf("'%s': %w", locale, ErrorUnsupportedLocale)
}

func GetAllLocales() []string {
	return supportedLocales
}
