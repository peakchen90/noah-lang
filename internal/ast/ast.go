package ast

type Position struct {
	Start int
	End   int
}

type Label struct {
	Name string
	Position
}

type Identifier struct {
	Name string
	Kind KindExpr
	Position
}

type ConstInt struct {
	Value int
	Position
}

type File struct {
	Body []Statement
	Position
}

type Statement struct {
	Node S
	Position
}

type Expression struct {
	Node E
	Position
}

type KindDecl struct {
	Node KD
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
	Name KindIdentifier
	Kind KindExpr
	Position
}

type Argument struct {
	Identifier
	Rest bool
}