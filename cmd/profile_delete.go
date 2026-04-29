package cmd

import (
	"mws/internal/service"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// NewProfileDeleteCmd builds the command that deletes a profile by name.
func NewProfileDeleteCmd(pm service.ProfileManager) *cobra.Command {
	var name string

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a profile",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := pm.Delete(name); err != nil {
				if errors.Is(err, service.ErrProfileNotFound) {
					cmd.Printf("Profile %q does not exist\n", name)
					return nil
				}

				return fmt.Errorf("delete profile %q: %w", name, err)
			}

			cmd.Printf("Profile %q deleted\n", name)
			return nil
		},
	}

	deleteCmd.Flags().StringVar(&name, "name", "", "Profile name")
	_ = deleteCmd.MarkFlagRequired("name")

	return deleteCmd
}
