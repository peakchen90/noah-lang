package compiler

import (
	"path/filepath"
)

/* compiler */

type Compiler struct {
	Main      *Module
	Modules   ModuleMap
	VirtualFS *VirtualFS
}

func NewCompiler(root string, isFileSystem bool) *Compiler {
	virtualFS := newVirtualFS(root, isFileSystem)
	mainFile := filepath.Join(virtualFS.Root, "main.noah")
	code, err := virtualFS.ReadFile(mainFile)
	if err != nil {
		panic(err)
	}

	c := &Compiler{
		Modules:   make(ModuleMap),
		VirtualFS: virtualFS,
	}

	c.Main = NewModule(c, string(code), "", ":main")
	c.Modules.add(c.Main)

	return c
}

func (c *Compiler) Compile() *Compiler {
	c.Main.compile()
	return c
}

/* module map */

type ModuleMap map[string]*Module

func (m *ModuleMap) get(moduleId string) *Module {
	return (*m)[moduleId]
}

func (m *ModuleMap) add(module *Module) {
	(*m)[module.moduleId] = module
}

func (m *ModuleMap) has(moduleId string) bool {
	return m.get(moduleId) != nil
}
