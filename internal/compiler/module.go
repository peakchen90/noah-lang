package compiler

import (
	"errors"
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/parser"
	"path/filepath"
	"strings"
)

type Module struct {
	Ast      *ast.File
	compiler *Compiler
	moduleId string
	path     string
	parser   *parser.Parser
	exports  *Scope
	scopes   *ScopeStack

	/* context flags */
	state       ModuleState
	allowImport bool
}

func NewModule(compiler *Compiler) *Module {
	return &Module{
		compiler: compiler,
		exports: &Scope{
			value: make(map[string]Value),
			kind:  make(map[string]*KindRef),
		},
		scopes:      newScopeStack(),
		state:       MSInit,
		allowImport: true,
	}
}

// 解析模块
func (m *Module) resolve(moduleId string) (*Module, error) {
	if m.state >= MSResolve {
		return m, nil
	}
	m.state = MSResolve

	packageName := ""
	pathIds := moduleId

	index := strings.IndexByte(moduleId, ':')
	if index >= 0 {
		packageName = moduleId[:index]
		pathIds = moduleId[index+1:]
	}

	relativePath := strings.Map(func(r rune) rune {
		if r == '.' {
			return '/'
		}
		return r
	}, pathIds)

	modulePath := ""
	virtualFS := m.compiler.VirtualFS
	if len(packageName) == 0 {
		modulePath = filepath.Join(virtualFS.Root, relativePath+".noah")
	} else {
		modulePath = filepath.Join(virtualFS.PackageRoot, packageName, relativePath+".noah")
	}

	if !virtualFS.ExistFile(modulePath) {
		return nil, errors.New("Module not found: " + moduleId)
	}

	m.moduleId = moduleId
	m.path = modulePath
	m.compiler.Modules.add(m)

	return m, nil
}

// 解析模块
func (m *Module) parse() (*Module, error) {
	if m.state >= MSParse {
		return m, nil
	}
	m.state = MSParse

	code, err := m.compiler.VirtualFS.ReadFile(m.path)
	if err != nil {
		return nil, err
	}

	m.parser = parser.NewParser(string(code), m.moduleId)
	m.Ast = m.parser.Parse()

	return m, nil
}

// 预编译模块
func (m *Module) precompile() {
	if m.state >= MSPrecompile {
		return
	}
	m.state = MSPrecompile

	// push scope : 将顶层定义全部放在这里
	m.scopes.push()

	for _, stmt := range m.Ast.Body {
		_, isImport := stmt.Node.(*ast.ImportDecl)
		if !isImport {
			m.allowImport = false
		}

		switch stmt.Node.(type) {
		case *ast.ImportDecl:
			if !m.allowImport {
				m.unexpectedPos(stmt.Start, "`import` should be at the top of the file")
			}
			m.compileImportDecl(stmt.Node.(*ast.ImportDecl), true)
		case *ast.FuncDecl:
			m.compileFuncDecl(stmt.Node.(*ast.FuncDecl), nil, true)
		case *ast.VarDecl:
			m.compileVarDecl(stmt.Node.(*ast.VarDecl), true)
		case *ast.TAliasDecl:
			m.compileTAliasDecl(stmt.Node.(*ast.TAliasDecl), true)
		case *ast.TInterfaceDecl:
			m.compileTInterfaceDecl(stmt.Node.(*ast.TInterfaceDecl), true)
		case *ast.TStructDecl:
			m.compileTStructDecl(stmt.Node.(*ast.TStructDecl), true)
		case *ast.TEnumDecl:
			m.compileTEnumDecl(stmt.Node.(*ast.TEnumDecl), true)
		}
	}
}

// 编译模块
func (m *Module) compile() {
	if m.state >= MSCompile {
		return
	}
	m.state = MSCompile

	for _, stmt := range m.Ast.Body {
		m.compileStmt(stmt)
	}
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

func (m *Module) putKind(name *ast.Identifier, scope *KindRef, isPanic bool) {
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

func (m *Module) putSelfKind(scope *KindRef) {
	last := m.scopes.last()
	if last != nil {
		last.setKind("self", scope)
	}
}

func (m *Module) putSelfValue(scope Value) {
	last := m.scopes.last()
	if last != nil {
		last.setValue("self", scope)
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

func (m *Module) findKind(name string) (*KindRef, error) {
	for i := m.scopes.size() - 1; i >= 0; i-- {
		scope := m.scopes.stack[i].getKind(name)
		if scope != nil {
			return scope, nil
		}
	}
	return nil, errors.New(name + " is not found")
}

func (m *Module) findIdentifierKind(name *ast.Identifier, isPanic bool) *KindRef {
	kind, err := m.findKind(name.Name)

	if err != nil && isPanic {
		m.unexpectedPos(name.Start, name.Name+" is not found")
	}

	return kind
}

func (m *Module) findMemberKind(kindExpr *ast.KindExpr, module *Module, isPanic bool) *KindRef {
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

	var kind *KindRef
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

func (m *Module) findSelfKind(kindExpr *ast.KindExpr, isPanic bool) *KindRef {
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
