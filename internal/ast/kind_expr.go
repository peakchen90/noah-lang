package ast

type KE interface{ isKindExpr() }

func (*TypeNumber) isKindExpr()   {}
func (*TypeByte) isKindExpr()     {}
func (*TypeChar) isKindExpr()     {}
func (*TypeString) isKindExpr()   {}
func (*TypeBool) isKindExpr()     {}
func (*TypeAny) isKindExpr()      {}
func (*TypeArray) isKindExpr()    {}
func (*TypeId) isKindExpr()       {}
func (*TypeMember) isKindExpr()   {}
func (*TypeFuncKind) isKindExpr() {}

type (
	TypeNumber struct{}

	TypeByte struct{}

	TypeChar struct{}

	TypeString struct{}

	TypeBool struct{}

	TypeAny struct{}

	TypeArray struct {
		Kind *KindExpr
		Len  *Expr
	}

	TypeId struct {
		Name *KindIdentifier
	}

	TypeMember struct {
		Parent *KindExpr
		Id     *KindExpr
	}

	TypeFuncKind struct {
		Arguments []*Argument
		Return    *KindExpr
	}
)

func (t *TypeArray) IsVector() bool {
	return t.Len == nil
}
