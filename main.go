package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"fhist/history"
	"fhist/shell"
	"fhist/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var (
		contextRadius = flag.Int("n", 5, "Number of commands before and after hit in context view")
		contextLong   = flag.Int("context", 5, "Number of commands before and after hit in context view")
		query         = flag.String("q", "", "Initial search query")
		queryLong     = flag.String("query", "", "Initial search query")
		source        = flag.String("s", "auto", "History source (auto, bash, zsh, fish, all)")
		sourceLong    = flag.String("source", "auto", "History source (auto, bash, zsh, fish, all)")
		historyFile   = flag.String("f", "", "Path to custom history file")
		fileLong      = flag.String("file", "", "Path to custom history file")
		initShell     = flag.String("init", "", "Output shell integration script (bash, zsh, fish)")
		printContext  = flag.Bool("print-context", false, "Print selected command along with surrounding context window on exit")
	)

	flag.Parse()

	if *initShell != "" {
		script, err := shell.GetInitScript(*initShell)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(script)
		return
	}

	N := 5
	if *contextRadius != 5 {
		N = *contextRadius
	} else if *contextLong != 5 {
		N = *contextLong
	}

	qStr := *query
	if *queryLong != "" {
		qStr = *queryLong
	}
	if qStr == "" && len(flag.Args()) > 0 {
		qStr = strings.Join(flag.Args(), " ")
	}

	srcStr := *source
	if *sourceLong != "auto" {
		srcStr = *sourceLong
	}

	filePath := *historyFile
	if *fileLong != "" {
		filePath = *fileLong
	}

	store, err := history.LoadAllHistory(filePath, srcStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading history: %v\n", err)
		os.Exit(1)
	}

	if len(store.GetCurrentItems()) == 0 {
		fmt.Fprintf(os.Stderr, "No history commands found.\n")
		os.Exit(1)
	}

	model := tui.NewModel(store, qStr, N)

	// Run bubbletea program rendering TUI onto os.Stderr so os.Stdout remains clean for shell evaluation
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithOutput(os.Stderr))
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	m, ok := finalModel.(tui.Model)
	if !ok {
		os.Exit(1)
	}

	if m.SelectedCommand != "" {
		if *printContext || m.PrintContextOnExit {
			if len(m.FilteredIndexes) > 0 && m.SelectedIndex < len(m.FilteredIndexes) {
				idx := m.FilteredIndexes[m.SelectedIndex]
				targetHist := m.CurrentItems[idx]

				start := max(0, targetHist.Index-m.ContextRadius)
				end := min(len(m.CurrentItems)-1, targetHist.Index+m.ContextRadius)

				fmt.Fprintf(os.Stderr, "=== Context Window (N=%d around #%d) ===\n", m.ContextRadius, targetHist.Index+1)
				for i := start; i <= end; i++ {
					item := m.CurrentItems[i]
					marker := "  "
					if i == targetHist.Index {
						marker = "> "
					}
					fmt.Fprintf(os.Stderr, "%s#%-5d [%s] %s\n", marker, item.Index+1, item.ShortTime(), item.Command)
				}
				fmt.Fprintln(os.Stderr, "======================================")
			}
		}

		// Clean output of ONLY the selected command string to stdout
		fmt.Print(m.SelectedCommand)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return a
	}
	return b
}
