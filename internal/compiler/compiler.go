package compiler

/* compiler */

type Compiler struct {
	Main      *Module
	Modules   ModuleMap
	VirtualFS *VirtualFS
}

func NewCompiler(root string, isFileSystem bool) *Compiler {
	virtualFS := newVirtualFS(root, isFileSystem)
	_compiler := &Compiler{
		Modules:   make(ModuleMap),
		VirtualFS: virtualFS,
	}

	module, err := NewModule(_compiler).resolve("main")
	if err != nil {
		panic(err)
	}
	_compiler.Main = module
	return _compiler
}

func (c *Compiler) Compile() *Compiler {
	_, err := c.Main.parse()
	if err != nil {
		panic(err)
	}
	c.Main.precompile()
	c.Main.compile()
	return c
}

/* module map */

type ModuleMap map[string]*Module

func (m *ModuleMap) find(moduleId string) (*Module, bool) {
	v, has := (*m)[moduleId]
	return v, has
}

func (m *ModuleMap) add(module *Module) {
	(*m)[module.moduleId] = module
}
