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
				stmt = p.parseVariableDeclaration(pubToken)
			case "const":
				stmt = p.parseVariableDeclaration(pubToken)
			case "type":
				stmt = p.parseTypeDeclaration(pubToken)
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
		p.unexpectedToken("a function name", p.current)
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

func (p *Parser) parseVariableDeclaration(pubToken *lexer.Token) *ast.Statement {
	stmt := ast.Statement{}

	return &stmt
}

func (p *Parser) parseTypeDeclaration(pubToken *lexer.Token) *ast.Statement {
	stmt := new(ast.Statement)
	kindDecl := &ast.KindDecl{}
	if pubToken != nil {
		kindDecl.Start = pubToken.Start
	} else {
		kindDecl.Start = p.current.Start
	}

	p.nextToken()

	defer func() {
		*stmt = ast.Statement{
			Node: &ast.TypeDeclaration{
				Decl:  kindDecl,
				Pubic: pubToken != nil,
			},
			Position: kindDecl.Position,
		}
	}()

	// maybe type name
	var nameToken *lexer.Token
	var Name *ast.KindIdentifier
	if p.isToken(lexer.TTIdentifier) {
		nameToken = p.current
		Name = NewKindIdentifier(p.current)
		p.nextToken()
	}

	// maybe interface name
	var interfaceToken *lexer.Token
	var Interface *ast.KindIdentifier
	if p.consume(lexer.TTColon, false) != nil {
		interfaceToken = p.lexer.LastToken
		token := p.consume(lexer.TTIdentifier, false)
		if token == nil {
			p.unexpectedToken("interface name", p.current)
		}
		Interface = NewKindIdentifier(token)
	}

	// nameToken
	finalNameToken := nameToken
	isInterface := false
	if nameToken == nil {
		finalNameToken = interfaceToken
		isInterface = true
	}
	if finalNameToken == nil {
		p.unexpectedPos(p.current.Start, "Unexpected type declaration")
	} else if IsReservedType(finalNameToken.Value) {
		p.unexpectedPos(finalNameToken.Start, "Cannot declare reserved type "+finalNameToken.Value)
	}

	// type alias
	if p.consume(lexer.TTAssign, false) != nil {
		kind := p.parseKindExpr()
		kindDecl.Node = &ast.TypeDeclAlias{
			Name: Name,
			Kind: kind,
		}
		kindDecl.End = kind.End
		return stmt
	}

	// maybe has extends
	var extendsToken *lexer.Token
	var Extends *ast.KindIdentifier
	if p.consumeKeyword("extends", false) != nil {
		extendsToken = p.lexer.LastToken
		token := p.consume(lexer.TTIdentifier, false)
		if token == nil {
			p.unexpectedToken("type identifier", p.current)
		}
		Extends = NewKindIdentifier(token)
	}

	// `{`
	p.consume(lexer.TTBraceL, true)
	Properties := make([]*ast.KindProperty, helper.DefaultCap)

	if !p.isToken(lexer.TTBraceR) {
		p.consume(lexer.TTIdentifier, true)
		if p.isToken(lexer.TTBraceR) || p.isToken(lexer.TTComma) { // 枚举类型
			if interfaceToken != nil {
				p.unexpectedPos(interfaceToken.Start, "Enum types cannot implement interfaces")
			} else if extendsToken != nil {
				p.unexpectedPos(extendsToken.Start, "Enum types cannot support extends")
			}

			p.revertLastToken()
			items := p.parseEnumItems()
			kindDecl.Node = &ast.TypeDeclEnum{
				Name:  Name,
				Items: items,
			}

			token := p.consume(lexer.TTBraceR, true)
			kindDecl.End = token.End
			return stmt

		} else { // 结构体或接口类型
			p.revertLastToken()
			Properties = p.parseKindProperties(isInterface)
		}
	}

	if isInterface {
		kindDecl.Node = &ast.TypeDeclInterface{
			Name:       Interface,
			Extends:    Extends,
			Properties: Properties,
		}
	} else {
		kindDecl.Node = &ast.TypeDeclStruct{
			Name:       Name,
			Interface:  Interface,
			Extends:    Extends,
			Properties: Properties,
		}
	}

	token := p.consume(lexer.TTBraceR, true)
	kindDecl.End = token.End

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
