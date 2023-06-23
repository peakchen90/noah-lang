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
	case *ast.TNumber:
		kind = &TNumber{Impl: newImpl()}
	case *ast.TByte:
		kind = &TByte{Impl: newImpl()}
	case *ast.TChar:
		kind = &TChar{Impl: newImpl()}
	case *ast.TString:
		kind = &TString{Impl: newImpl()}
	case *ast.TBool:
		kind = &TBool{Impl: newImpl()}
	case *ast.TAny:
		kind = &TAny{}
	case *ast.TSelf:
		kind = m.findIdentifierKind(kindExpr, true)
	case *ast.TArray:
		node := node.(*ast.TArray)
		kind = m.compileArrayKind(node)
	case *ast.TIdentifier:
		kind = m.findIdentifierKind(kindExpr, true)
	case *ast.TMemberKind:
		kind = m.findMemberKind(kindExpr, nil, true)
	case *ast.TFuncKind:
		node := node.(*ast.TFuncKind)
		kind = m.compileFuncKind(node)
	case *ast.TStructKind:
		node := node.(*ast.TStructKind)
		kind = m.compileStructKind(node)
	}

	return kind
}

func (m *Module) compileArrayKind(t *ast.TArray) Kind {
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

func (m *Module) compileFuncKind(t *ast.TFuncKind) Kind {
	rest := false
	arguments := make([]Kind, 0, helper.DefaultCap)

	for i, arg := range t.Arguments {
		if arg.Rest {
			if i == len(t.Arguments)-1 {
				rest = true
			} else {
				// TODO unexpected rest arg
				panic("unexpected rest arg")
			}
		}
		arguments = append(arguments, m.compileKindExpr(arg.Kind))
	}

	return &TFunc{
		Id:           getNextTypeId(),
		Arguments:    arguments,
		Return:       m.compileKindExpr(t.Return),
		RestArgument: rest,
		Impl:         newImpl(),
	}
}

func (m *Module) compileStructKind(t *ast.TStructKind) Kind {
	extends := make([]Kind, 0, helper.SmallCap)
	props := make(map[string]Kind)

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
		Id:         getNextTypeId(),
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
			expr.InferKind = &ast.TNumber{}
		case *ast.BoolLiteral:
			expr.InferKind = &ast.TBool{}
		case *ast.NullLiteral:
			// TODO
		case *ast.StringLiteral:
			expr.InferKind = &ast.TString{}
		case *ast.CharLiteral:
			expr.InferKind = &ast.TChar{}
		}
	}

	return &expr.InferKind
}
