package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	TasksListShellCommand = "task ls"
)

func init() {
	rootCmd.AddCommand(shellCmd)
}

var shellCmd = &cobra.Command{
	Use:   "shell [enable|disable]",
	Short: "Enable or disable task listing on shell startup",
	Long: `Configure whether tasks are automatically displayed when you open a new shell.

When enabled, your current task list will be shown automatically upon shell startup.
This is useful for keeping your tasks visible and top-of-mind.

The configuration is stored in your shell's configuration file (.bashrc, .zshrc, etc.).`,
	Example: `  # Enable task listing on shell startup
  task shell enable

  # Disable task listing on shell startup
  task shell disable`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly one argument is required (enable or disable)")
		}

		operation := args[0]
		if operation != "enable" && operation != "disable" {
			return errors.New("invalid argument: must be \"enable\" or \"disable\"")
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get user home dir failed: %w", err)
		}
		configFile, err := getShellConfigFileName()
		if err != nil {
			return fmt.Errorf("failed to determine shell config file: %w", err)
		}
		shdir := filepath.Join(home, configFile)

		switch operation {
		case "enable":
			if err := enable(shdir); err != nil {
				return fmt.Errorf("failed to enable shell integration: %w", err)
			}
			fmt.Println("✅ Shell integration enabled. Tasks will be shown on shell startup.")
		case "disable":
			if err := disable(shdir); err != nil {
				return fmt.Errorf("failed to disable shell integration: %w", err)
			}
			fmt.Println("✅ Shell integration disabled. Tasks will not be shown on shell startup.")
		}

		return nil
	},
}

func enable(shdir string) error {
	exists, err := checkConfigExists(shdir)
	if err != nil {
		return fmt.Errorf("check config exists failed: %w", err)
	}
	if exists {
		return fmt.Errorf("task already enabled in your shell: %s", shdir)
	}

	if err := addCommandToShellConfig(shdir); err != nil {
		return fmt.Errorf("add command to configs failed: %w", err)
	}

	return nil
}

func disable(shdir string) error {
	exists, err := checkConfigExists(shdir)
	if err != nil {
		return fmt.Errorf("check config exists failed: %w", err)
	}
	if !exists {
		return fmt.Errorf("task already disabled in your shell: %s", shdir)
	}

	if err := removeCommandFromShellConfig(shdir); err != nil {
		return fmt.Errorf("remove command from configs failed: %w", err)
	}

	return nil
}

func checkConfigExists(shdir string) (bool, error) {
	reg := fmt.Sprintf("^%s$", TasksListShellCommand)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("grep", reg, shdir)
	case "windows":
		cmd = exec.Command("sls", reg, shdir)
	default:
		return false, fmt.Errorf("not supported OS: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil
			} else if exitErr.ExitCode() == 2 {
				return false, fmt.Errorf("check config exists failed: %w\n", err)
			}
		}
		return false, fmt.Errorf("failed to run check command: %w", err)
	}
	return true, nil

}

func addCommandToShellConfig(shdir string) error {
	f, err := os.OpenFile(shdir, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("open config file at %s failed: %w", shdir, err)
	}
	defer f.Close()

	if _, err := f.WriteString(TasksListShellCommand + "\n"); err != nil {
		return fmt.Errorf("append command to config file failed: %w", err)
	}

	return nil
}

func removeCommandFromShellConfig(shdir string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("sed", "-i", fmt.Sprintf("/^%s$/d", TasksListShellCommand), shdir)
	case "windows":
		return errors.New("not implemented")
	default:
		return fmt.Errorf("not supported OS: %s", runtime.GOOS)
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run remove command: %w", err)
	}

	return nil
}

func getShellConfigFileName() (string, error) {
	sh := os.Getenv("SHELL")
	switch runtime.GOOS {
	case "darwin", "linux":
		switch filepath.Base(sh) {
		case "zsh":
			return ".zshrc", nil
		case "bash":
			return ".bashrc", nil
		case "sh", "dash", "ash":
			return ".profile", nil
		}

	case "windows":
		return filepath.Join(
			"Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1",
		), nil
	}

	return "", fmt.Errorf("could not determine shell config file (SHELL=%q)", sh)
}
