package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/lexer"
)

func (p *Parser) parseStatement() *ast.Statement {
	var stmt *ast.Statement

	omitTailingSemi := false
	switch p.current.Type {
	case lexer.TTKeyword:
		switch p.current.Value {
		case "pub":
			p.nextToken()
			p.expect(lexer.TTKeyword)

			switch p.current.Value {
			case "fn":
				omitTailingSemi = true
				stmt = p.parseFunctionDeclaration(true)
			case "let":
				stmt = p.parseVariableDeclaration(true)
			case "const":
				stmt = p.parseVariableDeclaration(true)
			case "type":
				stmt = p.parseTypeDeclaration(true)
			case "interface":
				stmt = p.parseInterfaceDeclaration(true)
			default:
				p.unexpected()
			}
		case "import":
			stmt = p.parseImportStatement()
		case "fn":
			omitTailingSemi = true
			stmt = p.parseFunctionDeclaration(false)
		case "let":
			stmt = p.parseVariableDeclaration(false)
		case "const":
			stmt = p.parseVariableDeclaration(false)
		case "type":
			stmt = p.parseTypeDeclaration(false)
		case "interface":
			stmt = p.parseInterfaceDeclaration(false)
		case "if":
			omitTailingSemi = true
			stmt = p.parseIfStatement()
		case "for":
			omitTailingSemi = true
			stmt = p.parseForStatement("")
		case "return":
			stmt = p.parseReturnStatement()
		case "break":
			stmt = p.parseBreakStatement()
		case "continue":
			stmt = p.parseContinueStatement()
		default:
			p.unexpected()
		}
	case lexer.TTIdentifier:
		maybeLabel := p.current.Value
		if p.lexer.LookNext() == ':' {
			p.nextToken()
			p.consume(lexer.TTColon, true)
			p.expect(lexer.TTKeyword)

			switch p.current.Value {
			case "for":
				omitTailingSemi = true
				stmt = p.parseForStatement(maybeLabel)
			default:
				p.unexpected()
			}

		} else {
			stmt = p.parseExpressionStatement()
		}
	case lexer.TTBraceL:
		omitTailingSemi = true
		stmt = p.parseBlockStatement()
	default:
		p.unexpected()
	}

	tailSemiCount := 0
	for p.consume(lexer.TTSemi, false) {
		tailSemiCount++
		if p.isEnd() {
			break
		}
	}

	if !omitTailingSemi && tailSemiCount == 0 && !p.isSeenNewline && !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		p.unexpected()
	}

	return stmt
}

func (p *Parser) parseImportStatement() *ast.Statement {
	stmt := ast.Statement{}
	stmt.Start = p.current.Start

	// source
	p.nextToken()
	p.expect(lexer.TTString)
	source := p.current.Value

	// as
	p.nextToken()
	p.consumeKeyword("as", true)

	// localName
	p.nextToken()
	p.expect(lexer.TTIdentifier)
	local := *NewIdentifier(p.current, nil)
	p.nextToken()

	stmt.Node = &ast.ImportDeclaration{
		Source: source,
		Local:  local,
	}
	stmt.End = local.End

	return &stmt
}

func (p *Parser) parseFunctionDeclaration(pub bool) *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseVariableDeclaration(pub bool) *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseInterfaceDeclaration(pub bool) *ast.Statement {
	kindDecl := ast.KindDecl{}
	kindDecl.Start = p.current.Start

	// type name
	p.nextToken()
	p.expect(lexer.TTIdentifier)
	Name := *NewKindIdentifier(p.current)
	p.nextToken()

	Properties := make([]ast.KindProperty, 0, 3)

	// maybe extends a struct
	var Extends ast.KindIdentifier
	if p.consumeKeyword("extends", false) {
		p.expect(lexer.TTIdentifier)
		Extends = *NewKindIdentifier(p.current)
		p.nextToken()
	}

	// {
	p.consume(lexer.TTBraceL, true)

	// properties
	Properties = p.parseKindProperties(Properties, true)

	// }
	p.consume(lexer.TTBraceR, true)

	kindDecl.Node = &ast.TypeInterface{
		Name:       Name,
		Extends:    Extends,
		Properties: Properties,
	}

	stmt := ast.Statement{
		Node: &ast.TypeDeclaration{
			Decl:  kindDecl,
			Pubic: pub,
		},
	}

	return &stmt
}

func (p *Parser) parseTypeDeclaration(pub bool) *ast.Statement {
	var stmt ast.Statement
	kindDecl := ast.KindDecl{}
	kindDecl.Start = p.current.Start
	shouldBeStruct := false

	defer func() {
		stmt = ast.Statement{
			Node: &ast.TypeDeclaration{
				Decl:  kindDecl,
				Pubic: pub,
			},
		}
		p.expect(lexer.TTBraceR)
		kindDecl.End = p.current.End
		stmt.Position = kindDecl.Position
		p.nextToken()
	}()

	// type name
	p.nextToken()
	p.expect(lexer.TTIdentifier)
	Name := *NewKindIdentifier(p.current)
	p.nextToken()

	// maybe has an interface
	var Interface ast.KindIdentifier
	if p.consume(lexer.TTColon, false) {
		shouldBeStruct = true
		p.expect(lexer.TTIdentifier)
		Interface = *NewKindIdentifier(p.current)
		p.nextToken()
	}

	// maybe has an extends
	var Extends ast.KindIdentifier
	if p.consumeKeyword("extends", false) {
		shouldBeStruct = true
		p.expect(lexer.TTIdentifier)
		Extends = *NewKindIdentifier(p.current)
		p.nextToken()
	}

	// {
	p.consume(lexer.TTBraceL, true)
	Properties := make([]ast.KindProperty, 0)

	if !p.isToken(lexer.TTBraceR) {
		p.expect(lexer.TTIdentifier)
		headToken := p.current
		p.nextToken()
		if p.isToken(lexer.TTBraceR) || p.isToken(lexer.TTComma) { // 枚举类型
			if shouldBeStruct {
				p.unexpected()
			}

			items := make([]ast.KindIdentifier, 0, 3)
			items = append(items, *NewKindIdentifier(headToken))
			hasComma := p.consume(lexer.TTComma, false)
			if !hasComma && !p.lexer.SeenNewline && !p.isToken(lexer.TTBraceR) {
				p.unexpected()
			}
			items = p.parseEnumItems(items)
			kindDecl.Node = &ast.TypeEnum{
				Name:  Name,
				Items: items,
			}
			return &stmt
		} else { // 结构体类型
			p.consume(lexer.TTColon, true)
			p.expect(lexer.TTIdentifier)
			firstPair := ast.KindProperty{
				Name: *NewKindIdentifier(headToken),
				Kind: *p.parseKindExpr(),
			}
			firstPair.Start = firstPair.Name.Start
			firstPair.End = firstPair.Kind.End
			Properties = make([]ast.KindProperty, 0, 3)
			Properties = append(Properties, firstPair)

			p.nextToken()
			hasComma := p.consume(lexer.TTComma, false)
			if !hasComma && !p.lexer.SeenNewline && !p.isToken(lexer.TTBraceR) {
				p.unexpected()
			}
			Properties = p.parseKindProperties(Properties, false)
		}
	}

	kindDecl.Node = &ast.TypeStruct{
		Name:       Name,
		Interface:  Interface,
		Extends:    Extends,
		Properties: Properties,
	}

	return &stmt
}

func (p *Parser) parseIfStatement() *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseForStatement(label string) *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseReturnStatement() *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseBreakStatement() *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseContinueStatement() *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseExpressionStatement() *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseBlockStatement() *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}
