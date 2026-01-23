package i18n

import (
	"embed"
	"encoding/json"
	"os"
	"strings"
	"text/template"
)

//go:embed locales/*.json
var localesFS embed.FS

var (
	currentLang   = "en"
	translations  = make(map[string]map[string]string)
	defaultLang   = "en"
)

func init() {
	// Load all locale files
	loadLocale("en")
	loadLocale("ja")

	// Detect language from environment
	DetectLanguage()
}

func loadLocale(lang string) {
	data, err := localesFS.ReadFile("locales/" + lang + ".json")
	if err != nil {
		return
	}

	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		return
	}

	translations[lang] = messages
}

func DetectLanguage() {
	// Check environment variable
	if lang := os.Getenv("BREW_DISCOVER_LANG"); lang != "" {
		SetLanguage(lang)
		return
	}

	// Check LANG environment variable
	if lang := os.Getenv("LANG"); lang != "" {
		// Extract language code from "ja_JP.UTF-8" -> "ja"
		parts := strings.Split(lang, "_")
		if len(parts) > 0 {
			langCode := strings.ToLower(parts[0])
			if _, ok := translations[langCode]; ok {
				SetLanguage(langCode)
				return
			}
		}
	}

	// Default to English
	SetLanguage(defaultLang)
}

func SetLanguage(lang string) {
	lang = strings.ToLower(lang)
	if _, ok := translations[lang]; ok {
		currentLang = lang
	} else {
		currentLang = defaultLang
	}
}

func GetLanguage() string {
	return currentLang
}

func T(key string, args ...any) string {
	messages, ok := translations[currentLang]
	if !ok {
		messages = translations[defaultLang]
	}

	msg, ok := messages[key]
	if !ok {
		// Fallback to default language
		if defaultMessages, ok := translations[defaultLang]; ok {
			if defaultMsg, ok := defaultMessages[key]; ok {
				msg = defaultMsg
			} else {
				return key
			}
		} else {
			return key
		}
	}

	// If args provided, use template
	if len(args) > 0 && len(args)%2 == 0 {
		data := make(map[string]any)
		for i := 0; i < len(args); i += 2 {
			if k, ok := args[i].(string); ok {
				data[k] = args[i+1]
			}
		}

		tmpl, err := template.New("msg").Parse(msg)
		if err != nil {
			return msg
		}

		var buf strings.Builder
		if err := tmpl.Execute(&buf, data); err != nil {
			return msg
		}
		return buf.String()
	}

	return msg
}
