package parser

import (
	"github.com/peakchen90/hera-lang/internal/ast"
	"github.com/peakchen90/hera-lang/internal/lexer"
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
	case lexer.TTIdentifier:
		maybeLabel := p.current.Value
		if p.lexer.LookNext() == ':' {
			p.nextToken()
			p.consume(lexer.TTColon)
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
	for p.consume(lexer.TTSemi) {
		tailSemiCount++
		if p.isToken(lexer.TTEof) {
			break
		}
	}

	if !omitTailingSemi && tailSemiCount == 0 && !p.isSeenNewline && !p.isToken(lexer.TTEof) && !p.isToken(lexer.TTBraceR) {
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
	Source := p.current.Value

	// as
	p.nextToken()
	if !(p.isToken(lexer.TTKeyword) && p.current.Value == "as") {
		p.unexpected()
	}

	// localName
	p.nextToken()
	p.expect(lexer.TTIdentifier)
	LocalName := p.current.Value

	stmt.Node = &ast.ImportDeclaration{
		Source,
		LocalName,
	}
	stmt.End = p.current.End

	p.nextToken()

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

func (p *Parser) parseTypeDeclaration(pub bool) *ast.Statement {
	stmt := ast.Statement{}

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
