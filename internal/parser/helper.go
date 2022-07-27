package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/lexer"
	"strconv"
)

func NewIdentifier(token *lexer.Token) *ast.Identifier {
	return &ast.Identifier{
		Name:     token.Value,
		Position: token.Position,
	}
}

func NewKindIdentifier(token *lexer.Token) *ast.KindIdentifier {
	return &ast.KindIdentifier{
		Name:     token.Value,
		Position: token.Position,
	}
}

func NewIdentifierExpr(token *lexer.Token) *ast.Expression {
	return &ast.Expression{
		Node:     &ast.IdentifierLiteral{Name: token.Value},
		Position: token.Position,
	}
}

func NewNumberExpr(token *lexer.Token, parser *Parser) *ast.Expression {
	value, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		parser.unexpectedPos(token.Start, err.Error())
	}
	return &ast.Expression{
		Node:     &ast.NumberLiteral{Value: value},
		Position: token.Position,
	}
}
