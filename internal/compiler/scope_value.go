package compiler

type Value interface{ isConst() bool }

func (*ModuleValue) isConst() bool { return false }
func (*FuncValue) isConst() bool   { return false }
func (v *VarValue) isConst() bool  { return !v.Const }
func (v *SelfValue) isConst() bool { return false }

type (
	ModuleValue struct {
		Name   string
		Module *Module
	}

	FuncValue struct {
		Name    string
		KindRef *KindRef
		Ptr     uintptr
	}

	VarValue struct {
		Name    string
		KindRef *KindRef
		Const   bool
		Ptr     uintptr
	}

	SelfValue struct {
		KindRef *KindRef
	}
)
