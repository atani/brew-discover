package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12"))

	RankStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Bold(true)

	NameStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))

	CountStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("14"))

	DescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13"))

	TipStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true)
)

type TableRow struct {
	Rank        int
	Name        string
	Count       string
	Description string
}

func PrintTable(title string, rows []TableRow, tip string) {
	fmt.Println()
	fmt.Println(TitleStyle.Render(title))
	fmt.Println()

	// Header
	header := fmt.Sprintf(" %-4s %-24s %-12s %s",
		HeaderStyle.Render("#"),
		HeaderStyle.Render("Name"),
		HeaderStyle.Render("Installs"),
		HeaderStyle.Render("Description"))
	fmt.Println(header)
	fmt.Println(strings.Repeat("─", 80))

	// Rows
	for _, row := range rows {
		desc := row.Description
		if len(desc) > 35 {
			desc = desc[:32] + "..."
		}

		line := fmt.Sprintf(" %-4s %-24s %-12s %s",
			RankStyle.Render(fmt.Sprintf("%d", row.Rank)),
			NameStyle.Render(truncate(row.Name, 22)),
			CountStyle.Render(row.Count),
			DescStyle.Render(desc))
		fmt.Println(line)
	}

	fmt.Println()
	if tip != "" {
		fmt.Println(TipStyle.Render(tip))
	}
	fmt.Println()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func FormatCount(count string) string {
	// Parse count string like "1,234,567" and format it
	return count
}
