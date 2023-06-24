package parser

// 保留类型
var reservedTypes = [...]string{
	"number",
	"byte",
	"char",
	"string",
	"bool",
	"any",
}

// 判断是否为保留类型
func isReservedType(value string) bool {
	for _, item := range reservedTypes {
		if item == value {
			return true
		}
	}
	return false
}

type AccessType uint8

const (
	_aoBase        AccessType = 0b00000001
	AccessDot                 = _aoBase << 0
	AccessComputed            = _aoBase << 1
	AccessCall                = _aoBase << 2
	AccessStruct              = _aoBase << 3
)
