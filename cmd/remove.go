package cmd

import (
	"errors"
	"fmt"
	"strconv"

	repository "github.com/arxenn/tasks/internal/repository/sqlite"
	"github.com/arxenn/tasks/internal/service"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove [task ID]",
	Aliases: []string{"rm", "del", "delete"},
	Short:   "Remove a task",
	Long: `Remove (delete) a task permanently by its ID.

This action cannot be undone. The task ID must be a positive integer.`,
	Example: `  task remove 42
  task rm 7
  task delete 15`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly one task ID is required")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %w", err)
		}

		if id <= 0 {
			return errors.New("task ID must be a positive integer")
		}

		repo, err := repository.NewSQLiteRepository()
		if err != nil {
			return fmt.Errorf("error initializing repository: %w", err)
		}
		defer repo.Close()

		svc := service.NewService(repo)

		if err := svc.Delete(id); err != nil {
			return fmt.Errorf("failed to remove task %d: %w", id, err)
		}

		fmt.Printf("Task %d removed successfully\n", id)
		return nil
	},
}
