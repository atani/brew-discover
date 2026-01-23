package api

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type AnalyticsResponse struct {
	Category   string `json:"category"`
	TotalItems int    `json:"total_items"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	TotalCount int    `json:"total_count"`
	Formulae   map[string][]FormulaInstall `json:"formulae"`
	Casks      map[string][]CaskInstall    `json:"casks"`
}

type FormulaInstall struct {
	Formula string `json:"formula"`
	Count   string `json:"count"`
}

type CaskInstall struct {
	Cask  string `json:"cask"`
	Count string `json:"count"`
}

type AnalyticsItem struct {
	Rank    int
	Name    string
	Count   int
	CountStr string
}

func (r *AnalyticsResponse) GetTopFormulae(limit int) []AnalyticsItem {
	items := make([]AnalyticsItem, 0, len(r.Formulae))

	for name, installs := range r.Formulae {
		totalCount := 0
		for _, install := range installs {
			count := parseCount(install.Count)
			totalCount += count
		}
		items = append(items, AnalyticsItem{
			Name:     name,
			Count:    totalCount,
			CountStr: formatCount(totalCount),
		})
	}

	// Sort by count descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	// Add ranks and limit
	if limit > len(items) {
		limit = len(items)
	}
	result := items[:limit]
	for i := range result {
		result[i].Rank = i + 1
	}

	return result
}

func (r *AnalyticsResponse) GetTopCasks(limit int) []AnalyticsItem {
	items := make([]AnalyticsItem, 0, len(r.Casks))

	for name, installs := range r.Casks {
		totalCount := 0
		for _, install := range installs {
			count := parseCount(install.Count)
			totalCount += count
		}
		items = append(items, AnalyticsItem{
			Name:     name,
			Count:    totalCount,
			CountStr: formatCount(totalCount),
		})
	}

	// Sort by count descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	// Add ranks and limit
	if limit > len(items) {
		limit = len(items)
	}
	result := items[:limit]
	for i := range result {
		result[i].Rank = i + 1
	}

	return result
}

func parseCount(s string) int {
	// Remove commas and parse
	s = strings.ReplaceAll(s, ",", "")
	count, _ := strconv.Atoi(s)
	return count
}

func formatCount(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return strconv.Itoa(n)
}

type Period string

const (
	Period30Days  Period = "30d"
	Period90Days  Period = "90d"
	Period365Days Period = "365d"
)

func (c *Client) GetFormulaAnalytics(period Period) (*AnalyticsResponse, error) {
	url := fmt.Sprintf("%s/analytics/install/homebrew-core/%s.json", BaseURL, period)
	data, err := c.getBytes(url)
	if err != nil {
		return nil, err
	}

	var analytics AnalyticsResponse
	if err := json.Unmarshal(data, &analytics); err != nil {
		return nil, err
	}
	return &analytics, nil
}

func (c *Client) GetFormulaAnalyticsBytes(period Period) ([]byte, error) {
	url := fmt.Sprintf("%s/analytics/install/homebrew-core/%s.json", BaseURL, period)
	return c.getBytes(url)
}

func (c *Client) GetCaskAnalytics(period Period) (*AnalyticsResponse, error) {
	url := fmt.Sprintf("%s/analytics/cask-install/homebrew-cask/%s.json", BaseURL, period)
	data, err := c.getBytes(url)
	if err != nil {
		return nil, err
	}

	var analytics AnalyticsResponse
	if err := json.Unmarshal(data, &analytics); err != nil {
		return nil, err
	}
	return &analytics, nil
}

func (c *Client) GetCaskAnalyticsBytes(period Period) ([]byte, error) {
	url := fmt.Sprintf("%s/analytics/cask-install/homebrew-cask/%s.json", BaseURL, period)
	return c.getBytes(url)
}
