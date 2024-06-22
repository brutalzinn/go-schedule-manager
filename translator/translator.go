package translator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type TranslatorBundle struct {
	localizer *i18n.Localizer
	bundle    *i18n.Bundle
	language  language.Tag
}

func New(language language.Tag) TranslatorBundle {
	bundle := i18n.NewBundle(language)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	localizer := i18n.NewLocalizer(bundle, language.String())
	if err := loadTranslations(bundle, "./translations"); err != nil {
		log.Fatalf("Failed to autoload translations: %v", err)
	}
	translator := TranslatorBundle{
		bundle:    bundle,
		language:  language,
		localizer: localizer,
	}
	return translator
}

// to translate messages. This because Home assistant can be in english too. Like us particular server
func (translator TranslatorBundle) T(messageID string) string {
	translated, err := translator.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID
	}
	return translated
}

func loadTranslations(bundle *i18n.Bundle, dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			_, err := bundle.LoadMessageFile(path)
			if err != nil {
				return fmt.Errorf("failed to load translation file %s: %v", path, err)
			}
			fmt.Printf("Loaded translation file: %s\n", path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking directory %s: %v", dirPath, err)
	}
	return nil
}
