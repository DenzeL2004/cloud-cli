/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"mws/internal/service"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewRootCmd builds the root Cobra command for the mws CLI.
func NewRootCmd(pm service.ProfileManager) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "mws",
		Short: "CLI for managing local YAML profiles",
		Args: cobra.NoArgs,
	}

	rootCmd.AddCommand(NewProfileCmd(pm))
	configureRootHelp(rootCmd)

	return rootCmd
}

func Execute(pm service.ProfileManager) {
	rootCmd := NewRootCmd(pm)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func configureRootHelp(rootCmd *cobra.Command) {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		w := cmd.OutOrStdout()
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

		fmt.Fprintf(w, "%s\n\n", cmd.Short)
		fmt.Fprintln(w, "Usage:")
		fmt.Fprintln(w, "  mws [command]")
		fmt.Fprintln(w)

		fmt.Fprintln(w, "Commands:")
		fmt.Fprintln(tw, "  help [command]\tHelp about any command")
		printCommandsRecursively(tw, cmd, "  ")
		tw.Flush()
		fmt.Fprintln(w)
		fmt.Fprintln(w, `Use "mws help [command]" for more information about a command.`)
	})
}

func printCommandsRecursively(w io.Writer, cmd *cobra.Command, indent string) {
	for _, child := range cmd.Commands() {
		if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
			continue
		}

		if child.Name() == "completion" || child.Name() == "help" {
			continue
		}

		commandWithFlags := child.Use + formatFlagsInline(child)
		fmt.Fprintf(w, "%s%s\t%s\n", indent, commandWithFlags, child.Short)
		printCommandsRecursively(w, child, indent+"  ")
	}
}

func formatFlagsInline(cmd *cobra.Command) string {
	var flags []string

	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "help" {
			return
		}
		
		flags = append(flags, "--"+flag.Name)
	})

	if len(flags) == 0 {
		return ""
	}

	return " " + strings.Join(flags, " ")
}
