# brew-discover

Discover new Homebrew packages through popularity rankings, category browsing, and random recommendations.

![Demo](assets/demo.gif)

## Features

- **Top Packages**: View the most popular packages by install count
- **Category Browsing**: Explore packages organized by category (dev, media, utils, network, security, data, games)
- **Random Discovery**: Get random package recommendations
- **Enhanced Search**: Search packages with popularity-sorted results
- **Detailed Info**: View package details with popularity stats
- **Bilingual**: Supports English and Japanese

## Installation

```bash
brew tap atani/tap
brew install brew-discover
```

Or download the binary from [Releases](https://github.com/atani/brew-discover/releases).

## Usage

### Top Packages

```bash
# Show top 20 formulae
brew-discover top

# Show top 10 casks
brew-discover top --cask -n 10
```

### Browse by Category

```bash
# List all categories
brew-discover browse

# Browse development tools
brew-discover browse dev

# Available categories: dev, media, utils, network, security, data, games
```

### Random Recommendations

```bash
# Get a random recommendation
brew-discover random

# Get 5 random picks
brew-discover random -n 5

# Random pick with install prompt
brew-discover random --lucky
```

### Search

```bash
# Search packages by name and description
brew-discover search editor

# Search casks
brew-discover search browser --cask
```

### Package Info

```bash
# Show detailed info with popularity stats
brew-discover info bat
```

### Language

```bash
# Use Japanese
brew-discover top --lang ja

# Or set environment variable
export BREW_DISCOVER_LANG=ja
```

## Development

```bash
# Build
go build -o brew-discover

# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## License

MIT
