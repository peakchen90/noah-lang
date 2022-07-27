package parser

import (
	"github.com/peakchen90/hera-lang/internal/ast"
	"github.com/peakchen90/hera-lang/internal/lexer"
)

func NewIdentifier(token *lexer.Token, kind *ast.KindExpr) *ast.Identifier {
	id := ast.Identifier{
		Name: token.Value,
		Kind: *kind,
	}
	id.Position = token.Position

	return &id
}

func NewKindIdentifier(token *lexer.Token) *ast.KindIdentifier {
	kindId := ast.KindIdentifier{Name: token.Value}
	kindId.Position = token.Position

	return &kindId
}
