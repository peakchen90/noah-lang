package lexer

import (
	"fmt"
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"strconv"
	"strings"
)

type Lexer struct {
	source       []rune // utf-8 字符
	index        int    // 光标位置
	allowExpr    bool   // 当前上下文是否允许表达式
	SeenNewline  bool   // 读取下一个 token 时前面是否遇到过换行符
	CurrentToken *Token // 当前的 token
	LastToken    *Token // 上一个 token
}

func NewLexer(source []rune) *Lexer {
	lexer := Lexer{
		source:    source,
		index:     0,
		allowExpr: true,
	}
	return &lexer
}

func (l *Lexer) Next() *Token {
	l.SeenNewline = false
	l.skipSpace()
	l.skipComment()

	var token *Token
	ch := l.Look(0)

	if helper.IsIdentifierChar(ch, true) {
		token = l.readAsIdentifier()
	} else if ch >= '0' && ch <= '9' {
		token = l.readAsNumber()
	} else if l.index == len(l.source) {
		token = l.createToken(TTEof, l.index, l.index)
	} else {
		switch ch {
		case '"':
			token = l.readAsString(false)
		case '`':
			token = l.readAsString(true)
		case '\'':
			token = l.readAsChar()
		case '=':
			if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTEq, l.index-2, l.index)
			} else {
				l.index++
				token = l.createToken(TTAssign, l.index-1, l.index)
			}
		case '+':
			if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTPlusAssign, l.index-2, l.index)
			} else if l.allowExpr || l.SeenNewline {
				if l.Look(1) == '+' {
					l.index += 2
					token = l.createToken(TTPrefixInc, l.index-2, l.index)
				} else {
					l.index++
					token = l.createToken(TTPrefixPlus, l.index-1, l.index)
				}
			} else {
				if l.Look(1) == '+' {
					l.index += 2
					token = l.createToken(TTPostfixInc, l.index-2, l.index)
				} else {
					l.index++
					token = l.createToken(TTPlus, l.index-1, l.index)
				}
			}
		case '-':
			if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTSubAssign, l.index-2, l.index)
			} else if l.Look(1) == '>' {
				l.index += 2
				token = l.createToken(TTReturnSym, l.index-2, l.index)
			} else if l.allowExpr || l.SeenNewline {
				if l.Look(1) == '-' {
					l.index += 2
					token = l.createToken(TTPrefixDec, l.index-2, l.index)
				} else {
					l.index++
					token = l.createToken(TTPrefixSub, l.index-1, l.index)
				}
			} else {
				if l.Look(1) == '-' {
					l.index += 2
					token = l.createToken(TTPostfixDec, l.index-2, l.index)
				} else {
					l.index++
					token = l.createToken(TTSub, l.index-1, l.index)
				}
			}
		case '*':
			if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTMulAssign, l.index-2, l.index)
			} else {
				l.index++
				token = l.createToken(TTMul, l.index-1, l.index)
			}
		case '/':
			if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTDivAssign, l.index-2, l.index)
			} else {
				l.index++
				token = l.createToken(TTDiv, l.index-1, l.index)
			}
		case '%':
			if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTRemAssign, l.index-2, l.index)
			} else {
				l.index++
				token = l.createToken(TTRem, l.index-1, l.index)
			}
		case '<':
			if l.Look(1) == '-' {
				l.index += 2
				token = l.createToken(TTExtendSym, l.index-2, l.index)
			} else if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTLe, l.index-2, l.index)
			} else if l.Look(1) == '<' {
				if l.Look(2) == '=' {
					l.index += 3
					token = l.createToken(TTBitLeftShiftAssign, l.index-3, l.index)
				} else {
					l.index += 2
					token = l.createToken(TTBitLeftShift, l.index-2, l.index)
				}
			} else {
				l.index++
				token = l.createToken(TTLt, l.index-1, l.index)
			}
		case '>':
			if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTGe, l.index-2, l.index)
			} else if l.Look(1) == '>' {
				if l.Look(2) == '=' {
					l.index += 3
					token = l.createToken(TTBitRightShiftAssign, l.index-3, l.index)
				} else {
					l.index += 2
					token = l.createToken(TTBitRightShift, l.index-2, l.index)
				}
			} else {
				l.index++
				token = l.createToken(TTGt, l.index-1, l.index)
			}
		case '&':
			if l.Look(1) == '&' {
				l.index += 2
				token = l.createToken(TTLogicAnd, l.index-2, l.index)
			} else if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTBitAndAssign, l.index-2, l.index)
			} else {
				l.index++
				token = l.createToken(TTBitAnd, l.index-1, l.index)
			}
		case '|':
			if l.Look(1) == '|' {
				l.index += 2
				token = l.createToken(TTLogicOr, l.index-2, l.index)
			} else if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTBitOrAssign, l.index-2, l.index)
			} else {
				l.index++
				token = l.createToken(TTBitOr, l.index-1, l.index)
			}
		case '!':
			if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTNe, l.index-2, l.index)
			} else {
				l.index++
				token = l.createToken(TTLogicNot, l.index-1, l.index)
			}
		case '~':
			l.index++
			token = l.createToken(TTBitNot, l.index-1, l.index)
		case '^':
			if l.Look(1) == '=' {
				l.index += 2
				token = l.createToken(TTBitXorAssign, l.index-2, l.index)
			} else {
				l.index++
				token = l.createToken(TTBitXor, l.index-1, l.index)
			}
		case '(':
			l.index++
			token = l.createToken(TTParenL, l.index-1, l.index)
		case ')':
			l.index++
			token = l.createToken(TTParenR, l.index-1, l.index)
		case '[':
			l.index++
			token = l.createToken(TTBracketL, l.index-1, l.index)
		case ']':
			l.index++
			token = l.createToken(TTBracketR, l.index-1, l.index)
		case '{':
			l.index++
			token = l.createToken(TTBraceL, l.index-1, l.index)
		case '}':
			l.index++
			token = l.createToken(TTBraceR, l.index-1, l.index)
		case ',':
			l.index++
			token = l.createToken(TTComma, l.index-1, l.index)
		case '.':
			if l.Look(1) == '.' && l.Look(2) == '.' {
				l.index += 3
				token = l.createToken(TTRest, l.index-3, l.index)
			} else {
				l.index++
				token = l.createToken(TTDot, l.index-1, l.index)
			}
		case ';':
			l.index++
			token = l.createToken(TTSemi, l.index-1, l.index)
		case ':':
			l.index++
			token = l.createToken(TTColon, l.index-1, l.index)
		default:
			l.unexpected(l.index, "")
		}
	}

	l.LastToken = l.CurrentToken
	l.CurrentToken = token

	return token
}

