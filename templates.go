package main

var cmdTpl = `
{{ define "bodycmd" }}
{{ end }}

{{ define "header" }}
package main

import "github.com/spf13/cobra"

{{ end }}

{{ define "command" }}
{{ if .IsRoot -}}
// NewRootCmd returns the top most *cobra.Command for the {{ .Name }} application.
{{ else -}}
// New{{ .Fullname }}Cmd returns a *cobra.Command for "{{ .Application }} {{ snakecase .Fullname | replace "_" " "}}".
{{ end -}}
func New{{ .Fullname }}Cmd() *cobra.Command {
{{ if .Flags -}}
    var flags {{ .Fullname}}Flags

{{ end -}}
    cmd := &cobra.Command{
        Use: "{{ .Use }}",
{{ if ne .Short "" -}}
        Short: "{{ .Short }}",
{{ end -}}
{{ if ne .Long "" -}}
        Long: "{{ .Long }}",
{{ end -}}
{{ if ne .Args "" -}}
        Args: {{ .Args }},
{{ end -}}

{{ if .PersistentPreRun -}}
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            ctx := cmd.Context()
            if err := RootPersistentPreRun(ctx, args); err != nil { return err }
            if cmd.Context() != ctx { cmd.SetContext(ctx) }
            return nil
        },
{{ end -}}
{{ if .PersistentPostRun -}}
        PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
            ctx := cmd.Context()
            if err := RootPersistentPostRun(ctx, args); err != nil { return err }
            if cmd.Context() != ctx { cmd.SetContext(ctx) }
            return nil
        },
{{ end -}}

{{ if not .NoRun -}}
        RunE: func(cmd *cobra.Command, args []string) error {
{{ if .Flags -}}
            return {{ .Fullname }}(cmd.Context(), args, flags)
{{ else -}}
            return {{ .Fullname }}(cmd.Context(), args)
{{ end -}}
        },

{{ end -}}
    }

{{ if .Flags -}}
    flags.Bind(cmd)
{{ end }}

{{ range .Commands -}}
    cmd.AddCommand(New{{ .Fullname }}Cmd())
{{ end }}

    return cmd
}
{{ end }}

{{ define "stubs" }}
{{ if .Flags -}}
{{ if .IsRoot -}}
// {{ .Fullname }}Flags holds the pflag variables for the {{ .Name }} application.
{{ else -}}
// {{ .Fullname }}Flags holds the pflag variables for "{{ snakecase .Fullname | replace "_" " "}}".
{{ end -}}
type {{ .Fullname }}Flags struct {
}

// Bind attaches f to the cobra.Command.
func (f *{{ .Fullname }}Flags) Bind(cmd *cobra.Command) {

}
{{ end }}
{{ if not .NoRun }}
{{ if .IsRoot -}}
// {{ .Fullname }} implements the cli for the {{ .Name }} application.
{{ else -}}
// {{ .Fullname }} implements the cli "{{ snakecase .Fullname | replace "_" " "}}" command.
{{ end -}}
{{ if .Flags -}}
func {{ .Fullname }}(ctx context.Context, args []string, flags {{ .Fullname}}Flags) error {
{{ else -}}
func {{ .Fullname }}(ctx context.Context, args []string) error {
{{ end -}}
    return nil
}
{{ end }}
{{ end }}
`
