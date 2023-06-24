package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/lexer"
	"strconv"
	"strings"
)

func newIdentifier(token *lexer.Token) *ast.Identifier {
	return &ast.Identifier{
		Name:     token.Value,
		Position: token.Position,
	}
}

func newKindIdentifier(token *lexer.Token) *ast.Identifier {
	return &ast.Identifier{
		Name:     token.Value,
		Position: token.Position,
	}
}

func newOperator(token *lexer.Token) *ast.Operator {
	return &ast.Operator{
		Value:    token.Value,
		Position: token.Position,
	}
}

func newIdentifierExpr(token *lexer.Token) *ast.Expr {
	return &ast.Expr{
		Node:     &ast.IdentifierLiteral{Name: newIdentifier(token)},
		Position: token.Position,
	}
}

func newNumberExpr(token *lexer.Token, parser *Parser) *ast.Expr {
	value, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		parser.UnexpectedPos(token.Start, err.Error())
	}
	return &ast.Expr{
		Node:     &ast.NumberLiteral{Value: value},
		Position: token.Position,
	}
}

func isUnsignedInt(value string) bool {
	first := value[0:1]
	return len(first) > 0 &&
		first[0] != '-' &&
		strings.IndexByte(value, '.') == -1
}

func getNumberExprValue(expr ast.Expr) float64 {
	switch node := expr.Node.(type) {
	case *ast.NumberLiteral:
		return node.Value
	default:
		panic("Internal Error")
	}
}
