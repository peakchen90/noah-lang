package ast

func (*TypeAliasDecl) isStmt()     {}
func (*TypeInterfaceDecl) isStmt() {}
func (*TypeStructDecl) isStmt()    {}
func (*TypeEnumDecl) isStmt()      {}

type (
	TypeAliasDecl struct {
		Name  *KindIdentifier
		Kind  *KindExpr
		Pubic bool
	}

	TypeInterfaceDecl struct {
		Name       *KindIdentifier
		Extends    *KindExpr
		Properties []*KindProperty
		Pubic      bool
	}

	TypeStructDecl struct {
		Name       *KindIdentifier
		Impl       *KindExpr
		Extends    *KindExpr
		Properties []*KindProperty
		Pubic      bool
	}

	TypeEnumDecl struct {
		Name  *KindIdentifier
		Items []*KindIdentifier
		Pubic bool
	}
)
