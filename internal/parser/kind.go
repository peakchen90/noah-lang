package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/lexer"
)

func (p *Parser) parseKindExpr() *ast.KindExpr {
	kindExpr := &ast.KindExpr{}

	if p.isToken(lexer.TTIdentifier) { // type refer
		token := p.current
		kindExpr.Position = token.Position
		p.nextToken()

		switch token.Value {
		case "number":
			kindExpr.Node = &ast.TNumber{}
		case "byte":
			kindExpr.Node = &ast.TByte{}
		case "char":
			kindExpr.Node = &ast.TChar{}
		case "string":
			kindExpr.Node = &ast.TString{}
		case "bool":
			kindExpr.Node = &ast.TBool{}
		case "any":
			kindExpr.Node = &ast.TAny{}
		case "self":
			kindExpr.Node = &ast.TSelf{}
		default:
			kindExpr.Node = &ast.TIdentifier{Name: newKindIdentifier(token)}
			return p.parseMaybeChainKindExpr(kindExpr)
		}
	} else if p.isToken(lexer.TTBracketL) {
		kindExpr.Start = p.current.Start
		p.nextToken()

		var Len *ast.Expr
		if p.isToken(lexer.TTNumber) { // [n]T
			Len = p.parseNumberExpr()
		}

		p.consume(lexer.TTBracketR, true)
		kind := p.parseKindExpr()
		kindExpr.Node = &ast.TArray{
			Kind: kind,
			Len:  Len,
		}
		kindExpr.End = kind.End
	} else if p.isKeyword("fn") { // fn(...params: []T) -> T
		start := p.current.Start
		p.nextToken()
		kindExpr = p.parseFuncKindExpr(start)
	} else if p.isKeyword("struct") { // struct{ a: number }
		start := p.current.Start
		p.nextToken()
		kindExpr = p.parseStructKindExpr(start)
	} else {
		p.unexpected()
	}

	return kindExpr
}

func (p *Parser) parseMaybeChainKindExpr(left *ast.KindExpr) *ast.KindExpr {
	if p.consume(lexer.TTDot, false) != nil {
		id := newKindIdentifier(p.consumeVarId(true))
		right := &ast.KindExpr{
			Node:     &ast.TIdentifier{Name: id},
			Position: id.Position,
		}
		nextLeft := &ast.KindExpr{
			Node:     &ast.TMemberKind{Left: left, Right: right},
			Position: *ast.NewPosition(left.Start, id.End),
		}
		return p.parseMaybeChainKindExpr(nextLeft)
	}

	return left
}

func (p *Parser) parseFuncKindExpr(start int) *ast.KindExpr {
	p.consume(lexer.TTParenL, true)
	kindExpr := &ast.KindExpr{}
	kindExpr.Start = start

	// params
	params := make([]*ast.Param, 0, helper.DefaultCap)
	var lastRestToken *lexer.Token
	for !p.isEnd() && !p.isToken(lexer.TTParenR) {
		if lastRestToken != nil {
			p.UnexpectedPos(lastRestToken.Start, "Only use `...` in the last parameter")
		}

		restToken := p.consume(lexer.TTRest, false)
		rest := restToken != nil
		if rest {
			lastRestToken = restToken
		}
		nameToken := p.consumeVarId(false)
		if nameToken == nil {
			p.unexpectedToken("identifier", p.current)
		}
		p.consume(lexer.TTColon, true)
		kind := p.parseKindExpr()
		param := &ast.Param{
			Name: newIdentifier(nameToken),
			Kind: kind,
			Rest: rest,
		}
		param.Start = nameToken.Start
		param.End = kind.End
		params = append(params, param)

		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}
	end := p.current.End
	p.consume(lexer.TTParenR, true)

	// return returnKind
	var returnKind *ast.KindExpr
	if p.consume(lexer.TTReturnSym, false) != nil {
		returnKind = p.parseKindExpr()
	}

	kindExpr.Node = &ast.TFuncKind{
		Params: params,
		Return: returnKind,
	}
	if returnKind != nil {
		end = returnKind.End
	}
	kindExpr.End = end

	return kindExpr
}

func (p *Parser) parseStructKindExpr(start int) *ast.KindExpr {
	extends := make([]*ast.KindExpr, 0, helper.SmallCap)
	if p.consume(lexer.TTExtendSym, false) != nil {
		for (p.isValidId() && !isReservedType(p.current.Value)) || p.isKeyword("struct") {
			extends = append(extends, p.parseKindExpr())
			if p.consume(lexer.TTComma, false) == nil {
				break
			}
		}
		if len(extends) == 0 {
			p.unexpected()
		}
	}

	// `{`
	p.consume(lexer.TTBraceL, true)
	properties := p.parseKindProperties(false)
	p.consume(lexer.TTBraceR, true)

	return &ast.KindExpr{
		Node: &ast.TStructKind{
			Extends:    extends,
			Properties: properties,
		},
		Position: *ast.NewPosition(start, p.lexer.LastToken.End),
	}
}

func (p *Parser) parseKindProperties(isFunc bool) []*ast.KindProperty {
	properties := make([]*ast.KindProperty, 0, helper.DefaultCap)

	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		pair := &ast.KindProperty{}

		if isFunc {
			start := p.current.Start
			p.consumeKeyword("fn", true)
			pair.Key = newIdentifier(p.consumeVarId(true))
			pair.Kind = p.parseFuncKindExpr(start)
		} else if !isFunc {
			token := p.consumeVarId(true)
			pair.Key = newIdentifier(token)
			p.consume(lexer.TTColon, true)
			pair.Kind = p.parseKindExpr()
		} else {
			p.unexpected()
		}

		pair.Start = pair.Key.Start
		pair.End = pair.Kind.End
		properties = append(properties, pair)

		tail := p.consume(lexer.TTComma, false)
		if tail == nil {
			tail = p.consume(lexer.TTSemi, false)
		}
		if tail == nil && !p.lexer.SeenNewline {
			break
		}
	}

	return properties
}

func (p *Parser) parseEnumItems() []*ast.Identifier {
	choices := make([]*ast.Identifier, 0, helper.DefaultCap)

	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		token := p.consumeVarId(true)
		choices = append(choices, newKindIdentifier(token))

		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}

	return choices
}
