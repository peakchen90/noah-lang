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
		kind.Ref = &TSelf{Kind: m.scopes.findSelfKind(kindExpr, true)}
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
		Kind: m.compileKindExpr(t.Kind),
		Len:  size,
		Impl: newImpl(),
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
		return m.compileKindExpr(expr.Node.(*ast.FuncExpr).FuncKind), nil
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

func (m *Module) inferCallExprKind(expr *ast.CallExpr) (kind *KindRef, err error) {
	callee := expr.Callee.Node

	switch callee.(type) {
	case *ast.IdentifierLiteral:
		value := m.scopes.findValue(callee.(*ast.IdentifierLiteral).Name, true)
		kind, err = m.getValueKind(value)
	case *ast.MemberExpr:
		value := m.scopes.findMemberValue(callee.(*ast.MemberExpr), true)
		kind, err = m.getValueKind(value)
	case *ast.FuncExpr:
		kind = m.compileKindExpr(callee.(*ast.FuncExpr).FuncKind)
	default:
		panic("Internal Err")
	}

	if kind != nil {
		funcKind, ok := kind.Ref.(*TFunc)
		if !ok {
			m.unexpectedPos(expr.Callee.Start, "not a function")
		}
		kind = funcKind.Return
	}

	return
}

func (m *Module) getValueKind(value Value) (*KindRef, error) {
	kind := &KindRef{}

	switch value.(type) {
	case *FuncValue:
		kind = value.(*FuncValue).Kind
	case *VarValue:
		kind = value.(*VarValue).Kind
	case *SelfValue:
		kind = value.(*SelfValue).Kind
	default:
		panic("Internal Error")
	}

	return kind, nil
}

func (m *Module) inferIdentifierLiteralKind(expr *ast.IdentifierLiteral) (*KindRef, error) {
	value := m.scopes.findValue(expr.Name, true)
	return m.getValueKind(value)
}

func (m *Module) inferMemberExprKind(expr *ast.MemberExpr) (*KindRef, error) {
	value := m.scopes.findMemberValue(expr, true)
	return m.getValueKind(value)
}

func (m *Module) inferBinaryExprKind(expr *ast.BinaryExpr) (*KindRef, error) {
	kind := &KindRef{}

	switch expr.Operator.Value {
	// assign
	case "=", "+=", "-=", "*=", "/=", "%=", "<<=", ">>=", "&=", "|=", "^=":
		return m.inferKind(expr.Left)

	// logic
	case "||", "&&", "==", "!=", "<", "<=", ">", ">=":
		kind.Ref = typeBool

	// bit op
	case "|", "^", "&", "<<", ">>":
		kind.Ref = typeNumber

	// decimal calc
	case "+", "-", "*", "/", "%":
		kind.Ref = typeNumber

	default:
		panic("Internal Err")
	}

	return kind, nil
}

func (m *Module) inferBinaryTypeExprKind(expr *ast.BinaryTypeExpr) (*KindRef, error) {
	kind := &KindRef{}

	switch expr.Operator.Value {
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
			m.unexpectedPos(expr.Operator.Start, "cannot use `as` on incompatible type: "+getKindExprString(expr.Right))
		}
	default:
		panic("Internal Err")
	}

	return kind, nil
}

func (m *Module) inferUnaryExprKind(expr *ast.UnaryExpr) (*KindRef, error) {
	kind := &KindRef{}

	switch expr.Operator.Value {
	// number op
	case "+", "-", "++", "--":
		kind.Ref = typeNumber

	// logic
	case "!":
		kind.Ref = typeBool

	// bit op
	case "~":
		kind.Ref = typeNumber

	default:
		panic("Internal Err")
	}

	return kind, nil
}

func (m *Module) inferStructExprKind(expr *ast.StructExpr) (*KindRef, error) {
	kind := &KindRef{}
	props := make(map[string]*KindRef)

	for _, pair := range expr.Properties {
		key := pair.Key.Node.(*ast.IdentifierLiteral).Name.Name
		_, has := props[key]
		if has {
			m.unexpectedPos(pair.Key.Start, "duplicate key: "+key)
		}
		inferKind, err := m.inferKind(pair.Value)
		if err != nil {
			return nil, err
		}
		props[key] = inferKind
	}

	kind.Ref = &TStruct{
		Extends:    make([]*KindRef, 0, 0),
		Properties: props,
		Impl:       newImpl(),
	}

	if expr.Ctor != nil {
		inferCtorKind, err := m.inferKind(expr.Ctor)
		if err != nil {
			return nil, err
		}
		_, ok := inferCtorKind.Ref.(*TStruct)
		if !ok {
			m.unexpectedPos(expr.Ctor.Start, "expect a struct")
		}

		// TODO check struct: k&v, extends

		kind = inferCtorKind
	}

	return kind, nil
}

func (m *Module) inferArrayExprKind(expr *ast.ArrayExpr) (*KindRef, error) {
	kind := &KindRef{}
	arr := &TArray{}
	kind.Ref = arr

	if len(expr.Items) > 0 {
		inferKind, err := m.inferKind(expr.Items[0])
		if err != nil {
			return nil, err
		}
		arr.Kind = inferKind
		arr.Len = len(expr.Items)
	}

	// TODO 没有元素如何推断？

	return kind, nil
}
