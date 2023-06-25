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
	token := p.current

	if token.OpType.IsOpUnaryPrefix() && precedence < token.Precedence {
		p.nextToken()
		argument := p.parseMaybePostfixUnaryExpr(p.parseMaybeBinaryExpr(token.Precedence), token.Precedence)
		return &ast.Expr{
			Node: &ast.UnaryExpr{
				Argument: argument,
				Operator: newOperator(token),
				Prefix:   true,
			},
			Position: *ast.NewPosition(token.Start, argument.End),
		}
	}

	return p.parseMaybePostfixUnaryExpr(p.parseAtomExpr(), precedence)
}

// maybe postfix operator, only once
func (p *Parser) parseMaybePostfixUnaryExpr(left *ast.Expr, precedence int8) *ast.Expr {
	token := p.current

	if token.OpType.IsOpUnaryPostfix() && precedence < token.Precedence {
		p.nextToken()
		return &ast.Expr{
			Node: &ast.UnaryExpr{
				Argument: left,
				Operator: newOperator(token),
				Prefix:   false,
			},
			Position: *ast.NewPosition(left.Start, token.End),
		}
	}

	return left
}

func (p *Parser) parseBinaryExprPrecedence(left *ast.Expr, precedence int8) *ast.Expr {
	token := p.current

	if (token.OpType.IsOpBinaryLTR() && precedence < token.Precedence) || (token.OpType.IsOpBinaryRTL() && precedence <= token.Precedence) {
		nextPrecedence := token.Precedence
		operator := newOperator(token)
		p.nextToken()

		var node *ast.Expr

		if token.OpType.IsOpBinaryType() {
			right := p.parseKindExpr()
			node = &ast.Expr{
				Node: &ast.BinaryTypeExpr{
					Left:     left,
					Right:    right,
					Operator: operator,
				},
				Position: *ast.NewPosition(left.Start, right.End),
			}
		} else {
			// 解析可能更高优先级的右侧表达式，如: `1 + 2 * 3` 将解析 `2 * 3` 作为右值
			nextExpr := p.parseMaybeBinaryExpr(nextPrecedence)
			maybeHigherPrecedenceExpr := p.parseBinaryExprPrecedence(nextExpr, nextPrecedence)
			node = &ast.Expr{
				Node: &ast.BinaryExpr{
					Left:     left,
					Right:    maybeHigherPrecedenceExpr,
					Operator: operator,
				},
				Position: *ast.NewPosition(left.Start, maybeHigherPrecedenceExpr.End),
			}
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
		return p.parseMaybeChainExpr(p.parseFuncExpr(), AccessCall)
	} else if p.isToken(lexer.TTConst) {
		value := p.current.Value
		if value == "true" || value == "false" {
			return p.parseMaybeChainExpr(p.parseBooleanExpr(), AccessDot)
		} else if value == "null" {
			return p.parseNullExpr()
		}
	} else if p.consume(lexer.TTIdentifier, false) != nil {
		return p.parseMaybeChainExpr(newIdentifierExpr(p.lexer.LastToken), AccessDot|AccessComputed|AccessCall|AccessStruct)
	}

	switch p.current.Type {
	case lexer.TTString:
		return p.parseMaybeChainExpr(p.parseStringExpr(), AccessDot|AccessComputed)
	case lexer.TTChar:
		return p.parseMaybeChainExpr(p.parseCharExpr(), AccessDot)
	case lexer.TTNumber:
		return p.parseMaybeChainExpr(p.parseNumberExpr(), AccessDot)
	case lexer.TTBraceL:
		return p.parseMaybeChainExpr(p.parseStructExpr(nil), AccessDot)
	case lexer.TTBracketL:
		return p.parseMaybeChainExpr(p.parseArrayExpr(), AccessComputed)
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
		Position: *ast.NewPosition(start, body.End),
	}
}

func (p *Parser) parseStructExpr(ctor *ast.Expr) *ast.Expr {
	ctorKind := exprToKindExpr(ctor)
	properties := make([]*ast.ValueProperty, 0, helper.DefaultCap)
	start := p.current.Start
	if ctorKind != nil {
		start = ctorKind.Start
	}

	p.consume(lexer.TTBraceL, true) // `{`

	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		p.expect(lexer.TTIdentifier)
		name := p.parseIdentifierExpr(nil)

		p.consume(lexer.TTColon, true)
		value := p.parseExpr()

		properties = append(properties, &ast.ValueProperty{
			Key:   name,
			Value: value,
		})

		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}

	p.consume(lexer.TTBraceR, true)

	return &ast.Expr{
		Node:     &ast.StructExpr{Ctor: ctorKind, Properties: properties},
		Position: *ast.NewPosition(start, p.lexer.LastToken.End),
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
		Position: *ast.NewPosition(start, p.lexer.LastToken.End),
	}
}

