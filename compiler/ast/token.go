package ast

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType = uint8

type TokenMeta struct {
	Type       TokenType
	Text       string
	Precedence int8
	AllowExpr  bool
}

// 更新 token 类型同时要更新 tokenMetaMap
const (
	TTEof        TokenType = iota + 1 // 结束 Token
	TTComment                         // 注释
	TTKeyword                         // 关键字
	TTConst                           // 内置常量（关键字）
	TTIdentifier                      // 标识符
	TTNumber                          // 数字字面量
	TTString                          // 字符串字面量
	TTReturnSym                       // ->
	TTParenL                          // (
	TTParenR                          // )
	TTBracketL                        // [
	TTBracketR                        // ]
	TTBraceL                          // {
	TTBraceR                          // }
	TTRest                            // ..
	TTSemi                            // ;
	TTColon                           // :
	TTComma                           // ,

	// operator

	TTAssign   // =
	TTPlus     // +
	TTSub      // -
	TTMul      // *
	TTDiv      // /
	TTRem      // %
	TTLt       // <
	TTLe       // <=
	TTGt       // >
	TTGe       // >=
	TTEq       // ==
	TTNe       // !=
	TTLogicAnd // &&
	TTLogicOr  // ||
	TTLogicNot // !
	TTBitAnd   // &
	TTBitOr    // |
	TTBitNot   // ~
	TTBitXor   // ^
	TTDot      // .
	TTRef      // &
	TTUnref    // *
)

var tokenMetaMap = map[TokenType]TokenMeta{
	TTEof:        {TTEof, "EOF", -1, false},
	TTComment:    {TTComment, "Comment", -1, false},
	TTKeyword:    {TTKeyword, "Keyword", -1, false},
	TTIdentifier: {TTIdentifier, "Identifier", -1, false},
	TTNumber:     {TTNumber, "Number", -1, false},
	TTString:     {TTString, "String", -1, false},
	TTConst:      {TTConst, "Const", -1, false},
	TTReturnSym:  {TTReturnSym, "->", -1, false},
	TTParenL:     {TTParenL, "(", -1, true},
	TTParenR:     {TTParenR, ")", -1, false},
	TTBracketL:   {TTBracketL, "[", -1, true},
	TTBracketR:   {TTBracketR, "]", -1, true},
	TTBraceL:     {TTBraceL, "{", -1, true},
	TTBraceR:     {TTBraceR, "}", -1, true},
	TTRest:       {TTRest, "..", -1, true},
	TTSemi:       {TTSemi, ";", -1, true},
	TTColon:      {TTColon, ":", -1, true},
	TTComma:      {TTComma, ",", -1, true},

	// operator
	// precedence see: https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Operators/Operator_Precedence

	TTAssign:   {TTAssign, "=", 2, true},
	TTPlus:     {TTPlus, "+", 12, true},
	TTSub:      {TTSub, "-", 12, true},
	TTMul:      {TTMul, "*", 13, true},
	TTDiv:      {TTDiv, "/", 13, true},
	TTRem:      {TTRem, "%", 13, true},
	TTLt:       {TTLt, "<", 10, true},
	TTLe:       {TTLe, "<=", 10, true},
	TTGt:       {TTGt, ">", 10, true},
	TTGe:       {TTGe, ">=", 10, true},
	TTEq:       {TTEq, "==", 9, true},
	TTNe:       {TTNe, "!=", 9, true},
	TTLogicAnd: {TTLogicAnd, "&&", 5, true},
	TTLogicOr:  {TTLogicOr, "||", 4, true},
	TTLogicNot: {TTLogicNot, "!", 15, true},
	TTBitAnd:   {TTBitAnd, "&", 8, true},
	TTBitOr:    {TTBitOr, "|", 6, true},
	TTBitNot:   {TTBitNot, "~", 15, true},
	TTBitXor:   {TTBitXor, "^", 7, true},
	TTDot:      {TTDot, ".", 18, true},
	TTRef:      {TTRef, "&", 17, true},
	TTUnref:    {TTUnref, "*", 17, true},
}

type Token struct {
	*TokenMeta
	Position
	Value string
	Ext   interface{}
}

func (t *Token) String() string {
	if t.Type == TTString {
		return fmt.Sprintf(`"%s"`, t.Value)
	}
	if len(t.Value) > 0 {
		return t.Value
	}
	return t.Text
}

func (p *Parser) newToken(tokenType TokenType, start int, end int) *Token {
	tokenMeta := tokenMetaMap[tokenType]
	token := Token{TokenMeta: &tokenMeta}
	token.Start = start
	token.End = end

	if tokenType != TTComment {
		p.allowExpr = tokenMeta.AllowExpr
	}

	return &token
}

