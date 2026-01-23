package i18n

import (
	"os"
	"testing"
)

func TestSetLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"en", "en"},
		{"ja", "ja"},
		{"EN", "en"},
		{"JA", "ja"},
		{"unknown", "en"}, // Falls back to default
	}

	for _, tt := range tests {
		SetLanguage(tt.input)
		if got := GetLanguage(); got != tt.expected {
			t.Errorf("SetLanguage(%q): got %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestT_English(t *testing.T) {
	SetLanguage("en")

	tests := []struct {
		key      string
		expected string
	}{
		{"top.tip", "Tip: Run `brew install <name>` to install"},
		{"random.title", "Random Formula Pick"},
		{"category.dev", "Development Tools"},
	}

	for _, tt := range tests {
		if got := T(tt.key); got != tt.expected {
			t.Errorf("T(%q): got %q, want %q", tt.key, got, tt.expected)
		}
	}
}

func TestT_Japanese(t *testing.T) {
	SetLanguage("ja")

	tests := []struct {
		key      string
		expected string
	}{
		{"top.tip", "Tip: `brew install <name>` でインストール"},
		{"random.title", "今日のおすすめ Formula"},
		{"category.dev", "開発ツール"},
	}

	for _, tt := range tests {
		if got := T(tt.key); got != tt.expected {
			t.Errorf("T(%q): got %q, want %q", tt.key, got, tt.expected)
		}
	}
}

func TestT_WithArgs(t *testing.T) {
	SetLanguage("en")

	got := T("top.title.formula", "Count", 10, "Days", 30)
	expected := "Homebrew Formula Top 10 (installs in last 30 days)"

	if got != expected {
		t.Errorf("T with args: got %q, want %q", got, expected)
	}
}

func TestT_UnknownKey(t *testing.T) {
	SetLanguage("en")

	key := "unknown.key.that.does.not.exist"
	if got := T(key); got != key {
		t.Errorf("T unknown key: got %q, want %q", got, key)
	}
}

func TestDetectLanguage_FromEnv(t *testing.T) {
	// Save original
	origLang := os.Getenv("LANG")
	origBrewLang := os.Getenv("BREW_DISCOVER_LANG")
	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("BREW_DISCOVER_LANG", origBrewLang)
	}()

	// Test BREW_DISCOVER_LANG takes priority
	os.Setenv("BREW_DISCOVER_LANG", "ja")
	os.Setenv("LANG", "en_US.UTF-8")
	DetectLanguage()
	if got := GetLanguage(); got != "ja" {
		t.Errorf("DetectLanguage with BREW_DISCOVER_LANG=ja: got %q, want %q", got, "ja")
	}

	// Test fallback to LANG
	os.Unsetenv("BREW_DISCOVER_LANG")
	os.Setenv("LANG", "ja_JP.UTF-8")
	DetectLanguage()
	if got := GetLanguage(); got != "ja" {
		t.Errorf("DetectLanguage with LANG=ja_JP.UTF-8: got %q, want %q", got, "ja")
	}
}
