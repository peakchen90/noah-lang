package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/lexer"
)

func (p *Parser) parseStmt() *ast.Stmt {
	var stmt *ast.Stmt

	switch p.current.Type {
	case lexer.TTKeyword:
		switch p.current.Value {
		case "pub":
			pubToken := p.current
			p.nextToken()
			p.expect(lexer.TTKeyword)

			switch p.current.Value {
			case "fn":
				stmt = p.parseFuncDecl(pubToken)
			case "let":
				stmt = p.parseVarDecl(pubToken, false)
			case "const":
				stmt = p.parseVarDecl(pubToken, true)
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
		case "fn":
			stmt = p.parseFuncDecl(nil)
		case "let":
			stmt = p.parseVarDecl(nil, false)
		case "const":
			stmt = p.parseVarDecl(nil, true)
		case "type":
			stmt = p.parseTypeDecl(nil)
		case "interface":
			stmt = p.parseInterfaceDecl(nil)
		case "struct":
			stmt = p.parseStructDecl(nil)
		case "enum":
			stmt = p.parseEnumDecl(nil)
		case "import":
			stmt = p.parseImportDecl()
		case "impl":
			stmt = p.parseImplDecl()
		case "if":
			stmt = p.parseIfStmt()
		case "for":
			stmt = p.parseForStmt(nil)
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
		maybeLabel := p.current
		if p.lexer.LookNext() == ':' {
			p.nextToken()
			p.consume(lexer.TTColon, true)
			if p.isKeyword("for") {
				stmt = p.parseForStmt(maybeLabel)
			} else {
				p.unexpected()
			}
		} else {
			stmt = p.parseExprStmt()
		}
	default:
		stmt = p.parseExprStmt()
	}

	tailSemiCount := 0
	for p.consume(lexer.TTSemi, false) != nil {
		tailSemiCount++
		if p.isEnd() {
			break
		}
	}

	if tailSemiCount == 0 && !p.lexer.SeenNewline && !p.isEnd() && !p.isToken(lexer.TTBraceR) {
		p.unexpected()
	}

	return stmt
}

func (p *Parser) parseImportDecl() *ast.Stmt {
	stmt := &ast.Stmt{}
	stmt.Start = p.current.Start
	p.nextToken() // skip `import`

	var Package *ast.Identifier
	headToken := p.consume(lexer.TTIdentifier, true)
	paths := make([]*ast.Identifier, 0, helper.DefaultCap)
	index := headToken.End

	if p.current.Start == index {
		if p.consume(lexer.TTColon, false) != nil {
			index++
			Package = newIdentifier(headToken)
			p.expect(lexer.TTIdentifier)
		} else {
			paths = append(paths, newIdentifier(headToken))
			if p.consume(lexer.TTDot, false) != nil {
				index++
			}
		}

		for p.current.Start == index {
			if p.isToken(lexer.TTIdentifier) {
				index = p.current.End
				paths = append(paths, newIdentifier(p.consume(lexer.TTIdentifier, true)))

				if p.current.Start == index && p.isToken(lexer.TTDot) {
					index = p.current.End
					p.nextToken()
					p.expect(lexer.TTIdentifier)
					continue
				}
			}

			break
		}
	} else {
		paths = append(paths, newIdentifier(headToken))
	}

	stmt.End = p.lexer.LastToken.End
	var local *ast.Identifier
	localId := paths[len(paths)-1]

	if p.consume(lexer.TTAsOp, false) != nil {
		local = newIdentifier(p.consume(lexer.TTIdentifier, true))
		localId = local
		stmt.End = local.End
	}

	if localId.Name == "self" {
		p.UnexpectedPos(localId.Start, "cannot use `self` as a local identifier")
	}

	stmt.Node = &ast.ImportDecl{
		Package: Package,
		Paths:   paths,
		Local:   local,
	}

	return stmt
}

func (p *Parser) parseFuncDecl(pubToken *lexer.Token) *ast.Stmt {
	stmt := &ast.Stmt{}
	if pubToken != nil {
		stmt.Start = pubToken.Start
	} else {
		stmt.Start = p.current.Start
	}
	p.nextToken() // skip `fn`

	nameToken := p.consumeVarId(true)
	funcKind := p.parseFuncKindExpr(stmt.Start)
	funcDecl := &ast.FuncDecl{
		Name: newIdentifier(nameToken),
		Kind: funcKind,
		Body: p.parseBlockStmt(),
		Pub:  pubToken != nil,
	}

	stmt.Node = funcDecl
	stmt.End = p.lexer.LastToken.End

	return stmt
}

func (p *Parser) parseVarDecl(pubToken *lexer.Token, isConst bool) *ast.Stmt {
	stmt := &ast.Stmt{}
	if pubToken != nil {
		stmt.Start = pubToken.Start
	} else {
		stmt.Start = p.current.Start
	}
	p.nextToken()

	// id
	token := p.consumeVarId(false)
	if token == nil {
		p.unexpectedMissing("variable name")
	}
	if isReservedType(token.Value) {
		p.UnexpectedPos(token.Start, "Reserved type cannot be used: "+token.Value)
	}
	id := newIdentifier(token)

	// maybe kind
	var kind *ast.KindExpr
	token = p.consume(lexer.TTColon, false)
	if token != nil {
		kind = p.parseKindExpr()
	}

	// maybe init
	token = p.consume(lexer.TTAssign, false)
	var init *ast.Expr
	if token != nil {
		init = p.parseExpr()
	}

	stmt.Node = &ast.VarDecl{
		Id:    id,
		Kind:  kind,
		Init:  init,
		Const: isConst,
		Pub:   pubToken != nil,
	}

	return stmt
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
	if isReservedType(p.current.Value) {
		p.UnexpectedPos(p.current.Start, "Reserved type cannot be used: "+p.current.Value)
	}
	name := newKindIdentifier(p.current)
	p.consumeVarId(true)

	kind := p.parseKindExpr()
	stmt.Node = &ast.TTypeDecl{
		Name: name,
		Kind: kind,
		Pub:  pubToken != nil,
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
	if isReservedType(p.current.Value) {
		p.UnexpectedPos(p.current.Start, "Reserved type cannot be used: "+p.current.Value)
	}
	name := newKindIdentifier(p.current)
	p.consumeVarId(true)

	// `{`
	p.consume(lexer.TTBraceL, true)
	properties := p.parseKindProperties(true)
	p.consume(lexer.TTBraceR, true)

	stmt.Node = &ast.TInterfaceDecl{
		Name:       name,
		Properties: properties,
		Pub:        pubToken != nil,
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

	start := p.current.Start
	p.nextToken() // skip `struct`

	nameToken := p.consumeVarId(true)
	if isReservedType(nameToken.Value) {
		p.UnexpectedPos(nameToken.Start, "Reserved type cannot be used: "+nameToken.Value)
	}

	stmt.Node = &ast.TStructDecl{
		Name: newKindIdentifier(nameToken),
		Kind: p.parseStructKindExpr(start),
		Pub:  pubToken != nil,
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
	if isReservedType(p.current.Value) {
		p.UnexpectedPos(p.current.Start, "Reserved type cannot be used: "+p.current.Value)
	}
	name := newKindIdentifier(p.current)
	p.consumeVarId(true)

	// `{`
	p.consume(lexer.TTBraceL, true)
	choices := p.parseEnumItems()
	p.consume(lexer.TTBraceR, true)

	stmt.Node = &ast.TEnumDecl{
		Name:    name,
		Choices: choices,
		Pub:     pubToken != nil,
	}
	stmt.End = p.lexer.LastToken.End
	return stmt
}

func (p *Parser) parseImplDecl() *ast.Stmt {
	stmt := &ast.Stmt{}
	stmt.Start = p.current.Start
	p.nextToken()

	var Interface *ast.KindExpr
	// maybe with interface
	if p.consume(lexer.TTParenL, false) != nil {
		p.expect(lexer.TTIdentifier)
		Interface = p.parseKindExpr()
		p.consume(lexer.TTParenR, true)
	}

	p.expect(lexer.TTIdentifier)
	target := p.parseKindExpr()

	body := make([]*ast.Stmt, 0, helper.DefaultCap)
	blockStmt := &ast.BlockStmt{}
	bodyStmt := &ast.Stmt{Node: blockStmt}
	bodyStmt.Start = p.current.Start

	p.consume(lexer.TTBraceL, true)
	for !p.isToken(lexer.TTBraceR) {
		if p.isKeyword("fn") {
			body = append(body, p.parseFuncDecl(nil))
			for {
				if p.consume(lexer.TTSemi, false) == nil {
					break
				}
			}
		} else {
			p.unexpected()
		}
	}
	bodyStmt.End = p.current.End
	p.consume(lexer.TTBraceR, true)

	blockStmt.Body = body

	stmt.Node = &ast.ImplDecl{
		Target:    target,
		Interface: Interface,
		Body:      bodyStmt,
	}
	stmt.End = p.lexer.LastToken.End
	return stmt
}

func (p *Parser) parseIfStmt() *ast.Stmt {
	stmt := ast.Stmt{}
	stmt.Start = p.current.Start
	p.nextToken()

	hasParentheses := p.consume(lexer.TTParenL, false) != nil // `(`
	condition := p.parseExpr()

	if hasParentheses {
		p.consume(lexer.TTParenR, true) // `)`
	}

	var alternate *ast.Stmt

	consequent := p.parseBlockStmt()
	if p.isKeyword("else") {
		p.nextToken()
		if p.isKeyword("if") {
			alternate = p.parseIfStmt()
		} else {
			alternate = p.parseBlockStmt()
		}
	}

	stmt.Node = &ast.IfStmt{
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
	}
	stmt.End = p.lexer.LastToken.End
	return &stmt
}

func (p *Parser) parseForStmt(labelToken *lexer.Token) *ast.Stmt {
	stmt := ast.Stmt{}
	stmt.Start = p.current.Start
	p.nextToken()

	var label *ast.Identifier
	var init *ast.Stmt
	var test *ast.Expr
	var update *ast.Expr
	var eachVisitor *ast.EachVisitor

	if labelToken != nil {
		label = newIdentifier(labelToken)
	}

	hasParentheses := p.consume(lexer.TTParenL, false) != nil // `(`

	if p.isToken(lexer.TTIdentifier) {
		headToken := p.current
		p.consumeVarId(true)

		if p.isToken(lexer.TTComma) || p.isToken(lexer.TTColon) { // for value, key: target {}
			var key *ast.Identifier
			if p.isToken(lexer.TTComma) {
				p.nextToken()
				key = newIdentifier(p.consumeVarId(true))
			}

			p.consume(lexer.TTColon, true)

			eachVisitor = &ast.EachVisitor{
				Value:  newIdentifier(headToken),
				Key:    key,
				Target: p.parseExpr(),
			}
		} else { // for value {}
			test = p.parseIdentifierExpr(headToken)
		}
	} else if p.isKeyword("let") || p.isKeyword("const") { // for (init; test; update) {}
		init = p.parseVarDecl(nil, p.current.Value == "const")
		p.consume(lexer.TTSemi, true)
		if !p.isToken(lexer.TTSemi) {
			test = p.parseExpr()
		}
		p.consume(lexer.TTSemi, true)
		if !p.isToken(lexer.TTParenR) && !p.isToken(lexer.TTBraceL) {
			update = p.parseExpr()
		}
	} else if !p.isToken(lexer.TTBraceL) { // for test {}
		test = p.parseExpr()
	}

	if hasParentheses {
		p.consume(lexer.TTParenR, true) // `)`
	}

	body := p.parseBlockStmt()

	stmt.Node = &ast.ForStmt{
		Label:       label,
		Init:        init,
		Test:        test,
		Update:      update,
		EachVisitor: eachVisitor,
		Body:        body,
	}
	stmt.End = p.lexer.LastToken.End
	return &stmt
}

func (p *Parser) parseReturnStmt() *ast.Stmt {
	stmt := ast.Stmt{}
	stmt.Start = p.current.Start
	p.nextToken()

	var argument *ast.Expr
	if !p.isToken(lexer.TTSemi) && !p.lexer.SeenNewline {
		argument = p.parseExpr()
	}

	stmt.Node = &ast.ReturnStmt{
		Argument: argument,
	}
	stmt.End = p.lexer.LastToken.End
	return &stmt
}

func (p *Parser) parseBreakStmt() *ast.Stmt {
	stmt := ast.Stmt{}
	stmt.Start = p.current.Start
	p.nextToken()

	var label *ast.Identifier
	if p.isToken(lexer.TTIdentifier) {
		label = newIdentifier(p.current)
		p.consumeVarId(true)
	}

	stmt.Node = &ast.BreakStmt{
		Label: label,
	}

	stmt.End = p.lexer.LastToken.End
	return &stmt
}

func (p *Parser) parseContinueStmt() *ast.Stmt {
	stmt := ast.Stmt{}
	stmt.Start = p.current.Start
	p.nextToken()

	var label *ast.Identifier
	if p.isToken(lexer.TTIdentifier) {
		label = newIdentifier(p.current)
		p.consumeVarId(true)
	}

	stmt.Node = &ast.ContinueStmt{
		Label: label,
	}

	stmt.End = p.lexer.LastToken.End
	return &stmt
}

func (p *Parser) parseExprStmt() *ast.Stmt {
	expr := p.parseExpr()

	return &ast.Stmt{
		Node: &ast.ExprStmt{
			Expression: expr,
		},
		Position: expr.Position,
	}
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
