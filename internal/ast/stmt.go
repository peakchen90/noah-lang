package ast

// S statements
type S interface{ isStmt() }

func (*ImportDecl) isStmt()   {}
func (*FuncDecl) isStmt()     {}
func (*VarDecl) isStmt()      {}
func (*BlockStmt) isStmt()    {}
func (*ReturnStmt) isStmt()   {}
func (*ExprStmt) isStmt()     {}
func (*IfStmt) isStmt()       {}
func (*ForStmt) isStmt()      {}
func (*ForOfStmt) isStmt()    {}
func (*BreakStmt) isStmt()    {}
func (*ContinueStmt) isStmt() {}

type (
	ImportDecl struct {
		Source string
		Local  *Identifier
	}

	FuncDecl struct {
		Name     *Identifier
		Impl     *KindExpr
		FuncKind *KindExpr
		Body     *Stmt
		Pubic    bool
	}

	VarDecl struct {
		Id    *Identifier
		Kind  *KindExpr
		Init  *Expr
		Const bool
		Pubic bool
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

	ForOfStmt struct {
		Label     *Identifier
		IterIndex *Identifier
		IterName  *Identifier
		Target    *Expr
	}

	BreakStmt struct {
		Label *Identifier
	}

	ContinueStmt struct {
		Label *Identifier
	}
)
