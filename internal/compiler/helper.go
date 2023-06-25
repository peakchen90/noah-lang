package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"strings"
)

type ModuleState uint

const (
	ModuleInit ModuleState = iota
	ModuleResolve
	ModuleParse
	ModulePrecompile
	ModuleCompile
)

// isLooseStruct 为 true 表示 expected >= received (属性)
func matchKind(expected *KindRef, received *KindRef, isLooseStruct bool) bool {
	if expected == nil || received == nil {
		return false
	}

	// 引用 self 指向
	_e, ok := expected.current.(*TSelf)
	if ok {
		expected = _e.Kind
	}
	_r, ok := received.current.(*TSelf)
	if ok {
		received = _r.Kind
	}

	// equal
	if expected.current == received.current {
		return true
	}

	// any
	if expected.current == typeAny || received.current == typeAny {
		return true
	}

	switch expected.current.(type) {
	case *TArray:
		r, ok := received.current.(*TArray)
		if !ok {
			return false
		}

		e := expected.current.(*TArray)
		return e.Len == r.Len && matchKind(e.Kind, r.Kind, isLooseStruct)
	case *TFunc:
		r, ok := received.current.(*TFunc)
		if !ok {
			return false
		}
		e := expected.current.(*TFunc)

		if e.RestParam != r.RestParam || len(e.Params) != len(r.Params) {
			return false
		}

		for i, param := range e.Params {
			if !matchKind(param, r.Params[i], isLooseStruct) {
				return false
			}
		}
		return matchKind(e.Return, r.Return, isLooseStruct)
	case *TStruct:
		_, ok := received.current.(*TStruct)
		if !ok {
			return false
		}

		// extends
		expectedProps := getStructProperties(expected)
		receivedProps := getStructProperties(received)

		if !isLooseStruct && len(expectedProps) != len(receivedProps) {
			return false
		}
		for key, kind := range receivedProps {
			if !matchKind(kind, expectedProps[key], false) {
				return false
			}
		}
		return true
	case *TInterface:
		r, ok := received.current.(*TInterface)
		e := expected.current.(*TInterface)
		if !ok {
			for _, ref := range expected.refs {
				if matchKind(ref, received, isLooseStruct) {
					return true
				}
			}
			return false
		}

		if len(e.Properties) != len(r.Properties) {
			return false
		}
		for key, kind := range e.Properties {
			if !matchKind(kind, r.Properties[key], isLooseStruct) {
				return false
			}
		}
		return true
	case *TCustom:
		e := expected.current.(*TCustom)
		return matchKind(e.Kind, received, isLooseStruct)
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
	switch kind.current {
	case typeNumber, typeByte, typeChar, typeString, typeBool:
		return false
	}

	switch kind.current.(type) {
	case *TCustom:
		return isReferenceKind(kind.current.(*TCustom).Kind)
	case *TSelf:
		return isReferenceKind(kind.current.(*TSelf).Kind)
	}

	return true
}

func walkStruct(kind *KindRef, callback func(*KindRef)) {
	node := kind.current.(*TStruct)
	callback(kind)

	for _, extend := range node.Extends {
		walkStruct(extend, callback)
	}
}

func getStructProperties(kind *KindRef) map[string]*KindRef {
	node := kind.current.(*TStruct)
	if len(node.Extends) == 0 {
		return node.Properties
	}

	properties := make(map[string]*KindRef)

	var walkProps func(k *KindRef)
	walkProps = func(_kind *KindRef) {
		_node := _kind.current.(*TStruct)
		for i := len(_node.Extends) - 1; i >= 0; i-- {
			walkProps(_node.Extends[i])
		}

		for k, v := range _node.Properties {
			if kind.module != v.module && k[0] == '_' {
				continue
			}
			properties[k] = v
		}
	}

	walkProps(kind)

	return properties
}
