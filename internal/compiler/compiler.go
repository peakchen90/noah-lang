package compiler

type Compiler struct {
	MainModule *Module
	ModuleMap  ModuleMap
}

func NewCompiler(code string, path string) *Compiler {
	c := &Compiler{
		ModuleMap: make(ModuleMap),
	}

	c.MainModule = NewModule(c, code, path)
	c.ModuleMap.add(c.MainModule)

	return c
}

func (c *Compiler) Compile() {
	c.MainModule.compile()
}
