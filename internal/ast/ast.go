package ast

type Position struct {
	Start int
	End   int
}

type Identifier struct {
	Name string
	Position
}

type File struct {
	Body []*Stmt
	Position
}

type Stmt struct {
	Node S
	Position
}

type Expr struct {
	Node E
	Position
}

type KindExpr struct {
	Node KE
	Position
}

type KindIdentifier struct {
	Name string
	Position
}

type KindProperty struct {
	Name *Identifier
	Kind *KindExpr
	Position
}

type Argument struct {
	Name *Identifier
	Kind *KindExpr
	Rest bool
	Position
}

type EachVisitor struct {
	Value  *Identifier
	Key    *Identifier
	Target *Expr
}
