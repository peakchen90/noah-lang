package lexer

import (
	"fmt"
	"github.com/peakchen90/noah-lang/internal/ast"
)

type TokenType = uint8

type TokenMeta struct {
	Type       TokenType
	Name       string
	Text       string
	Precedence int8
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

	// operator

	TTAssign    // =
	TTPlus      // +
	TTSub       // -
	TTMul       // *
	TTDiv       // /
	TTRem       // %
	TTLt        // <
	TTLe        // <=
	TTGt        // >
	TTGe        // >=
	TTEq        // ==
	TTNe        // !=
	TTLogicAnd  // &&
	TTLogicOr   // ||
	TTLogicNot  // !
	TTBitAnd    // &
	TTBitOr     // |
	TTBitNot    // ~
	TTBitXor    // ^
	TTDot       // .
	TTUnaryPlus // +
	TTUnarySub  // -
)

var tokenMetaTable = [...]TokenMeta{
	TTEof:        {TTEof, "TTEof", "", -1, false},
	TTComment:    {TTComment, "TTComment", "", -1, false},
	TTKeyword:    {TTKeyword, "TTKeyword", "", -1, false},
	TTIdentifier: {TTIdentifier, "TTIdentifier", "", -1, false},
	TTNumber:     {TTNumber, "TTNumber", "", -1, false},
	TTString:     {TTString, "TTString", "", -1, false},
	TTChar:       {TTChar, "TTChar", "", -1, false},
	TTConst:      {TTConst, "TTConst", "", -1, false},
	TTReturnSym:  {TTReturnSym, "TTReturnSym", "->", -1, false},
	TTExtendSym:  {TTExtendSym, "TTExtendSym", "<-", -1, false},
	TTParenL:     {TTParenL, "TTParenL", "(", -1, true},
	TTParenR:     {TTParenR, "TTParenR", ")", -1, false},
	TTBracketL:   {TTBracketL, "TTBracketL", "[", -1, true},
	TTBracketR:   {TTBracketR, "TTBracketR", "]", -1, false},
	TTBraceL:     {TTBraceL, "TTBraceL", "{", -1, true},
	TTBraceR:     {TTBraceR, "TTBraceR", "}", -1, false},
	TTRest:       {TTRest, "TTRest", "...", -1, false},
	TTSemi:       {TTSemi, "TTSemi", ";", -1, true},
	TTColon:      {TTColon, "TTColon", ":", -1, true},
	TTComma:      {TTComma, "TTComma", ",", -1, true},
	TTDot:        {TTDot, "TTDot", ".", -1, true},

	// operator
	// precedence see: https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Operators/Operator_Precedence

	TTAssign:   {TTAssign, "TTAssign", "=", 2, true},
	TTPlus:     {TTPlus, "TTPlus", "+", 12, true},
	TTSub:      {TTSub, "TTSub", "-", 12, true},
	TTMul:      {TTMul, "TTMul", "*", 13, true},
	TTDiv:      {TTDiv, "TTDiv", "/", 13, true},
	TTRem:      {TTRem, "TTRem", "%", 13, true},
	TTLt:       {TTLt, "TTLt", "<", 10, true},
	TTLe:       {TTLe, "TTLe", "<=", 10, true},
	TTGt:       {TTGt, "TTGt", ">", 10, true},
	TTGe:       {TTGe, "TTGe", ">=", 10, true},
	TTEq:       {TTEq, "TTEq", "==", 9, true},
	TTNe:       {TTNe, "TTNe", "!=", 9, true},
	TTLogicAnd: {TTLogicAnd, "TTLogicAnd", "&&", 5, true},
	TTLogicOr:  {TTLogicOr, "TTLogicOr", "||", 4, true},
	TTLogicNot: {TTLogicNot, "TTLogicNot", "!", 15, true},
	TTBitAnd:   {TTBitAnd, "TTBitAnd", "&", 8, true},
	TTBitOr:    {TTBitOr, "TTBitOr", "|", 6, true},
	TTBitNot:   {TTBitNot, "TTBitNot", "~", 15, true},
	TTBitXor:   {TTBitXor, "TTBitXor", "^", 7, true},

	TTUnaryPlus: {TTUnaryPlus, "TTUnaryPlus", "+ ", 14, true},
	TTUnarySub:  {TTUnarySub, "TTUnarySub", "- ", 14, true},
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
