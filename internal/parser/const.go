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

type ChainOp uint8

const (
	_copBase          ChainOp = 0b00000001
	ChainTypeDot              = _copBase << 0
	ChainTypeComputed         = _copBase << 1
	ChainTypeCall             = _copBase << 2
	ChainTypeStruct           = _copBase << 3
)
