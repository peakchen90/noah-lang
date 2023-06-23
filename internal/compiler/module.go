package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/parser"
)

/* module */

type Module struct {
	Compiler    *Compiler
	PathId      string
	Ast         *ast.File
	PublicScope *Scope
	ScopeStack  *ScopeStack
}

func NewModule(compiler *Compiler, code string, path string) *Module {
	return &Module{
		Compiler: compiler,
		PathId:   normalizePathId(path),
		Ast:      parser.NewParser(code),
		PublicScope: &Scope{
			value: make(map[string]Value),
			kind:  make(map[string]Kind),
		},
		ScopeStack: newScopeStack(),
	}
}

// 执行编译
func (m *Module) compile() {
	m.preCompile()
}

/* module map */

type ModuleMap map[string]*Module

func (m *ModuleMap) get(path string) *Module {
	return (*m)[path]
}

func (m *ModuleMap) add(module *Module) {
	(*m)[module.PathId] = module
}

func (m *ModuleMap) has(path string) bool {
	return m.get(path) != nil
}

func normalizePathId(path string) string {
	// TODO
	return path
}
