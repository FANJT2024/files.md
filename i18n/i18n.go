package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	lang            *i18n.Bundle
	emojisByKeyword map[string]string
)

// LoadLangFile only supports single language for now
func LoadLangFile(path string) error {
	lang = i18n.NewBundle(language.English)
	//_, err := lang.LoadMessageFile(path)
	//if err != nil {
	//	return fmt.Errorf("i18n.Load: %w", err)
	//}

	return nil
}

func Tr(str string) string {
	return str
}
