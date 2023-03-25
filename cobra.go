// Generated code - Do Not Edit -

package main

import (
	"github.com/endobit/cgen/internal/gen"
	"github.com/spf13/cobra"
)

// NewRootCmd returns the top most *cobra.Command for the cgen application.
func NewRootCmd() *cobra.Command {
	var flags gen.RootFlags

	cmd := &cobra.Command{
		Use:   "cgen spec.yaml",
		Short: "cobra generator",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return gen.RootRun(cmd, args, flags)
		},
	}

	flags.Bind(cmd)

	return cmd
}
