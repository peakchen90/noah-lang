package ast

type Parser struct {
	source            []rune // utf-8 字符
	index             int    // 光标位置
	isStart           bool   // 光标是否在开始位置
	isSeenNewline     bool   // 读取下一个 token 时是否遇到过换行
	currentToken      *Token // 当前 token
	allowExpr         bool   // 当前上下文是否允许表达式
	currentBlockLevel int    // 当前进入到第几层块级作用域
	currentLoopLevel  int    // 当前进入到第几层循环块
}

func NewParser(input string) *Statement {
	parser := Parser{
		source: []rune(input),
	}
	return parser.parse()
}

func (p *Parser) parse() *Statement {
	body := make([]Statement, 0, 5)

	p.readNextToken()
	for p.checkIndexRange() {
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

func (p *Parser) checkIndexRange() bool {
	return p.index < len(p.source)
}

func (p *Parser) unexpectedPos(index int) {
}

func (p *Parser) unexpected() {
}
