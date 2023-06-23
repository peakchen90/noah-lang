package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"math"
)

func (m *Module) compileKindExpr(kindExpr *ast.KindExpr) Kind {
	var kind Kind
	if kindExpr == nil {
		return kind
	}

	node := kindExpr.Node

	switch node.(type) {
	case *ast.TypeNumber:
		kind = &TNumber{Impl: newImpl()}
	case *ast.TypeByte:
		kind = &TByte{Impl: newImpl()}
	case *ast.TypeChar:
		kind = &TChar{Impl: newImpl()}
	case *ast.TypeString:
		kind = &TString{Impl: newImpl()}
	case *ast.TypeBool:
		kind = &TBool{Impl: newImpl()}
	case *ast.TypeAny:
		kind = &TAny{Impl: newImpl()}
	case *ast.TypeArray:
		node := node.(*ast.TypeArray)
		kind = m.compileArrayKind(node)
	case *ast.TypeIdentifier:
		node := node.(*ast.TypeIdentifier)
		kind = m.ScopeStack.findKind(node.Name.Name)
	case *ast.TypeMemberKind:
		node := node.(*ast.TypeMemberKind)
		kind = m.ScopeStack.findMemberKind(node.ToMemberIds())
	case *ast.TypeFuncKind:
		node := node.(*ast.TypeFuncKind)
		kind = m.compileFuncKind(node, false)
	case *ast.TypeStructKind:
		node := node.(*ast.TypeStructKind)
		kind = m.compileStructKind(node, false)
	}

	return kind
}

func (m *Module) compileArrayKind(t *ast.TypeArray) Kind {
	size := -1 // vector array

	if t.Len != nil {
		rawVal := t.Len.Node.(*ast.NumberLiteral).Value
		if rawVal < 0 || math.Floor(rawVal) != rawVal {
			// TODO unexpected len
			panic("unexpected len")
		}
		size = int(rawVal)
	}

	return &TArray{
		Kind: m.compileKindExpr(t.Kind),
		Len:  size,
		Impl: newImpl(),
	}
}

func (m *Module) compileFuncKind(t *ast.TypeFuncKind, isDecl bool) Kind {
	id := -1
	rest := false
	size := len(t.Arguments)
	arguments := make([]Kind, 0, helper.DefaultCap)

	if isDecl {
		id = getNextTypeId()
	}

	for i, arg := range t.Arguments {
		if arg.Rest {
			if i == size-1 {
				rest = true
			} else {
				// TODO unexpected rest arg
				panic("unexpected rest arg")
			}
		}
		arguments = append(arguments, m.compileKindExpr(arg.Kind))
	}

	return &TFunc{
		Id:           id,
		Arguments:    arguments,
		Return:       m.compileKindExpr(t.Return),
		RestArgument: rest,
		Impl:         newImpl(),
	}
}

func (m *Module) compileStructKind(t *ast.TypeStructKind, isDecl bool) Kind {
	id := -1
	extends := make([]Kind, 0, helper.SmallCap)
	props := make(map[string]Kind)

	if isDecl {
		id = getNextTypeId()
	}

	for _, item := range t.Properties {
		key := item.Key.Name
		_, has := props[key]
		if has {
			// TODO duplicate
			panic("duplicate " + key)
		}
		props[key] = m.compileKindExpr(item.Kind)
	}

	for _, item := range t.Extends {
		extends = append(extends, m.compileKindExpr(item))
	}

	return &TStruct{
		Id:         id,
		Extends:    extends,
		Properties: props,
		Impl:       newImpl(),
	}
}

// TODO
func (c *Compiler) inferKind(expr *ast.Expr) *ast.KE {
	if expr.InferKind == nil {
		switch expr.Node.(type) {
		case *ast.CallExpr:
			// TODO
		case *ast.MemberExpr:
			// TODO
		case *ast.BinaryExpr:
			// TODO
		case *ast.BinaryTypeExpr:
			// TODO
		case *ast.UnaryExpr:
			// TODO
		case *ast.FuncExpr:
			// TODO
		case *ast.StructExpr:
			// TODO
		case *ast.ArrayExpr:
			// TODO
		case *ast.IdentifierLiteral:
			// TODO
		case *ast.NumberLiteral:
			expr.InferKind = &ast.TypeNumber{}
		case *ast.BoolLiteral:
			expr.InferKind = &ast.TypeBool{}
		case *ast.NullLiteral:
			// TODO
		case *ast.SelfLiteral:
			// TODO
		case *ast.StringLiteral:
			expr.InferKind = &ast.TypeString{}
		case *ast.CharLiteral:
			expr.InferKind = &ast.TypeChar{}
		}
	}

	return &expr.InferKind
}
