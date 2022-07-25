package ast

import (
	"fmt"
)

type Parser struct {
	source            []rune // utf-8 字符
	index             int    // 光标位置
	isSeenNewline     bool   // 读取下一个 token 时是否遇到过换行
	currentToken      *Token // 当前 token
	allowExpr         bool   // 当前上下文是否允许表达式
	currentBlockLevel int    // 当前进入到第几层块级作用域
	currentLoopLevel  int    // 当前进入到第几层循环块
}

func NewParser(input string) *Statement {
	parser := Parser{
		source:    []rune(input),
		allowExpr: true,
	}
	return parser.parse()
}

func (p *Parser) parse() *Statement {
	body := make([]Statement, 0, 5)

	p.readNextToken()

	for p.checkIndex() {
		stmt := p.parseStatement()
		body = append(body, *stmt)
	}

	node := Statement{
		Data: &Program{
			Body: body,
		},
	}

	if len(body) > 0 {
		node.Start = body[0].Start
		node.End = body[len(body)-1].End
	}

	return &node
}

func (p *Parser) checkIndex() bool {
	return p.index < len(p.source)
}

func (p *Parser) panicWithError(pos int, msg string) {
	source := string(p.source)
	line, column := PrintErrorFrame(&source, pos, msg)
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
	} else if token.Type == TTEof {
		message = "Unexpected end of file"
	} else {
		message = "Unexpected token " + token.String()
	}
	p.panicWithError(token.Start, message)
}

func (p *Parser) unexpected() {
	p.unexpectedToken(p.currentToken, "")
}
