package constant

import "fmt"

const (
	LocaleEN = "en"
	LocaleRU = "ru"
	LocaleUK = "uk"
)

var supportedLocales = []string{
	LocaleEN,
	LocaleRU,
	LocaleUK,
}

var ErrorUnsupportedLocale = fmt.Errorf("unsupported locale, supported locales are: %v", supportedLocales)

// IsValidLocaleSupported checks if the given locale is supported
func IsValidLocaleSupported(locale string) error {
	for _, l := range supportedLocales {
		if l == locale {
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrorUnsupportedLocale, locale)
}

func GetAllLocales() []string {
	return supportedLocales
}
