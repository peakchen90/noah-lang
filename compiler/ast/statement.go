package ast

func (p *Parser) parseStatement() *Statement {
	var stmt *Statement

	omitTailingSemi := false
	switch p.current.Type {
	case TTKeyword:
		switch p.current.Value {
		case "pub":
			p.nextToken()
			p.expect(TTKeyword)

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
				stmt = p.parseTypeDeclaration(true)
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
			stmt = p.parseTypeDeclaration(false)
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
	case TTIdentifier:
		maybeLabel := p.current.Value
		if p.lexer.lookNext() == ':' {
			p.nextToken()
			p.consume(TTColon)
			p.expect(TTKeyword)

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
	case TTBraceL:
		omitTailingSemi = true
		stmt = p.parseBlockStatement()
	default:
		p.unexpected()
	}

	tailSemiCount := 0
	for p.consume(TTSemi) {
		tailSemiCount++
		if p.isToken(TTEof) {
			break
		}
	}

	if !omitTailingSemi && tailSemiCount == 0 && !p.isSeenNewline && !p.isToken(TTEof) && !p.isToken(TTBraceR) {
		p.unexpected()
	}

	return stmt
}

func (p *Parser) parseImportStatement() *Statement {
	stmt := Statement{}
	stmt.Start = p.current.Start

	// source
	p.nextToken()
	p.expect(TTString)
	Source := p.current.Value

	// as
	p.nextToken()
	if !(p.isToken(TTKeyword) && p.current.Value == "as") {
		p.unexpected()
	}

	// localName
	p.nextToken()
	p.expect(TTIdentifier)
	LocalName := p.current.Value

	stmt.Node = &ImportDeclaration{
		Source,
		LocalName,
	}
	stmt.End = p.current.End

	p.nextToken()

	return &stmt
}

func (p *Parser) parseFunctionDeclaration(pub bool) *Statement {
	stmt := Statement{}

	return &stmt
}

func (p *Parser) parseVariableDeclaration(pub bool) *Statement {
	stmt := Statement{}

	return &stmt
}

func (p *Parser) parseTypeDeclaration(pub bool) *Statement {
	stmt := Statement{}

	return &stmt
}
func (p *Parser) parseIfStatement() *Statement {
	stmt := Statement{}

	return &stmt
}

func (p *Parser) parseForStatement(label string) *Statement {
	stmt := Statement{}

	return &stmt
}

func (p *Parser) parseReturnStatement() *Statement {
	stmt := Statement{}

	return &stmt
}

func (p *Parser) parseBreakStatement() *Statement {
	stmt := Statement{}

	return &stmt
}

func (p *Parser) parseContinueStatement() *Statement {
	stmt := Statement{}

	return &stmt
}

func (p *Parser) parseExpressionStatement() *Statement {
	stmt := Statement{}

	return &stmt
}

func (p *Parser) parseBlockStatement() *Statement {
	stmt := Statement{}

	return &stmt
}
