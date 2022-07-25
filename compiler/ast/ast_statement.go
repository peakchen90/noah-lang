package ast

type Statement struct {
	Node S
	Position
}

func (*Statement) isNode() {}

// S statements
type S interface{ isStmt() }

func (*Program) isStmt()             {}
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
	Program struct {
		Body []Statement
	}

	ImportDeclaration struct {
		Source    string
		LocalName string
	}

	FunctionDeclaration struct {
		Id        string
		Arguments []Expression
		Body      []Statement
		Kind      Kind
		Pubic     bool
	}

	VariableDeclaration struct {
		Id    Expression
		Init  Expression
		Pubic bool
	}

	TypeDeclaration struct {
		Name  Expression
		Value Kind
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
		Label     string
		Init      Statement
		Condition Expression
		Update    Statement
	}

	ForOfStatement struct {
		Label     string
		IterIndex string
		IterName  string
		Target    Expression
	}

	BreakStatement struct {
		Label string
	}

	ContinueStatement struct {
		Label string
	}
)
