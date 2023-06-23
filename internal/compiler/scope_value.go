package compiler

type Value interface{ isConst() bool }

func (*ModuleValue) isConst() bool { return false }
func (*FuncValue) isConst() bool   { return false }
func (v *VarValue) isConst() bool  { return !v.Const }

type (
	ModuleValue struct {
		Name   string
		Module *Module
	}

	FuncValue struct {
		Name string
		Kind Kind
		Ptr  uintptr
	}

	VarValue struct {
		Name  string
		Kind  Kind
		Const bool
		Ptr   uintptr
	}
)
