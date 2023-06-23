package ast

// S statements
type S interface{ isStmt() }

func (*ImportDecl) isStmt()   {}
func (*FuncDecl) isStmt()     {}
func (*ImplDecl) isStmt()     {}
func (*VarDecl) isStmt()      {}
func (*BlockStmt) isStmt()    {}
func (*ReturnStmt) isStmt()   {}
func (*ExprStmt) isStmt()     {}
func (*IfStmt) isStmt()       {}
func (*ForStmt) isStmt()      {}
func (*BreakStmt) isStmt()    {}
func (*ContinueStmt) isStmt() {}

func (*TAliasDecl) isStmt()     {}
func (*TInterfaceDecl) isStmt() {}
func (*TStructDecl) isStmt()    {}
func (*TEnumDecl) isStmt()      {}

/* statements */
type (
	ImportDecl struct {
		Package *Identifier
		Paths   []*Identifier
		Local   *Identifier
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
	TAliasDecl struct {
		Name *Identifier
		Kind *KindExpr
		Pub  bool
	}

	TInterfaceDecl struct {
		Name       *Identifier
		Properties []*KindProperty
		Pub        bool
	}

	TStructDecl struct {
		Name *Identifier
		Kind *KindExpr
		Pub  bool
	}

	TEnumDecl struct {
		Name    *Identifier
		Choices []*Identifier
		Pub     bool
	}
)
