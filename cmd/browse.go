package cmd

import (
	"fmt"
	"sort"

	"github.com/atani/brew-discover/internal/api"
	"github.com/atani/brew-discover/internal/cache"
	"github.com/atani/brew-discover/internal/category"
	"github.com/atani/brew-discover/internal/i18n"
	"github.com/atani/brew-discover/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	browseCask  bool
	browseLimit int
)

var browseCmd = &cobra.Command{
	Use:   "browse [category]",
	Short: "Browse packages by category",
	Long: `Browse Homebrew packages organized by category.
Available categories: dev, media, utils, network, security, data, games

Without a category argument, shows all categories with package counts.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBrowse,
}

func init() {
	rootCmd.AddCommand(browseCmd)
	browseCmd.Flags().BoolVar(&browseCask, "cask", false, "Browse casks instead of formulae")
	browseCmd.Flags().IntVarP(&browseLimit, "limit", "l", 20, "Number of packages to show per category")
}

type categoryInfo struct {
	cat   category.Category
	count int
	items []categoryItem
}

type categoryItem struct {
	name     string
	desc     string
	count    int
	countStr string
}

func runBrowse(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	client := api.NewClient()
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Build category index
	categories, err := buildCategoryIndex(client, c, browseCask)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		// Show category overview
		printCategoryOverview(categories)
	} else {
		// Show specific category
		catName := category.Category(args[0])
		if info, ok := categories[catName]; ok {
			printCategoryDetail(catName, info)
		} else {
			return fmt.Errorf("unknown category: %s", args[0])
		}
	}

	return nil
}

func buildCategoryIndex(client *api.Client, c *cache.Cache, isCask bool) (map[category.Category]*categoryInfo, error) {
	result := make(map[category.Category]*categoryInfo)
	for _, cat := range category.AllCategories {
		result[cat] = &categoryInfo{cat: cat}
	}

	if isCask {
		casks, err := getCaskMap(client, c)
		if err != nil {
			return nil, err
		}

		analytics, _ := getCaskAnalytics(client, c)
		countMap := make(map[string]api.AnalyticsItem)
		if analytics != nil {
			for _, item := range analytics.GetTopCasks(10000) {
				countMap[item.Name] = item
			}
		}

		for token, cask := range casks {
			cat := category.Classify(token, cask.Desc)
			info := result[cat]
			info.count++

			count := 0
			countStr := ""
			if item, ok := countMap[token]; ok {
				count = item.Count
				countStr = item.CountStr
			}

			info.items = append(info.items, categoryItem{
				name:     token,
				desc:     cask.Desc,
				count:    count,
				countStr: countStr,
			})
		}
	} else {
		formulae, err := getFormulaMap(client, c)
		if err != nil {
			return nil, err
		}

		analytics, _ := getFormulaAnalytics(client, c)
		countMap := make(map[string]api.AnalyticsItem)
		if analytics != nil {
			for _, item := range analytics.GetTopFormulae(10000) {
				countMap[item.Name] = item
			}
		}

		for name, formula := range formulae {
			cat := category.Classify(name, formula.Desc)
			info := result[cat]
			info.count++

			count := 0
			countStr := ""
			if item, ok := countMap[name]; ok {
				count = item.Count
				countStr = item.CountStr
			}

			info.items = append(info.items, categoryItem{
				name:     name,
				desc:     formula.Desc,
				count:    count,
				countStr: countStr,
			})
		}
	}

	// Sort items in each category by count
	for _, info := range result {
		sort.Slice(info.items, func(i, j int) bool {
			return info.items[i].count > info.items[j].count
		})
	}

	return result, nil
}

func printCategoryOverview(categories map[category.Category]*categoryInfo) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13"))
	catStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	countStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	fmt.Println()
	fmt.Println(titleStyle.Render("📂 " + i18n.T("browse.select")))
	fmt.Println()

	// Order categories
	orderedCats := []category.Category{
		category.CategoryDev,
		category.CategoryMedia,
		category.CategoryUtils,
		category.CategoryNetwork,
		category.CategorySecurity,
		category.CategoryData,
		category.CategoryGames,
		category.CategoryOther,
	}

	for _, cat := range orderedCats {
		info := categories[cat]
		emoji := category.GetEmoji(cat)
		name := i18n.T("category." + string(cat))
		countStr := fmt.Sprintf("(%d packages)", info.count)

		fmt.Printf("  %s %s %s\n",
			emoji,
			catStyle.Render(fmt.Sprintf("%-10s", cat)),
			descStyle.Render(name)+" "+countStyle.Render(countStr))
	}

	fmt.Println()
	fmt.Println(countStyle.Render("Usage: brew-discover browse <category>"))
	fmt.Println()
}

func printCategoryDetail(cat category.Category, info *categoryInfo) {
	title := category.GetEmoji(cat) + " " + i18n.T("category."+string(cat))

	limit := browseLimit
	if limit > len(info.items) {
		limit = len(info.items)
	}

	rows := make([]ui.TableRow, limit)
	for i := 0; i < limit; i++ {
		item := info.items[i]
		rows[i] = ui.TableRow{
			Rank:        i + 1,
			Name:        item.name,
			Count:       item.countStr,
			Description: item.desc,
		}
	}

	ui.PrintTable(title, rows, "")
}
