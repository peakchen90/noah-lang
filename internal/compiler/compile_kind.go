package compiler

import (
	"errors"
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"math"
)

func (m *Module) compileKindExpr(kindExpr *ast.KindExpr) *KindRef {
	kind := newKindRef(m, -1)
	if kindExpr == nil {
		return kind
	}

	node := kindExpr.Node

	switch node.(type) {
	case *ast.TNumber:
		kind.current = typeNumber
	case *ast.TByte:
		kind.current = typeByte
	case *ast.TChar:
		kind.current = typeChar
	case *ast.TString:
		kind.current = typeString
	case *ast.TBool:
		kind.current = typeBool
	case *ast.TAny:
		kind.current = typeAny
	case *ast.TSelf:
		kind.current = &TSelf{Kind: m.scopes.findSelfKind(kindExpr, true)}
	case *ast.TArray:
		return m.compileArrayKind(kindExpr)
	case *ast.TIdentifier:
		node := node.(*ast.TIdentifier)
		return m.scopes.findIdentifierKind(node.Name, true)
	case *ast.TMemberKind:
		return m.scopes.findMemberKind(kindExpr, true)
	case *ast.TFuncKind:
		return m.compileFuncKind(kindExpr)
	case *ast.TStructKind:
		return m.compileStructKind(nil, kindExpr)
	}

	return kind
}

func (m *Module) compileArrayKind(kindExpr *ast.KindExpr) *KindRef {
	node := kindExpr.Node.(*ast.TArray)

	kind := newKindRef(m, -1)
	size := -1 // vector array

	if node.Len != nil {
		rawVal := node.Len.Node.(*ast.NumberLiteral).Value
		if rawVal < 0 || math.Floor(rawVal) != rawVal {
			m.unexpectedPos(node.Len.Start, "expect be a positive integer")
		}
		size = int(rawVal)
	}

	kind.current = &TArray{
		Kind: m.compileKindExpr(node.Kind),
		Len:  size,
		Impl: newImpl(),
	}
	return kind
}

func (m *Module) compileFuncKind(kindExpr *ast.KindExpr) *KindRef {
	node := kindExpr.Node.(*ast.TFuncKind)

	kind := newKindRef(m, -1)
	hasRest := false
	arguments := make([]*KindRef, 0, helper.DefaultCap)

	for i, arg := range node.Arguments {
		if arg.Rest {
			if i < len(node.Arguments)-1 {
				m.unexpectedPos(arg.Start, "the rest argument should be placed last")
			}
			hasRest = true
		}
		arguments = append(arguments, m.compileKindExpr(arg.Kind))
	}

	kind.current = &TFunc{
		Arguments: arguments,
		Return:    m.compileKindExpr(node.Return),
		HasRest:   hasRest,
		Impl:      newImpl(),
	}
	return kind
}

func (m *Module) compileStructKind(kind *KindRef, kindExpr *ast.KindExpr) *KindRef {
	node := kindExpr.Node.(*ast.TStructKind)

	if kind == nil {
		kind = newKindRef(m, helper.SmallCap)
	}
	extends := make([]*KindRef, 0, helper.SmallCap)
	props := make(map[string]*KindRef)

	for _, pair := range node.Properties {
		key := pair.Key.Name
		_, has := props[key]
		if has {
			m.unexpectedPos(pair.Start, "duplicate key: "+key)
		}
		props[key] = m.compileKindExpr(pair.Kind)
	}

	for _, item := range node.Extends {
		extendKind := m.compileKindExpr(item)
		if extendKind == kind {
			m.unexpectedPos(item.Start, "cannot extend itself")
		}
		for _, ref := range kind.refs {
			if ref == extendKind {
				m.unexpectedPos(item.Start, "cannot extends cycle")
			}
		}
		_, is := extendKind.current.(*TStruct)
		if !is {
			m.unexpectedPos(item.Start, "expect a struct")
		}
		extends = append(extends, extendKind)

		walkStruct(extendKind, func(ref *KindRef) {
			ref.refs = append(ref.refs, kind)
		}, false)
	}

	kind.current = &TStruct{
		Extends:    extends,
		Properties: props,
		Impl:       newImpl(),
	}
	return kind
}

