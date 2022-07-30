package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/lexer"
)

func (p *Parser) parseStmt() *ast.Stmt {
	var stmt *ast.Stmt

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
				stmt = p.parseFuncDecl(pubToken)
			case "let":
				stmt = p.parseVarDecl(false, pubToken)
			case "const":
				stmt = p.parseVarDecl(true, pubToken)
			case "type":
				stmt = p.parseTypeDecl(pubToken)
			case "interface":
				stmt = p.parseInterfaceDecl(pubToken)
			case "struct":
				stmt = p.parseStructDecl(pubToken)
			case "enum":
				stmt = p.parseEnumDecl(pubToken)
			default:
				p.unexpected()
			}
		case "import":
			stmt = p.parseImportStmt()
		case "fn":
			omitTailingSemi = true
			stmt = p.parseFuncDecl(nil)
		case "let":
			stmt = p.parseVarDecl(false, nil)
		case "const":
			stmt = p.parseVarDecl(true, nil)
		case "type":
			stmt = p.parseTypeDecl(nil)
		case "interface":
			stmt = p.parseInterfaceDecl(nil)
		case "struct":
			stmt = p.parseStructDecl(nil)
		case "enum":
			stmt = p.parseEnumDecl(nil)
		case "if":
			omitTailingSemi = true
			stmt = p.parseIfStmt()
		case "for":
			omitTailingSemi = true
			stmt = p.parseForStmt("")
		case "return":
			stmt = p.parseReturnStmt()
		case "break":
			stmt = p.parseBreakStmt()
		case "continue":
			stmt = p.parseContinueStmt()
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
				stmt = p.parseForStmt(maybeLabel)
			default:
				p.unexpected()
			}

		} else {
			stmt = p.parseExprStmt()
		}
	case lexer.TTBraceL:
		omitTailingSemi = true
		stmt = p.parseBlockStmt()
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

func (p *Parser) parseImportStmt() *ast.Stmt {
	stmt := ast.Stmt{}
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

	stmt.Node = &ast.ImportDecl{
		Source: source,
		Local:  local,
	}
	stmt.End = local.End

	return &stmt
}

func (p *Parser) parseFuncDecl(pubToken *lexer.Token) *ast.Stmt {
	stmt := ast.Stmt{}
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

	funcDecl := &ast.FuncDecl{
		Name:     NewIdentifier(nameToken),
		FuncSign: funcSign,
		Body:     p.parseBlockStmt(),
		Pubic:    pubToken != nil,
	}
	if implToken != nil {
		funcDecl.Impl = NewKindIdentifier(implToken)
	}
	stmt.Node = funcDecl

	return &stmt
}

func (p *Parser) parseVarDecl(isConst bool, pubToken *lexer.Token) *ast.Stmt {
	stmt := ast.Stmt{}

	return &stmt
}

func (p *Parser) parseTypeDecl(pubToken *lexer.Token) *ast.Stmt {
	stmt := &ast.Stmt{}
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

func (p *Parser) parseInterfaceDecl(pubToken *lexer.Token) *ast.Stmt {
	stmt := &ast.Stmt{}
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
	properties := p.parseKindProps(true)
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

func (p *Parser) parseStructDecl(pubToken *lexer.Token) *ast.Stmt {
	stmt := &ast.Stmt{}
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
	properties := p.parseKindProps(false)
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

func (p *Parser) parseEnumDecl(pubToken *lexer.Token) *ast.Stmt {
	stmt := &ast.Stmt{}
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

func (p *Parser) parseIfStmt() *ast.Stmt {
	stmt := ast.Stmt{}

	return &stmt
}

func (p *Parser) parseForStmt(label string) *ast.Stmt {
	stmt := ast.Stmt{}

	return &stmt
}

func (p *Parser) parseReturnStmt() *ast.Stmt {
	stmt := ast.Stmt{}

	return &stmt
}

func (p *Parser) parseBreakStmt() *ast.Stmt {
	stmt := ast.Stmt{}

	return &stmt
}

func (p *Parser) parseContinueStmt() *ast.Stmt {
	stmt := ast.Stmt{}

	return &stmt
}

func (p *Parser) parseExprStmt() *ast.Stmt {
	stmt := ast.Stmt{}

	return &stmt
}

func (p *Parser) parseBlockStmt() *ast.Stmt {
	stmt := ast.Stmt{}
	stmt.Start = p.current.Start
	p.consume(lexer.TTBraceL, true)

	body := make([]*ast.Stmt, 0, helper.DefaultCap)

	for !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		body = append(body, p.parseStmt())
	}

	stmt.Node = &ast.BlockStmt{Body: body}
	stmt.End = p.current.End
	p.consume(lexer.TTBraceR, true)

	return &stmt
}
