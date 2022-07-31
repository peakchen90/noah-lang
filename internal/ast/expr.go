package ast

// E expression
type E interface{ isExpr() }

func (*CallExpr) isExpr()          {}
func (*MemberExpr) isExpr()        {}
func (*BinaryExpr) isExpr()        {}
func (*UnaryExpr) isExpr()         {}
func (*AssignmentExpr) isExpr()    {}
func (*IdentifierLiteral) isExpr() {}
func (*NumberLiteral) isExpr()     {}
func (*BooleanLiteral) isExpr()    {}
func (*NullLiteral) isExpr()       {}
func (*StringLiteral) isExpr()     {}

type (
	FuncExpr struct {
		FuncSign *KindExpr
		Body     *Stmt
	}

	StructExpr struct {
		Properties []*KindProperty
	}

	CallExpr struct {
		Callee    *Expr
		Arguments []*Expr
	}

	MemberExpr struct {
		Object   *Expr
		Property *Expr
		Computed bool
	}

	BinaryExpr struct {
		Left     *Expr
		Right    *Expr
		Operator string
	}

	UnaryExpr struct {
		Argument *Expr
		Operator string
		Prefix   bool
	}

	AssignmentExpr struct {
		Left     *Expr
		Right    *Expr
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
	}
)
