package ast

type (
	File struct {
		Body []*Stmt
		Position
	}

	Stmt struct {
		Node S
		Position
	}

	Expr struct {
		Node E
		Position
	}

	KindExpr struct {
		Node KE
		Position
	}

	Identifier struct {
		Name string
		Position
	}

	Operator struct {
		Value string
		Position
	}

	KindProperty struct {
		Key  *Identifier
		Kind *KindExpr
		Position
	}

	ValueProperty struct {
		Key   *Expr
		Value *Expr
	}

	Param struct {
		Name *Identifier
		Kind *KindExpr
		Rest bool
		Position
	}

	EachVisitor struct {
		Key    *Identifier
		Value  *Identifier
		Target *Expr
	}
)

type Position struct {
	Start int
	End   int
}

func NewPosition(start int, end int) *Position {
	return &Position{
		Start: start,
		End:   end,
	}
}
