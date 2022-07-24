package ast

type Kind uint8

const (
	KNumber    Kind = iota + 1 // `num`
	KString                    // `str`
	KBool                      // `bool`
	KArray                     // `[]str`、`[5]num`
	KStruct                    // `type T {a: str, b: num}`
	KEnum                      // `type T {A, B}`
	KInterface                 // `interface I {}`
	KCustom                    // `type foo = num`
)

var Keyword = [...]string{
	// 变量、类型声明
	"fn", "let", "const", "type", "interface",
	// 逻辑控制
	"if", "else", "for", "return", "break", "continue",
	// 内置常量
	"true", "false", "self", "null",
	// 其他修饰符
	"pub", "extends", "import", "as",
}
