// Copyright 2023 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/cli/go-gh/pkg/jsonpretty"
	"github.com/cli/go-gh/pkg/term"
	"github.com/heaths/gh-merge-json/internal/merge"
)

func main() {
	t := term.FromEnv()
	if term.IsTerminal(os.Stdin) {
		fmt.Fprintf(t.ErrOut(), "missing input\n")
		os.Exit(1)
	}

	r := bufio.NewReader(t.In())
	b, err := merge.MergeJSON(r)
	if err != nil {
		fmt.Fprintf(t.ErrOut(), "failed to merge JSON: %s\n", err)
		os.Exit(1)
	}

	if t.IsTerminalOutput() {
		buf := bytes.NewBuffer(b)
		err = jsonpretty.Format(t.Out(), buf, "  ", t.IsColorEnabled())
	} else {
		_, err = t.Out().Write(b)
	}
	if err != nil {
		fmt.Fprintf(t.ErrOut(), "failed to format JSON: %s\n", err)
		os.Exit(1)
	}
}
