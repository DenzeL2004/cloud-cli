package cmd

import (
	"mws/internal/models"
	"mws/internal/service"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// NewProfileCreateCmd builds the command that creates a new profile.
func NewProfileCreateCmd(pm service.ProfileManager) *cobra.Command {
	var name string
	var user string
	var project string

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a profile",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			profile := &models.Profile{
				User:    user,
				Project: project,
			}

			if err := pm.Create(name, profile); err != nil {
				if errors.Is(err, service.ErrProfileAlreadyExists) {
					cmd.Printf("Profile %q already exists\n", name)
					return nil
				}

				return fmt.Errorf("create profile %q: %w", name, err)
			}

			cmd.Printf("Profile %q created\n", name)
			return nil
		},
	}

	createCmd.Flags().StringVar(&name, "name", "", "Profile name")
	createCmd.Flags().StringVar(&user, "user", "", "user name")
	createCmd.Flags().StringVar(&project, "project", "", "Project name")

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("user")
	_ = createCmd.MarkFlagRequired("project")

	return createCmd
}
