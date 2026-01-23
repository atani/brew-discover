package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/atani/brew-discover/internal/api"
	"github.com/atani/brew-discover/internal/cache"
	"github.com/atani/brew-discover/internal/i18n"
	"github.com/atani/brew-discover/internal/ui"
	"github.com/spf13/cobra"
)

var (
	topCount int
	topCask  bool
)

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Show top packages by install count",
	Long:  `Display the most popular Homebrew packages sorted by install count over the last 30 days.`,
	RunE:  runTop,
}

func init() {
	rootCmd.AddCommand(topCmd)
	topCmd.Flags().IntVarP(&topCount, "number", "n", 20, "Number of packages to show")
	topCmd.Flags().BoolVar(&topCask, "cask", false, "Show top casks instead of formulae")
}

func runTop(cmd *cobra.Command, args []string) error {
	// Set language if specified
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	client := api.NewClient()
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	var analytics *api.AnalyticsResponse
	var items []api.AnalyticsItem
	var formulae map[string]*api.Formula
	var casks map[string]*api.Cask

	if topCask {
		// Get cask analytics
		cacheKey := fmt.Sprintf(cache.CaskAnalytics, api.Period30Days)
		if data, ok := c.Get(cacheKey, cache.AnalyticsTTL); ok && !refresh {
			if err := json.Unmarshal(data, &analytics); err == nil {
				goto gotCaskAnalytics
			}
		}

		analytics, err = client.GetCaskAnalytics(api.Period30Days)
		if err != nil {
			return fmt.Errorf("failed to get cask analytics: %w", err)
		}

		if data, err := json.Marshal(analytics); err == nil {
			c.Set(cacheKey, data)
		}

	gotCaskAnalytics:
		items = analytics.GetTopCasks(topCount)

		// Get cask details for descriptions
		casks, err = getCaskMap(client, c)
		if err != nil {
			return fmt.Errorf("failed to get casks: %w", err)
		}
	} else {
		// Get formula analytics
		cacheKey := fmt.Sprintf(cache.FormulaAnalytics, api.Period30Days)
		if data, ok := c.Get(cacheKey, cache.AnalyticsTTL); ok && !refresh {
			if err := json.Unmarshal(data, &analytics); err == nil {
				goto gotFormulaAnalytics
			}
		}

		analytics, err = client.GetFormulaAnalytics(api.Period30Days)
		if err != nil {
			return fmt.Errorf("failed to get formula analytics: %w", err)
		}

		if data, err := json.Marshal(analytics); err == nil {
			c.Set(cacheKey, data)
		}

	gotFormulaAnalytics:
		items = analytics.GetTopFormulae(topCount)

		// Get formula details for descriptions
		formulae, err = getFormulaMap(client, c)
		if err != nil {
			return fmt.Errorf("failed to get formulae: %w", err)
		}
	}

	// Build table rows
	rows := make([]ui.TableRow, 0, len(items))
	for _, item := range items {
		desc := ""

		if topCask {
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

	// Print table
	var title string
	if topCask {
		title = i18n.T("top.title.cask", "Count", len(rows), "Days", 30)
	} else {
		title = i18n.T("top.title.formula", "Count", len(rows), "Days", 30)
	}

	ui.PrintTable(title, rows, i18n.T("top.tip"))

	return nil
}

func getFormulaMap(client *api.Client, c *cache.Cache) (map[string]*api.Formula, error) {
	var formulae []api.Formula

	if data, ok := c.Get(cache.FormulaeFile, cache.FormulaTTL); ok && !refresh {
		if err := json.Unmarshal(data, &formulae); err == nil {
			goto buildMap
		}
	}

	{
		data, err := client.GetFormulaeBytes()
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &formulae); err != nil {
			return nil, err
		}

		c.Set(cache.FormulaeFile, data)
	}

buildMap:
	result := make(map[string]*api.Formula)
	for i := range formulae {
		result[formulae[i].Name] = &formulae[i]
	}
	return result, nil
}

func getCaskMap(client *api.Client, c *cache.Cache) (map[string]*api.Cask, error) {
	var casks []api.Cask

	if data, ok := c.Get(cache.CasksFile, cache.FormulaTTL); ok && !refresh {
		if err := json.Unmarshal(data, &casks); err == nil {
			goto buildMap
		}
	}

	{
		data, err := client.GetCasksBytes()
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &casks); err != nil {
			return nil, err
		}

		c.Set(cache.CasksFile, data)
	}

buildMap:
	result := make(map[string]*api.Cask)
	for i := range casks {
		result[casks[i].Token] = &casks[i]
	}
	return result, nil
}
