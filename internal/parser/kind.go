package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/lexer"
)

func (p *Parser) parseKindExpr() *ast.KindExpr {
	kindExpr := ast.KindExpr{}

	if p.isToken(lexer.TTIdentifier) { // type alias
		switch p.current.Value {
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
			kindExpr.Node = &ast.TypeId{Name: p.current.Value}
		}

		kindExpr.Position = p.current.Position
		p.nextToken()
	} else if p.isToken(lexer.TTBracketL) {
		kindExpr.Start = p.current.Start
		p.nextToken()

		switch p.current.Type {
		case lexer.TTRest: // [..]T
			p.nextToken()
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
		default: // [n]T
			var expr *ast.Expression
			if p.isToken(lexer.TTNumber) {
				if !IsUnsignedInt(p.current.Value) {
					p.unexpectedToken("constant integer", p.current)
				}
				expr = NewNumberExpr(p.current, p)
			} else if p.isToken(lexer.TTIdentifier) {
				expr = NewIdentifierExpr(p.current)
			} else {
				p.unexpectedToken("constant integer", p.current)
			}

			p.nextToken()
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

func (p *Parser) parseKindProperties(allowFunc bool) []ast.KindProperty {
	properties := make([]ast.KindProperty, 0, helper.DefaultCap)

	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		pair := ast.KindProperty{}

		// name
		token := p.consume(lexer.TTIdentifier, true)
		pair.Name = *NewKindIdentifier(token)

		if p.consume(lexer.TTColon, false) != nil {
			pair.Kind = *p.parseKindExpr()
		} else if allowFunc && p.consume(lexer.TTParenL, false) != nil {
			// TODO interface func
		} else {
			p.unexpected()
		}

		if p.consume(lexer.TTComma, false) == nil && !p.isToken(lexer.TTBraceR) {
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

func (p *Parser) parseEnumItems() []ast.KindIdentifier {
	items := make([]ast.KindIdentifier, 0, helper.DefaultCap)

	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		token := p.consume(lexer.TTIdentifier, true)
		items = append(items, *NewKindIdentifier(token))

		if p.consume(lexer.TTComma, false) == nil {
			break
		}
	}

	return items
}