func (p *Parser) readNextToken() *Token {
	p.skipSpace()
	p.skipComment()

	var token *Token
	ch := p.look(0)

	if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' || ch == '$' {
		token = p.readAsIdentifier()
	} else if ch >= '0' && ch <= '9' {
		token = p.readAsNumber()
	} else if p.index == len(p.source) {
		token = p.newToken(TTEof, p.index, p.index)
	} else {
		switch ch {
		case '"':
			token = p.readAsString()
		case '=':
			if p.look(1) == '=' {
				p.index += 2
				token = p.newToken(TTEq, p.index-2, p.index)
			} else {
				p.index++
				token = p.newToken(TTAssign, p.index-1, p.index)
			}
		case '+':
			p.index++
			token = p.newToken(TTPlus, p.index-1, p.index)
		case '-':
			if p.look(1) == '>' {
				p.index += 2
				token = p.newToken(TTReturnSym, p.index-2, p.index)
			} else if p.allowExpr {
				token = p.readAsNumber()
			} else {
				p.index++
				token = p.newToken(TTSub, p.index-1, p.index)
			}
		case '*':
			p.index++
			if p.allowExpr {
				token = p.newToken(TTUnref, p.index-1, p.index)
			} else {
				token = p.newToken(TTMul, p.index-1, p.index)
			}
		case '/':
			p.index++
			token = p.newToken(TTDiv, p.index-1, p.index)
		case '%':
			p.index++
			token = p.newToken(TTRem, p.index-1, p.index)
		case '<':
			if p.look(1) == '=' {
				p.index += 2
				token = p.newToken(TTLe, p.index-2, p.index)
			} else {
				p.index++
				token = p.newToken(TTLt, p.index-1, p.index)
			}
		case '>':
			if p.look(1) == '=' {
				p.index += 2
				token = p.newToken(TTGe, p.index-2, p.index)
			} else {
				p.index++
				token = p.newToken(TTGt, p.index-1, p.index)
			}
		case '&':
			if p.look(1) == '&' {
				p.index += 2
				token = p.newToken(TTLogicAnd, p.index-2, p.index)
			} else if p.allowExpr {
				p.index++
				token = p.newToken(TTRef, p.index-1, p.index)
			} else {
				p.index++
				token = p.newToken(TTBitAnd, p.index-1, p.index)
			}
		case '|':
			if p.look(1) == '|' {
				p.index += 2
				token = p.newToken(TTLogicOr, p.index-2, p.index)
			} else {
				p.index++
				token = p.newToken(TTBitOr, p.index-1, p.index)
			}
		case '!':
			if p.look(1) == '=' {
				p.index += 2
				token = p.newToken(TTNe, p.index-2, p.index)
			} else {
				p.index++
				token = p.newToken(TTLogicNot, p.index-1, p.index)
			}
		case '~':
			p.index++
			token = p.newToken(TTBitNot, p.index-1, p.index)
		case '^':
			p.index++
			token = p.newToken(TTBitXor, p.index-1, p.index)
		case '(':
			p.index++
			token = p.newToken(TTParenL, p.index-1, p.index)
		case ')':
			p.index++
			token = p.newToken(TTParenR, p.index-1, p.index)
		case '[':
			p.index++
			token = p.newToken(TTBracketL, p.index-1, p.index)
		case ']':
			p.index++
			token = p.newToken(TTBracketR, p.index-1, p.index)
		case '{':
			p.index++
			token = p.newToken(TTBraceL, p.index-1, p.index)
		case '}':
			p.index++
			token = p.newToken(TTBraceR, p.index-1, p.index)
		case ',':
			p.index++
			token = p.newToken(TTComma, p.index-1, p.index)
		case '.':
			if p.look(1) == '.' {
				p.index += 2
				token = p.newToken(TTRest, p.index-2, p.index)
			} else {
				p.index++
				token = p.newToken(TTDot, p.index-1, p.index)
			}
		case ';':
			p.index++
			token = p.newToken(TTSemi, p.index-1, p.index)
		case ':':
			p.index++
			token = p.newToken(TTColon, p.index-1, p.index)
		default:
			p.unexpectedPos(p.index)
		}
	}

	p.currentToken = token
	return token
}

func (p *Parser) look(n int) rune {
	next := p.index + n
	if next < len(p.source) {
		return p.source[next]
	}
	return 0
}

func (p *Parser) skipSpace() {
	for p.checkIndex() {
		ch := p.look(0)
		if ch == '\r' || ch == '\n' || ch == '\t' || ch == ' ' {
			if ch == '\r' || ch == '\n' {
				p.isSeenNewline = true
			}
			p.index++
		} else {
			break
		}
	}
}

func (p *Parser) skipComment() {
	if p.look(0) == '/' && p.look(1) == '/' {
		p.index += 2
		for p.checkIndex() && p.look(0) != '\n' {
			p.index++
		}
		p.index++
		p.skipSpace()
		p.skipComment()
	} else if p.look(0) == '/' && p.look(1) == '*' {
		p.index += 2
		for p.checkIndex() && !(p.look(0) == '*' && p.look(1) == '/') {
			p.index++
		}
		p.index += 2
		p.skipSpace()
		p.skipComment()
	}
}

