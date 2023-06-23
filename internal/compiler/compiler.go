package compiler

import (
	"os"
	"path/filepath"
)

/* compiler */

type Compiler struct {
	Root      string
	Main      *Module
	ModuleMap ModuleMap
}

func NewCompiler(root string) *Compiler {
	root, _ = filepath.Abs(root)
	mainPath := filepath.Join(root, "main.noah")
	code, err := os.ReadFile(mainPath)
	if err != nil {
		panic(err)
	}
	return NewPureCompiler(string(code), mainPath)
}

func NewPureCompiler(code string, modulePath string) *Compiler {
	absPath, err := filepath.Abs(modulePath)
	if err == nil {
		modulePath = absPath
	}
	c := &Compiler{
		Root:      filepath.Dir(modulePath),
		ModuleMap: make(ModuleMap),
	}

	c.Main = NewModule(c, code, modulePath)
	c.ModuleMap.add(c.Main)

	return c
}

func (c *Compiler) Compile() *Compiler {
	c.Main.compile()
	return c
}

/* module map */

type ModuleMap map[string]*Module

func (m *ModuleMap) get(pathId string) *Module {
	return (*m)[pathId]
}

func (m *ModuleMap) add(module *Module) {
	(*m)[module.PathId] = module
}

func (m *ModuleMap) has(pathId string) bool {
	return m.get(pathId) != nil
}
