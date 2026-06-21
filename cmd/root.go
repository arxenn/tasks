/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	memrepo "github.com/arxenn/tasks/internal/repository/memory"
	"github.com/arxenn/tasks/internal/service"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Always watch your priorioties",
	Long: `tasks is a simple todo list manager, for managing your daily tasks
	with a priority first aproach.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new task (doesn't need double quotes).",
	RunE: func(cmd *cobra.Command, args []string) error {
		priority, err := cmd.Flags().GetString("priority")
		if err != nil {
			return err
		}
		content := strings.Join(args, " ")

		memRepo := memrepo.NewInMemoryRepository()
		svc := service.NewService(memRepo)

		if _, err := svc.Add(content, priority); err != nil {
			return fmt.Errorf("add task failed: %w", err)
		}

		return nil
	},
}
