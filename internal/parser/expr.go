package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/lexer"
)

func (p *Parser) parseExpr() *ast.Expr {
	if p.isKeyword("fn") {
		return p.parseFuncExpr()
	} else if p.isKeyword("struct") {
		return p.parseStructExpr()
	}

	switch p.current.Type {
	case lexer.TTString:

	default:

		p.unexpected()
	}

	return nil
}

func (p *Parser) parseFuncExpr() *ast.Expr {
	expr := ast.Expr{}

	return &expr
}

func (p *Parser) parseStructExpr() *ast.Expr {
	expr := ast.Expr{}

	return &expr
}
