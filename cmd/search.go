package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/atani/brew-discover/internal/api"
	"github.com/atani/brew-discover/internal/cache"
	"github.com/atani/brew-discover/internal/i18n"
	"github.com/atani/brew-discover/internal/ui"
	"github.com/spf13/cobra"
)

var (
	searchCask  bool
	searchLimit int
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search packages with enhanced results",
	Long:  `Search Homebrew packages by name and description, sorted by popularity.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVar(&searchCask, "cask", false, "Search casks instead of formulae")
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 20, "Maximum number of results")
}

type searchResult struct {
	name     string
	desc     string
	count    int
	countStr string
	matchIn  string // "name" or "desc"
}

func runSearch(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	query := strings.ToLower(args[0])
	client := api.NewClient()
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	var results []searchResult

	if searchCask {
		results, err = searchCasks(client, c, query)
	} else {
		results, err = searchFormulae(client, c, query)
	}

	if err != nil {
		return err
	}

	// Sort by install count (popularity)
	sort.Slice(results, func(i, j int) bool {
		return results[i].count > results[j].count
	})

	// Limit results
	if len(results) > searchLimit {
		results = results[:searchLimit]
	}

	// Print results
	title := i18n.T("search.results", "Query", args[0])

	rows := make([]ui.TableRow, len(results))
	for i, r := range results {
		rows[i] = ui.TableRow{
			Rank:        i + 1,
			Name:        r.name,
			Count:       r.countStr,
			Description: r.desc,
		}
	}

	totalMsg := i18n.T("search.count", "Count", len(results))
	ui.PrintTable(title, rows, totalMsg)

	return nil
}

func searchFormulae(client *api.Client, c *cache.Cache, query string) ([]searchResult, error) {
	formulae, err := getFormulaMap(client, c)
	if err != nil {
		return nil, err
	}

	analytics, _ := getFormulaAnalytics(client, c)

	// Build count map
	countMap := make(map[string]api.AnalyticsItem)
	if analytics != nil {
		for _, item := range analytics.GetTopFormulae(10000) {
			countMap[item.Name] = item
		}
	}

	var results []searchResult
	for name, formula := range formulae {
		nameLower := strings.ToLower(name)
		descLower := strings.ToLower(formula.Desc)

		var matchIn string
		if strings.Contains(nameLower, query) {
			matchIn = "name"
		} else if strings.Contains(descLower, query) {
			matchIn = "desc"
		} else {
			continue
		}

		count := 0
		countStr := ""
		if item, ok := countMap[name]; ok {
			count = item.Count
			countStr = item.CountStr
		}

		results = append(results, searchResult{
			name:     name,
			desc:     formula.Desc,
			count:    count,
			countStr: countStr,
			matchIn:  matchIn,
		})
	}

	return results, nil
}

func searchCasks(client *api.Client, c *cache.Cache, query string) ([]searchResult, error) {
	casks, err := getCaskMap(client, c)
	if err != nil {
		return nil, err
	}

	analytics, _ := getCaskAnalytics(client, c)

	// Build count map
	countMap := make(map[string]api.AnalyticsItem)
	if analytics != nil {
		for _, item := range analytics.GetTopCasks(10000) {
			countMap[item.Name] = item
		}
	}

	var results []searchResult
	for token, cask := range casks {
		tokenLower := strings.ToLower(token)
		descLower := strings.ToLower(cask.Desc)

		var matchIn string
		if strings.Contains(tokenLower, query) {
			matchIn = "name"
		} else if strings.Contains(descLower, query) {
			matchIn = "desc"
		} else {
			continue
		}

		count := 0
		countStr := ""
		if item, ok := countMap[token]; ok {
			count = item.Count
			countStr = item.CountStr
		}

		results = append(results, searchResult{
			name:     token,
			desc:     cask.Desc,
			count:    count,
			countStr: countStr,
			matchIn:  matchIn,
		})
	}

	return results, nil
}
