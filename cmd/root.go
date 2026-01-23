package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	lang    string
	refresh bool
)

var rootCmd = &cobra.Command{
	Use:   "brew-discover",
	Short: "Discover new Homebrew packages",
	Long: `brew-discover helps you discover new and popular Homebrew packages.

Features:
  - View top packages by install count
  - Browse packages by category
  - Get random package recommendations
  - Search with enhanced results`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func SetVersion(v string) {
	rootCmd.Version = v
}

func init() {
	rootCmd.PersistentFlags().StringVar(&lang, "lang", "", "Language (en, ja)")
	rootCmd.PersistentFlags().BoolVar(&refresh, "refresh", false, "Refresh cache")
}
