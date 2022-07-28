package lexer

// Keywords 关键字
var Keywords = [...]string{
	// 变量、类型声明
	"fn", "let", "const", "type", "interface",
	// 逻辑控制
	"if", "else", "for", "return", "break", "continue",
	// 内置常量
	"true", "false", "null", "self",
	// 其他修饰符
	"pub", "extends", "import", "as",
}

// IsKeyword 判断是否为关键字
func IsKeyword(value string) bool {
	for _, item := range Keywords {
		if item == value {
			return true
		}
	}

	return false
}
