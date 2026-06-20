// Package nonnil provides a go/analysis analyzer that enforces the golang-cop
// Null-Object invariant: a function whose result is a "Null-Object interface"
// (any interface declaring IsNull() bool) must never return a bare nil. Returning
// nil reintroduces exactly the nil-dereference the Null-Object pattern exists to
// prevent — callers expect a real object (e.g. Null.New()/NullError.New()) they
// can always send messages to.
//
// It turns the "always return a Null object, never nil" convention from a matter
// of human discipline into a check that fails CI.
//
// Scope: it flags a bare `return nil` (the realistic mistake). A deliberate typed
// conversion such as `return (Thing)(nil)`, or returning a nil interface variable,
// is left alone — those are explicit and uncommon.
package nonnil

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Doc is the analyzer documentation shown by `go vet`/unitchecker.
const Doc = `report nil returns for Null-Object interfaces (interfaces with IsNull() bool)

A function returning such an interface must hand back a real Null object
(e.g. Null.New()) rather than a bare nil, so callers never dereference an
undefined value. This is the compiler-checkable form of the golang-cop
"never nil" invariant.`

// Analyzer is the nonnil analyzer.
var Analyzer = &analysis.Analyzer{
	Name: "nonnil",
	Doc:  Doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// process walks one function body and reports nil returned in any result
	// position whose declared type is a Null-Object interface. Nested function
	// literals are skipped here because they are visited as their own nodes.
	process := func(body *ast.BlockStmt, sig *types.Signature) {
		if body == nil || sig == nil || sig.Results().Len() == 0 {
			return
		}
		results := sig.Results()
		ast.Inspect(body, func(n ast.Node) bool {
			switch s := n.(type) {
			case *ast.FuncLit:
				return false // its own signature; handled separately
			case *ast.ReturnStmt:
				checkReturn(pass, results, s)
			}
			return true
		})
	}

	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch fn := n.(type) {
			case *ast.FuncDecl:
				if obj, ok := pass.TypesInfo.Defs[fn.Name]; ok {
					if sig, ok := obj.Type().(*types.Signature); ok {
						process(fn.Body, sig)
					}
				}
			case *ast.FuncLit:
				if tv, ok := pass.TypesInfo.Types[fn]; ok {
					if sig, ok := tv.Type.(*types.Signature); ok {
						process(fn.Body, sig)
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func checkReturn(pass *analysis.Pass, results *types.Tuple, ret *ast.ReturnStmt) {
	// Only the 1:1 form (return a, b, c) maps expressions to result positions.
	// A naked return or `return f()` (single call feeding multiple results) is
	// left alone to avoid false positives.
	if len(ret.Results) != results.Len() {
		return
	}
	for i, expr := range ret.Results {
		rt := results.At(i).Type()
		if !isNullObjectInterface(rt) {
			continue
		}
		if tv, ok := pass.TypesInfo.Types[expr]; ok && tv.IsNil() {
			pass.Reportf(expr.Pos(),
				"return a Null object (e.g. Null.New()) instead of nil for %s: it is a Null-Object interface (IsNull() bool), and nil reintroduces the nil-dereference the pattern prevents",
				rt)
		}
	}
}

// isNullObjectInterface reports whether t is an interface type that declares
// IsNull() bool — the marker of the golang-cop Null-Object family.
func isNullObjectInterface(t types.Type) bool {
	if _, ok := t.Underlying().(*types.Interface); !ok {
		return false
	}
	obj, _, _ := types.LookupFieldOrMethod(t, true, nil, "IsNull")
	fn, ok := obj.(*types.Func)
	if !ok {
		return false
	}
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}
	b, ok := sig.Results().At(0).Type().(*types.Basic)
	return ok && b.Kind() == types.Bool
}
