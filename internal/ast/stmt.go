package ast

// S statements
type S interface{ isStmt() }

func (*UseModuleStmt) isStmt() {}
func (*FuncDecl) isStmt()      {}
func (*ImplDecl) isStmt()      {}
func (*VarDecl) isStmt()       {}
func (*BlockStmt) isStmt()     {}
func (*ReturnStmt) isStmt()    {}
func (*ExprStmt) isStmt()      {}
func (*IfStmt) isStmt()        {}
func (*ForStmt) isStmt()       {}
func (*BreakStmt) isStmt()     {}
func (*ContinueStmt) isStmt()  {}

func (*TypeAliasDecl) isStmt()     {}
func (*TypeInterfaceDecl) isStmt() {}
func (*TypeStructDecl) isStmt()    {}
func (*TypeEnumDecl) isStmt()      {}

/* statements */
type (
	UseModuleStmt struct {
		Source *Expr
		Local  *Identifier
		Pub    bool
	}

	FuncDecl struct {
		Name *Identifier
		Kind *KindExpr
		Body *Stmt
		Pub  bool
	}

	ImplDecl struct {
		Target    *KindExpr
		Interface *KindExpr
		Body      []*Stmt
	}

	VarDecl struct {
		Id    *Identifier
		Kind  *KindExpr
		Init  *Expr
		Const bool
		Pub   bool
	}

	BlockStmt struct {
		Body []*Stmt
	}

	ReturnStmt struct {
		Argument *Expr
	}

	ExprStmt struct {
		Expression *Expr
	}

	IfStmt struct {
		Condition  *Expr
		Consequent *Stmt
		Alternate  *Stmt
	}

	ForStmt struct {
		Label       *Identifier
		Init        *Stmt
		Test        *Expr
		Update      *Expr
		EachVisitor *EachVisitor
		Body        *Stmt
	}

	BreakStmt struct {
		Label *Identifier
	}

	ContinueStmt struct {
		Label *Identifier
	}
)

/* kind decl */
type (
	TypeAliasDecl struct {
		Name *KindIdentifier
		Kind *KindExpr
		Pub  bool
	}

	TypeInterfaceDecl struct {
		Name       *KindIdentifier
		Properties []*KindProperty
		Pub        bool
	}

	TypeStructDecl struct {
		Name *KindIdentifier
		Kind *KindExpr
		Pub  bool
	}

	TypeEnumDecl struct {
		Name    *KindIdentifier
		Choices []*KindIdentifier
		Pub     bool
	}
)
