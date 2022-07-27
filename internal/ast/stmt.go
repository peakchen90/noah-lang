package ast

func (*Statement) isNode() {}

// S statements
type S interface{ isStmt() }

func (*ImportDeclaration) isStmt()   {}
func (*FunctionDeclaration) isStmt() {}
func (*VariableDeclaration) isStmt() {}
func (*TypeDeclaration) isStmt()     {}
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
		Local  Identifier
	}

	FunctionDeclaration struct {
		TypeFunction
		Impl  KindIdentifier
		Body  []Statement
		Pubic bool
	}

	VariableDeclaration struct {
		Id    Identifier
		Init  Expression
		Const bool
		Pubic bool
	}

	TypeDeclaration struct {
		Decl  KindDecl
		Pubic bool
	}

	BlockStatement struct {
		Body []Statement
	}

	ReturnStatement struct {
		Argument Expression
	}

	ExpressionStatement struct {
		Expression Expression
	}

	IfStatement struct {
		Condition  Expression
		Consequent Statement
		Alternate  Statement
	}

	ForStatement struct {
		Label     Label
		Init      Statement
		Condition Expression
		Update    Statement
	}

	ForOfStatement struct {
		Label     Label
		IterIndex Identifier
		IterName  Identifier
		Target    Expression
	}

	BreakStatement struct {
		Label Label
	}

	ContinueStatement struct {
		Label Label
	}
)
