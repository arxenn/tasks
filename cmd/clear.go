package cmd

import (
	"errors"
	"fmt"

	repository "github.com/arxenn/tasks/internal/repository/sqlite"
	"github.com/arxenn/tasks/internal/service"
	"github.com/spf13/cobra"
)

func init() {
	clearCmd.Flags().BoolP("all", "a", false, "clears ALL tasks")

	rootCmd.AddCommand(clearCmd)
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear completed tasks from the list",
	Long: `Remove tasks from your todo list.
By default, this command removes only tasks that are marked as "done".
Use the --all (-a) flag to remove ALL tasks, including todo ones.`,
	Example: `  # Clear all completed tasks
  task clear

  # Clear all tasks (including todos)
  task clear --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("invalid usage")
		}

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			return cmdError(cmd, err, "failed to read all flag")
		}

		repo, err := repository.NewSQLiteRepository()
		if err != nil {
			return cmdError(cmd, err, "failed to connect to the database")
		}
		defer repo.Close()

		svc := service.NewService(repo)

		if err := svc.Clear(all); err != nil {
			return cmdError(cmd, err, "could not clear tasks")
		}

		fmt.Print("Tasks cleared successfully\n")
		return nil
	},
}
