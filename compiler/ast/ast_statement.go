package ast

type Statement struct {
	Data S
	Position
}

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

type Program struct {
	Body []Statement
}

type ImportDeclaration struct {
	Source     string
	specifiers []Expression
}

type FunctionDeclaration struct {
	Id        string
	Arguments []Expression
	Body      []Statement
	Kind      KindMeta
	Pubic     bool
}

type VariableDeclaration struct {
	Id    Expression
	Init  Expression
	Pubic bool
}

type TypeDeclaration struct {
	Name  Expression
	Kind  KindMeta
	Pubic bool
}

type BlockStatement struct {
	Body []Statement
}

type ReturnStatement struct {
	Argument Expression
}

type ExpressionStatement struct {
	Expression Expression
}

type IfStatement struct {
	Condition  Expression
	Consequent Statement
	Alternate  Statement
}

type ForStatement struct {
	Label     string
	Init      Statement
	Condition Expression
	Update    Statement
}

type ForOfStatement struct {
	Label     string
	IterIndex string
	IterName  string
	Target    Expression
}

type BreakStatement struct {
	Label string
}

type ContinueStatement struct {
	Label string
}
