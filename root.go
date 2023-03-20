package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type RootFlags struct {
	OutputFile string
	StubsFile  string
}

func (f *RootFlags) Bind(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.OutputFile, "output", "o", "", "file name for generated commands (default is stdout)")
	cmd.Flags().StringVar(&f.StubsFile, "output-stubs", "", "file name for optionally generated stubs")
}

func Root(ctx context.Context, args []string, flags RootFlags) error {
	spec, err := readSpec(args[0])
	if err != nil {
		return err
	}

	var code, stubs io.Writer

	if flags.OutputFile != "" {
		fout, err := os.Create(flags.OutputFile)
		if err != nil {
			return err
		}
		defer fout.Close()
		code = fout
	} else {
		code = os.Stdout
	}

	if flags.StubsFile != "" {
		fout, err := os.Create(flags.StubsFile)
		if err != nil {
			return err
		}
		defer fout.Close()
		stubs = fout
	}

	process(spec)

	buf := new(bytes.Buffer)
	if err := generateCode(buf, spec); err != nil {
		return err
	}

	if err := gofmt(buf.Bytes(), code); err != nil {
		fmt.Println(buf.String())
		return err
	}

	if stubs == nil {
		return nil
	}

	buf = new(bytes.Buffer)
	if err := generateStubs(buf, spec); err != nil {
		return err
	}

	if err := gofmt(buf.Bytes(), stubs); err != nil {
		fmt.Println(buf.String())
		return err
	}

	return nil
}
