package main

import (
	"go/format"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type cmd struct {
	IsRoot            bool
	Application       string
	Name              string
	Fullname          string
	Use               string   `yaml:"use"`
	Short             string   `yaml:"short"`
	Long              string   `yaml:"long"`
	Example           string   `yaml:"example"`
	Aliases           []string `yaml:"aliases"`
	Flags             bool     `yaml:"flags"`
	Args              string   `yaml:"args"`
	NoRun             bool     `yaml:"norun"`
	PersistentPreRun  bool     `yaml:"persistent_pre_run"`
	PersistentPostRun bool     `yaml:"persistent_post_run"`
	PreRun            bool     `yaml:"pre_run"`
	PostRun           bool     `yaml:"post_run"`
	Commands          []cmd    `yaml:"commands"`
}

var version string

func main() {
	root := NewRootCmd()
	root.Version = version

	if err := root.Execute(); err != nil {
		os.Exit(-1)
	}
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

func process(cmd *cmd) {
	cmd.IsRoot = true
	cmd.Name = strings.Fields(cmd.Use)[0]
	cmd.Application = cmd.Name
	cmd.Fullname = "Root"

	for i := range cmd.Commands {
		processCommand(cmd.Application, "", &cmd.Commands[i])
	}
}

func processCommand(app, basename string, cmd *cmd) {
	cmd.Application = app
	cmd.Name = strings.Fields(cmd.Use)[0]
	cmd.Fullname = fullname(basename, cmd.Name)

	for i := range cmd.Commands {
		processCommand(app, cmd.Fullname, &cmd.Commands[i])
	}
}

func generateStubs(w io.Writer, cmd *cmd) error {
	tmpl, err := template.New("cli").Funcs(sprig.TxtFuncMap()).Parse(cmdTpl)
	if err != nil {
		return err
	}

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

func generateCode(w io.Writer, cmd *cmd) error {
	tmpl, err := template.New("cli").Funcs(sprig.TxtFuncMap()).Parse(cmdTpl)
	if err != nil {
		return err
	}

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

func genStub(w io.Writer, tmpl *template.Template, cmd cmd) error {
	return tmpl.ExecuteTemplate(w, "stubs", cmd)
}

func genCommand(w io.Writer, tmpl *template.Template, cmd cmd) error {
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

func readSpec(file string) (*cmd, error) {
	fin, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	var cmd cmd

	d := yaml.NewDecoder(fin)
	if err := d.Decode(&cmd); err != nil {
		return nil, err
	}

	return &cmd, nil
}
