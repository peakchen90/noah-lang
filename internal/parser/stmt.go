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
			token := p.consume(lexer.TTKeyword, true)

			switch token.Value {
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
			token := p.consume(lexer.TTKeyword, true)

			switch token.Value {
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
	for p.consume(lexer.TTSemi, false) != nil {
		tailSemiCount++
		if p.isEnd() {
			break
		}
	}

	if !omitTailingSemi && tailSemiCount == 0 && !p.lexer.SeenNewline && !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		p.unexpected()
	}

	return stmt
}

func (p *Parser) parseImportStatement() *ast.Statement {
	stmt := ast.Statement{}
	stmt.Start = p.current.Start
	p.nextToken()

	// source
	if !p.isToken(lexer.TTString) {
		p.unexpectedToken("string literal", p.current)
	}
	source := p.current.Value
	p.nextToken()

	// as
	if p.consumeKeyword("as", false) == nil {
		p.unexpectedToken("`as` keyword", p.current)
	}

	// localName
	if !p.isToken(lexer.TTIdentifier) {
		p.unexpectedToken("local name", p.current)
	}
	local := NewIdentifier(p.current)
	p.nextToken()

	stmt.Node = &ast.ImportDeclaration{
		Source: source,
		Local:  *local,
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
	token := p.consume(lexer.TTIdentifier, true)
	Name := *NewKindIdentifier(token)

	Properties := make([]ast.KindProperty, 0, 3)

	// maybe extends a struct
	var Extends ast.KindIdentifier
	if p.consumeKeyword("extends", false) != nil {
		token := p.consume(lexer.TTIdentifier, true)
		Extends = *NewKindIdentifier(token)
	}

	// {
	p.consume(lexer.TTBraceL, true)

	// properties
	Properties = p.parseKindProperties(true)

	// }
	token = p.consume(lexer.TTBraceR, true)
	kindDecl.End = token.End

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
		Position: kindDecl.Position,
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
			Position: kindDecl.Position,
		}
	}()

	// type name
	p.nextToken()
	token := p.consume(lexer.TTIdentifier, true)
	Name := NewKindIdentifier(token)

	// type alias
	if p.isToken(lexer.TTIdentifier) || p.isToken(lexer.TTBracketL) {
		kind := p.parseKindExpr()
		kindDecl.Node = &ast.TypeAlias{
			Name: *Name,
			Kind: *kind,
		}
		kindDecl.End = kind.End
		return &stmt
	}

	// maybe has an interface
	var Interface *ast.KindIdentifier
	if p.consume(lexer.TTColon, false) != nil {
		shouldBeStruct = true
		token := p.consume(lexer.TTIdentifier, true)
		Interface = NewKindIdentifier(token)
	}

	// maybe has an extends
	var Extends *ast.KindIdentifier
	if p.consumeKeyword("extends", false) != nil {
		shouldBeStruct = true
		token := p.consume(lexer.TTIdentifier, true)
		Extends = NewKindIdentifier(token)
	}

	// {
	p.consume(lexer.TTBraceL, true)
	Properties := make([]ast.KindProperty, 0)

	if !p.isToken(lexer.TTBraceR) {
		p.consume(lexer.TTIdentifier, true)
		if p.isToken(lexer.TTBraceR) || p.isToken(lexer.TTComma) { // 枚举类型
			if shouldBeStruct {
				p.unexpected()
			}

			p.revertLastToken()
			items := p.parseEnumItems()
			kindDecl.Node = &ast.TypeEnum{
				Name:  *Name,
				Items: items,
			}

			token := p.consume(lexer.TTBraceR, true)
			kindDecl.End = token.End
			return &stmt

		} else { // 结构体类型
			p.revertLastToken()
			Properties = p.parseKindProperties(false)
		}
	}

	kindDecl.Node = &ast.TypeStruct{
		Name:       *Name,
		Interface:  *Interface,
		Extends:    *Extends,
		Properties: Properties,
	}
	p.nextToken()
	p.expect(lexer.TTBraceR)
	kindDecl.End = p.current.End

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
