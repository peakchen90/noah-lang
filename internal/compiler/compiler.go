package compiler

/* compiler */

type Compiler struct {
	Main      *Module
	Modules   ModuleMap
	VirtualFS *VirtualFS
}

func NewCompiler(root string, isFileSystem bool) *Compiler {
	virtualFS := newVirtualFS(root, isFileSystem)
	return &Compiler{
		Modules:   make(ModuleMap),
		VirtualFS: virtualFS,
	}
}

func (c *Compiler) Compile() *Compiler {
	module, err := NewModule(c).resolve("main")
	if err != nil {
		panic(err)
	}

	err = module.parse()
	if err != nil {
		panic(err)
	}

	c.Main = module
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
