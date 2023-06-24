package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"math"
)

func (m *Module) compileKindExpr(kindExpr *ast.KindExpr) *KindRef {
	kind := &KindRef{}
	if kindExpr == nil {
		return kind
	}

	node := kindExpr.Node

	switch node.(type) {
	case *ast.TNumber:
		kind.Ref = TNumberConst
	case *ast.TByte:
		kind.Ref = TByteConst
	case *ast.TChar:
		kind.Ref = TCharConst
	case *ast.TString:
		kind.Ref = TStringConst
	case *ast.TBool:
		kind.Ref = TBoolConst
	case *ast.TAny:
		kind.Ref = TAnyConst
	case *ast.TSelf:
		kind.Ref = &TSelf{KindRef: m.scopes.findSelfKind(kindExpr, true)}
	case *ast.TArray:
		node := node.(*ast.TArray)
		return m.compileArrayKind(node)
	case *ast.TIdentifier:
		node := node.(*ast.TIdentifier)
		return m.scopes.findIdentifierKind(node.Name, true)
	case *ast.TMemberKind:
		return m.scopes.findMemberKind(kindExpr, true)
	case *ast.TFuncKind:
		node := node.(*ast.TFuncKind)
		return m.compileFuncKind(node)
	case *ast.TStructKind:
		node := node.(*ast.TStructKind)
		return m.compileStructKind(node)
	}

	return kind
}

func (m *Module) compileArrayKind(t *ast.TArray) *KindRef {
	kind := &KindRef{}
	size := -1 // vector array

	if t.Len != nil {
		rawVal := t.Len.Node.(*ast.NumberLiteral).Value
		if rawVal < 0 || math.Floor(rawVal) != rawVal {
			m.unexpectedPos(t.Len.Start, "Expect be a positive integer")
		}
		size = int(rawVal)
	}

	kind.Ref = &TArray{
		KindRef: m.compileKindExpr(t.Kind),
		Len:     size,
		Impl:    newImpl(),
	}
	return kind
}

func (m *Module) compileFuncKind(t *ast.TFuncKind) *KindRef {
	kind := &KindRef{}
	rest := false
	params := make([]*KindRef, 0, helper.DefaultCap)

	for i, param := range t.Params {
		if param.Rest {
			if i < len(t.Params)-1 {
				m.unexpectedPos(param.Start, "The rest parameter should be placed last")
			}
			rest = true
		}
		params = append(params, m.compileKindExpr(param.Kind))
	}

	kind.Ref = &TFunc{
		Params:    params,
		Return:    m.compileKindExpr(t.Return),
		RestParam: rest,
		Impl:      newImpl(),
	}
	return kind
}

func (m *Module) compileStructKind(t *ast.TStructKind) *KindRef {
	kind := &KindRef{}
	extends := make([]*KindRef, 0, helper.SmallCap)
	props := make(map[string]*KindRef)

	for _, pair := range t.Properties {
		key := pair.Key.Name
		_, has := props[key]
		if has {
			m.unexpectedPos(pair.Start, "Duplicate key: "+key)
		}
		props[key] = m.compileKindExpr(pair.Kind)
	}

	for _, item := range t.Extends {
		extends = append(extends, m.compileKindExpr(item))
	}

	kind.Ref = &TStruct{
		Extends:    extends,
		Properties: props,
		Impl:       newImpl(),
	}
	return kind
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
