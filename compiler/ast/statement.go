package ast

func (p *Parser) parseStatement() *Statement {
	var stmt *Statement

	omitTailingSemi := false
	switch p.currentToken.Type {
	case TTKeyword:
		switch p.currentToken.Value {
		case "pub":
			p.readNextToken()
			if !p.isToken(TTKeyword) {
				// TODO unexpected
			}

			switch p.currentToken.Value {
			case "fn":
				omitTailingSemi = true
				stmt = p.parseFunctionDeclaration(true)
			case "let":
			case "const":
			case "type":
			case "interface":
				// TODO parse pub vars & kinds
			default:
				// TODO unexpected
			}
		case "import":
		case "fn":
		case "let":
		case "const":
		case "type":
		case "interface":
		case "if":
		case "for":
		case "return":
		case "break":
		case "continue":
		}
	case TTIdentifier:
	case TTNumber:
	case TTParenL:
	case TTBraceL:

	default:
		// TODO unexpected
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

func (p *Parser) parseFunctionDeclaration(pub bool) *Statement {
	stmt := Statement{}

	return &stmt
}