func (p *Parser) parseBooleanExpr() *ast.Expr {
	text := p.current.Value
	expr := ast.Expr{
		Node: &ast.BoolLiteral{
			Value: text == "true", // `true`
			Text:  text,
		},
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

func (p *Parser) parseStringExpr() *ast.Expr {
	expr := ast.Expr{
		Node:     &ast.StringLiteral{Value: p.current.Value},
		Position: p.current.Position,
	}
	p.nextToken()
	return &expr
}

func (p *Parser) parseCharExpr() *ast.Expr {
	text := p.current.Value
	expr := ast.Expr{
		Node: &ast.CharLiteral{
			Value: rune(text[0]),
			Text:  text,
		},
		Position: p.current.Position,
	}
	p.nextToken()
	return &expr
}

func (p *Parser) parseNumberExpr() *ast.Expr {
	text := p.current.Value
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		p.unexpected()
	}

	expr := ast.Expr{
		Node: &ast.NumberLiteral{
			Value: value,
			Text:  text,
		},
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
		Node:     newIdentifierExpr(token).Node,
		Position: token.Position,
	}
	p.nextToken()
	return &expr
}

// 解析可能的链式调用表达式，如：`expr.b.c`, `expr[n]`, `expr()`
func (p *Parser) parseMaybeChainExpr(parent *ast.Expr, access AccessType) *ast.Expr {
	if (access&AccessDot > 0) && p.isToken(lexer.TTDot) { // `.`
		p.nextToken()
		property := newIdentifierExpr(p.consume(lexer.TTIdentifier, true))
		memberExpr := &ast.Expr{
			Node: &ast.MemberExpr{
				Object:   parent,
				Property: property,
				Computed: false,
			},
			Position: *ast.NewPosition(parent.Start, property.End),
		}
		return p.parseMaybeChainExpr(memberExpr, AccessDot|AccessComputed|AccessCall|AccessStruct)
	} else if (access&AccessComputed > 0) && p.isToken(lexer.TTBracketL) { // `[`
		p.nextToken()
		property := p.parseExpr()
		p.consume(lexer.TTBracketR, true)

		computedMemberExpr := &ast.Expr{
			Node: &ast.MemberExpr{
				Object:   parent,
				Property: property,
				Computed: true,
			},
			Position: *ast.NewPosition(parent.Start, p.lexer.LastToken.End),
		}
		return p.parseMaybeChainExpr(computedMemberExpr, AccessDot|AccessComputed|AccessCall)
	} else if (access&AccessCall > 0) && p.isToken(lexer.TTParenL) { // `(`
		callExpr := p.parseCallExpr(parent)
		return p.parseMaybeChainExpr(callExpr, AccessDot|AccessComputed|AccessCall)
	} else if (access&AccessStruct > 0) && p.isToken(lexer.TTBraceL) { // `{`
		structExpr := p.parseStructExpr(parent)
		return p.parseMaybeChainExpr(structExpr, AccessDot)
	}

	return parent
}

func (p *Parser) parseCallExpr(callee *ast.Expr) *ast.Expr {
	p.nextToken()
	params := make([]*ast.Expr, 0, helper.DefaultCap)

	for !p.isToken(lexer.TTParenR) {
		params = append(params, p.parseExpr())
		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}

	p.consume(lexer.TTParenR, true)

	return &ast.Expr{
		Node: &ast.CallExpr{
			Callee: callee,
			Params: params,
		},
		Position: *ast.NewPosition(
			callee.Start,
			p.lexer.LastToken.End,
		),
	}
}
