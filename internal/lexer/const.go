package lexer

// Keywords 关键字
var Keywords = [...]string{
	// 变量声明
	"fn", "let", "const",
	// 类型声明
	"type", "interface", "struct", "enum",
	// 逻辑控制
	"if", "else", "for", "return", "break", "continue",
	// 其他修饰符
	"pub", "extends", "import", "as", "new",
}

// Constants 内置常量
var Constants = [...]string{
	"true", "false", "null", "self",
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

// IsConstant 判断是否为内置常量
func IsConstant(value string) bool {
	for _, item := range Constants {
		if item == value {
			return true
		}
	}
	return false
}
