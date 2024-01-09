// Copyright 2023 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/jsonpretty"
	ghtemplate "github.com/cli/go-gh/v2/pkg/template"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/heaths/gh-merge-json/internal/merge"
	"github.com/spf13/cobra"
)

var (
	template string
	tee      string

	cmd = &cobra.Command{
		Use:   "gh merge-json",
		Short: "Merges multiple pages of JSON",
		Long:  "Merges multiple pages of JSON since older versions of `gh api --paginate` do not emit proper JSON with multiple pages. This extension merges them together for both REST and GraphQL responses.",
		RunE:  run,
	}
)

func init() {
	cmd.Flags().StringVarP(&template, "template", "t", "", `Format JSON output using a Go template; see "gh help formatting"`)
	cmd.Flags().StringVar(&tee, "tee", "", "Write JSON to `file` as well formatted JSON to standard output")
}

func main() {
	cmd.Execute()
}

func run(_ *cobra.Command, _ []string) error {
	t := term.FromEnv()
	if term.IsTerminal(os.Stdin) {
		return errors.New("missing input")
	}

	r := bufio.NewReader(t.In())
	b, err := merge.MergeJSON(r)
	if err != nil {
		return fmt.Errorf("merge JSON: %s", err)
	}

	if tee != "" {
		w, err := os.Create(tee)
		if err != nil {
			return err
		}

		defer w.Close()

		if _, err = w.Write(b); err != nil {
			return err
		}
	}

	if template != "" {
		width, _, _ := t.Size()
		templ := ghtemplate.New(t.Out(), width, t.IsColorEnabled())
		if err = templ.Parse(template); err != nil {
			return err
		}

		buf := bytes.NewBuffer(b)
		if err = templ.Execute(buf); err != nil {
			return err
		}

		return templ.Flush()
	}

	if t.IsTerminalOutput() {
		buf := bytes.NewBuffer(b)
		err = jsonpretty.Format(t.Out(), buf, "  ", t.IsColorEnabled())
	} else {
		_, err = t.Out().Write(b)
	}
	if err != nil {
		return fmt.Errorf("format JSON: %s", err)
	}

	return nil
}
