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

	left := p.parseMaybeUnaryExpr(precedence)
	return p.parseBinaryExprPrecedence(left, precedence)
}

func (p *Parser) parseMaybeUnaryExpr(precedence int8) *ast.Expr {
	if precedence < p.current.Precedence {
		switch p.current.Type {
		case lexer.TTBitNot, lexer.TTLogicNot, lexer.TTSub, lexer.TTPlus:
			operator := p.current.Text
			start := p.current.Start
			p.nextToken()
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
	if precedence < p.current.Precedence || (p.isAssignToken() && precedence == p.current.Precedence) {
		nextPrecedence := p.current.Precedence
		operator := p.current.Text
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

// 解析一个原子表达式，如: `foo()`, `3.14`, `a.b`, `var2 = expr`, `true`, `"str"`, `fn() {}`, `A{}`
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
	} else if p.consume(lexer.TTIdentifier, false) != nil {
		return p.parseMaybeChainExpr(NewIdentifierExpr(p.lexer.LastToken))
	}

	switch p.current.Type {
	case lexer.TTString:
		return p.parseStringExpr()
	case lexer.TTChar:
		return p.parseCharExpr()
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
	start := p.current.Start
	p.nextToken()
	funcKind := p.parseFuncKindExpr(start)
	body := p.parseBlockStmt()

	return &ast.Expr{
		Node:     &ast.FuncExpr{FuncKind: funcKind, Body: body},
		Position: ast.Position{Start: start, End: body.End},
	}
}

func (p *Parser) parseStructExpr(ctor *ast.Expr) *ast.Expr {
	properties := make([]*ast.StructProperty, 0, helper.DefaultCap)
	start := p.current.Start
	if ctor != nil {
		start = ctor.Start
	}

	p.consume(lexer.TTBraceL, true) // `{`

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

	p.consume(lexer.TTBraceR, true)

	return &ast.Expr{
		Node:     &ast.StructExpr{Ctor: ctor, Properties: properties},
		Position: ast.Position{Start: start, End: p.lexer.LastToken.End},
	}
}

func (p *Parser) parseArrayExpr() *ast.Expr {
	items := make([]*ast.Expr, 0, helper.DefaultCap)
	start := p.current.Start
	p.consume(lexer.TTBracketL, true) // `[`

	for !p.isEnd() && !p.isToken(lexer.TTBracketR) {
		items = append(items, p.parseExpr())
		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}

	p.consume(lexer.TTBracketR, true) // `]`

	return &ast.Expr{
		Node:     &ast.ArrayExpr{Items: items},
		Position: ast.Position{Start: start, End: p.lexer.LastToken.End},
	}
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

func (p *Parser) parseCharExpr() *ast.Expr {
	expr := ast.Expr{
		Node:     &ast.CharLiteral{Value: rune(p.current.Value[0])},
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

// 解析可能的链式调用表达式，如：`expr.b.c`, `expr[n]`, `expr()`
func (p *Parser) parseMaybeChainExpr(parent *ast.Expr) *ast.Expr {
	if p.isToken(lexer.TTDot) { // `.`
		p.nextToken()
		property := NewIdentifierExpr(p.consume(lexer.TTIdentifier, true))
		memberExpr := &ast.Expr{
			Node: &ast.MemberExpr{
				Object:   parent,
				Property: property,
				Computed: false,
			},
			Position: ast.Position{
				Start: parent.Start,
				End:   property.End,
			},
		}
		return p.parseMaybeChainExpr(memberExpr)
	} else if p.isToken(lexer.TTBracketL) { // `[`
		p.nextToken()
		property := p.parseExpr()
		p.consume(lexer.TTBracketR, true)

		computedMemberExpr := &ast.Expr{
			Node: &ast.MemberExpr{
				Object:   parent,
				Property: property,
				Computed: true,
			},
			Position: ast.Position{
				Start: parent.Start,
				End:   p.lexer.LastToken.End,
			},
		}
		return p.parseMaybeChainExpr(computedMemberExpr)
	} else if p.isToken(lexer.TTParenL) { // `(`
		callExpr := p.parseCallExpr(parent)
		return p.parseMaybeChainExpr(callExpr)
	} else if p.isToken(lexer.TTBraceL) { // `{`
		structExpr := p.parseStructExpr(parent)
		return p.parseMaybeChainExpr(structExpr)
	}

	return parent
}

func (p *Parser) parseCallExpr(callee *ast.Expr) *ast.Expr {
	p.nextToken()
	arguments := make([]*ast.Expr, 0, helper.DefaultCap)

	for !p.isToken(lexer.TTParenR) {
		arguments = append(arguments, p.parseExpr())
		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}

	p.consume(lexer.TTParenR, true)

	return &ast.Expr{
		Node: &ast.CallExpr{
			Callee:    callee,
			Arguments: arguments,
		},
		Position: ast.Position{
			Start: callee.Start,
			End:   p.lexer.LastToken.End,
		},
	}
}
