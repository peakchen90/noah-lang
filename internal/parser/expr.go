package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/lexer"
	"strconv"
)

func (p *Parser) parseExpr() *ast.Expr {
	if p.isKeyword("fn") {
		return p.parseFuncExpr()
	} else if p.isToken(lexer.TTConst) {
		value := p.current.Value
		if value == "true" || value == "false" {
			return p.parseBooleanExpr()
		} else if value == "null" {
			return p.parseNullExpr()
		} else if value == "self" {
			return p.parseSelfExpr()
		}
	} else if p.isToken(lexer.TTIdentifier) {
		return p.parseMaybeMemberExpr()
	}

	switch p.current.Type {
	case lexer.TTString:
		return p.parseStringExpr()
	case lexer.TTNumber:
		return p.parseNumberExpr()
	case lexer.TTBraceL:
		return p.parseStructExpr()
	case lexer.TTBracketL:
		return p.parseArrayExpr()
	default:
		p.unexpected()
	}

	return nil // NEVER
}

func (p *Parser) parseFuncExpr() *ast.Expr {
	expr := ast.Expr{}

	return &expr
}

func (p *Parser) parseStructExpr() *ast.Expr {
	kind := ast.KindExpr{}
	properties := make([]*ast.StructProperty, 0, helper.DefaultCap)
	start := p.current.Start
	p.nextToken() // skip `{`

	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		p.expect(lexer.TTIdentifier)
		name := p.parseIdentifierExpr(nil)

		p.consume(lexer.TTColon, true)
		value := p.parseExpr()

		properties = append(properties, &ast.StructProperty{
			Name:  name,
			Value: value,
		})

		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}

	if !p.isToken(lexer.TTBraceR) {
		p.unexpected()
	}

	expr := ast.Expr{
		Node:     &ast.StructExpr{Kind: &kind, Properties: properties},
		Position: ast.Position{Start: start, End: p.current.End},
	}
	p.nextToken() // skip `}`
	return &expr
}

func (p *Parser) parseArrayExpr() *ast.Expr {
	items := make([]*ast.Expr, 0, helper.DefaultCap)
	start := p.current.Start
	p.nextToken() // skip `[`

	for !p.isEnd() && !p.isToken(lexer.TTBracketR) {
		items = append(items, p.parseExpr())
		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}

	if !p.isToken(lexer.TTBracketR) {
		p.unexpected()
	}

	expr := ast.Expr{
		Node:     &ast.ArrayExpr{Items: items},
		Position: ast.Position{Start: start, End: p.current.End},
	}
	p.nextToken() // skip `]`
	return &expr
}

func (p *Parser) parseBooleanExpr() *ast.Expr {
	expr := ast.Expr{
		Node:     &ast.BooleanLiteral{Value: p.current.Value == "true"},
		Position: p.current.Position,
	}
	p.nextToken()
	return &expr
}

func (p *Parser) parseNullExpr() *ast.Expr {
	expr := ast.Expr{
		Node:     &ast.NullLiteral{},
		Position: p.current.Position,
	}
	p.nextToken()
	return &expr
}

func (p *Parser) parseSelfExpr() *ast.Expr {
	expr := ast.Expr{
		Node:     &ast.SelfLiteral{},
		Position: p.current.Position,
	}
	p.nextToken()
	return &expr
}

func (p *Parser) parseStringExpr() *ast.Expr {
	expr := ast.Expr{
		Node:     &ast.StringLiteral{Value: p.current.Value},
		Position: p.current.Position,
	}
	p.nextToken()
	return &expr
}

func (p *Parser) parseNumberExpr() *ast.Expr {
	value, err := strconv.ParseFloat(p.current.Value, 64)
	if err != nil {
		p.unexpected()
	}

	expr := ast.Expr{
		Node:     &ast.NumberLiteral{Value: value},
		Position: p.current.Position,
	}
	p.nextToken()
	return &expr
}

func (p *Parser) parseIdentifierExpr(token *lexer.Token) *ast.Expr {
	if token == nil {
		token = p.current
	}
	expr := ast.Expr{
		Node:     &ast.IdentifierLiteral{Name: token.Value},
		Position: token.Position,
	}
	p.nextToken()
	return &expr
}

func (p *Parser) parseMaybeMemberExpr() *ast.Expr {
	expr := ast.Expr{
		Node:     &ast.NumberLiteral{},
		Position: p.current.Position,
	}
	p.nextToken()
	return &expr
}
