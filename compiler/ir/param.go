package ir

import (
	"fmt"
	"go/ast"
	"go/types"
)

// Parameter represents a single variable that can be passed into a block.
type Param struct {
	// Name is the identifier for this parameter.
	//
	// This identifier needs to have a types.Info entry to get the object.
	//
	// Unnamed parameters (parameter lists which only contain types) have
	// a nil name. This is not allowed here since these are only for blocks.
	// The function call to kick off the first block in a function will
	// only pass named parameters into the block. The function will keep
	// the unnamed parameters so that its signature remained the same as
	// it was defined in the AST.
	Name *ast.Ident // TODO: REPLACE

	// Expr is the source-level type expression for this parameter.
	//
	// May be nil for params synthesized by the blocker for variables
	// that were defined inside the function body (no source type expression
	// is available). In that case Type must be set.
	Expr ast.Expr // TODO: REPLACE

	// Type is the resolved type for this parameter.
	//
	// Used as a fallback when Type is nil. The blocker sets this when
	// synthesizing block params for variables defined inside the function
	// body.
	Type types.Type // TODO: REPLACE
}

func (p *Param) String() string {
	if p.Expr != nil {
		return fmt.Sprintf(`%s %s`, p.Name.String(), nodeString(p.Expr))
	}
	return fmt.Sprintf(`%s %v`, p.Name.String(), p.Type)
}

func ExpandParams(ft *ast.FuncType) []*Param {
	if ft == nil {
		return nil
	}
	return expandFieldList(ft.Params)
}

func ExpandResults(ft *ast.FuncType) []*Param {
	if ft == nil {
		return nil
	}
	return expandFieldList(ft.Results)
}

func expandFieldList(fl *ast.FieldList) []*Param {
	if fl == nil || len(fl.List) <= 0 {
		return nil
	}
	params := make([]*Param, 0, len(fl.List))
	for _, field := range fl.List {
		// If field.Names is nil then this is an unnamed field
		// that should be skipped over when defining block params.
		for _, name := range field.Names {
			if name == nil {
				continue
			}
			param := &Param{
				Name: name,
				Expr: field.Type,
			}
			params = append(params, param)
		}
	}
	return params
}
