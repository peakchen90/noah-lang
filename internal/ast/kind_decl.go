package ast

type KD interface{ isKindDecl() }

func (*TypeDeclAlias) isKindDecl()     {}
func (*TypeDeclStruct) isKindDecl()    {}
func (*TypeDeclEnum) isKindDecl()      {}
func (*TypeDeclInterface) isKindDecl() {}

type (
	TypeDeclAlias struct {
		Name KindIdentifier
		Kind KindExpr
	}

	TypeDeclStruct struct {
		Name       KindIdentifier
		Interface  KindIdentifier
		Extends    KindIdentifier
		Properties []KindProperty
	}

	TypeDeclEnum struct {
		Name  KindIdentifier
		Items []KindIdentifier
	}

	TypeDeclInterface struct {
		Name       KindIdentifier
		Extends    KindIdentifier
		Properties []KindProperty
	}
)
