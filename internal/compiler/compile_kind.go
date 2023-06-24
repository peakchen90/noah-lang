package compiler

import (
	"errors"
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
		kind.Ref = typeNumber
	case *ast.TByte:
		kind.Ref = typeByte
	case *ast.TChar:
		kind.Ref = typeChar
	case *ast.TString:
		kind.Ref = typeString
	case *ast.TBool:
		kind.Ref = typeBool
	case *ast.TAny:
		kind.Ref = typeAny
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
			m.unexpectedPos(t.Len.Start, "expect be a positive integer")
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
				m.unexpectedPos(param.Start, "the rest parameter should be placed last")
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
			m.unexpectedPos(pair.Start, "duplicate key: "+key)
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

func (m *Module) inferKind(expr *ast.Expr) (*KindRef, error) {
	kind := &KindRef{}

	switch expr.Node.(type) {
	case *ast.CallExpr:
		return m.inferCallExprKind(expr.Node.(*ast.CallExpr))
	case *ast.MemberExpr:
		return m.inferMemberExprKind(expr.Node.(*ast.MemberExpr))
	case *ast.BinaryExpr:
		return m.inferBinaryExprKind(expr.Node.(*ast.BinaryExpr))
	case *ast.BinaryTypeExpr:
		return m.inferBinaryTypeExprKind(expr.Node.(*ast.BinaryTypeExpr))
	case *ast.UnaryExpr:
		return m.inferUnaryExprKind(expr.Node.(*ast.UnaryExpr))
	case *ast.FuncExpr:
		return m.inferFuncExprKind(expr.Node.(*ast.FuncExpr))
	case *ast.StructExpr:
		return m.inferStructExprKind(expr.Node.(*ast.StructExpr))
	case *ast.ArrayExpr:
		return m.inferArrayExprKind(expr.Node.(*ast.ArrayExpr))
	case *ast.IdentifierLiteral:
		return m.inferIdentifierLiteralKind(expr.Node.(*ast.IdentifierLiteral))
	case *ast.NumberLiteral:
		kind.Ref = typeNumber
	case *ast.BoolLiteral:
		kind.Ref = typeBool
	case *ast.NullLiteral:
		return nil, errors.New("cannot infer the type of null")
	case *ast.StringLiteral:
		kind.Ref = typeString
	case *ast.CharLiteral:
		kind.Ref = typeChar
	default:
		panic("Internal Err")
	}

	return kind, nil
}

func (m *Module) inferCallExprKind(expr *ast.CallExpr) (*KindRef, error) {
	kind := &KindRef{}
	// TODO

	return kind, nil
}

func (m *Module) inferMemberExprKind(expr *ast.MemberExpr) (*KindRef, error) {
	kind := &KindRef{}
	// TODO

	return kind, nil
}

func (m *Module) inferBinaryExprKind(expr *ast.BinaryExpr) (*KindRef, error) {
	kind := &KindRef{}
	// TODO

	return kind, nil
}

func (m *Module) inferBinaryTypeExprKind(expr *ast.BinaryTypeExpr) (*KindRef, error) {
	kind := &KindRef{}

	switch expr.Operator {
	case "is":
		kind.Ref = typeBool
	case "as":
		kind = m.compileKindExpr(expr.Right)
		leftKind, err := m.inferKind(expr.Left)
		if err != nil {
			_, isNull := expr.Left.Node.(*ast.NullLiteral)
			if isNull {
				if !isReferenceKind(kind) {
					m.unexpectedPos(expr.Right.Start, "expect a reference type, but found: "+getKindExprString(expr.Right))
				}
			} else {
				return nil, err
			}
		} else if !matchKind(leftKind, kind) {
			m.unexpectedPos(expr.Left.Start, "cannot use `as` on incompatible type: "+getKindExprString(expr.Right))
		}
	default:
		panic("Internal Err")
	}

	return kind, nil
}

func (m *Module) inferUnaryExprKind(expr *ast.UnaryExpr) (*KindRef, error) {
	kind := &KindRef{}
	// TODO

	return kind, nil
}

func (m *Module) inferFuncExprKind(expr *ast.FuncExpr) (*KindRef, error) {
	kind := &KindRef{}
	// TODO

	return kind, nil
}

func (m *Module) inferStructExprKind(expr *ast.StructExpr) (*KindRef, error) {
	kind := &KindRef{}
	// TODO

	return kind, nil
}

func (m *Module) inferArrayExprKind(expr *ast.ArrayExpr) (*KindRef, error) {
	kind := &KindRef{}
	// TODO

	return kind, nil
}

func (m *Module) inferIdentifierLiteralKind(expr *ast.IdentifierLiteral) (*KindRef, error) {
	//name := expr.Node.(*ast.IdentifierLiteral).Name
	//value := m.scopes.findValue(name, true)
	// TODO

	kind := &KindRef{}
	// TODO

	return kind, nil
}
