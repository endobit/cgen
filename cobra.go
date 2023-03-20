package main

import "github.com/spf13/cobra"

// NewRootCmd returns the top most *cobra.Command for the cgen application.
func NewRootCmd() *cobra.Command {
	var flags RootFlags

	cmd := &cobra.Command{
		Use:   "cgen spec.yaml",
		Short: "cobra generator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Root(cmd.Context(), args, flags)
		},
	}

	flags.Bind(cmd)

	return cmd
}
