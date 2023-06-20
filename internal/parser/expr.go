package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/lexer"
	"strconv"
)

func (p *Parser) parseExpr() *ast.Expr {
	return p.parseMaybeBinaryExpr(-1)
}

func (p *Parser) parseMaybeBinaryExpr(precedence int8) *ast.Expr {
	if p.isToken(lexer.TTParenL) { // `(`
		p.nextToken()
		expr := p.parseExpr()
		p.consume(lexer.TTParenR, true)
		return p.parseBinaryExprPrecedence(expr, precedence)
	}

	return p.parseBinaryExprPrecedence(p.parseMaybeUnaryExpr(precedence), precedence)
}

func (p *Parser) parseMaybeUnaryExpr(precedence int8) *ast.Expr {
	if precedence < p.current.Precedence {
		if p.isToken(lexer.TTBitNot) || p.isToken(lexer.TTLogicNot) {
			operator := p.current.Value
			start := p.current.Start
			argument := p.parseExpr()
			return &ast.Expr{
				Node: &ast.UnaryExpr{
					Argument: argument,
					Operator: operator,
					Prefix:   true,
				},
				Position: ast.Position{Start: start, End: argument.End},
			}
		}
	}

	return p.parseAtomExpr()
}

func (p *Parser) parseBinaryExprPrecedence(left *ast.Expr, precedence int8) *ast.Expr {
	// 当前 token 拥有更高优先级
	if precedence < p.current.Precedence {
		// TODO validate operator

		nextPrecedence := p.current.Precedence
		operator := p.current.Value
		p.nextToken()

		// 解析可能更高优先级的右侧表达式，如: `1 + 2 * 3` 将解析 `2 * 3` 作为右值
		maybeHigherPrecedenceExpr := p.parseMaybeBinaryExpr(nextPrecedence)
		right := p.parseBinaryExprPrecedence(maybeHigherPrecedenceExpr, nextPrecedence)

		node := &ast.Expr{
			Node: &ast.BinaryExpr{
				Left:     left,
				Right:    right,
				Operator: operator,
			},
			Position: ast.Position{Start: left.Start, End: right.End},
		}

		// 将已经解析的二元表达式作为左值，然后递归解析后面可能的同等优先级或低优先级的表达式作为右值
		// 如: `1 + 2 + 3`, 当前已经解析 `1 + 2`, 然后将该节点作为左值递归解析表达式优先级
		return p.parseBinaryExprPrecedence(node, precedence)
	}

	return left
}

// 解析一个原子表达式，如: `foo()`, `3.14`, `foo`, `var2 = expr`, `true`, `"str"`, `fn() {}`, `A{}`
func (p *Parser) parseAtomExpr() *ast.Expr {
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
		left := p.parseMaybeMemberExpr()

		switch p.current.Type {
		case lexer.TTParenL:
			return p.parseCallExpr(left)
		case lexer.TTAssign:
			return p.parseAssignExpr(left)
		case lexer.TTBraceL:
			return p.parseStructExpr(left)
		default:
			return left
		}
	}

	switch p.current.Type {
	case lexer.TTString:
		return p.parseStringExpr()
	case lexer.TTNumber:
		return p.parseNumberExpr()
	case lexer.TTBraceL:
		return p.parseStructExpr(nil)
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

func (p *Parser) parseStructExpr(ctor *ast.Expr) *ast.Expr {
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
		Node:     &ast.StructExpr{Ctor: ctor, Properties: properties},
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
		Node:     &ast.BooleanLiteral{Value: len(p.current.Value) == 4},
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

// a.b.c.d
func (p *Parser) parseMaybeMemberExpr() *ast.Expr {
	object := NewIdentifierExpr(p.consume(lexer.TTIdentifier, true))

	for p.consume(lexer.TTDot, false) != nil {
		property := NewIdentifierExpr(p.consume(lexer.TTIdentifier, true))
		object = &ast.Expr{
			Node: &ast.MemberExpr{
				Object:   object,
				Property: property,
			},
			Position: ast.Position{
				Start: object.Start,
				End:   property.End,
			},
		}
	}

	return object
}

func (p *Parser) parseCallExpr(callee *ast.Expr) *ast.Expr {
	// TODO
	expr := ast.Expr{}
	return &expr
}

func (p *Parser) parseAssignExpr(left *ast.Expr) *ast.Expr {
	// TODO
	expr := ast.Expr{}
	return &expr
}
