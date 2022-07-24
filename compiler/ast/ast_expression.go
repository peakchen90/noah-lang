package ast

// E expression
type E interface{ isExpr() }

func (*ImportSpecifier) isExpr()      {}
func (*CallExpression) isExpr()       {}
func (*BinaryExpression) isExpr()     {}
func (*UnaryExpression) isExpr()      {}
func (*AssignmentExpression) isExpr() {}
func (*Identifier) isExpr()           {}
func (*NumberLiteral) isExpr()        {}
func (*BooleanLiteral) isExpr()       {}
func (*StringLiteral) isExpr()        {}

type Expression struct {
	Data E
	Position
}

type ImportSpecifier struct {
	Imported string
	Local    string
}

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator string
}

type UnaryExpression struct {
	Argument Expression
	Operator string
}

type AssignmentExpression struct {
	Left     Expression
	Right    Expression
	Operator string
}

type Identifier struct {
	Name string
	Kind
	Refer bool
}

type NumberLiteral struct {
	Value float64
}

type BooleanLiteral struct {
	Value bool
}

type StringLiteral struct {
	Value string
	Raw   bool // 多行字符串
}
