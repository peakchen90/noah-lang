package parser

import (
	"fmt"
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"github.com/peakchen90/noah-lang/internal/lexer"
)

type Parser struct {
	source     []rune       // utf-8 字符
	lexer      *lexer.Lexer // 词法分析器
	current    *lexer.Token // 当前 token
	seenToken  *lexer.Token // 缓存的后一个 token
	blockLevel int          // 当前进入到第几层块级作用域
	loopLevel  int          // 当前进入到第几层循环块
}

func NewParser(input string) *ast.File {
	source := []rune(input)
	parser := Parser{
		source: source,
		lexer:  lexer.NewLexer(source),
	}
	return parser.parse()
}

func (p *Parser) parse() *ast.File {
	body := make([]*ast.Stmt, 0, helper.DefaultCap)
	p.nextToken()

	for !p.isEnd() {
		stmt := p.parseStmt()
		body = append(body, stmt)
	}

	node := ast.File{Body: body}

	if len(body) > 0 {
		node.Start = body[0].Start
		node.End = body[len(body)-1].End
	}

	return &node
}

func (p *Parser) nextToken() *lexer.Token {
	if p.seenToken != nil {
		p.current = p.seenToken
		p.seenToken = nil
	} else {
		p.current = p.lexer.Next()
	}
	return p.current
}

func (p *Parser) revertLastToken() {
	if p.seenToken == nil {
		p.seenToken = p.current
		p.current = p.lexer.LastToken
	}
}

func (p *Parser) unexpectedPos(index int, msg string) {
	line, column := helper.PrintErrorFrame(p.source, index, msg)
	panic(fmt.Sprintf("%s (%d:%d)", msg, line, column))
}

func (p *Parser) unexpectedToken(expectHint string, receiveToken *lexer.Token) {
	p.unexpectedPos(
		receiveToken.Start,
		fmt.Sprintf("Expected %s, found %s", expectHint, receiveToken.String()),
	)
}

func (p *Parser) unexpectedMissing(missingHint string) {
	p.unexpectedPos(p.lexer.LastToken.End, fmt.Sprintf("Missing %s", missingHint))
}

func (p *Parser) unexpected() {
	var message string
	token := p.current

	switch token.Type {
	case lexer.TTEof:
		message = "unexpected end of file"
	case lexer.TTString, lexer.TTNumber, lexer.TTConst:
		message = "unexpected literal: " + token.String()
	default:
		message = "unexpected token " + token.String()
	}

	p.unexpectedPos(token.Start, message)
}

// 判断当前是否为指定 token 类型
func (p *Parser) isToken(tokenType lexer.TokenType) bool {
	return p.current.Type == tokenType
}

// 判断当前是否是指定名称的关键字
func (p *Parser) isKeyword(name string) bool {
	if p.isToken(lexer.TTKeyword) && p.current.Value == name {
		return true
	}
	return false
}

// 是否结束
func (p *Parser) isEnd() bool {
	return p.isToken(lexer.TTEof)
}

// 消费一个 token 类型，如果消费成功返回 token 并读取下一个 token，否则返回 nil
func (p *Parser) consume(tokenType lexer.TokenType, isPanic bool) *lexer.Token {
	if p.isToken(tokenType) {
		token := p.current
		p.nextToken()
		return token
	}
	if isPanic {
		p.unexpected()
	}
	return nil
}

// 消费一个 keyword token 类型，如果消费成功返回 token 并读取下一个 token，否则返回 nil
func (p *Parser) consumeKeyword(name string, isPanic bool) *lexer.Token {
	if p.isKeyword(name) {
		token := p.current
		p.nextToken()
		return token
	}
	if isPanic {
		p.unexpected()
	}
	return nil
}

// 期待当前 token 为指定类型，否则抛错
func (p *Parser) expect(tokenType lexer.TokenType) {
	if !p.isToken(tokenType) {
		p.unexpected()
	}
}
