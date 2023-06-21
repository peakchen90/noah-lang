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
	//case lexer.TTBraceL:
	//	stmt = p.parseBlockStmt()
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
	stmt := &ast.Stmt{}
	if pubToken != nil {
		stmt.Start = pubToken.Start
	} else {
		stmt.Start = p.current.Start
	}
	p.nextToken()

	firstToken := p.consume(lexer.TTIdentifier, false)
	if firstToken == nil {
		p.unexpectedMissing("function name")
	}

	var impl *ast.KindExpr

	if !p.isToken(lexer.TTParenL) {
		impl = p.parseMaybeChainKindExpr(&ast.KindExpr{
			Node:     &ast.TypeId{Name: NewKindIdentifier(firstToken)},
			Position: firstToken.Position,
		})
		p.consume(lexer.TTColon, true)
	}

	nameToken := p.consume(lexer.TTIdentifier, false)
	if nameToken == nil {
		nameToken = firstToken
		firstToken = nil
	}

	funcKind := p.parseFuncKindExpr(stmt.Start)

	funcDecl := &ast.FuncDecl{
		Name:     NewIdentifier(nameToken),
		Impl:     impl,
		FuncKind: funcKind,
		Body:     p.parseBlockStmt(),
		Pubic:    pubToken != nil,
	}

	stmt.Node = funcDecl
	stmt.End = p.lexer.LastToken.End

	return stmt
}

func (p *Parser) parseVarDecl(isConst bool, pubToken *lexer.Token) *ast.Stmt {
	stmt := &ast.Stmt{}
	if pubToken != nil {
		stmt.Start = pubToken.Start
	} else {
		stmt.Start = p.current.Start
	}
	p.nextToken()

	// id
	token := p.consume(lexer.TTIdentifier, false)
	if token == nil {
		p.unexpectedMissing("variable name")
	}
	if IsReservedType(token.Value) {
		p.unexpectedPos(token.Start, "Reserved type cannot be used: "+token.Value)
	}
	id := NewIdentifier(token)

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
		Pubic: pubToken != nil,
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
	if IsReservedType(p.current.Value) {
		p.unexpectedPos(p.current.Start, "Reserved type cannot be used: "+p.current.Value)
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
	if IsReservedType(p.current.Value) {
		p.unexpectedPos(p.current.Start, "Reserved type cannot be used: "+p.current.Value)
	}
	name := NewKindIdentifier(p.current)
	p.nextToken()

	// extends
	var Extends *ast.KindExpr
	if p.consume(lexer.TTExtendSym, false) != nil {
		p.expect(lexer.TTIdentifier)
		Extends = p.parseKindExpr()
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
	if IsReservedType(p.current.Value) {
		p.unexpectedPos(p.current.Start, "Reserved type cannot be used: "+p.current.Value)
	}
	name := NewKindIdentifier(p.current)
	p.nextToken()

	// extends
	var extends *ast.KindExpr
	if p.consume(lexer.TTExtendSym, false) != nil {
		p.expect(lexer.TTIdentifier)
		extends = p.parseKindExpr()
	}

	// impl
	var impl *ast.KindExpr
	if p.consume(lexer.TTColon, false) != nil {
		p.expect(lexer.TTIdentifier)
		impl = p.parseKindExpr()
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
	if IsReservedType(p.current.Value) {
		p.unexpectedPos(p.current.Start, "Reserved type cannot be used: "+p.current.Value)
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
		label = NewIdentifier(labelToken)
	}

	hasParentheses := p.consume(lexer.TTParenL, false) != nil // `(`

	if p.isToken(lexer.TTIdentifier) {
		headToken := p.current
		p.nextToken()

		if p.isToken(lexer.TTComma) || p.isToken(lexer.TTColon) { // for value, key: target {}
			var key *ast.Identifier
			if p.isToken(lexer.TTComma) {
				p.nextToken()
				key = NewIdentifier(p.consume(lexer.TTIdentifier, true))
			}

			p.consume(lexer.TTColon, true)

			eachVisitor = &ast.EachVisitor{
				Value:  NewIdentifier(headToken),
				Key:    key,
				Target: p.parseExpr(),
			}
		} else { // for value {}
			test = p.parseIdentifierExpr(headToken)
		}
	} else if p.isKeyword("let") || p.isKeyword("const") { // for (init; test; update) {}
		init = p.parseVarDecl(len(p.current.Value) == 5, nil)
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
		label = NewIdentifier(p.current)
		p.nextToken()
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
		label = NewIdentifier(p.current)
		p.nextToken()
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
