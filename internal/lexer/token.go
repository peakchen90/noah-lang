package lexer

import (
	"fmt"
	"github.com/peakchen90/noah-lang/internal/ast"
)

type TokenType = uint8

type OpType = uint8

const (
	OpNone = 0b00000000

	OpBinary     = 0b00001111
	OpBinaryLTR  = 0b00000001 // left to right
	OpBinaryRTL  = 0b00000010 // right to left
	OpBinaryType = 0b00000100 // typeof value

	OpUnary        = 0b11110000
	OpUnaryPrefix  = 0b00010000
	OpUnaryPostfix = 0b00100000
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

	TTAssign              // =
	TTPlusAssign          // +=
	TTSubAssign           // -=
	TTMulAssign           // *=
	TTDivAssign           // /=
	TTRemAssign           // %=
	TTBitLeftShiftAssign  // <<=
	TTBitRightShiftAssign // >>=
	TTBitAndAssign        // &=
	TTBitOrAssign         // |=
	TTBitXorAssign        // ^=

	TTLogicOr       // ||
	TTLogicAnd      // &&
	TTBitOr         // |
	TTBitXor        // ^
	TTBitAnd        // &
	TTEq            // ==
	TTNe            // !=
	TTLt            // <
	TTLe            // <=
	TTGt            // >
	TTGe            // >=
	TTIsOp          // `is`
	TTBitLeftShift  // <<
	TTBitRightShift // >>
	TTPlus          // +
	TTSub           // -
	TTMul           // *
	TTDiv           // /
	TTRem           // %

	TTPrefixPlus // + ...
	TTPrefixSub  // - ...
	TTLogicNot   // ! ...
	TTBitNot     // ~ ...
	TTPrefixInc  // ++ ...
	TTPrefixDec  // -- ...
	TTPostfixInc // ... ++
	TTPostfixDec // ... --
)

// precedence see: https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Operators/Operator_Precedence
var tokenMetaTable = [59]TokenMeta{
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
	TTDot:        {TTDot, "TTDot", ".", -1, OpNone, false},

	// binary operator
	TTAssign:              {TTAssign, "TTAssign", "=", 2, OpBinaryRTL, true},
	TTPlusAssign:          {TTPlusAssign, "TTPlusAssign", "+=", 2, OpBinaryRTL, true},
	TTSubAssign:           {TTSubAssign, "TTSubAssign", "-=", 2, OpBinaryRTL, true},
	TTMulAssign:           {TTMulAssign, "TTMulAssign", "*=", 2, OpBinaryRTL, true},
	TTDivAssign:           {TTDivAssign, "TTDivAssign", "/=", 2, OpBinaryRTL, true},
	TTRemAssign:           {TTRemAssign, "TTRemAssign", "%=", 2, OpBinaryRTL, true},
	TTBitLeftShiftAssign:  {TTBitLeftShiftAssign, "TTBitLeftShiftAssign", "<<=", 2, OpBinaryRTL, true},
	TTBitRightShiftAssign: {TTBitRightShiftAssign, "TTBitRightShiftAssign", ">>=", 2, OpBinaryRTL, true},
	TTBitAndAssign:        {TTBitAndAssign, "TTBitAndAssign", "&=", 2, OpBinaryRTL, true},
	TTBitOrAssign:         {TTBitOrAssign, "TTBitOrAssign", "|=", 2, OpBinaryRTL, true},
	TTBitXorAssign:        {TTBitXorAssign, "TTBitXorAssign", "^=", 2, OpBinaryRTL, true},

	TTLogicOr:       {TTLogicOr, "TTLogicOr", "||", 3, OpBinaryLTR, true},
	TTLogicAnd:      {TTLogicAnd, "TTLogicAnd", "&&", 4, OpBinaryLTR, true},
	TTBitOr:         {TTBitOr, "TTBitOr", "|", 5, OpBinaryLTR, true},
	TTBitXor:        {TTBitXor, "TTBitXor", "^", 6, OpBinaryLTR, true},
	TTBitAnd:        {TTBitAnd, "TTBitAnd", "&", 7, OpBinaryLTR, true},
	TTEq:            {TTEq, "TTEq", "==", 8, OpBinaryLTR, true},
	TTNe:            {TTNe, "TTNe", "!=", 8, OpBinaryLTR, true},
	TTLt:            {TTLt, "TTLt", "<", 9, OpBinaryLTR, true},
	TTLe:            {TTLe, "TTLe", "<=", 9, OpBinaryLTR, true},
	TTGt:            {TTGt, "TTGt", ">", 9, OpBinaryLTR, true},
	TTGe:            {TTGe, "TTGe", ">=", 9, OpBinaryLTR, true},
	TTIsOp:          {TTIsOp, "TTIsOp", "is", 9, OpBinaryLTR | OpBinaryType, false},
	TTBitLeftShift:  {TTBitLeftShift, "TTBitLeftShift", "<<", 10, OpBinaryLTR, true},
	TTBitRightShift: {TTBitRightShift, "TTBitRightShift", ">>", 10, OpBinaryLTR, true},
	TTPlus:          {TTPlus, "TTPlus", "+", 11, OpBinaryLTR, true},
	TTSub:           {TTSub, "TTSub", "-", 11, OpBinaryLTR, true},
	TTMul:           {TTMul, "TTMul", "*", 12, OpBinaryLTR, true},
	TTDiv:           {TTDiv, "TTDiv", "/", 12, OpBinaryLTR, true},
	TTRem:           {TTRem, "TTRem", "%", 12, OpBinaryLTR, true},

	// unary operator
	TTPrefixPlus: {TTPrefixPlus, "TTPrefixPlus", "+", 14, OpUnaryPrefix, true},
	TTPrefixSub:  {TTPrefixSub, "TTPrefixSub", "-", 14, OpUnaryPrefix, true},
	TTLogicNot:   {TTLogicNot, "TTLogicNot", "!", 14, OpUnaryPrefix, true},
	TTBitNot:     {TTBitNot, "TTBitNot", "~", 14, OpUnaryPrefix, true},
	TTPrefixInc:  {TTPrefixInc, "TTPrefixInc", "++", 14, OpUnaryPrefix, true},
	TTPrefixDec:  {TTPrefixDec, "TTPrefixDec", "--", 14, OpUnaryPrefix, true},
	TTPostfixInc: {TTPostfixInc, "TTPostfixInc", "++", 15, OpUnaryPostfix, false},
	TTPostfixDec: {TTPostfixDec, "TTPostfixDec", "--", 15, OpUnaryPostfix, false},
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
