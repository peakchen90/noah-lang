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
			pubToken := p.current
			p.nextToken()
			p.expect(lexer.TTKeyword)

			switch p.current.Value {
			case "fn":
				omitTailingSemi = true
				stmt = p.parseFunctionDeclaration(pubToken)
			case "let":
				stmt = p.parseVariableDeclaration(pubToken)
			case "const":
				stmt = p.parseVariableDeclaration(pubToken)
			case "type":
				stmt = p.parseTypeDeclaration(pubToken)
			case "interface":
				stmt = p.parseInterfaceDeclaration(pubToken)
			default:
				p.unexpected()
			}
		case "import":
			stmt = p.parseImportStatement()
		case "fn":
			omitTailingSemi = true
			stmt = p.parseFunctionDeclaration(nil)
		case "let":
			stmt = p.parseVariableDeclaration(nil)
		case "const":
			stmt = p.parseVariableDeclaration(nil)
		case "type":
			stmt = p.parseTypeDeclaration(nil)
		case "interface":
			stmt = p.parseInterfaceDeclaration(nil)
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

func (p *Parser) parseFunctionDeclaration(pubToken *lexer.Token) *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseVariableDeclaration(pubToken *lexer.Token) *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseInterfaceDeclaration(pubToken *lexer.Token) *ast.Statement {
	kindDecl := ast.KindDecl{}
	kindDecl.Start = p.current.Start

	// type name
	p.nextToken()
	token := p.consume(lexer.TTIdentifier, true)
	if IsReservedType(token.Value) {
		p.unexpectedPos(token.Start, "Cannot declare reserved type "+token.Value)
	}
	Name := *NewKindIdentifier(token)

	// maybe extends a struct
	var Extends ast.KindIdentifier
	if p.consumeKeyword("extends", false) != nil {
		token := p.consume(lexer.TTIdentifier, true)
		Extends = *NewKindIdentifier(token)
	}

	// {
	p.consume(lexer.TTBraceL, true)

	// properties
	Properties := p.parseKindProperties(true)

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
			Pubic: pubToken != nil,
		},
		Position: kindDecl.Position,
	}

	return &stmt
}

func (p *Parser) parseTypeDeclaration(pubToken *lexer.Token) *ast.Statement {
	var stmt ast.Statement
	kindDecl := ast.KindDecl{}
	kindDecl.Start = p.current.Start

	defer func() {
		stmt = ast.Statement{
			Node: &ast.TypeDeclaration{
				Decl:  kindDecl,
				Pubic: pubToken != nil,
			},
			Position: kindDecl.Position,
		}
	}()

	// type name
	p.nextToken()
	token := p.consume(lexer.TTIdentifier, true)
	if IsReservedType(token.Value) {
		p.unexpectedPos(token.Start, "Cannot declare reserved type "+token.Value)
	}
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
	var interfaceToken *lexer.Token
	if p.consume(lexer.TTColon, false) != nil {
		interfaceToken = p.lexer.LastToken
		token := p.consume(lexer.TTIdentifier, false)
		if token == nil {
			p.unexpectedToken("interface name", p.current)
		}
		Interface = NewKindIdentifier(token)
	}

	// maybe has an extends
	var extendsToken *lexer.Token
	var Extends *ast.KindIdentifier
	if p.consumeKeyword("extends", false) != nil {
		extendsToken = p.lexer.LastToken
		token := p.consume(lexer.TTIdentifier, false)
		if token == nil {
			p.unexpectedToken("struct type identifier", p.current)
		}
		Extends = NewKindIdentifier(token)
	}

	// {
	p.consume(lexer.TTBraceL, true)
	Properties := make([]ast.KindProperty, 0)

	if !p.isToken(lexer.TTBraceR) {
		p.consume(lexer.TTIdentifier, true)
		if p.isToken(lexer.TTBraceR) || p.isToken(lexer.TTComma) { // 枚举类型
			if interfaceToken != nil {
				p.unexpectedPos(interfaceToken.Start, "Enumeration type does not support interface")
			} else if extendsToken != nil {
				p.unexpectedPos(extendsToken.Start, "Enumeration type does not support extends")
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

	p.consume(lexer.TTBraceR, true)
	kindDecl.Node = &ast.TypeStruct{
		Name:       *Name,
		Interface:  *Interface,
		Extends:    *Extends,
		Properties: Properties,
	}
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
