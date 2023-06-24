package compiler

type Value interface{ isConst() bool }

func (*FuncValue) isConst() bool   { return false }
func (v *VarValue) isConst() bool  { return !v.Const }
func (v *SelfValue) isConst() bool { return false }

type (
	FuncValue struct {
		Name string
		Kind *KindRef
		Ptr  uintptr
	}

	VarValue struct {
		Name  string
		Kind  *KindRef
		Const bool
		Ptr   uintptr
	}

	SelfValue struct {
		Kind *KindRef
	}
)
