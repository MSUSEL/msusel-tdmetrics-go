package usages

import (
	"go/ast"
	"go/token"
	"go/types"
	"iter"

	"github.com/Snow-Gremlin/goToolbox/collections/stack"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type posReader interface {
	Pos() token.Pos
}

func isLocal(root ast.Node, query posReader) bool {
	pos := query.Pos()
	return root.Pos() <= pos && pos <= root.End()
}

func unparen(node ast.Node) ast.Node {
	if p, ok := node.(*ast.ParenExpr); ok {
		return p.X
	}
	return node
}

func getPkgPath(o types.Object) string {
	if o.Pkg() != nil {
		return o.Pkg().Path()
	}
	return ``
}

func getInstTypes(o types.Object) *types.TypeList {
	type TA interface{ TypeArgs() *types.TypeList }
	if t, ok := o.Type().(TA); ok {
		return t.TypeArgs()
	}
	return nil
}

func getName(fSet *token.FileSet, expr ast.Expr) string {
	exp := unparen(expr)
	if id, ok := exp.(*ast.Ident); ok {
		return id.Name
	}
	if sel, ok := exp.(*ast.SelectorExpr); ok {
		src := unparen(sel.X)
		if id, ok := src.(*ast.Ident); ok {
			return id.Name + `.` + sel.Sel.Name
		}
		panic(terror.New(`unexpected expression in selection for name`).
			WithType(`type`, src).
			With(`expression`, src).
			With(`selection`, sel).
			With(`position`, fSet.Position(expr.Pos())))
	}
	panic(terror.New(`unexpected expression for name`).
		WithType(`type`, exp).
		With(`expression`, exp).
		With(`position`, fSet.Position(expr.Pos())))
}

func getNamed(t types.Type) *types.Named {
	for named := range where[*types.Named](walkType(t)) {
		// return first type hit.
		return named
	}
	return nil
}

func where[TOut, TIn any](it iter.Seq[TIn]) iter.Seq[TOut] {
	return func(yield func(TOut) bool) {
		for v := range it {
			if out, ok := any(v).(TOut); ok {
				if !yield(out) {
					return
				}
			}
		}
	}
}

// walkType walks the tree of types. This will only return unique
// types and skip any types already outputted. This will output a type
// followed walking the children. Siblings in a type are output in reverse
// order, such that `map[T]S` will output `S` then `T`.
func walkType(start types.Type) iter.Seq[types.Type] {
	return func(yield func(types.Type) bool) {
		s := stack.With(start)
		touched := map[types.Type]struct{}{}
		for !s.Empty() {
			cur := s.Pop()
			if utils.IsNil(cur) {
				continue
			}
			if !yield(cur) {
				return
			}
			if _, has := touched[cur]; has {
				continue
			}
			touched[cur] = struct{}{}
			switch t := cur.(type) {
			case *types.Alias:
				s.Push(t.Rhs())
			case *types.Array:
				s.Push(t.Elem())
			case *types.Basic:
				// Do Nothing
			case *types.Chan:
				s.Push(t.Elem())
			case *types.Interface:
				for i := range t.NumEmbeddeds() {
					s.Push(t.EmbeddedType(i))
				}
				for i := range t.NumExplicitMethods() {
					s.Push(t.ExplicitMethod(i).Type())
				}
			case *types.Map:
				s.Push(t.Key(), t.Elem())
			case *types.Named:
				if tp := t.TypeParams(); tp != nil {
					for i := range tp.Len() {
						s.Push(tp.At(i))
					}
				}
				if ta := t.TypeArgs(); ta != nil {
					for i := range ta.Len() {
						s.Push(ta.At(i))
					}
				}
				s.Push(t.Underlying())
				for i := range t.NumMethods() {
					s.Push(t.Method(i).Type())
				}
			case *types.Pointer:
				s.Push(t.Elem())
			case *types.Signature:
				if tp := t.TypeParams(); tp != nil {
					for i := range tp.Len() {
						s.Push(tp.At(i))
					}
				}
				s.Push(t.Params(), t.Results())
			case *types.Slice:
				s.Push(t.Elem())
			case *types.Struct:
				for i := range t.NumFields() {
					s.Push(t.Field(i).Type())
				}
			case *types.Tuple:
				for i := range t.Len() {
					s.Push(t.At(i).Type())
				}
			case *types.TypeParam:
				s.Push(t.Constraint())
			case *types.Union:
				for i := range t.Len() {
					s.Push(t.Term(i).Type())
				}
			default:
				panic(terror.New(`encountered unhandled type during walk`).
					WithType(`type`, t).
					With(`value`, t))
			}
		}
	}
}
