package compiler

import (
	"errors"
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/parser"
	"strings"
)

/* module */

type Module struct {
	Ast         *ast.File
	compiler    *Compiler
	parser      *parser.Parser
	packageName string
	moduleId    string
	exports     *Scope
	scopes      *ScopeStack
}

func NewModule(compiler *Compiler, code string, packageName string, moduleId string) *Module {
	p := parser.NewParser(code)

	return &Module{
		compiler:    compiler,
		parser:      p,
		packageName: packageName,
		moduleId:    moduleId,
		Ast:         p.Parse(),
		exports: &Scope{
			value: make(map[string]Value),
			kind:  make(map[string]Kind),
		},
		scopes: newScopeStack(),
	}
}

// 执行编译
func (m *Module) compile() {
	m.preCompile()
}

func (m *Module) putValue(name *ast.Identifier, scope Value, isPanic bool) {
	last := m.scopes.last()
	if last != nil {
		if last.has(name.Name) {
			if isPanic {
				m.unexpectedPos(name.Start, "Identifier has already been declared: "+name.Name)
			}
		}

		last.setValue(name.Name, scope)
	}
}

func (m *Module) putModule(name *ast.Identifier, scope Value, isPanic bool) {
	last := m.scopes.last()
	if last != nil {
		if last.has(name.Name) {
			if isPanic {
				m.unexpectedPos(name.Start, "Identifier has already been declared: "+name.Name)
			}
		}

		last.setValue(name.Name, scope)
	}
}

func (m *Module) putKind(name *ast.Identifier, scope Kind, isPanic bool) {
	last := m.scopes.last()
	if last != nil {
		if last.has(name.Name) {
			if isPanic {
				m.unexpectedPos(name.Start, "Identifier has already been declared: "+name.Name)
			}
		}

		last.setKind(name.Name, scope)
	}
}

func (m *Module) putSelfKind(scope Kind) {
	last := m.scopes.last()
	if last != nil {
		last.setKind("self", scope)
	}
}

func (m *Module) findValue(name *ast.Identifier, isPanic bool) Value {
	for i := m.scopes.size() - 1; i >= 0; i-- {
		scope := m.scopes.stack[i].getValue(name.Name)
		if scope != nil {
			return scope
		}
	}

	if isPanic {
		m.unexpectedPos(name.Start, name.Name+" is not defined")
	}

	return nil
}

func (m *Module) findKind(name string) (Kind, error) {
	for i := m.scopes.size() - 1; i >= 0; i-- {
		scope := m.scopes.stack[i].getKind(name)
		if scope != nil {
			return scope, nil
		}
	}
	return nil, errors.New(name + " is not found")
}

func (m *Module) findIdentifierKind(name *ast.Identifier, isPanic bool) Kind {
	kind, err := m.findKind(name.Name)

	if err != nil && isPanic {
		m.unexpectedPos(name.Start, name.Name+" is not found")
	}

	return kind
}

func (m *Module) findMemberKind(kindExpr *ast.KindExpr, module *Module, isPanic bool) Kind {
	memberIdStack := make([]*ast.KindExpr, 0, helper.SmallCap)
	current := kindExpr

outer:
	for {
		switch current.Node.(type) {
		case *ast.TMemberKind:
			node := current.Node.(*ast.TMemberKind)
			memberIdStack = append(memberIdStack, node.Right)
			current = node.Left
		case *ast.TIdentifier:
			memberIdStack = append(memberIdStack, current)
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

	for i := len(memberIdStack) - 1; i >= 0; i-- {
		item := memberIdStack[i]
		node := item.Node.(*ast.TIdentifier)
		builder.WriteString(node.Name.Name)

		if i > 0 {
			value := module.findModule(node.Name, isPanic)
			if value == nil {
				break
			}
			module = value.Module
			builder.WriteByte('.')
		} else {
			kind = module.findIdentifierKind(node.Name, isPanic)
		}
	}

	if kind == nil {
		if isPanic {
			m.unexpectedPos(kindExpr.Start, builder.String()+" is not found")
		}
	}

	return kind
}

func (m *Module) findSelfKind(kindExpr *ast.KindExpr, isPanic bool) Kind {
	_, ok := kindExpr.Node.(*ast.TSelf)
	if !ok {
		panic("Internal Err")
	}

	kind, err := m.findKind("self")

	if err != nil && isPanic {
		m.unexpectedPos(kindExpr.Start, "Cannot use self here")
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
		m.unexpectedPos(name.Start, name.Name+" is not found")
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
		m.unexpectedPos(name.Start, name.Name+" is not found")
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
		m.unexpectedPos(name.Start, name.Name+" is not found")
	}

	return nil
}

func (m *Module) unexpectedPos(index int, msg string) {
	m.parser.UnexpectedPos(index, msg)
}
