package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"

	"github.com/atani/brew-discover/internal/api"
	"github.com/atani/brew-discover/internal/cache"
	"github.com/atani/brew-discover/internal/i18n"
	"github.com/atani/brew-discover/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	randomCount int
	randomCask  bool
	randomLucky bool
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get random package recommendations",
	Long:  `Discover new packages through random recommendations from popular Homebrew packages.`,
	RunE:  runRandom,
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.Flags().IntVarP(&randomCount, "number", "n", 1, "Number of random packages to show")
	randomCmd.Flags().BoolVar(&randomCask, "cask", false, "Pick from casks instead of formulae")
	randomCmd.Flags().BoolVar(&randomLucky, "lucky", false, "Immediately prompt to install")
}

func runRandom(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	client := api.NewClient()
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	var items []api.AnalyticsItem
	var formulae map[string]*api.Formula
	var casks map[string]*api.Cask

	if randomCask {
		analytics, err := getCaskAnalytics(client, c)
		if err != nil {
			return err
		}
		allItems := analytics.GetTopCasks(500) // Pick from top 500
		items = pickRandom(allItems, randomCount)

		casks, err = getCaskMap(client, c)
		if err != nil {
			return err
		}
	} else {
		analytics, err := getFormulaAnalytics(client, c)
		if err != nil {
			return err
		}
		allItems := analytics.GetTopFormulae(500) // Pick from top 500
		items = pickRandom(allItems, randomCount)

		formulae, err = getFormulaMap(client, c)
		if err != nil {
			return err
		}
	}

	if randomCount == 1 {
		// Single package - show detailed view
		item := items[0]
		printRandomDetail(item, formulae, casks, randomCask)

		if randomLucky {
			promptInstall(item.Name, randomCask)
		}
	} else {
		// Multiple packages - show table
		title := i18n.T("random.title.plural", "Count", randomCount)
		rows := make([]ui.TableRow, 0, len(items))

		for _, item := range items {
			desc := ""
			if randomCask {
				if cask, ok := casks[item.Name]; ok {
					desc = cask.Desc
				}
			} else {
				if formula, ok := formulae[item.Name]; ok {
					desc = formula.Desc
				}
			}

			rows = append(rows, ui.TableRow{
				Rank:        item.Rank,
				Name:        item.Name,
				Count:       item.CountStr,
				Description: desc,
			})
		}

		ui.PrintTable(title, rows, i18n.T("random.tip"))
	}

	return nil
}

func pickRandom(items []api.AnalyticsItem, count int) []api.AnalyticsItem {
	if count >= len(items) {
		return items
	}

	// Fisher-Yates shuffle
	shuffled := make([]api.AnalyticsItem, len(items))
	copy(shuffled, items)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:count]
}

func printRandomDetail(item api.AnalyticsItem, formulae map[string]*api.Formula, casks map[string]*api.Cask, isCask bool) {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2)

	var content strings.Builder

	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))

	content.WriteString(nameStyle.Render(item.Name) + "\n")
	content.WriteString(strings.Repeat("─", 50) + "\n")

	var desc, homepage string
	if isCask {
		if cask, ok := casks[item.Name]; ok {
			desc = cask.Desc
			homepage = cask.Homepage
		}
	} else {
		if formula, ok := formulae[item.Name]; ok {
			desc = formula.Desc
			homepage = formula.Homepage
		}
	}

	content.WriteString(desc + "\n\n")
	content.WriteString(labelStyle.Render(i18n.T("info.installs")+": ") + item.CountStr + " (30d)\n")
	content.WriteString(labelStyle.Render(i18n.T("info.ranking")+": ") + fmt.Sprintf("#%d", item.Rank) + "\n")
	if homepage != "" {
		content.WriteString(labelStyle.Render(i18n.T("info.homepage")+": ") + homepage + "\n")
	}

	fmt.Println()
	fmt.Println(boxStyle.Render(content.String()))
	fmt.Println()
}

func promptInstall(name string, isCask bool) {
	fmt.Print(i18n.T("random.install_prompt") + " ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		var cmd *exec.Cmd
		if isCask {
			cmd = exec.Command("brew", "install", "--cask", name)
		} else {
			cmd = exec.Command("brew", "install", name)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}
}

func getFormulaAnalytics(client *api.Client, c *cache.Cache) (*api.AnalyticsResponse, error) {
	cacheKey := fmt.Sprintf(cache.FormulaAnalytics, api.Period30Days)
	if data, ok := c.Get(cacheKey, cache.AnalyticsTTL); ok && !refresh {
		var analytics api.AnalyticsResponse
		if err := json.Unmarshal(data, &analytics); err == nil {
			return &analytics, nil
		}
	}

	analytics, err := client.GetFormulaAnalytics(api.Period30Days)
	if err != nil {
		return nil, fmt.Errorf("failed to get formula analytics: %w", err)
	}

	if data, err := json.Marshal(analytics); err == nil {
		_ = c.Set(cacheKey, data)
	}

	return analytics, nil
}

func getCaskAnalytics(client *api.Client, c *cache.Cache) (*api.AnalyticsResponse, error) {
	cacheKey := fmt.Sprintf(cache.CaskAnalytics, api.Period30Days)
	if data, ok := c.Get(cacheKey, cache.AnalyticsTTL); ok && !refresh {
		var analytics api.AnalyticsResponse
		if err := json.Unmarshal(data, &analytics); err == nil {
			return &analytics, nil
		}
	}

	analytics, err := client.GetCaskAnalytics(api.Period30Days)
	if err != nil {
		return nil, fmt.Errorf("failed to get cask analytics: %w", err)
	}

	if data, err := json.Marshal(analytics); err == nil {
		_ = c.Set(cacheKey, data)
	}

	return analytics, nil
}
