package gen

import (
	"embed"
	"io"
	"io/fs"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

//go:embed templates
var templates embed.FS

// NewTemplate parses and initializes the embedded templates.
func NewTemplate() (*template.Template, error) {
	disk, err := fs.Sub(templates, "templates")
	if err != nil {
		return nil, err
	}

	tmpl := template.New("cgen").Funcs(sprig.FuncMap())

	return tmpl.ParseFS(disk, "*.tmpl")
}

func generateCode(w io.Writer, tmpl *template.Template, cmd *Command) error {
	if err := tmpl.ExecuteTemplate(w, "header", cmd); err != nil {
		return err
	}

	if err := genCommand(w, tmpl, *cmd); err != nil {
		return err
	}

	for _, c := range cmd.Commands {
		if err := genCommand(w, tmpl, c); err != nil {
			return err
		}
	}

	return nil
}

func genCommand(w io.Writer, tmpl *template.Template, cmd Command) error {
	if err := tmpl.ExecuteTemplate(w, "command", cmd); err != nil {
		return err
	}

	for _, c := range cmd.Commands {
		if err := genCommand(w, tmpl, c); err != nil {
			return err
		}
	}

	return nil
}

func generateStubs(w io.Writer, tmpl *template.Template, cmd *Command) error {
	if err := tmpl.ExecuteTemplate(w, "header", cmd); err != nil {
		return err
	}

	if err := genStub(w, tmpl, *cmd); err != nil {
		return err
	}

	for _, c := range cmd.Commands {
		if err := genStub(w, tmpl, c); err != nil {
			return err
		}
	}

	return nil
}

func genStub(w io.Writer, tmpl *template.Template, cmd Command) error {
	return tmpl.ExecuteTemplate(w, "stubs", cmd)
}
