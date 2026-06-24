package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// cmdError checks the verbose flag and returns either the full wrapped error
// (verbose) or a simple user-facing message.
func cmdError(cmd *cobra.Command, err error, fallback string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		return err
	}
	return errors.New(fallback)
}
