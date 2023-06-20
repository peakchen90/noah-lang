package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/lexer"
)

func (p *Parser) parseKindExpr() *ast.KindExpr {
	kindExpr := &ast.KindExpr{}

	if p.isToken(lexer.TTIdentifier) { // type alias
		token := p.current
		kindExpr.Position = token.Position
		p.nextToken()

		switch token.Value {
		case "num":
			kindExpr.Node = &ast.TypeNumber{}
		case "byte":
			kindExpr.Node = &ast.TypeByte{}
		case "char":
			kindExpr.Node = &ast.TypeChar{}
		case "str":
			kindExpr.Node = &ast.TypeString{}
		case "bool":
			kindExpr.Node = &ast.TypeBool{}
		case "any":
			kindExpr.Node = &ast.TypeAny{}
		default:
			kindExpr.Node = &ast.TypeId{Name: NewKindIdentifier(token)}
			return p.parseChainKindExpr(kindExpr)
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
		kindExpr.Node = &ast.TypeArray{
			Kind: kind,
			Len:  Len,
		}
		kindExpr.End = kind.End
	} else if p.isKeyword("fn") { // fn(...arg: []T) -> T
		kindExpr = p.parseFuncKindExpr(p.current.Start)
	} else {
		p.unexpected()
	}

	return kindExpr
}

func (p *Parser) parseChainKindExpr(parent *ast.KindExpr) *ast.KindExpr {
	if p.consume(lexer.TTDot, false) != nil {
		id := NewKindIdentifier(p.consume(lexer.TTIdentifier, true))
		idKindExpr := &ast.KindExpr{
			Node:     &ast.TypeId{Name: id},
			Position: id.Position,
		}
		nextParent := &ast.KindExpr{
			Node: &ast.TypeMember{
				Parent: parent,
				Id:     idKindExpr,
			},
			Position: ast.Position{
				Start: parent.Start,
				End:   id.End,
			},
		}
		return p.parseChainKindExpr(nextParent)
	}

	return parent
}

func (p *Parser) parseFuncKindExpr(start int) *ast.KindExpr {
	p.consume(lexer.TTParenL, true)
	kindExpr := &ast.KindExpr{}
	kindExpr.Start = start

	// arguments
	args := make([]*ast.Argument, 0, helper.DefaultCap)
	var lastRestToken *lexer.Token
	for !p.isEnd() && !p.isToken(lexer.TTParenR) {
		if lastRestToken != nil {
			p.unexpectedPos(lastRestToken.Start, "Can only use '...' as the final argument")
		}

		restToken := p.consume(lexer.TTRest, false)
		rest := restToken != nil
		if rest {
			lastRestToken = restToken
		}
		nameToken := p.consume(lexer.TTIdentifier, false)
		if nameToken == nil {
			p.unexpectedToken("identifier", p.current)
		}
		p.consume(lexer.TTColon, true)
		kind := p.parseKindExpr()
		argument := &ast.Argument{
			Name: NewIdentifier(nameToken),
			Kind: kind,
			Rest: rest,
		}
		argument.Start = nameToken.Start
		argument.End = kind.End
		args = append(args, argument)

		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}
	end := p.current.End
	p.consume(lexer.TTParenR, true)

	// return kind
	var kind *ast.KindExpr
	if p.consume(lexer.TTReturnSym, false) != nil {
		kind = p.parseKindExpr()
	}

	kindExpr.Node = &ast.TypeFuncKind{
		Arguments: args,
		Return:    kind,
	}
	if kind != nil {
		end = kind.End
	}
	kindExpr.End = end

	return kindExpr
}

func (p *Parser) parseKindProps(isFunc bool) []*ast.KindProperty {
	properties := make([]*ast.KindProperty, 0, helper.DefaultCap)

	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		pair := &ast.KindProperty{}

		if isFunc && p.isKeyword("fn") {
			start := p.current.Start
			p.nextToken()
			if !p.isToken(lexer.TTIdentifier) {
				p.unexpectedMissing("function name")
			}
			pair.Name = NewIdentifier(p.current)
			p.nextToken()
			pair.Kind = p.parseFuncKindExpr(start)
		} else if !isFunc {
			token := p.consume(lexer.TTIdentifier, true)
			pair.Name = NewIdentifier(token)
			p.consume(lexer.TTColon, true)
			pair.Kind = p.parseKindExpr()
		} else {
			p.unexpected()
		}

		pair.Start = pair.Name.Start
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

func (p *Parser) parseEnumItems() []*ast.KindIdentifier {
	items := make([]*ast.KindIdentifier, 0, helper.DefaultCap)

	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		token := p.consume(lexer.TTIdentifier, true)
		items = append(items, NewKindIdentifier(token))

		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}

	return items
}
