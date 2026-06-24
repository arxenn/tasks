package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/arxenn/tasks/internal/domain"
	repository "github.com/arxenn/tasks/internal/repository/sqlite"
	"github.com/arxenn/tasks/internal/service"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().StringP("priority", "p", "", "Filter by priority (low, medium, high, block)")
	listCmd.Flags().IntP("number", "n", domain.DefaulListCountNumber, "number of shown tasks (default is 3)")
	listCmd.Flags().BoolP("done", "d", false, "show done tasks")

	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all tasks",
	Long: `List tasks with optional filtering by status and priority.

Tasks are displayed ordered by priority (block > high > medium > low).
Use the -d or --done flag to list completed tasks.
Use the -n or --number flag to limit the number of shown tasks (default is 3).
Use the -p or --priority flag to filter by priority.`,
	Example: `  task list
  task list --done --number 5
  task list --priority high
  task ls -d -p low`,
	RunE: func(cmd *cobra.Command, args []string) error {
		priority, err := cmd.Flags().GetString("priority")
		if err != nil {
			return cmdError(cmd, err, "failed to read priority flag")
		}

		done, err := cmd.Flags().GetBool("done")
		if err != nil {
			return cmdError(cmd, err, "failed to read done flag")
		}

		num, err := cmd.Flags().GetInt("number")
		if err != nil {
			return cmdError(cmd, err, "failed to read number flag")
		}

		repo, err := repository.NewSQLiteRepository()
		if err != nil {
			return cmdError(cmd, err, "failed to connect to the database")
		}
		defer repo.Close()

		svc := service.NewService(repo)

		tasks, err := svc.List(priority, done, num)
		if err != nil {
			return cmdError(cmd, err, "could not list tasks")
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return nil
		}

		printTasks(tasks, done)
		return nil
	},
}

func printTasks(tasks []domain.Task, done bool) {
	f := " [%d] %-*s - %s %s\n"

	var maxContentLen int
	for i := range tasks {
		maxContentLen = max(maxContentLen, len(tasks[i].Content))
	}

	for _, task := range tasks {
		priority := fmt.Sprintf("[%s%s%s]", priorityColor(task.Priority), task.Priority, ColorReset)
		var timeStr string
		if done {
			timeStr = "took " + formatDuration(task.DoneAt.Sub(task.CreatedAt))
		} else {
			timeStr = task.CreatedAt.Format(domain.ListTimeDisplayFormat)
		}

		fmt.Printf(
			f,
			task.ID,
			maxContentLen,
			task.Content,
			timeStr,
			priority,
		)
	}
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	parts := []string{}

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}

// ANSI colors
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

func priorityColor(p domain.TaskPriority) string {
	var color string
	switch p {
	case "block":
		color = ColorRed + ColorBold
	case "high":
		color = ColorRed
	case "medium":
		color = ColorYellow
	case "low":
		color = ColorGreen
	default:
		color = ColorReset
	}

	return color
}
