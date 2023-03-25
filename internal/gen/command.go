package gen

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type meta struct {
	Application   string
	Package       string
	ImportPath    string
	ImportPackage string
}

type command struct {
	Meta              meta
	IsRoot            bool
	Name              string
	Fullname          string
	Use               string    `yaml:"use"`
	Short             string    `yaml:"short"`
	Long              string    `yaml:"long"`
	Example           string    `yaml:"example"`
	Aliases           []string  `yaml:"aliases"`
	Flags             bool      `yaml:"flags"`
	Args              string    `yaml:"args"`
	NoRun             bool      `yaml:"norun"`
	PersistentPreRun  bool      `yaml:"persistent_pre_run"`
	PersistentPostRun bool      `yaml:"persistent_post_run"`
	PreRun            bool      `yaml:"pre_run"`
	PostRun           bool      `yaml:"post_run"`
	Commands          []command `yaml:"commands"`
}

// ParseCommand parses a yaml spec file and return the command.
func ParseCommand(file string) (*command, error) {
	fin, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	var cmd command

	d := yaml.NewDecoder(fin)
	if err := d.Decode(&cmd); err != nil {
		return nil, err
	}

	return &cmd, nil
}

// ForAll executed f on c and all its sub-commands.
func (c *command) ForAll(f func(*command)) {
	f(c)

	for i := range c.Commands {
		c.Commands[i].ForAll(f)
	}
}

// Configure traverses c and populates all the meta data to prepare for template
// execution.
func (c *command) Configure(pkg, imp string) {
	var impbase string

	if imp != "" {
		a := strings.Split(imp, "/")
		impbase = a[len(a)-1]
	}

	c.IsRoot = true
	c.Fullname = "Root"
	c.Name = strings.Fields(c.Use)[0]

	meta := meta{
		Application:   c.Name,
		Package:       pkg,
		ImportPath:    imp,
		ImportPackage: impbase,
	}

	c.ForAll(func(c *command) {
		c.Meta = meta
	})

	for i := range c.Commands {
		c.Commands[i].configure("")
	}
}

func (c *command) configure(basename string) {
	c.Name = strings.Fields(c.Use)[0]
	c.Fullname = fullname(basename, c.Name)

	for i := range c.Commands {
		c.Commands[i].configure(c.Fullname)
	}
}
