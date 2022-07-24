package ast

type TokenType = uint8

const (
	TTBegin      TokenType = iota + 1 // 初始Token
	TTEof                             // 结束 Token
	TTComment                         // 注释
	TTKeyword                         // 关键字
	TTIdentifier                      // 标识符
	TTNumber                          // 数字
	TTString                          // 字符串
	TTBool                            // 布尔值
	TTAssign                          // =
	TTPlus                            // +
	TTSub                             // -
	TTMul                             // *
	TTDiv                             // /
	TTRem                             // %
	TTLt                              // <
	TTLe                              // <=
	TTGt                              // >
	TTGe                              // >=
	TTEq                              // ==
	TTNe                              // !=
	TTLogicAnd                        // &&
	TTLogicOr                         // ||
	TTLogicNot                        // !
	TTBitAnd                          // &
	TTBitOr                           // |
	TTBitNot                          // ~
	TTBitXor                          // ^
	TTParenL                          // (
	TTParenR                          // )
	TTBracketL                        // [
	TTBracketR                        // ]
	TTBraceL                          // {
	TTBraceR                          // }
	TTComma                           // ,
	TTDot                             // .
	TTRest                            // ..
	TTSemi                            // ;
	TTColon                           // :
	TTStar                            // *
	TTReturnSym                       // ->
)

type Token struct {
	Type       TokenType
	Value      string
	Precedence int8
	Start      int
	End        int
}

func (p *Parser) newToken(tokenType TokenType, value string, start int, end int) *Token {
	token := Token{
		Type:       tokenType,
		Value:      value,
		Precedence: -1,
		Start:      start,
		End:        end,
	}
	return &token
}

func (p *Parser) readNextToken() *Token {
	p.isStart = false
	p.skipSpace(true)
	p.skipComment()

	char := p.source[p.index]
	var token *Token

	if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || char == '_' || char == '$' {
		token = p.readAsIdentifier()
	} else if char >= '0' && char <= '9' {
		token = p.readAsNumber()
	} else if p.index == len(p.source) {
		token = p.newToken(TTEof, "EOF", p.index, p.index)
	} else {
		switch char {
		case '"':
			token = p.readAsString()
		case '=':
			if p.lookBehind(1) == '=' {
				p.index += 2
				token = p.newToken(TTEq, "==", p.index-2, p.index)
				token.Precedence = 11
			} else {
				p.index++
				token = p.newToken(TTAssign, "=", p.index-1, p.index)
				token.Precedence = 1
			}
		case '+':
			p.index++
			token = p.newToken(TTPlus, "+", p.index-1, p.index)
			token.Precedence = 14
		case '-':
			if p.lookBehind(1) == '>' {
				p.index += 2
				token = p.newToken(TTReturnSym, "->", p.index-2, p.index)
			} else if p.allowExpr {
				token = p.readAsNumber()
			} else {
				p.index++
				token = p.newToken(TTSub, "-", p.index-1, p.index)
				token.Precedence = 14
			}
		case '*':
			p.index++
			if p.allowExpr {
				token = p.newToken(TTStar, "*", p.index-1, p.index)
			} else {
				token = p.newToken(TTMul, "*", p.index-1, p.index)
				token.Precedence = 15
			}
		case '/':
			p.index++
			token = p.newToken(TTDiv, "/", p.index-1, p.index)
			token.Precedence = 15
		case '%':
			p.index++
			token = p.newToken(TTRem, "%", p.index-1, p.index)
			token.Precedence = 15
		case '<':
			if p.lookBehind(1) == '=' {
				p.index += 2
				token = p.newToken(TTLe, "<=", p.index-2, p.index)
				token.Precedence = 12
			} else {
				p.index++
				token = p.newToken(TTLt, "<", p.index-1, p.index)
				token.Precedence = 12
			}
		case '>':
			if p.lookBehind(1) == '=' {
				p.index += 2
				token = p.newToken(TTGe, ">=", p.index-2, p.index)
				token.Precedence = 12
			} else {
				p.index++
				token = p.newToken(TTGt, ">", p.index-1, p.index)
				token.Precedence = 12
			}
		case '&':
			if p.lookBehind(1) == '&' {
				p.index += 2
				token = p.newToken(TTLogicAnd, "&&", p.index-2, p.index)
				token.Precedence = 7
			} else if p.allowExpr {
				// TODO: object reference
			} else {
				p.index++
				token = p.newToken(TTBitAnd, "&", p.index-1, p.index)
				token.Precedence = 10
			}
		case '|':
			if p.lookBehind(1) == '|' {
				p.index += 2
				token = p.newToken(TTLogicOr, "||", p.index-2, p.index)
				token.Precedence = 6
			} else {
				p.index++
				token = p.newToken(TTBitOr, "|", p.index-1, p.index)
				token.Precedence = 8
			}
		case '!':
			if p.lookBehind(1) == '=' {
				p.index += 2
				token = p.newToken(TTNe, "!=", p.index-2, p.index)
				token.Precedence = 11
			} else {
				p.index++
				token = p.newToken(TTLogicNot, "!", p.index-1, p.index)
				token.Precedence = 17
			}
		case '~':
			p.index++
			token = p.newToken(TTBitNot, "~", p.index-1, p.index)
			token.Precedence = 17
		case '^':
			p.index++
			token = p.newToken(TTBitXor, "^", p.index-1, p.index)
			token.Precedence = 9
		case '(':
			p.index++
			token = p.newToken(TTParenL, "(", p.index-1, p.index)
		case ')':
			p.index++
			token = p.newToken(TTParenR, ")", p.index-1, p.index)
		case '[':
			p.index++
			token = p.newToken(TTBracketL, "[", p.index-1, p.index)
		case ']':
			p.index++
			token = p.newToken(TTBracketR, "]", p.index-1, p.index)
		case '{':
			p.index++
			token = p.newToken(TTBraceL, "{", p.index-1, p.index)
		case '}':
			p.index++
			token = p.newToken(TTBraceR, "}", p.index-1, p.index)
		case ',':
			p.index++
			token = p.newToken(TTComma, ",", p.index-1, p.index)
		case '.':
			if p.lookBehind(1) == '.' {
				p.index += 2
				token = p.newToken(TTRest, "..", p.index-2, p.index)
			} else {
				p.index++
				token = p.newToken(TTDot, ".", p.index-1, p.index)
			}
		case ';':
			p.index++
			token = p.newToken(TTSemi, ";", p.index-1, p.index)
		case ':':
			p.index++
			token = p.newToken(TTColon, ":", p.index-1, p.index)
		default:
			p.unexpectedPos(p.index)
		}
	}

	p.currentToken = token
	return token
}

func (p *Parser) lookBehind(n int) rune {
	return p.source[p.index+n]
}

func (p *Parser) skipSpace(isSkipNewline bool) {

}

func (p *Parser) skipComment() {

}

func (p *Parser) readAsString() *Token {
	token := Token{}
	return &token
}

func (p *Parser) readAsNumber() *Token {
	token := Token{}
	return &token
}

func (p *Parser) readAsIdentifier() *Token {
	token := Token{}
	return &token
}

func (p *Parser) isToken(tokenType TokenType) bool {
	return p.currentToken.Type == tokenType
}

// 消费一个 token 类型，如果消费成功，返回 true 并读取下一个 token，否则返回 false
func (p *Parser) consume(tokenType TokenType) bool {
	if p.isToken(tokenType) {
		p.readNextToken()
		return true
	}
	return false
}
