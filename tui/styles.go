package tui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	HeaderTitle    lipgloss.Style
	SearchPrompt   lipgloss.Style
	SearchText     lipgloss.Style
	BadgeBox       lipgloss.Style
	BadgeBash      lipgloss.Style
	BadgeZsh       lipgloss.Style
	BadgeFish      lipgloss.Style
	BadgeAll       lipgloss.Style
	CounterText    lipgloss.Style

	// Top Pane (Search List)
	TableBorder    lipgloss.Style
	SelectedItem   lipgloss.Style
	NormalItem     lipgloss.Style
	MatchHighlight lipgloss.Style
	DimmedText     lipgloss.Style

	// Bottom Pane (Context Window)
	ContextBox    lipgloss.Style
	ContextTitle  lipgloss.Style
	ContextTarget lipgloss.Style
	ContextBefore lipgloss.Style
	ContextAfter  lipgloss.Style
	ContextBadge  lipgloss.Style
	ContextTag    lipgloss.Style

	// Status & Help Bar
	StatusBar lipgloss.Style
	KeyHint   lipgloss.Style
	KeyDesc   lipgloss.Style
}

func DefaultStyles() Styles {
	s := Styles{}

	// Color palette
	cyan := lipgloss.Color("#89B4FA")
	purple := lipgloss.Color("#CBA6F7")
	pink := lipgloss.Color("#F5C2E7")
	green := lipgloss.Color("#A6E3A1")
	yellow := lipgloss.Color("#F9E2AF")
	subtext := lipgloss.Color("#A6ADC8")
	surface0 := lipgloss.Color("#313244")
	overlay0 := lipgloss.Color("#6C7086")

	s.HeaderTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(pink).
		Padding(0, 1)

	s.SearchPrompt = lipgloss.NewStyle().
		Bold(true).
		Foreground(cyan)

	s.SearchText = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CDD6F4"))

	s.BadgeBox = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		MarginRight(1)

	s.BadgeBash = s.BadgeBox.
		Foreground(lipgloss.Color("#11111B")).
		Background(cyan)

	s.BadgeZsh = s.BadgeBox.
		Foreground(lipgloss.Color("#11111B")).
		Background(yellow)

	s.BadgeFish = s.BadgeBox.
		Foreground(lipgloss.Color("#11111B")).
		Background(purple)

	s.BadgeAll = s.BadgeBox.
		Foreground(lipgloss.Color("#11111B")).
		Background(green)

	s.CounterText = lipgloss.NewStyle().
		Foreground(subtext).
		Italic(true)

	s.TableBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(surface0)

	s.SelectedItem = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#11111B")).
		Background(cyan)

	s.NormalItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CDD6F4"))

	s.MatchHighlight = lipgloss.NewStyle().
		Bold(true).
		Foreground(yellow)

	s.DimmedText = lipgloss.NewStyle().
		Foreground(overlay0)

	// Context Window
	s.ContextBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purple).
		Padding(0, 1)

	s.ContextTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(purple)

	s.ContextTarget = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#11111B")).
		Background(pink).
		Padding(0, 1)

	s.ContextBefore = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BAC2DE"))

	s.ContextAfter = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BAC2DE"))

	s.ContextBadge = lipgloss.NewStyle().
		Bold(true).
		Foreground(green)

	s.ContextTag = lipgloss.NewStyle().
		Bold(true).
		Foreground(yellow)

	// Status Bar
	s.StatusBar = lipgloss.NewStyle().
		Foreground(subtext).
		Padding(0, 1)

	s.KeyHint = lipgloss.NewStyle().
		Bold(true).
		Foreground(cyan)

	s.KeyDesc = lipgloss.NewStyle().
		Foreground(subtext)

	return s
}
