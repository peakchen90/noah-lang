package ast

type KD interface{ isKind() }

func (*TypeAlias) isKind()     {}
func (*TypeStruct) isKind()    {}
func (*TypeEnum) isKind()      {}
func (*TypeInterface) isKind() {}
func (*TypeFunction) isKind()  {}

type (
	TypeAlias struct {
		Name KindIdentifier
		Kind KindExpr
	}

	TypeStruct struct {
		Name       KindIdentifier
		Interface  KindIdentifier
		Extends    KindIdentifier
		Properties []KindProperty
	}

	TypeEnum struct {
		Name  KindIdentifier
		Items []KindIdentifier
	}

	TypeInterface struct {
		Name       KindIdentifier
		Extends    KindIdentifier
		Properties []KindProperty
	}

	TypeFunction struct {
		Name      KindIdentifier
		Arguments []Argument
		Kind      KindExpr
	}
)
