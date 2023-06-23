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
		Node      E
		InferKind KE
		Position
	}

	KindExpr struct {
		Node KE
		Position
	}

	KindIdentifier struct {
		Name string
		Position
	}

	Identifier struct {
		Name string
		Position
	}

	KindProperty struct {
		Key  *Identifier
		Kind *KindExpr
		Position
	}

	Argument struct {
		Name *Identifier
		Kind *KindExpr
		Rest bool
		Position
	}

	StructProperty struct {
		Key   *Expr
		Value *Expr
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
