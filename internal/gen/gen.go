// Package gen contains the internals of the cgen application.
package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// RootFlags belong to the Root command.
type RootFlags struct {
	OutputFile string
	StubsFile  string
	Package    string
	Import     string
}

// Bind binds f to the cobra.Command.
func (f *RootFlags) Bind(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.OutputFile, "output", "o", "", "file name for generated commands (default is stdout)")
	cmd.Flags().StringVar(&f.StubsFile, "output-stubs", "", "file name for optionally generated stubs")
	cmd.Flags().StringVar(&f.Package, "package", "main", "generated code package")
	cmd.Flags().StringVar(&f.Import, "import", "", "import path for implementation code")
}

// RootRun is the "cgen" command.
func RootRun(_ *cobra.Command, args []string, flags RootFlags) error {
	input := "cobra.yaml"
	if len(args) > 0 {
		input = args[0]
	}

	spec, err := ParseCommand(input)
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

	spec.Configure(flags.Package, flags.Import)

	tmpl, err := NewTemplate()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := generateCode(buf, tmpl, spec); err != nil {
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
	if err := generateStubs(buf, tmpl, spec); err != nil {
		return err
	}

	if err := gofmt(buf.Bytes(), stubs); err != nil {
		fmt.Println(buf.String())
		return err
	}

	return nil
}

func gofmt(buf []byte, w io.Writer) error {
	code, err := format.Source(buf)
	if err != nil {
		return err
	}

	if _, err := w.Write(code); err != nil {
		return err
	}

	return nil
}

func fullname(base, use string) string {
	name := use

	f := strings.Fields(name)
	if len(f) > 0 {
		name = f[0]
	}

	return base + cases.Title(language.English).String(name)
}
