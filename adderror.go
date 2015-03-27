// Copyright 2015 Auburn University. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package demo

import (
	"fmt"
	"strings"

	"go/ast"
	"go/token"

	"github.com/godoctor/godoctor/analysis/cfg"
	"github.com/godoctor/godoctor/refactoring"
	"github.com/godoctor/godoctor/text"
)

type AddError struct {
	base refactoring.RefactoringBase
}

func (r *AddError) Description() *refactoring.Description {
	return &refactoring.Description{
		Name:      "Add Error Return",
		Synopsis:  "Adds an error to a function's return tuple",
		Usage:     "",
		Multifile: false,
		Params:    []refactoring.Parameter{},
		Hidden:    false,
	}
}

func (r *AddError) Run(config *refactoring.Config) *refactoring.Result {
	r.base.Run(config)
	if !refactoring.ValidateArgs(config, r.Description(), r.base.Log) ||
		r.base.Log.ContainsErrors() {
		return &r.base.Result
	}

	f := r.findEnclosingFuncDecl()
	if f == nil || f.Body == nil {
		r.base.Log.Error("Please select a function definition.")
		r.base.Log.AssociatePos(r.base.SelectionStart, r.base.SelectionEnd)
		return &r.base.Result
	}

	r.addErrorReturn(f)
	r.editReturnStmts(f)
	r.addMissingReturnStmts(f)
	// FIXME does not handle anonymous functions correctly at all
	// FIXME update call sites
	r.base.FormatFileInEditor()
	return &r.base.Result
}

func (r *AddError) findEnclosingFuncDecl() *ast.FuncDecl {
	for _, node := range r.base.PathEnclosingSelection {
		if f, ok := node.(*ast.FuncDecl); ok {
			return f
		}
	}
	return nil
}

func (r *AddError) addErrorReturn(f *ast.FuncDecl) {
	var replaceRange *text.Extent
	var results string
	if f.Type.Results != nil {
		// Function has at least one existing result; replace it
		replaceRange = r.base.Extent(f.Type.Results)
		if f.Type.Results.List != nil {
			startPos := f.Type.Results.List[0].Pos()
			endPos := f.Type.Results.List[len(f.Type.Results.List)-1].End()
			results = r.base.TextFromPosRange(startPos, endPos)
		}
	} else {
		// Function has no results; insert after type
		replaceRange = &text.Extent{r.base.OffsetOfPos(f.Type.End()), 0}
	}

	results = strings.TrimSpace(results)
	if results == "" {
		results = "error"
	} else {
		results = fmt.Sprintf("(%s, error)", results)
	}
	r.base.Edits[r.base.Filename].Add(replaceRange, results)
}

func (r *AddError) editReturnStmts(f *ast.FuncDecl) {
	ast.Inspect(r.base.File, func(node ast.Node) bool {
		if ret, ok := node.(*ast.ReturnStmt); ok {
			insert := &text.Extent{r.base.OffsetOfPos(ret.End()), 0}
			r.base.Edits[r.base.Filename].Add(insert, ", nil")
		}
		return true
	})
}

func (r *AddError) addMissingReturnStmts(f *ast.FuncDecl) {
	cfg := cfg.FromFunc(f)
	exitPoints := cfg.Preds(cfg.Exit)
	insertPos := []token.Pos{}
	for _, node := range exitPoints {
		if node == cfg.Entry {
			insertPos = append(insertPos, f.Body.Lbrace+1)
		} else if _, isReturn := node.(*ast.ReturnStmt); !isReturn {
			insertPos = append(insertPos, node.End())
		}
	}

	for _, pos := range insertPos {
		insert := &text.Extent{r.base.OffsetOfPos(pos), 0}
		r.base.Edits[r.base.Filename].Add(insert, "\n\treturn nil")
	}
}
