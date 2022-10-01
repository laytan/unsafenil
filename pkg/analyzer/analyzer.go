package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "unafenil"
	doc  = "Checks that there is no return of a nil error or false, and a nil/default value before it"

	reportMsg = "returns both a nil error or false, and a nil/default value before it"
)

// New returns new nilnil analyzer.
func New() *analysis.Analyzer {
	n := newNilNil()

	a := &analysis.Analyzer{
		Name:     name,
		Doc:      doc,
		Run:      n.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	a.Flags.Var(&n.checkedTypes, "checked-types", "coma separated list")

	return a
}

type nilNil struct {
	checkedTypes checkedTypes
}

func newNilNil() *nilNil {
	return &nilNil{
		checkedTypes: newDefaultCheckedTypes(),
	}
}

var (
	types = []ast.Node{(*ast.TypeSpec)(nil)}

	funcAndReturns = []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
		(*ast.ReturnStmt)(nil),
	}
)

type typeSpecByName map[string]*ast.TypeSpec

func (n *nilNil) run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	typeSpecs := typeSpecByName{}
	insp.Preorder(types, func(node ast.Node) {
		t := node.(*ast.TypeSpec)
		typeSpecs[t.Name.Name] = t
	})

	var fs funcTypeStack
	insp.Nodes(funcAndReturns, func(node ast.Node, push bool) (proceed bool) {
		switch v := node.(type) {
		case *ast.FuncLit:
			if push {
				fs.Push(v.Type)
			} else {
				fs.Pop()
			}

		case *ast.FuncDecl:
			if push {
				fs.Push(v.Type)
			} else {
				fs.Pop()
			}

		case *ast.ReturnStmt:
			ft := fs.Top() // Current function.

			if !push || ft == nil || ft.Results == nil {
				return
			}

			if len(v.Results) < 2 || len(ft.Results.List) < 2 {
				return
			}

			fResLast := ft.Results.List[len(ft.Results.List)-1]
			if !n.isErrorField(fResLast) {
				return
			}

			rResLast := v.Results[len(v.Results)-1]
			if !isNil(rResLast) {
				return
			}

			for i, res := range ft.Results.List[0 : len(ft.Results.List)-1] {
				if n.isDangerNilField(res, typeSpecs) {
					if isNil(v.Results[i]) {
						pass.Reportf(v.Pos(), reportMsg)
					}
				}
			}
		}

		return true
	})

	return nil, nil //nolint:nilnil
}

func (n *nilNil) isDangerNilField(f *ast.Field, typeSpecs typeSpecByName) bool {
	return n.isDangerNilType(f.Type, typeSpecs)
}

func (n *nilNil) isDangerNilType(t ast.Expr, typeSpecs typeSpecByName) bool {
	switch v := t.(type) {
	case *ast.StarExpr:
		return n.checkedTypes.Contains(ptrType)

	case *ast.FuncType:
		return n.checkedTypes.Contains(funcType)

	case *ast.InterfaceType:
		return n.checkedTypes.Contains(ifaceType)

	case *ast.MapType:
		return n.checkedTypes.Contains(mapType)

	case *ast.ChanType:
		return n.checkedTypes.Contains(chanType)

	case *ast.Ident:
		if t, ok := typeSpecs[v.Name]; ok {
			return n.isDangerNilType(t.Type, nil)
		}
	}
	return false
}

func (n *nilNil) isErrorField(f *ast.Field) bool {
	return isIdent(f.Type, "error") || isIdent(f.Type, "bool")
}

func isNil(e ast.Expr) bool {
	return isIdent(e, "nil") || isIdent(e, "true")
}

func isIdent(n ast.Node, name string) bool {
	i, ok := n.(*ast.Ident)
	if !ok {
		return false
	}
	return i.Name == name
}
