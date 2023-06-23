package lexer

// 关键字
var keywords = [...]string{
	// 变量声明
	"fn", "let", "const",
	// 类型声明
	"type", "interface", "struct", "enum",
	// 逻辑控制
	"if", "else", "for", "return", "break", "continue",
	// 运算符
	"as", "is",
	// 其他修饰符
	"pub", "import", "impl",
}

var reservedKeywords = [...]string{
	"async", "await", "switch", "case", "default", "try", "catch", "throw", "new",
}

// 内置常量
var builtInConstants = [...]string{
	"true", "false", "null", "self",
}

// isKeyword 判断是否为关键字
func isKeyword(value string) bool {
	for _, item := range keywords {
		if item == value {
			return true
		}
	}
	return false
}

// isReversedKeyword 判断是否为保留关键字
func isReversedKeyword(value string) bool {
	for _, item := range reservedKeywords {
		if item == value {
			return true
		}
	}
	return false
}

// isBuiltInConstant 判断是否为内置常量
func isBuiltInConstant(value string) bool {
	for _, item := range builtInConstants {
		if item == value {
			return true
		}
	}
	return false
}
