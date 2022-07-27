package parser

// ReservedTypes 保留类型
var ReservedTypes = [...]string{
	"num",
	"byte",
	"char",
	"str",
	"bool",
	"any",
}

// IsReservedType 判断是否为保留类型
func IsReservedType(value string) bool {
	for _, item := range ReservedTypes {
		if item == value {
			return true
		}
	}
	return false
}