func (l *Lexer) Look(n int) rune {
	next := l.index + n
	if next < len(l.source) {
		return l.source[next]
	}
	return 0
}

func (l *Lexer) LookNext() rune {
	l.skipSpace()
	l.skipComment()
	if l.checkIndex() {
		return l.source[l.index]
	}
	return 0
}

func (l *Lexer) createToken(tokenType TokenType, start int, end int) *Token {
	tokenMeta := &tokenMetaTable[tokenType]
	token := Token{TokenMeta: tokenMeta}
	token.Position = *ast.NewPosition(start, end)

	if tokenType != TTComment {
		l.allowExpr = tokenMeta.AllowExpr
	}

	return &token
}

func (l *Lexer) checkIndex() bool {
	return l.index < len(l.source)
}

func (l *Lexer) skipSpace() {
	for l.checkIndex() {
		ch := l.Look(0)
		if ch == '\r' || ch == '\n' || ch == '\t' || ch == ' ' {
			if ch == '\r' {
				if l.Look(1) == '\n' { // \r\n
					l.index++
				}
				l.SeenNewline = true
			} else if ch == '\n' {
				l.SeenNewline = true
			}
			l.index++
		} else {
			break
		}
	}
}

func (l *Lexer) skipComment() {
	if l.Look(0) == '/' && l.Look(1) == '/' {
		l.index += 2
		for l.checkIndex() && l.Look(0) != '\n' {
			l.index++
		}
		l.skipSpace()
		l.skipComment()
	} else if l.Look(0) == '/' && l.Look(1) == '*' {
		l.index += 2
		for l.checkIndex() && !(l.Look(0) == '*' && l.Look(1) == '/') {
			ch := l.Look(0)
			if ch == '\r' || ch == '\n' {
				l.SeenNewline = true
			}
			l.index++
		}
		l.index += 2
		l.skipSpace()
		l.skipComment()
	}
}

