// Copyright 2015 Auburn University. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/godoctor/demo"
	"github.com/godoctor/godoctor/engine"
	"github.com/godoctor/godoctor/engine/cli"
)

func main() {
	engine.AddRefactoring("adderror", new(demo.AddError))
	os.Exit(cli.Run(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
