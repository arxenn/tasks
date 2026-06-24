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
	rootCmd.AddCommand(doneCmd)
}

var doneCmd = &cobra.Command{
	Use:   "done [task ID]",
	Short: "Mark a task as completed",
	Long: `Mark a task as completed by setting its status to "done".

The task ID must be a positive integer.`,
	Example: `  task done 42
  task done 7`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly one task ID is required")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return cmdError(cmd, err, "invalid task ID")
		}

		if id <= 0 {
			return errors.New("task ID must be a positive integer")
		}

		repo, err := repository.NewSQLiteRepository()
		if err != nil {
			return cmdError(cmd, err, "failed to connect to the database")
		}
		defer repo.Close()

		svc := service.NewService(repo)

		if err := svc.Done(id); err != nil {
			return cmdError(cmd, err, "could not mark task as done")
		}

		fmt.Printf("Task %d marked as done\n", id)
		return nil
	},
}
