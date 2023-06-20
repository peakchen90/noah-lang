package ast

// E expression
type E interface{ isExpr() }

func (*CallExpr) isExpr()          {}
func (*MemberExpr) isExpr()        {}
func (*BinaryExpr) isExpr()        {}
func (*UnaryExpr) isExpr()         {}
func (*FuncExpr) isExpr()          {}
func (*StructExpr) isExpr()        {}
func (*ArrayExpr) isExpr()         {}
func (*IdentifierLiteral) isExpr() {}
func (*NumberLiteral) isExpr()     {}
func (*BooleanLiteral) isExpr()    {}
func (*NullLiteral) isExpr()       {}
func (*SelfLiteral) isExpr()       {}
func (*StringLiteral) isExpr()     {}
func (*CharLiteral) isExpr()       {}

// expr
type (
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

	FuncExpr struct {
		FuncKind *KindExpr
		Body     *Stmt
	}

	StructExpr struct {
		Ctor       *Expr
		Properties []*StructProperty
	}

	ArrayExpr struct {
		Items []*Expr
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

	SelfLiteral struct {
	}

	StringLiteral struct {
		Value string
	}

	CharLiteral struct {
		Value rune
	}
)