func (l *Lexer) readAsString(raw bool) *Token {
	start := l.index
	valid := false
	value := strings.Builder{}
	l.index++

	for l.checkIndex() {
		ch := l.Look(0)
		if raw && ch == '`' {
			valid = true
			break
		} else if !raw && ch == '"' {
			valid = true
			break
		}

		// 换行
		if ch == '\n' && !raw {
			l.unexpected(
				l.index,
				"The string literal cannot wrap. Tip: You can use the template string `...`",
			)
		}

		// escape char
		// see: https://baike.baidu.com/item/%E8%BD%AC%E4%B9%89%E5%AD%97%E7%AC%A6/86397
		if ch == '\\' {
			l.index++
			switch l.Look(0) {
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
			case '?':
				value.WriteByte(63)
			case 'x': // \xhh 2位十六进制字符
				l.index++
				ch1 := l.Look(0)
				ch2 := l.Look(1)
				if ((ch1 >= '0' && ch1 <= '9') || (ch1 >= 'a' && ch1 <= 'f') || (ch1 >= 'A' && ch1 <= 'F')) &&
					(ch2 >= '0' && ch2 <= '9') || (ch2 >= 'a' && ch2 <= 'f') || (ch2 >= 'A' && ch2 <= 'F') {
					l.index++
					code, _ := strconv.ParseUint(string([]rune{ch1, ch2}), 16, 8)
					value.WriteRune(rune(code))
				} else {
					l.unexpected(l.index, "Invalid hexadecimal escape sequence")
				}
			default:
				// \ddd 1~3位八进制字符
				str := strings.Builder{}
				for i := 0; i < 3; i++ {
					ch := l.Look(i)
					if ch >= '0' && ch <= '7' {
						str.WriteRune(ch)
					} else {
						break
					}
				}
				if str.Len() > 0 {
					l.index += str.Len() - 1
					code, _ := strconv.ParseUint(str.String(), 8, 16)
					value.WriteRune(rune(code))
				} else {
					value.WriteRune(l.Look(0))
				}
			}

			l.index++
			continue
		}

		value.WriteRune(l.Look(0))
		l.index++
	}

	if !valid {
		if raw {
			l.unexpected(l.index, "Invalid template string")
		} else {
			l.unexpected(l.index, "Invalid string literal")
		}
	}

	l.index++

	if !valid {
		l.unexpected(l.index, "Invalid string literal")
	}
	token := l.createToken(TTString, start, l.index)
	token.Value = value.String()
	if raw {
		token.Flag = "raw"
	}
	return token
}

func (l *Lexer) readAsChar() *Token {
	start := l.index
	value := ""
	l.index++

	// escape char
	if l.Look(0) == '\\' {
		l.index++
		switch l.Look(0) {
		case 'a':
			value = "\a"
		case 'b':
			value = "\b"
		case 'f':
			value = "\f"
		case 'n':
			value = "\n"
		case 'r':
			value = "\r"
		case 't':
			value = "\t"
		case 'v':
			value = "\v"
		case '?':
			value = string(rune(63))
		default:
			value = string(l.Look(0))
		}
	} else {
		value = string(l.Look(0))
	}

	l.index++
	if l.Look(0) != '\'' {
		l.unexpected(start, "Invalid char literal")
	}

	token := l.createToken(TTChar, start, l.index)
	token.Value = value

	l.index++
	return token
}

func (l *Lexer) readAsNumber() *Token {
	start := l.index
	valid := true
	seenDot := false
	consumeNum := true
	value := strings.Builder{}

	for l.checkIndex() {
		ch := l.Look(0)
		if ch >= '0' && ch <= '9' {
			value.WriteRune(ch)
			consumeNum = false
			l.index++
		} else if !consumeNum && ch == '.' {
			if seenDot {
				valid = false
				break
			}
			seenDot = true
			consumeNum = true
			value.WriteRune(ch)
			l.index++
		} else if consumeNum { // 再次检查是否是消费掉一个数字
			valid = false
			break
		} else {
			break
		}
	}

	if !valid || consumeNum {
		l.unexpected(l.index, "unexpected number")
	}

	token := l.createToken(TTNumber, start, l.index)
	token.Value = value.String()
	return token
}

func (l *Lexer) readAsIdentifier() *Token {
	start := l.index
	value := strings.Builder{}

	for l.checkIndex() {
		ch := l.Look(0)
		if helper.IsIdentifierChar(ch, false) {
			value.WriteRune(ch)
			l.index++
		} else {
			break
		}
	}

	valueStr := value.String()
	var token *Token

	if valueStr == "is" {
		token = l.createToken(TTIsOp, start, l.index)
	} else if valueStr == "as" {
		token = l.createToken(TTAsOp, start, l.index)
	} else if isKeyword(valueStr) {
		token = l.createToken(TTKeyword, start, l.index)
		if valueStr == "return" {
			l.allowExpr = true
		}
	} else if isBuiltInConstant(valueStr) {
		token = l.createToken(TTConst, start, l.index)
	} else {
		token = l.createToken(TTIdentifier, start, l.index)
	}

	token.Value = valueStr

	return token
}

func (l *Lexer) unexpected(index int, msg string) {
	var message string
	if len(msg) > 0 {
		message = msg
	} else if index < len(l.source) {
		message = "unexpected token " + string(l.source[index])
	} else {
		message = "unexpected end of file"
	}
	line, column := helper.PrintErrorFrame(l.source, index, message)
	panic(fmt.Sprintf("%s (%d:%d)", msg, line, column))
}
