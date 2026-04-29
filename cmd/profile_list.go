package cmd

import (
	"mws/internal/service"
	"fmt"

	"github.com/spf13/cobra"
)

// NewProfileListCmd builds the command that lists available profiles.
func NewProfileListCmd(pm service.ProfileManager) *cobra.Command {

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			profiles, err := pm.List()

			if err != nil {
				return fmt.Errorf("list profiles: %w", err)
			}

			cmd.Printf("Profiles:\n")
			for _, profile := range profiles {
				cmd.Printf("  %s\n", profile)
			}
			return nil
		},
	}

	return listCmd
}
