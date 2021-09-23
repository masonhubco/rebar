package rebar

import (
	"errors"

	internationalization "github.com/qor/i18n"
	"github.com/qor/i18n/backends/yaml"
)

// Multiple nested lookups to allow for deeply packaged tests to have access to i18n
var i18n = internationalization.New(yaml.NewWithWalk("languages", "../languages", "../../languages"))

type String struct {
	i18nKey string
	args    []interface{}
}

// I18n will construct a new internationalized ValidationError
func I18n(i18nKey string, args ...interface{}) String {
	return String{
		i18nKey: i18nKey,
		args:    args,
	}
}

type LanguageScoped struct {
	Language string
}

func (l LanguageScoped) String(str String) string {
	return string(i18n.T(l.Language, str.i18nKey, str.args...))
}

func (l LanguageScoped) Raw(key string, args ...interface{}) string {
	return string(i18n.T(l.Language, key, args...))
}

func (l LanguageScoped) Error(key string, args ...interface{}) error {
	return errors.New(string(i18n.T(l.Language, key, args...)))
}

func WithLanguage(language string) LanguageScoped {
	if language == "" {
		language = "en"
	}
	return LanguageScoped{
		Language: language,
	}
}

// Translate will return the translated key based on the accept string (i.e. en-US, es)
func Translate(language, key string, args ...interface{}) string {
	if language == "" {
		language = "en"
	}
	return string(i18n.T(language, key, args...))
}