func (m *Module) inferKind(expr *ast.Expr) (*KindRef, error) {
	kind := newKindRef(m, -1)

	switch expr.Node.(type) {
	case *ast.CallExpr:
		return m.inferCallExprKind(expr.Node.(*ast.CallExpr))
	case *ast.MemberExpr:
		return m.inferMemberExprKind(expr)
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
		kind.current = typeNumber
	case *ast.BoolLiteral:
		kind.current = typeBool
	case *ast.NullLiteral:
		return nil, errors.New("cannot infer the type of null")
	case *ast.StringLiteral:
		kind.current = typeString
	case *ast.CharLiteral:
		kind.current = typeChar
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
		value := m.scopes.findMemberValue(expr.Callee, true)
		kind, err = m.getValueKind(value)
	case *ast.FuncExpr:
		kind = m.compileKindExpr(callee.(*ast.FuncExpr).FuncKind)
	default:
		panic("Internal Err")
	}

	if kind != nil {
		funcKind, ok := kind.current.(*TFunc)
		if !ok {
			m.unexpectedPos(expr.Callee.Start, "not a function")
		}
		kind = funcKind.Return
	}

	return
}

func (m *Module) getValueKind(value Value) (kind *KindRef, err error) {
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

	return
}

func (m *Module) inferIdentifierLiteralKind(expr *ast.IdentifierLiteral) (*KindRef, error) {
	value := m.scopes.findValue(expr.Name, false)
	if value != nil {
		return m.getValueKind(value)
	}

	kind := m.scopes.findIdentifierKind(expr.Name, true)
	return kind, nil
}

func (m *Module) inferMemberExprKind(expr *ast.Expr) (*KindRef, error) {
	value := m.scopes.findMemberValue(expr, false)
	if value != nil {
		return m.getValueKind(value)
	}

	// TODO
	return nil, nil
}

func (m *Module) inferBinaryExprKind(expr *ast.BinaryExpr) (*KindRef, error) {
	kind := newKindRef(m, -1)

	switch expr.Operator.Value {
	// assign
	case "=", "+=", "-=", "*=", "/=", "%=", "<<=", ">>=", "&=", "|=", "^=":
		return m.inferKind(expr.Left)

	// logic
	case "||", "&&", "==", "!=", "<", "<=", ">", ">=":
		kind.current = typeBool

	// bit op
	case "|", "^", "&", "<<", ">>":
		kind.current = typeNumber

	// decimal calc
	case "+", "-", "*", "/", "%":
		kind.current = typeNumber

	default:
		panic("Internal Err")
	}

	return kind, nil
}

func (m *Module) inferBinaryTypeExprKind(expr *ast.BinaryTypeExpr) (*KindRef, error) {
	kind := newKindRef(m, -1)

	switch expr.Operator.Value {
	case "is":
		kind.current = typeBool
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
		} else if !matchKind(leftKind, kind, true) {
			m.unexpectedPos(expr.Operator.Start, "cannot use `as` on incompatible type: "+getKindExprString(expr.Right))
		}
	default:
		panic("Internal Err")
	}

	return kind, nil
}

func (m *Module) inferUnaryExprKind(expr *ast.UnaryExpr) (*KindRef, error) {
	kind := newKindRef(m, -1)

	switch expr.Operator.Value {
	// number op
	case "+", "-", "++", "--":
		kind.current = typeNumber

	// logic
	case "!":
		kind.current = typeBool

	// bit op
	case "~":
		kind.current = typeNumber

	default:
		panic("Internal Err")
	}

	return kind, nil
}

func (m *Module) inferStructExprKind(expr *ast.StructExpr) (*KindRef, error) {
	kind := newKindRef(m, 0)
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

	kind.current = &TStruct{
		Extends:    make([]*KindRef, 0, 0),
		Properties: props,
		Impl:       newImpl(),
	}

	if expr.Ctor != nil {
		ctorKind := m.compileKindExpr(expr.Ctor)
		_, ok := ctorKind.current.(*TStruct)
		if !ok {
			m.unexpectedPos(expr.Ctor.Start, "expect a struct")
		}

		if !matchKind(ctorKind, kind, true) {
			m.unexpectedPos(expr.Ctor.End, "cannot match struct: "+getKindExprString(expr.Ctor))
		}

		kind = ctorKind
	}

	return kind, nil
}

func (m *Module) inferArrayExprKind(expr *ast.ArrayExpr) (*KindRef, error) {
	kind := newKindRef(m, -1)
	arr := &TArray{}
	kind.current = arr

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
