package cmd

import (
	"errors"
	"fmt"
	"strings"

	repository "github.com/arxenn/tasks/internal/repository/sqlite"
	"github.com/arxenn/tasks/internal/service"
	"github.com/spf13/cobra"
)

func init() {
	addCmd.Flags().StringP("priority", "p", "medium", "Task priority (low, medium, high, block)")

	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new task",
	Long: `Add a new task with an optional priority level.

The task description can be provided without quotes.
Priority can be specified using the --priority flag.
Valid priority values: low, medium, high, block (default: medium)`,
	Example: `  task add Complete project report --priority high
  task add low Fix bug #123 -p low`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("task description is required")
		}
		content := strings.Join(args, " ")

		priority, err := cmd.Flags().GetString("priority")
		if err != nil {
			return cmdError(cmd, err, "failed to read priority flag")
		}

		repo, err := repository.NewSQLiteRepository()
		if err != nil {
			return cmdError(cmd, err, "failed to connect to the database")
		}
		defer repo.Close()

		svc := service.NewService(repo)

		id, err := svc.Add(content, priority)
		if err != nil {
			return cmdError(cmd, err, "could not add task")
		}

		fmt.Printf("Task added successfully (ID: %d)\n", id)
		return nil
	},
}
