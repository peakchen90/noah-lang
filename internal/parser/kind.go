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
				Kind: kind,
			}
			kindExpr.End = kind.End
		case lexer.TTBracketR: // []T
			p.nextToken()
			kind := p.parseKindExpr()
			kindExpr.Node = &ast.TypeArray{
				Kind: kind,
				Len:  nil,
			}
			kindExpr.End = kind.End
		default: // [n]T
			// TODO 直接解析表达式
			var expr *ast.Expr
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
				Kind: kind,
				Len:  expr,
			}
			kindExpr.End = kind.End
		}
	} else if p.isKeyword("fn") { // fn(..arg: [..]T) -> T
		kindExpr = *p.parseFuncSignExpr(p.current.Start)
	} else {
		p.unexpected()
	}

	return &kindExpr
}

func (p *Parser) parseFuncSignExpr(start int) *ast.KindExpr {
	p.consume(lexer.TTParenL, true)
	kindExpr := &ast.KindExpr{}
	kindExpr.Start = start

	// arguments
	args := make([]*ast.Argument, 0, helper.DefaultCap)
	var lastRestToken *lexer.Token
	for !p.isEnd() && !p.isToken(lexer.TTParenR) {
		if lastRestToken != nil {
			p.unexpectedPos(lastRestToken.Start, "Can only use '..' as the final argument")
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

	kindExpr.Node = &ast.TypeFuncSign{
		Arguments: args,
		Kind:      kind,
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
			pair.Kind = p.parseFuncSignExpr(start)
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
