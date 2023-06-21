package lexer

import (
	"fmt"
	"github.com/peakchen90/noah-lang/internal/ast"
)

type TokenType = uint8

type OpType = uint8

const (
	OpNone = 0b00000000

	OpBinary       = 0b00000001
	OpBinaryAssign = 0b00000011
	OpBinaryType   = 0b00000101

	OpUnary       = 0b00010000
	OpUnaryPrefix = 0b00110000
	OpUnarySuffix = 0b01010000
)

type TokenMeta struct {
	Type       TokenType
	Name       string
	Text       string
	Precedence int8
	OpType     OpType
	AllowExpr  bool
}

// 更新 token 类型同时要更新 tokenMetaTable
const (
	TTEof        TokenType = iota // 结束 Token
	TTComment                     // 注释
	TTKeyword                     // 关键字
	TTConst                       // 内置常量（关键字）
	TTIdentifier                  // 标识符
	TTNumber                      // 数字字面量
	TTString                      // 字符串字面量
	TTChar                        // 字符字面量
	TTReturnSym                   // ->
	TTExtendSym                   // <-
	TTParenL                      // (
	TTParenR                      // )
	TTBracketL                    // [
	TTBracketR                    // ]
	TTBraceL                      // {
	TTBraceR                      // }
	TTRest                        // ...
	TTSemi                        // ;
	TTColon                       // :
	TTComma                       // ,
	TTDot                         // .

	TTAssign   // =
	TTLogicOr  // ||
	TTLogicAnd // &&
	TTBitOr    // |
	TTBitXor   // ^
	TTBitAnd   // &
	TTEq       // ==
	TTNe       // !=
	TTLt       // <
	TTLe       // <=
	TTGt       // >
	TTGe       // >=
	TTIsOp     // `is`
	TTPlus     // +
	TTSub      // -
	TTMul      // *
	TTDiv      // /
	TTRem      // %

	TTUnaryPlus // +
	TTUnarySub  // -
	TTLogicNot  // !
	TTBitNot    // ~
)

// precedence see: https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Operators/Operator_Precedence
var tokenMetaTable = [43]TokenMeta{
	TTEof:        {TTEof, "TTEof", "", -1, OpNone, false},
	TTComment:    {TTComment, "TTComment", "", -1, OpNone, false},
	TTKeyword:    {TTKeyword, "TTKeyword", "", -1, OpNone, false},
	TTIdentifier: {TTIdentifier, "TTIdentifier", "", -1, OpNone, false},
	TTNumber:     {TTNumber, "TTNumber", "", -1, OpNone, false},
	TTString:     {TTString, "TTString", "", -1, OpNone, false},
	TTChar:       {TTChar, "TTChar", "", -1, OpNone, false},
	TTConst:      {TTConst, "TTConst", "", -1, OpNone, false},
	TTReturnSym:  {TTReturnSym, "TTReturnSym", "->", -1, OpNone, false},
	TTExtendSym:  {TTExtendSym, "TTExtendSym", "<-", -1, OpNone, false},
	TTParenL:     {TTParenL, "TTParenL", "(", -1, OpNone, true},
	TTParenR:     {TTParenR, "TTParenR", ")", -1, OpNone, false},
	TTBracketL:   {TTBracketL, "TTBracketL", "[", -1, OpNone, true},
	TTBracketR:   {TTBracketR, "TTBracketR", "]", -1, OpNone, false},
	TTBraceL:     {TTBraceL, "TTBraceL", "{", -1, OpNone, true},
	TTBraceR:     {TTBraceR, "TTBraceR", "}", -1, OpNone, false},
	TTRest:       {TTRest, "TTRest", "...", -1, OpNone, false},
	TTSemi:       {TTSemi, "TTSemi", ";", -1, OpNone, true},
	TTColon:      {TTColon, "TTColon", ":", -1, OpNone, true},
	TTComma:      {TTComma, "TTComma", ",", -1, OpNone, true},
	TTDot:        {TTDot, "TTDot", ".", -1, OpNone, true},

	// binary operator (precedence 第二位为 1)
	TTAssign:   {TTAssign, "TTAssign", "=", 2, OpBinaryAssign, true},
	TTLogicOr:  {TTLogicOr, "TTLogicOr", "||", 3, OpBinary, true},
	TTLogicAnd: {TTLogicAnd, "TTLogicAnd", "&&", 4, OpBinary, true},
	TTBitOr:    {TTBitOr, "TTBitOr", "|", 5, OpBinary, true},
	TTBitXor:   {TTBitXor, "TTBitXor", "^", 6, OpBinary, true},
	TTBitAnd:   {TTBitAnd, "TTBitAnd", "&", 7, OpBinary, true},
	TTEq:       {TTEq, "TTEq", "==", 8, OpBinary, true},
	TTNe:       {TTNe, "TTNe", "!=", 8, OpBinary, true},
	TTLt:       {TTLt, "TTLt", "<", 9, OpBinary, true},
	TTLe:       {TTLe, "TTLe", "<=", 9, OpBinary, true},
	TTGt:       {TTGt, "TTGt", ">", 9, OpBinary, true},
	TTGe:       {TTGe, "TTGe", ">=", 9, OpBinary, true},
	TTIsOp:     {TTIsOp, "TTIsOp", "is", 9, OpBinaryType, false},
	TTPlus:     {TTPlus, "TTPlus", "+", 11, OpBinary, true},
	TTSub:      {TTSub, "TTSub", "-", 11, OpBinary, true},
	TTMul:      {TTMul, "TTMul", "*", 12, OpBinary, true},
	TTDiv:      {TTDiv, "TTDiv", "/", 12, OpBinary, true},
	TTRem:      {TTRem, "TTRem", "%", 12, OpBinary, true},

	// unary operator (precedence 第二位为 1)
	TTUnaryPlus: {TTUnaryPlus, "TTUnaryPlus", "+", 14, OpUnaryPrefix, true},
	TTUnarySub:  {TTUnarySub, "TTUnarySub", "-", 14, OpUnaryPrefix, true},
	TTLogicNot:  {TTLogicNot, "TTLogicNot", "!", 14, OpUnaryPrefix, true},
	TTBitNot:    {TTBitNot, "TTBitNot", "~", 14, OpUnaryPrefix, true},
}

type Token struct {
	*TokenMeta
	Value string
	Flag  string
	ast.Position
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