func (p *Parser) readAsString() *Token {
	start := p.index
	raw := false
	valid := false
	value := strings.Builder{}

	if p.look(1) == '"' && p.look(2) == '"' {
		p.index += 3
		raw = true
	} else {
		p.index++
	}

	for p.checkIndex() {
		ch := p.look(0)
		if ch == '"' {
			if !raw || (p.look(1) == '"' && p.look(2) == '"') {
				valid = true
				break
			}
		}
		// 换行
		if ch == '\n' && !raw {
			p.panicWithError(
				p.index,
				"String literals cannot wrap. Tip: You can use the raw string `\"\"\"...\"\"\"`",
			)
		}

		// escape char
		// see: https://baike.baidu.com/item/%E8%BD%AC%E4%B9%89%E5%AD%97%E7%AC%A6/86397
		if ch == '\\' {
			p.index++
			switch p.look(0) {
			case 'a':
				value.WriteByte('\a')
			case 'b':
				value.WriteByte('\b')
			case 'f':
				value.WriteByte('\f')
			case 'n':
				value.WriteByte('\n')
			case 'r':
				value.WriteByte('\r')
			case 't':
				value.WriteByte('\t')
			case 'v':
				value.WriteByte('\v')
			case '\\':
				value.WriteByte('\\')
			case '\'':
				value.WriteByte('\'')
			case '"':
				value.WriteByte('"')
			case '?':
				value.WriteByte(63)
			case 'x': // \xhh 2位十六进制字符
				p.index++
				ch1 := p.look(0)
				ch2 := p.look(1)
				if ((ch1 >= '0' && ch1 <= '9') || (ch1 >= 'a' && ch1 <= 'f') || (ch1 >= 'A' && ch1 <= 'F')) &&
					(ch2 >= '0' && ch2 <= '9') || (ch2 >= 'a' && ch2 <= 'f') || (ch2 >= 'A' && ch2 <= 'F') {
					p.index++
					code, _ := strconv.ParseUint(string([]rune{ch1, ch2}), 16, 8)
					value.WriteByte(byte(code))
				} else {
					p.panicWithError(p.index, "Invalid hexadecimal escape sequence")
				}
			default:
				// \ddd 1~3位八进制字符
				str := strings.Builder{}
				for i := 0; i < 3; i++ {
					ch := p.look(i)
					if ch >= '0' && ch <= '7' {
						str.WriteRune(ch)
					} else {
						break
					}
				}
				if str.Len() > 0 {
					p.index += str.Len() - 1
					code, _ := strconv.ParseUint(str.String(), 8, 8)
					value.WriteByte(byte(code))
				} else {
					value.WriteRune(p.look(0))
				}
			}

			p.index++
			continue
		}

		value.WriteRune(p.look(0))
		p.index++
	}

	if !valid {
		if raw {
			p.panicWithError(p.index, "The string literal is missing the terminator `\"\"\"`")
		} else {
			p.panicWithError(p.index, "The string literal is missing the terminator `\"`")
		}
	}

	if raw {
		p.index += 3
	} else {
		p.index++
	}

	token := p.newToken(TTString, start, p.index)
	token.Value = value.String()
	token.Ext = raw
	return token
}

func (p *Parser) readAsNumber() *Token {
	start := p.index
	valid := true
	seenDot := false
	consumeNum := true
	value := strings.Builder{}

	if p.look(0) == '-' {
		value.WriteByte('-')
		p.index++
	}

	for p.checkIndex() {
		ch := p.look(0)
		if ch >= '0' && ch <= '9' {
			value.WriteRune(ch)
			consumeNum = false
			p.index++
		} else if !consumeNum && ch == '.' {
			if seenDot {
				valid = false
				break
			}
			seenDot = true
			consumeNum = true
			value.WriteRune(ch)
			p.index++
		} else if consumeNum { // 再次检查是否是消费掉一个数字
			valid = false
			break
		} else {
			break
		}
	}

	if !valid || consumeNum {
		p.panicWithError(p.index, "Unexpected number")
	}

	token := p.newToken(TTNumber, start, p.index)
	token.Value = value.String()
	return token
}

func (p *Parser) readAsIdentifier() *Token {
	start := p.index
	value := strings.Builder{}

	for p.checkIndex() {
		ch := p.look(0)
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' || ch == '$' || (ch >= '0' && ch <= '9') {
			value.WriteRune(ch)
			p.index++
		}
	}

	valueStr := value.String()
	var token *Token

	if IsKeyword(valueStr) {
		switch valueStr {
		case "true", "false", "null", "self":
			token = p.newToken(TTConst, start, p.index)
		default:
			token = p.newToken(TTKeyword, start, p.index)
		}

	} else {
		token = p.newToken(TTIdentifier, start, p.index)
	}

	token.Value = valueStr

	return token
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
