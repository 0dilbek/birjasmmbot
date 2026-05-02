package bot

import (
	"strings"

	"github.com/birjasmm/bot/internal/models"
)

func splitCallback(data string) (action, param string) {
	parts := strings.SplitN(data, ":", 2)
	action = parts[0]
	if len(parts) > 1 {
		param = parts[1]
	}
	return
}

func taskEmoji(s models.TaskStatus) string {
	switch s {
	case models.TaskOpen:
		return "🟢"
	case models.TaskInProgress:
		return "🟡"
	case models.TaskCompleted:
		return "🔵"
	case models.TaskCancelled:
		return "🔴"
	}
	return "⚪"
}

func ifEmpty(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}
