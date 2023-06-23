package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/parser"
	"strings"
)

/* module */

type Module struct {
	Compiler    *Compiler
	PathId      string
	Ast         *ast.File
	PublicScope *Scope
	Scopes      *ScopeStack
}

func NewModule(compiler *Compiler, code string, pathId string) *Module {
	return &Module{
		Compiler: compiler,
		PathId:   pathId,
		Ast:      parser.NewParser(code).Parse(),
		PublicScope: &Scope{
			value: make(map[string]Value),
			kind:  make(map[string]Kind),
		},
		Scopes: newScopeStack(),
	}
}

// 执行编译
func (m *Module) compile() {
	m.preCompile()
}

func (m *Module) putValue(name *ast.Identifier, scope Value, isPanic bool) {
	last := m.Scopes.last()
	if last != nil {
		if last.has(name.Name) {
			if isPanic {
				// TODO
				panic("exist: " + name.Name)
			}
		}

		last.setValue(name.Name, scope)
	}
}

func (m *Module) putKind(name *ast.Identifier, scope Kind, isPanic bool) {
	last := m.Scopes.last()
	if last != nil {
		if last.has(name.Name) {
			if isPanic {
				// TODO
				panic("exist: " + name.Name)
			}
		}

		last.setKind(name.Name, scope)
	}
}

func (m *Module) findValue(name *ast.Identifier, isPanic bool) Value {
	for i := m.Scopes.size() - 1; i >= 0; i-- {
		scope := m.Scopes.stack[i].getValue(name.Name)
		if scope != nil {
			return scope
		}
	}

	if isPanic {
		// TODO
		panic("can not found: " + name.Name)
	}

	return nil
}

func (m *Module) findIdentifierKind(kindExpr *ast.KindExpr, isPanic bool) Kind {
	node, _ := kindExpr.Node.(*ast.TIdentifier)

	for i := m.Scopes.size() - 1; i >= 0; i-- {
		scope := m.Scopes.stack[i].getKind(node.Name.Name)
		if scope != nil {
			return scope
		}
	}

	if isPanic {
		// TODO
		panic("can not found: " + node.Name.Name)
	}

	return nil
}

func (m *Module) findMemberKind(kindExpr *ast.KindExpr, module *Module, isPanic bool) Kind {
	memberStack := make([]*ast.TIdentifier, 0, helper.SmallCap)
	current := kindExpr.Node

outer:
	for {
		switch current.(type) {
		case *ast.TMemberKind:
			node := current.(*ast.TMemberKind)
			memberStack = append(memberStack, node.Right.Node.(*ast.TIdentifier))
			current = node.Left.Node
		case *ast.TIdentifier:
			memberStack = append(memberStack, current.(*ast.TIdentifier))
			break outer
		default:
			panic("Internal Err")
		}
	}

	if module == nil {
		module = m
	}

	var kind Kind
	builder := strings.Builder{}

	for i := len(memberStack) - 1; i >= 0; i-- {
		node := memberStack[i]
		builder.WriteString(node.Name.Name)

		if i > 0 {
			value := module.findModule(node.Name, isPanic)
			if value == nil {
				break
			}
			module = value.Module
			builder.WriteString(".")
		} else {
			kind = module.findIdentifierKind(kindExpr, isPanic)
		}
	}

	if kind == nil {
		if isPanic {
			// TODO
			panic("can not found: " + builder.String())
		}
	}

	return kind
}

func (m *Module) findModule(name *ast.Identifier, isPanic bool) *ModuleValue {
	scope := m.findValue(name, isPanic)
	value, ok := scope.(*ModuleValue)
	if scope != nil && ok {
		return value
	}

	if isPanic {
		// TODO
		panic("can not found: " + name.Name)
	}

	return nil
}

func (m *Module) findFunc(name *ast.Identifier, isPanic bool) *FuncValue {
	scope := m.findValue(name, isPanic)
	value, ok := scope.(*FuncValue)
	if scope != nil && ok {
		return value
	}

	if isPanic {
		// TODO
		panic("can not found: " + name.Name)
	}

	return nil
}

func (m *Module) findVar(name *ast.Identifier, isPanic bool) *VarValue {
	scope := m.findValue(name, isPanic)
	value, ok := scope.(*VarValue)
	if scope != nil && ok {
		return value
	}

	if isPanic {
		// TODO
		panic("can not found: " + name.Name)
	}

	return nil
}
