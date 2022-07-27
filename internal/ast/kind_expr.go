package ast

type KE interface{ isKind() }

func (*TypeNumber) isKind()      {}
func (*TypeByte) isKind()        {}
func (*TypeChar) isKind()        {}
func (*TypeString) isKind()      {}
func (*TypeBool) isKind()        {}
func (*TypeAny) isKind()         {}
func (*TypeArray) isKind()       {}
func (*TypeVectorArray) isKind() {}
func (*TypeId) isKind()          {}

type (
	TypeNumber struct{}

	TypeByte struct{}

	TypeChar struct{}

	TypeString struct{}

	TypeBool struct{}

	TypeAny struct{}

	TypeArray struct {
		Kind KindExpr
		Len  Expression
	}

	TypeVectorArray struct {
		Kind KindExpr
	}

	TypeId struct {
		Name string
	}
)
