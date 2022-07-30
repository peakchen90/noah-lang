package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/lexer"
	"strconv"
	"strings"
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

func NewIdentifierExpr(token *lexer.Token) *ast.Expr {
	return &ast.Expr{
		Node:     &ast.IdentifierLiteral{Name: token.Value},
		Position: token.Position,
	}
}

func NewNumberExpr(token *lexer.Token, parser *Parser) *ast.Expr {
	value, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		parser.unexpectedPos(token.Start, err.Error())
	}
	return &ast.Expr{
		Node:     &ast.NumberLiteral{Value: value},
		Position: token.Position,
	}
}

func IsUnsignedInt(value string) bool {
	first := value[0:1]
	return len(first) > 0 &&
		first[0] != '-' &&
		strings.IndexByte(value, '.') == -1
}

func GetNumberExprValue(expr ast.Expr) float64 {
	switch node := expr.Node.(type) {
	case *ast.NumberLiteral:
		return node.Value
	default:
		panic("Internal Error")
	}
}
