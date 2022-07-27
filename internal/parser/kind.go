package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/lexer"
)

func (p *Parser) parseKindProperties(properties []ast.KindProperty, allowFunc bool) []ast.KindProperty {
	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		pair := ast.KindProperty{}

		// name
		p.expect(lexer.TTIdentifier)
		pair.Name = *NewKindIdentifier(p.current)
		p.nextToken()

		if p.consume(lexer.TTColon, false) {
			pair.Kind = *p.parseKindExpr()
		} else if allowFunc && p.consume(lexer.TTParenL, false) {
			// TODO interface func
		} else {
			p.unexpected()
		}

		p.nextToken()
		if !p.consume(lexer.TTComma, false) && !p.isToken(lexer.TTBraceR) {
			p.nextToken()
			if !p.lexer.SeenNewline {
				break
			}
		}

		pair.Start = pair.Name.Start
		pair.End = pair.Kind.End
		properties = append(properties, pair)
	}

	return properties
}

func (p *Parser) parseKindExpr() *ast.KindExpr {
	kindExpr := ast.KindExpr{}

	if p.isToken(lexer.TTIdentifier) {
		kindExpr.Node = &ast.KindId{Name: p.current.Value}
		kindExpr.Position = p.current.Position
	} else if p.isToken(lexer.TTBracketL) {
		kindExpr.Start = p.current.Start
		p.nextToken()

		switch p.current.Type {
		case lexer.TTRest: // [..]T
			p.consume(lexer.TTBracketR, true)
			kind := p.parseKindExpr()
			kindExpr.Node = &ast.TypeVectorArray{
				Kind: *kind,
			}
			kindExpr.End = kind.End
		case lexer.TTBracketR: // []T
			p.nextToken()
			kind := p.parseKindExpr()
			kindExpr.Node = &ast.TypeArray{
				Kind: *kind,
				Len:  *new(ast.Expression),
			}
			kindExpr.End = kind.End
		default: // [expr]T
			expr := p.parseExpression()
			p.consume(lexer.TTBracketR, true)
			kind := p.parseKindExpr()
			kindExpr.Node = &ast.TypeArray{
				Kind: *kind,
				Len:  *expr,
			}
			kindExpr.End = kind.End
		}
	} else {
		p.unexpected()
	}

	return &kindExpr
}

func (p *Parser) parseEnumItems(items []ast.KindIdentifier) []ast.KindIdentifier {
	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		p.expect(lexer.TTIdentifier)
		items = append(items, *NewKindIdentifier(p.current))
		p.nextToken()

		if !p.consume(lexer.TTComma, false) {
			break
		}
	}
	return items
}
