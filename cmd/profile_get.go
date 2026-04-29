package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"mws/internal/service"

	"github.com/spf13/cobra"
)

// NewProfileGetCmd builds the command that reads a profile by name.
func NewProfileGetCmd(pm service.ProfileManager) *cobra.Command {
	var name string

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a profile",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			profile, err := pm.Get(name)

			if err != nil {
				if errors.Is(err, service.ErrProfileNotFound) {
					cmd.Printf("Profile %q does not exist\n", name)
					return nil
				}

				return fmt.Errorf("get profile %q: %w", name, err)
			}

			cmd.Printf("Profile %s:\n", name)
			data, _ := json.MarshalIndent(profile, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}

	getCmd.Flags().StringVar(&name, "name", "", "Profile name")
	_ = getCmd.MarkFlagRequired("name")

	return getCmd
}
