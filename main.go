// Copyright 2023 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"bufio"
	"fmt"
	"os"

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
	err := merge.MergeJSON(r, t.Out(), t.IsColorEnabled())
	if err != nil {
		fmt.Fprintf(t.ErrOut(), "failed to merge JSON: %s\n", err)
		os.Exit(1)
	}
}
