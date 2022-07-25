package ast

type Expression struct {
	Node E
	Position
}

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

type (
	ImportSpecifier struct {
		Imported string
		Local    string
	}

	CallExpression struct {
		Callee    Expression
		Arguments []Expression
	}

	BinaryExpression struct {
		Left     Expression
		Right    Expression
		Operator string
	}

	UnaryExpression struct {
		Argument Expression
		Operator string
	}

	AssignmentExpression struct {
		Left     Expression
		Right    Expression
		Operator string
	}

	Identifier struct {
		Name  string
		Kind  Kind
		Refer bool
	}

	NumberLiteral struct {
		Value float64
	}

	BooleanLiteral struct {
		Value bool
	}

	StringLiteral struct {
		Value string
		Raw   bool // 原始字符串（多行）
	}
)
