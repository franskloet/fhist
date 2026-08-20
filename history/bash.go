package history

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadBashHistory(path string) ([]HistoryItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var items []HistoryItem
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var currentTs time.Time

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		// Check for timestamp line `#1724128500`
		if strings.HasPrefix(line, "#") && len(line) > 1 {
			if tsVal, err := strconv.ParseInt(line[1:], 10, 64); err == nil && tsVal > 1000000000 {
				currentTs = time.Unix(tsVal, 0)
				continue
			}
		}

		items = append(items, HistoryItem{
			Index:     len(items),
			Command:   line,
			Timestamp: currentTs,
			ExitCode:  -1,
			Source:    "bash",
		})
	}

	for i := range items {
		items[i].Index = i
	}

	return items, scanner.Err()
}
