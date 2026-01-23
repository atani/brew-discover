package category

import (
	"strings"
)

type Category string

const (
	CategoryDev      Category = "dev"
	CategoryMedia    Category = "media"
	CategoryUtils    Category = "utils"
	CategoryNetwork  Category = "network"
	CategorySecurity Category = "security"
	CategoryData     Category = "data"
	CategoryGames    Category = "games"
	CategoryOther    Category = "other"
)

var AllCategories = []Category{
	CategoryDev,
	CategoryMedia,
	CategoryUtils,
	CategoryNetwork,
	CategorySecurity,
	CategoryData,
	CategoryGames,
	CategoryOther,
}

var categoryEmojis = map[Category]string{
	CategoryDev:      "🛠️",
	CategoryMedia:    "🎬",
	CategoryUtils:    "🔧",
	CategoryNetwork:  "🌐",
	CategorySecurity: "🔒",
	CategoryData:     "📊",
	CategoryGames:    "🎮",
	CategoryOther:    "📦",
}

var categoryKeywords = map[Category][]string{
	CategoryDev: {
		"compiler", "build", "git", "docker", "kubernetes", "k8s",
		"python", "node", "ruby", "rust", "go", "java", "kotlin", "swift",
		"typescript", "javascript", "php", "perl", "lua", "clojure", "elixir",
		"haskell", "scala", "erlang", "ocaml", "nim", "zig", "crystal",
		"cmake", "make", "ninja", "gradle", "maven", "cargo", "npm", "yarn", "pnpm",
		"debug", "lint", "format", "test", "coverage", "profil",
		"ide", "editor", "vim", "neovim", "emacs", "code",
		"sdk", "jdk", "runtime", "interpreter", "repl",
		"terraform", "ansible", "puppet", "chef",
		"ci", "cd", "pipeline", "devops",
	},
	CategoryMedia: {
		"video", "audio", "image", "photo", "music", "sound",
		"ffmpeg", "vlc", "media", "stream", "player", "record",
		"youtube", "spotify", "podcast",
		"encode", "decode", "transcode", "convert",
		"mp3", "mp4", "mkv", "avi", "mov", "flac", "wav", "ogg",
		"png", "jpg", "jpeg", "gif", "webp", "svg", "heic",
		"camera", "screen", "capture", "screenshot",
		"subtitle", "caption",
	},
	CategoryUtils: {
		"file", "disk", "backup", "sync", "archive", "compress",
		"zip", "tar", "gzip", "bzip", "xz", "7z", "rar",
		"monitor", "system", "process", "memory", "cpu",
		"terminal", "shell", "console", "cli",
		"text", "string", "regex", "search", "find", "grep",
		"time", "date", "calendar", "clock",
		"calculator", "math", "convert",
		"clipboard", "paste", "copy",
		"notification", "alert", "reminder",
		"pdf", "document", "office",
		"font", "color", "theme",
	},
	CategoryNetwork: {
		"http", "https", "ftp", "sftp", "ssh", "ssl", "tls",
		"vpn", "proxy", "tunnel", "firewall",
		"dns", "dhcp", "ip", "tcp", "udp", "socket",
		"curl", "wget", "fetch", "download", "upload",
		"api", "rest", "graphql", "grpc", "websocket",
		"email", "smtp", "imap", "pop3", "mail",
		"chat", "message", "irc", "slack", "discord",
		"browser", "web", "internet",
		"network", "wifi", "bluetooth", "wireless",
		"server", "client", "host", "remote",
	},
	CategorySecurity: {
		"encrypt", "decrypt", "crypto", "cipher",
		"password", "secret", "key", "token", "auth",
		"hash", "md5", "sha", "bcrypt", "argon",
		"certificate", "ssl", "tls", "gpg", "pgp",
		"security", "secure", "protect", "guard",
		"scan", "audit", "vulnerability", "exploit",
		"firewall", "ids", "ips", "waf",
		"antivirus", "malware", "virus",
		"2fa", "mfa", "otp", "totp",
		"vault", "keychain", "keyring",
	},
	CategoryData: {
		"database", "sql", "mysql", "postgres", "sqlite", "mariadb",
		"mongodb", "redis", "memcache", "elasticsearch",
		"data", "analytics", "statistics", "machine learning", "ml",
		"csv", "json", "xml", "yaml", "toml",
		"etl", "pipeline", "transform",
		"graph", "chart", "plot", "visualization",
		"pandas", "numpy", "scipy", "tensorflow", "pytorch",
		"bigdata", "hadoop", "spark", "kafka",
		"warehouse", "lake", "bi",
	},
	CategoryGames: {
		"game", "gaming", "play", "entertainment",
		"emulator", "rom", "retro",
		"steam", "epic",
		"puzzle", "chess", "tetris",
		"roguelike", "rpg", "fps", "mmo",
	},
}

func Classify(name, desc string) Category {
	combined := strings.ToLower(name + " " + desc)

	// Score each category
	scores := make(map[Category]int)
	for cat, keywords := range categoryKeywords {
		for _, kw := range keywords {
			if strings.Contains(combined, kw) {
				scores[cat]++
			}
		}
	}

	// Find highest scoring category
	bestCat := CategoryOther
	bestScore := 0
	for cat, score := range scores {
		if score > bestScore {
			bestScore = score
			bestCat = cat
		}
	}

	return bestCat
}

func GetEmoji(cat Category) string {
	if emoji, ok := categoryEmojis[cat]; ok {
		return emoji
	}
	return "📦"
}

func GetCategories(name, desc string) []Category {
	combined := strings.ToLower(name + " " + desc)

	var cats []Category
	for cat, keywords := range categoryKeywords {
		for _, kw := range keywords {
			if strings.Contains(combined, kw) {
				cats = append(cats, cat)
				break
			}
		}
	}

	if len(cats) == 0 {
		cats = append(cats, CategoryOther)
	}

	return cats
}
