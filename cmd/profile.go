package cmd

import (
	"mws/internal/service"
	"github.com/spf13/cobra"
)

// NewProfileCmd builds the parent command for profile operations.
func NewProfileCmd(pm service.ProfileManager) *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles",
	}

	profileCmd.AddCommand(NewProfileCreateCmd(pm))
	profileCmd.AddCommand(NewProfileGetCmd(pm))
	profileCmd.AddCommand(NewProfileDeleteCmd(pm))
	profileCmd.AddCommand(NewProfileListCmd(pm))

	return profileCmd
}