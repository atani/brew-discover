package cmd

import (
	"fmt"
	"strings"

	"github.com/atani/brew-discover/internal/api"
	"github.com/atani/brew-discover/internal/cache"
	"github.com/atani/brew-discover/internal/i18n"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var infoCask bool

var infoCmd = &cobra.Command{
	Use:   "info <package>",
	Short: "Show detailed package information",
	Long:  `Display detailed information about a Homebrew package including popularity stats.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().BoolVar(&infoCask, "cask", false, "Show info for a cask")
}

func runInfo(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	name := args[0]
	client := api.NewClient()
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	if infoCask {
		return showCaskInfo(client, c, name)
	}
	return showFormulaInfo(client, c, name)
}

func showFormulaInfo(client *api.Client, c *cache.Cache, name string) error {
	formula, err := client.GetFormula(name)
	if err != nil {
		return fmt.Errorf("%s", i18n.T("error.not_found", "Name", name))
	}

	// Get analytics for ranking
	analytics, _ := getFormulaAnalytics(client, c)
	var rank int
	var installs string
	if analytics != nil {
		allItems := analytics.GetTopFormulae(10000)
		for _, item := range allItems {
			if item.Name == name {
				rank = item.Rank
				installs = item.CountStr
				break
			}
		}
	}

	printFormulaInfo(formula, rank, installs, len(analytics.Formulae))
	return nil
}

func showCaskInfo(client *api.Client, c *cache.Cache, name string) error {
	cask, err := client.GetCask(name)
	if err != nil {
		return fmt.Errorf("%s", i18n.T("error.not_found", "Name", name))
	}

	// Get analytics for ranking
	analytics, _ := getCaskAnalytics(client, c)
	var rank int
	var installs string
	if analytics != nil {
		allItems := analytics.GetTopCasks(10000)
		for _, item := range allItems {
			if item.Name == name {
				rank = item.Rank
				installs = item.CountStr
				break
			}
		}
	}

	printCaskInfo(cask, rank, installs, len(analytics.Casks))
	return nil
}

func printFormulaInfo(f *api.Formula, rank int, installs string, total int) {
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	divider := strings.Repeat("─", 60)

	fmt.Println()
	fmt.Println(nameStyle.Render("📦 " + f.Name))
	fmt.Println()
	fmt.Println(f.Desc)
	fmt.Println()
	fmt.Println(divider)

	// Popularity section
	fmt.Println(sectionStyle.Render("📊 " + i18n.T("info.popularity")))
	if installs != "" {
		fmt.Printf("   %s %s (30d)\n", labelStyle.Render(i18n.T("info.installs")+":"), valueStyle.Render(installs))
	}
	if rank > 0 {
		rankStr := fmt.Sprintf("#%d / %d+ packages", rank, total)
		fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.ranking")+":"), valueStyle.Render(rankStr))
		fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.popularity")+":"), valueStyle.Render(getPopularityStars(rank, total)))
	}
	fmt.Println()

	// Details section
	fmt.Println(sectionStyle.Render("📋 " + i18n.T("info.details")))
	fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.version")+":"), valueStyle.Render(f.Versions.Stable))
	if f.License != "" {
		fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.license")+":"), valueStyle.Render(f.License))
	}
	fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.homepage")+":"), valueStyle.Render(f.Homepage))
	fmt.Println()

	// Dependencies section
	if len(f.Dependencies) > 0 {
		fmt.Println(sectionStyle.Render("📎 " + i18n.T("info.dependencies")))
		fmt.Printf("   %s\n", valueStyle.Render(strings.Join(f.Dependencies, ", ")))
		fmt.Println()
	}

	// Install command
	fmt.Println(sectionStyle.Render("💻 " + i18n.T("info.install")))
	fmt.Printf("   brew install %s\n", f.Name)
	fmt.Println(divider)
	fmt.Println()
}

func printCaskInfo(c *api.Cask, rank int, installs string, total int) {
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	divider := strings.Repeat("─", 60)

	fmt.Println()
	fmt.Println(nameStyle.Render("📦 " + c.Token))
	fmt.Println()
	fmt.Println(c.Desc)
	fmt.Println()
	fmt.Println(divider)

	// Popularity section
	fmt.Println(sectionStyle.Render("📊 " + i18n.T("info.popularity")))
	if installs != "" {
		fmt.Printf("   %s %s (30d)\n", labelStyle.Render(i18n.T("info.installs")+":"), valueStyle.Render(installs))
	}
	if rank > 0 {
		rankStr := fmt.Sprintf("#%d / %d+ packages", rank, total)
		fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.ranking")+":"), valueStyle.Render(rankStr))
		fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.popularity")+":"), valueStyle.Render(getPopularityStars(rank, total)))
	}
	fmt.Println()

	// Details section
	fmt.Println(sectionStyle.Render("📋 " + i18n.T("info.details")))
	fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.version")+":"), valueStyle.Render(c.Version))
	fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.homepage")+":"), valueStyle.Render(c.Homepage))
	fmt.Println()

	// Install command
	fmt.Println(sectionStyle.Render("💻 " + i18n.T("info.install")))
	fmt.Printf("   brew install --cask %s\n", c.Token)
	fmt.Println(divider)
	fmt.Println()
}

func getPopularityStars(rank, total int) string {
	percent := float64(rank) / float64(total) * 100

	var stars string
	switch {
	case percent <= 1:
		stars = "⭐⭐⭐⭐⭐ Top 1%"
	case percent <= 5:
		stars = "⭐⭐⭐⭐ Top 5%"
	case percent <= 10:
		stars = "⭐⭐⭐ Top 10%"
	case percent <= 25:
		stars = "⭐⭐ Top 25%"
	case percent <= 50:
		stars = "⭐ Top 50%"
	default:
		stars = "Below 50%"
	}

	return stars
}
