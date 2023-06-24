package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"strings"
)

type ModuleState uint

const (
	MSInit ModuleState = iota
	MSResolve
	MSParse
	MSPrecompile
	MSCompile
)

func matchKind(expected *KindRef, received *KindRef) bool {
	// 引用 self 指向
	_e, ok := expected.Ref.(*TSelf)
	if ok {
		expected = _e.Kind
	}
	_r, ok := received.Ref.(*TSelf)
	if ok {
		received = _r.Kind
	}

	// equal
	if expected.Ref == received.Ref {
		return true
	}

	// any
	if expected.Ref == typeAny || received.Ref == typeAny {
		return true
	}

	switch expected.Ref.(type) {
	case *TArray:
		r, ok := received.Ref.(*TArray)
		if !ok {
			return false
		}

		e := expected.Ref.(*TArray)
		return e.Len == r.Len && matchKind(e.Kind, r.Kind)
	case *TFunc:
		r, ok := received.Ref.(*TFunc)
		if !ok {
			return false
		}
		e := expected.Ref.(*TFunc)

		if e.RestParam != r.RestParam || len(e.Params) != len(r.Params) {
			return false
		}

		for i, param := range e.Params {
			if !matchKind(param, r.Params[i]) {
				return false
			}
		}
		return matchKind(e.Return, r.Return)
	case *TStruct:
		r, ok := received.Ref.(*TStruct)
		if !ok {
			return false
		}
		e := expected.Ref.(*TStruct)

		// TODO think about extends

		if len(e.Properties) != len(r.Properties) {
			return false
		}
		for key, kind := range e.Properties {
			if !matchKind(kind, r.Properties[key]) {
				return false
			}
		}
		return true
	case *TInterface:
		r, ok := received.Ref.(*TInterface)
		e := expected.Ref.(*TInterface)
		if !ok {
			for _, ref := range e.Refs {
				if matchKind(ref, received) {
					return true
				}
			}
			return false
		}

		if len(e.Properties) != len(r.Properties) {
			return false
		}
		for key, kind := range e.Properties {
			if !matchKind(kind, r.Properties[key]) {
				return false
			}
		}
		return true
	case *TEnum:
		r, ok := received.Ref.(*TEnum)
		if !ok {
			return false
		}
		e := expected.Ref.(*TEnum)

		if len(e.Choices) != len(r.Choices) {
			return false
		}
		for i, v := range e.Choices {
			if v != r.Choices[i] {
				return false
			}
		}
		return true
	case *TCustom:
		e := expected.Ref.(*TCustom)
		return matchKind(e.Kind, received)
	}

	return false
}

func getKindExprString(expr *ast.KindExpr) string {
	if expr == nil {
		return ""
	}

	builder := strings.Builder{}

	switch expr.Node.(type) {
	case *ast.TNumber:
		builder.WriteString("number")
	case *ast.TByte:
		builder.WriteString("byte")
	case *ast.TChar:
		builder.WriteString("char")
	case *ast.TString:
		builder.WriteString("string")
	case *ast.TBool:
		builder.WriteString("bool")
	case *ast.TAny:
		builder.WriteString("any")
	case *ast.TSelf:
		builder.WriteString("self")
	case *ast.TArray:
		node := expr.Node.(*ast.TArray)
		builder.WriteString("[")
		if node.Len != nil {
			lenNode := node.Len.Node.(*ast.NumberLiteral)
			builder.WriteString(lenNode.Text)
		}
		builder.WriteString("]")
		builder.WriteString(getKindExprString(node.Kind))
	case *ast.TIdentifier:
		node := expr.Node.(*ast.TIdentifier)
		builder.WriteString(node.Name.Name)
	case *ast.TMemberKind:
		node := expr.Node.(*ast.TMemberKind)
		builder.WriteString(getKindExprString(node.Left))
		builder.WriteString(".")
		builder.WriteString(getKindExprString(node.Right))
	case *ast.TFuncKind:
		node := expr.Node.(*ast.TFuncKind)
		builder.WriteString("fn(")
		for i, param := range node.Params {
			builder.WriteString(param.Name.Name)
			builder.WriteString(": ")
			builder.WriteString(getKindExprString(param.Kind))
			if i < len(node.Params)-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString(")")
		if len(node.Params) > 0 {
			builder.WriteString(" -> ")
			builder.WriteString(getKindExprString(node.Return))
		}
	case *ast.TStructKind:
		node := expr.Node.(*ast.TStructKind)
		builder.WriteString("struct")
		// struct<-abc{ }
		if len(node.Extends) > 0 {
			builder.WriteString("<-")
			for i, kind := range node.Extends {
				builder.WriteString(getKindExprString(kind))
				if i < len(node.Extends)-1 {
					builder.WriteString(", ")
				}
			}
		}
		builder.WriteString("{")
		if len(node.Properties) > 0 {
			builder.WriteString(" ")
			for i, pair := range node.Properties {
				builder.WriteString(pair.Key.Name)
				builder.WriteString(": ")
				builder.WriteString(getKindExprString(pair.Kind))
				if i < len(node.Properties)-1 {
					builder.WriteString(", ")
				}
			}
			builder.WriteString(" ")
		}
		builder.WriteString("}")

	}

	return builder.String()
}

func isReferenceKind(kind *KindRef) bool {
	switch kind.Ref {
	case typeNumber, typeByte, typeChar, typeString, typeBool:
		return false
	}

	switch kind.Ref.(type) {
	case *TCustom:
		return isReferenceKind(kind.Ref.(*TCustom).Kind)
	case *TSelf:
		return isReferenceKind(kind.Ref.(*TSelf).Kind)
	}

	return true
}
