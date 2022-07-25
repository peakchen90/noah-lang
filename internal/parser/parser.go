package parser

import (
	"github.com/peakchen90/hera-lang/internal/ast"
	"github.com/peakchen90/hera-lang/internal/lexer"
)

type Parser struct {
	source        []rune       // utf-8 字符
	lexer         *lexer.Lexer // 词法分析器
	current       *lexer.Token // 当前 token
	isSeenNewline bool         // 读取下一个 token 时是否遇到过换行
	blockLevel    int          // 当前进入到第几层块级作用域
	loopLevel     int          // 当前进入到第几层循环块
}

func NewParser(input string) *ast.Statement {
	source := []rune(input)
	parser := Parser{
		source: source,
		lexer:  lexer.NewLexer(source),
	}
	return parser.parse()
}

func (p *Parser) parse() *ast.Statement {
	body := make([]ast.Statement, 0, 5)

	p.nextToken()

	for !p.isToken(lexer.TTEof) {
		stmt := p.parseStatement()
		body = append(body, *stmt)
	}

	node := ast.Statement{
		Node: &ast.Program{
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
	p.current = p.lexer.Next()
}

func (p *Parser) unexpectedToken(token *lexer.Token, msg string) {
	var message string
	if len(msg) > 0 {
		message = msg
	} else {
		switch token.Type {
		case lexer.TTEof:
			message = "Unexpected end of file"
		case lexer.TTString, lexer.TTNumber, lexer.TTConst:
			message = "Unexpected literal: " + token.String()
		default:
			message = "Unexpected token " + token.String()
		}
	}
	p.lexer.Unexpected(token.Start, message)
}

func (p *Parser) unexpected() {
	p.unexpectedToken(p.current, "")
}

func (p *Parser) isToken(tokenType lexer.TokenType) bool {
	return p.current.Type == tokenType
}

// 消费一个 token 类型，如果消费成功，返回 true 并读取下一个 token，否则返回 false
func (p *Parser) consume(tokenType lexer.TokenType) bool {
	if p.isToken(tokenType) {
		p.nextToken()
		return true
	}
	return false
}

// 期待当前 token 为指定类型，否则抛错
func (p *Parser) expect(tokenType lexer.TokenType) {
	if !p.isToken(tokenType) {
		p.unexpected()
	}
}
