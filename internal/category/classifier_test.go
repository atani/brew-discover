package category

import (
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		name     string
		desc     string
		expected Category
	}{
		{"python", "Interpreted, interactive, object-oriented programming language", CategoryDev},
		{"node", "Platform for server-side JavaScript applications", CategoryDev},
		{"cmake", "Cross-platform make", CategoryDev},
		{"ffmpeg", "Play, record, convert, and stream audio and video", CategoryMedia},
		{"yt-dlp", "Fork of youtube-dl with additional features", CategoryMedia},
		{"curl", "Get a file from an HTTP, HTTPS or FTP server", CategoryNetwork},
		{"wget", "Internet file retriever", CategoryNetwork},
		{"htop", "Improved top (interactive process viewer)", CategoryUtils},
		{"zip", "Compression and file packaging/archive utility", CategoryUtils},
		{"age", "Simple, modern and secure encryption tool", CategorySecurity},
		{"openssl", "Cryptography and SSL/TLS Toolkit", CategorySecurity},
		{"postgresql", "Object-relational database system", CategoryData},
		{"redis", "Persistent key-value database", CategoryData},
		{"xyzabc", "A completely random thing", CategoryOther},
	}

	for _, tt := range tests {
		got := Classify(tt.name, tt.desc)
		if got != tt.expected {
			t.Errorf("Classify(%q, %q): got %q, want %q", tt.name, tt.desc, got, tt.expected)
		}
	}
}

func TestGetEmoji(t *testing.T) {
	tests := []struct {
		cat      Category
		expected string
	}{
		{CategoryDev, "🛠️"},
		{CategoryMedia, "🎬"},
		{CategoryUtils, "🔧"},
		{CategoryNetwork, "🌐"},
		{CategorySecurity, "🔒"},
		{CategoryData, "📊"},
		{CategoryGames, "🎮"},
		{CategoryOther, "📦"},
	}

	for _, tt := range tests {
		got := GetEmoji(tt.cat)
		if got != tt.expected {
			t.Errorf("GetEmoji(%q): got %q, want %q", tt.cat, got, tt.expected)
		}
	}
}

func TestGetCategories(t *testing.T) {
	// Package that matches multiple categories
	cats := GetCategories("docker", "Platform to build, run, and share containerized applications with kubernetes")

	if len(cats) == 0 {
		t.Error("GetCategories should return at least one category")
	}

	// Check that dev is included (matches docker, kubernetes, build)
	hasDev := false
	for _, c := range cats {
		if c == CategoryDev {
			hasDev = true
			break
		}
	}
	if !hasDev {
		t.Error("GetCategories for docker should include CategoryDev")
	}
}

func TestAllCategories(t *testing.T) {
	expected := 8 // dev, media, utils, network, security, data, games, other
	if len(AllCategories) != expected {
		t.Errorf("AllCategories length: got %d, want %d", len(AllCategories), expected)
	}
}
