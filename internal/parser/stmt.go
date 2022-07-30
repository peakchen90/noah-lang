package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
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
				stmt = p.parseVariableDeclaration(false, pubToken)
			case "const":
				stmt = p.parseVariableDeclaration(true, pubToken)
			case "type":
				stmt = p.parseTypeDeclaration(pubToken)
			case "interface":
				stmt = p.parseInterfaceDeclaration(pubToken)
			case "struct":
				stmt = p.parseStructDeclaration(pubToken)
			case "enum":
				stmt = p.parseEnumDeclaration(pubToken)
			default:
				p.unexpected()
			}
		case "import":
			stmt = p.parseImportStatement()
		case "fn":
			omitTailingSemi = true
			stmt = p.parseFunctionDeclaration(nil)
		case "let":
			stmt = p.parseVariableDeclaration(false, nil)
		case "const":
			stmt = p.parseVariableDeclaration(true, nil)
		case "type":
			stmt = p.parseTypeDeclaration(nil)
		case "interface":
			stmt = p.parseInterfaceDeclaration(nil)
		case "struct":
			stmt = p.parseStructDeclaration(nil)
		case "enum":
			stmt = p.parseEnumDeclaration(nil)
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
		p.unexpectedMissing("import source")
	}
	source := p.current.Value
	p.nextToken()

	// as
	if p.consumeKeyword("as", false) == nil {
		p.unexpectedMissing("keyword `as`")
	}

	// localName
	if !p.isToken(lexer.TTIdentifier) {
		p.unexpectedMissing("local name")
	}
	local := NewIdentifier(p.current)
	p.nextToken()

	stmt.Node = &ast.ImportDeclaration{
		Source: source,
		Local:  local,
	}
	stmt.End = local.End

	return &stmt
}

func (p *Parser) parseFunctionDeclaration(pubToken *lexer.Token) *ast.Statement {
	stmt := ast.Statement{}
	if pubToken != nil {
		stmt.Start = pubToken.Start
	} else {
		stmt.Start = p.current.Start
	}

	p.nextToken()
	implToken := p.consume(lexer.TTIdentifier, false)
	if implToken == nil {
		p.unexpectedMissing("function name")
	}

	nameToken := p.consume(lexer.TTIdentifier, false)
	if nameToken == nil {
		nameToken = implToken
		implToken = nil
	}

	funcSign := p.parseFuncSignExpr(stmt.Start)

	funcDecl := &ast.FunctionDeclaration{
		Name:     NewIdentifier(nameToken),
		FuncSign: funcSign,
		Body:     p.parseBlockStatement(),
		Pubic:    pubToken != nil,
	}
	if implToken != nil {
		funcDecl.Impl = NewKindIdentifier(implToken)
	}
	stmt.Node = funcDecl

	return &stmt
}

func (p *Parser) parseVariableDeclaration(isConst bool, pubToken *lexer.Token) *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseTypeDeclaration(pubToken *lexer.Token) *ast.Statement {
	stmt := &ast.Statement{}
	if pubToken != nil {
		stmt.Start = pubToken.Start
	} else {
		stmt.Start = p.current.Start
	}
	p.nextToken()

	if !p.isToken(lexer.TTIdentifier) {
		p.unexpectedMissing("type name")
	}
	name := NewKindIdentifier(p.current)
	p.nextToken()

	if !p.isToken(lexer.TTAssign) {
		p.unexpectedMissing("=")
	}
	p.nextToken()

	// alias T
	kind := p.parseKindExpr()
	stmt.Node = &ast.TypeAliasDecl{
		Name:  name,
		Kind:  kind,
		Pubic: pubToken != nil,
	}
	stmt.End = kind.End
	return stmt
}

func (p *Parser) parseInterfaceDeclaration(pubToken *lexer.Token) *ast.Statement {
	stmt := &ast.Statement{}
	if pubToken != nil {
		stmt.Start = pubToken.Start
	} else {
		stmt.Start = p.current.Start
	}
	p.nextToken()

	// name
	if !p.isToken(lexer.TTIdentifier) {
		p.unexpectedMissing("interface name")
	}
	name := NewKindIdentifier(p.current)
	p.nextToken()

	// extends
	var Extends *ast.KindIdentifier
	if p.consumeKeyword("extends", false) != nil {
		token := p.consume(lexer.TTIdentifier, false)
		if token == nil {
			p.unexpectedMissing("extends type")
		}
		Extends = NewKindIdentifier(token)
	}

	// `{`
	p.consume(lexer.TTBraceL, true)
	properties := p.parseKindProperties(true)
	p.consume(lexer.TTBraceR, true)

	stmt.Node = &ast.TypeInterfaceDecl{
		Name:       name,
		Extends:    Extends,
		Properties: properties,
		Pubic:      pubToken != nil,
	}
	stmt.End = p.lexer.LastToken.End
	return stmt
}

func (p *Parser) parseStructDeclaration(pubToken *lexer.Token) *ast.Statement {
	stmt := &ast.Statement{}
	if pubToken != nil {
		stmt.Start = pubToken.Start
	} else {
		stmt.Start = p.current.Start
	}
	p.nextToken()

	// name
	if !p.isToken(lexer.TTIdentifier) {
		p.unexpectedMissing("struct name")
	}
	name := NewKindIdentifier(p.current)
	p.nextToken()

	// impl
	var impl *ast.KindIdentifier
	if p.consume(lexer.TTColon, false) != nil {
		token := p.consume(lexer.TTIdentifier, false)
		if token == nil {
			p.unexpectedMissing("extends type")
		}
		impl = NewKindIdentifier(token)
	}

	// extends
	var extends *ast.KindIdentifier
	if p.consumeKeyword("extends", false) != nil {
		token := p.consume(lexer.TTIdentifier, false)
		if token == nil {
			p.unexpectedMissing("extends type")
		}
		extends = NewKindIdentifier(token)
	}

	// `{`
	p.consume(lexer.TTBraceL, true)
	properties := p.parseKindProperties(false)
	p.consume(lexer.TTBraceR, true)

	stmt.Node = &ast.TypeStructDecl{
		Name:       name,
		Impl:       impl,
		Extends:    extends,
		Properties: properties,
		Pubic:      pubToken != nil,
	}
	stmt.End = p.lexer.LastToken.End
	return stmt
}

func (p *Parser) parseEnumDeclaration(pubToken *lexer.Token) *ast.Statement {
	stmt := &ast.Statement{}
	if pubToken != nil {
		stmt.Start = pubToken.Start
	} else {
		stmt.Start = p.current.Start
	}
	p.nextToken()

	// name
	if !p.isToken(lexer.TTIdentifier) {
		p.unexpectedMissing("enum name")
	}
	name := NewKindIdentifier(p.current)
	p.nextToken()

	// `{`
	p.consume(lexer.TTBraceL, true)
	items := p.parseEnumItems()
	p.consume(lexer.TTBraceR, true)

	stmt.Node = &ast.TypeEnumDecl{
		Name:  name,
		Items: items,
		Pubic: pubToken != nil,
	}
	stmt.End = p.lexer.LastToken.End
	return stmt
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
	stmt.Start = p.current.Start
	p.consume(lexer.TTBraceL, true)

	body := make([]*ast.Statement, 0, helper.DefaultCap)

	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		body = append(body, p.parseStatement())
	}

	stmt.Node = &ast.BlockStatement{Body: body}
	stmt.End = p.current.End
	p.consume(lexer.TTBraceR, true)

	return &stmt
}
