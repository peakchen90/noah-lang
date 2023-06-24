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
	module := &Module{
		compiler:    compiler,
		exports:     newScope(),
		state:       MSInit,
		allowImport: true,
	}
	module.scopes = newScopeStack(module)
	return module
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
func (m *Module) parse() error {
	if m.state >= MSParse {
		return nil
	}
	m.state = MSParse

	code, err := m.compiler.VirtualFS.ReadFile(m.path)
	if err != nil {
		return err
	}

	m.parser = parser.NewParser(string(code), m.moduleId)
	m.Ast = m.parser.Parse()

	return nil
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

	fns := make([]*ast.Stmt, 0, helper.DefaultCap)
	vars := make([]*ast.Stmt, 0, helper.DefaultCap)

	// 1. 优先编译 类型声明、模块
	for _, stmt := range m.Ast.Body {
		switch stmt.Node.(type) {
		case *ast.FuncDecl, *ast.ImplDecl:
			fns = append(fns, stmt)
		case *ast.VarDecl:
			vars = append(vars, stmt)
		default:
			m.compileStmt(stmt)
		}
	}

	// 2. 其次编译函数
	for _, stmt := range fns {
		m.compileStmt(stmt)
	}

	// 3. 编译全局变量（变量可能依赖类型定义、函数返回值等）
	for _, stmt := range vars {
		m.compileStmt(stmt)
	}
}

func (m *Module) unexpectedPos(index int, msg string) {
	m.parser.UnexpectedPos(index, msg)
}
