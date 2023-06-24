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

func compareKind(expected *KindRef, received *KindRef, isMatch bool) bool {
	if expected == nil && received == nil {
		return true
	}

	_, ok := received.Ref.(*TAny)
	if ok {
		return true
	}

	_e, ok := expected.Ref.(*TSelf)
	if ok {
		expected = _e.KindRef
	}
	_r, ok := received.Ref.(*TSelf)
	if ok {
		received = _r.KindRef
	}

	switch expected.Ref.(type) {
	case *TNumber:
		_, ok := received.Ref.(*TNumber)
		return ok
	case *TByte:
		_, ok := received.Ref.(*TByte)
		return ok
	case *TChar:
		_, ok := received.Ref.(*TChar)
		return ok
	case *TString:
		_, ok := received.Ref.(*TString)
		return ok
	case *TBool:
		_, ok := received.Ref.(*TBool)
		return ok
	case *TAny:
		return true
	case *TArray:
		r, ok := received.Ref.(*TArray)
		if !ok {
			return false
		}

		e := expected.Ref.(*TArray)
		return e.Len == r.Len && compareKind(e.KindRef, r.KindRef, isMatch)
	case *TFunc:
		r, ok := received.Ref.(*TFunc)
		if !ok {
			return false
		}
		e := expected.Ref.(*TFunc)

		if isMatch {
			if e.RestArgument != r.RestArgument || len(e.Arguments) != len(r.Arguments) {
				return false
			}

			for i, arg := range e.Arguments {
				if !compareKind(arg, r.Arguments[i], isMatch) {
					return false
				}
			}
			return compareKind(e.Return, r.Return, isMatch)
		} else {
			return r == e
		}
	case *TStruct:
		r, ok := received.Ref.(*TStruct)
		if !ok {
			return false
		}
		e := expected.Ref.(*TStruct)

		if isMatch {
			// TODO think about extends

			if len(e.Properties) != len(r.Properties) {
				return false
			}
			for key, kind := range e.Properties {
				if !compareKind(kind, r.Properties[key], isMatch) {
					return false
				}
			}
			return true
		} else {
			return r == e
		}
	case *TInterface:
		r, ok := received.Ref.(*TInterface)
		e := expected.Ref.(*TInterface)
		if !ok {
			for _, ref := range e.Refs {
				if isMatch && compareKind(ref, received, isMatch) {
					return true
				}
			}
			return false
		}

		if isMatch {
			if len(e.Properties) != len(r.Properties) {
				return false
			}
			for key, kind := range e.Properties {
				if !compareKind(kind, r.Properties[key], isMatch) {
					return false
				}
			}
			return true
		} else {
			return r == e
		}
	case *TEnum:
		r, ok := received.Ref.(*TEnum)
		if !ok {
			return false
		}
		e := expected.Ref.(*TEnum)

		if isMatch {
			if len(e.Choices) != len(r.Choices) {
				return false
			}
			for i, v := range e.Choices {
				if v != r.Choices[i] {
					return false
				}
			}
			return true
		} else {
			return r == e
		}
	case *TCustom:
		e := expected.Ref.(*TCustom)
		if isMatch {
			return compareKind(e.KindRef, received, true)
		} else {
			r, ok := received.Ref.(*TCustom)
			return ok && r == e
		}
	}

	return false
}

func getKindExprId(expr *ast.KindExpr) string {
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
	case *ast.TArray:
		node := expr.Node.(*ast.TArray)
		builder.WriteString("[")
		if node.Len != nil {
			lenNode := node.Len.Node.(*ast.NumberLiteral)
			builder.WriteString(lenNode.Text)
		}
		builder.WriteString("]")
		builder.WriteString(getKindExprId(node.Kind))
	case *ast.TIdentifier:
		node := expr.Node.(*ast.TIdentifier)
		builder.WriteString(node.Name.Name)
	case *ast.TMemberKind:
		node := expr.Node.(*ast.TMemberKind)
		builder.WriteString(getKindExprId(node.Left))
		builder.WriteString(".")
		builder.WriteString(getKindExprId(node.Right))
	case *ast.TFuncKind:
		node := expr.Node.(*ast.TFuncKind)
		builder.WriteString("fn(")
		for i, arg := range node.Arguments {
			builder.WriteString(arg.Name.Name)
			builder.WriteString(": ")
			builder.WriteString(getKindExprId(arg.Kind))
			if i < len(node.Arguments)-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString(")")
		if len(node.Arguments) > 0 {
			builder.WriteString(" -> ")
			builder.WriteString(getKindExprId(node.Return))
		}
	case *ast.TStructKind:
		node := expr.Node.(*ast.TStructKind)
		builder.WriteString("struct")
		// struct<-abc{ }
		if len(node.Extends) > 0 {
			builder.WriteString("<-")
			for i, kind := range node.Extends {
				builder.WriteString(getKindExprId(kind))
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
				builder.WriteString(getKindExprId(pair.Kind))
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
