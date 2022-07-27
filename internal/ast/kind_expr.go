package ast

type KE interface{ isKind() }

func (*TypeNumber) isKind()      {}
func (*TypeString) isKind()      {}
func (*TypeBool) isKind()        {}
func (*TypeChar) isKind()        {}
func (*TypeArray) isKind()       {}
func (*TypeVectorArray) isKind() {}
func (*TypeAny) isKind()         {}
func (*KindId) isKind()          {}

type (
	TypeNumber struct{}

	TypeString struct{}

	TypeBool struct{}

	TypeChar struct{}

	TypeArray struct {
		Kind KindExpr
		Len  Expression
	}

	TypeVectorArray struct {
		Kind KindExpr
	}

	TypeAny struct{}

	KindId struct {
		Name string
	}
)
