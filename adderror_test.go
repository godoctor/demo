// Copyright 2015 Auburn University. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package demo_test

import (
	"testing"

	"github.com/godoctor/demo"
	"github.com/godoctor/godoctor/engine"
	"github.com/godoctor/godoctor/refactoring/testutil"
)

const directory = "testdata/"

func TestRefactorings(t *testing.T) {
	engine.AddRefactoring("adderror", new(demo.AddError))
	testutil.TestRefactorings(directory, t)
}
