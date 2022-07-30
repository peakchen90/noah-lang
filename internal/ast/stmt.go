package ast

// S statements
type S interface{ isStmt() }

func (*ImportDeclaration) isStmt()   {}
func (*FunctionDeclaration) isStmt() {}
func (*VariableDeclaration) isStmt() {}
func (*BlockStatement) isStmt()      {}
func (*ReturnStatement) isStmt()     {}
func (*ExpressionStatement) isStmt() {}
func (*IfStatement) isStmt()         {}
func (*ForStatement) isStmt()        {}
func (*ForOfStatement) isStmt()      {}
func (*BreakStatement) isStmt()      {}
func (*ContinueStatement) isStmt()   {}

type (
	ImportDeclaration struct {
		Source string
		Local  *Identifier
	}

	FunctionDeclaration struct {
		Name     *Identifier
		Impl     *KindIdentifier
		FuncSign *KindExpr
		Body     *Statement
		Pubic    bool
	}

	VariableDeclaration struct {
		Id    *Identifier
		Init  *Expression
		Const bool
		Pubic bool
	}

	BlockStatement struct {
		Body []*Statement
	}

	ReturnStatement struct {
		Argument *Expression
	}

	ExpressionStatement struct {
		Expression *Expression
	}

	IfStatement struct {
		Condition  *Expression
		Consequent *Statement
		Alternate  *Statement
	}

	ForStatement struct {
		Label     *Identifier
		Init      *Statement
		Condition *Expression
		Update    *Statement
	}

	ForOfStatement struct {
		Label     *Identifier
		IterIndex *Identifier
		IterName  *Identifier
		Target    *Expression
	}

	BreakStatement struct {
		Label *Identifier
	}

	ContinueStatement struct {
		Label *Identifier
	}
)
