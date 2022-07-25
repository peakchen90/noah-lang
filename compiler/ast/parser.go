package ast

import (
	"fmt"
)

type Parser struct {
	source        []rune // utf-8 字符
	lexer         *Lexer // 词法分析器
	current       *Token // 当前 token
	isSeenNewline bool   // 读取下一个 token 时是否遇到过换行
	blockLevel    int    // 当前进入到第几层块级作用域
	loopLevel     int    // 当前进入到第几层循环块
}

func NewParser(input string) *Statement {
	source := []rune(input)
	parser := Parser{
		source: source,
		lexer:  NewLexer(source),
	}
	return parser.parse()
}

func (p *Parser) parse() *Statement {
	body := make([]Statement, 0, 5)

	p.nextToken()

	for !p.isToken(TTEof) {
		stmt := p.parseStatement()
		body = append(body, *stmt)
	}

	node := Statement{
		Node: &Program{
			Body: body,
		},
	}

	if len(body) > 0 {
		node.Start = body[0].Start
		node.End = body[len(body)-1].End
	}

	return &node
}

func (p *Parser) nextToken() {
	p.current = p.lexer.readNext()
}

func (p *Parser) panicWithError(pos int, msg string) {
	line, column := PrintErrorFrame(p.source, pos, msg)
	panic(fmt.Sprintf("%s (%d:%d)", msg, line, column))
}

func (p *Parser) unexpectedPos(index int) {
	var message string
	if index < len(p.source) {
		message = "Unexpected token " + string(p.source[index])
	} else {
		message = "Unexpected end of file"
	}
	p.panicWithError(index, message)
}

func (p *Parser) unexpectedToken(token *Token, msg string) {
	var message string
	if len(msg) > 0 {
		message = msg
	} else {
		switch token.Type {
		case TTEof:
			message = "Unexpected end of file"
		case TTString, TTNumber, TTConst:
			message = "Unexpected literal: " + token.String()
		default:
			message = "Unexpected token " + token.String()
		}
	}
	p.panicWithError(token.Start, message)
}

func (p *Parser) unexpected() {
	p.unexpectedToken(p.current, "")
}

func (p *Parser) isToken(tokenType TokenType) bool {
	return p.current.Type == tokenType
}

// 消费一个 token 类型，如果消费成功，返回 true 并读取下一个 token，否则返回 false
func (p *Parser) consume(tokenType TokenType) bool {
	if p.isToken(tokenType) {
		p.nextToken()
		return true
	}
	return false
}

// 期待当前 token 为指定类型，否则抛错
func (p *Parser) expect(tokenType TokenType) {
	if !p.isToken(tokenType) {
		p.unexpected()
	}
}
