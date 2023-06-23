package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/parser"
	"strings"
)

/* module */

type Module struct {
	Ast         *ast.File
	compiler    *Compiler
	packageName string
	moduleId    string
	exports     *Scope
	scopes      *ScopeStack
}

func NewModule(compiler *Compiler, code string, packageName string, moduleId string) *Module {
	return &Module{
		compiler:    compiler,
		packageName: packageName,
		moduleId:    moduleId,
		Ast:         parser.NewParser(code).Parse(),
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
				// TODO
				panic("exist: " + name.Name)
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
				// TODO
				panic("exist: " + name.Name)
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
				// TODO
				panic("exist: " + name.Name)
			}
		}

		last.setKind(name.Name, scope)
	}
}

func (m *Module) putSelfKind(name string, scope Kind) {
	last := m.scopes.last()
	if last != nil {
		last.setKind(name, scope)
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
		// TODO
		panic("can not found: " + name.Name)
	}

	return nil
}

func (m *Module) findIdentifierKind(kindExpr *ast.KindExpr, isPanic bool) Kind {
	node, ok := kindExpr.Node.(*ast.TIdentifier)
	if !ok {
		panic("Internal Err")
	}

	for i := m.scopes.size() - 1; i >= 0; i-- {
		scope := m.scopes.stack[i].getKind(node.Name.Name)
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
			kind = module.findIdentifierKind(item, isPanic)
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
