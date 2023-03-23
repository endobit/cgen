package main

import (
	"context"

	"github.com/spf13/cobra"
)

type cmdKey struct{}

func storeCmd(ctx context.Context, cmd *cobra.Command) context.Context {
	return context.WithValue(ctx, cmdKey{}, cmd)
}

// Ctx returns the *cobra.Commad associated with the ctx.
func Ctx(ctx context.Context) *cobra.Command {
	if c, ok := ctx.Value(cmdKey{}).(*cobra.Command); ok {
		return c
	}

	return nil
}

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
