package ast

// E expression
type E interface{ isExpr() }

func (*CallExpression) isExpr()       {}
func (*MemberExpression) isExpr()     {}
func (*BinaryExpression) isExpr()     {}
func (*UnaryExpression) isExpr()      {}
func (*AssignmentExpression) isExpr() {}
func (*IdentifierLiteral) isExpr()    {}
func (*NumberLiteral) isExpr()        {}
func (*BooleanLiteral) isExpr()       {}
func (*NullLiteral) isExpr()          {}
func (*StringLiteral) isExpr()        {}

type (
	CallExpression struct {
		Callee    Expression
		Arguments []Expression
	}

	MemberExpression struct {
		Object   Expression
		Property Expression
		Computed bool
	}

	BinaryExpression struct {
		Left     Expression
		Right    Expression
		Operator string
	}

	UnaryExpression struct {
		Argument Expression
		Operator string
		Prefix   bool
	}

	AssignmentExpression struct {
		Left     Expression
		Right    Expression
		Operator string
	}

	IdentifierLiteral struct {
		Name string
	}

	NumberLiteral struct {
		Value float64
	}

	BooleanLiteral struct {
		Value bool
	}

	NullLiteral struct {
	}

	StringLiteral struct {
		Value string
		Raw   bool // 原始字符串（多行）
	}
)
