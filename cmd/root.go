package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task [command] [flags]",
	Short: "A priority-first task manager for your daily todos",
	Long: `task is a simple todo list manager that helps you stay organized
with a priority-first approach. Tasks are automatically ordered by priority
(block > high > medium > low) to ensure you always focus on what matters most.
All data is persisted in a local SQLite database.`,
	Example: `  # Add a new task
  task add Review pull requests -p high

  # List all pending tasks
  task list -s pending

  # Mark a task as done
  task done 42

  # Remove a task
  task remove 15

 # Get help for a specific command
  task help add`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
