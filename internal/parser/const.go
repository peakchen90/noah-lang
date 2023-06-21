package parser

// ReservedTypes 保留类型
var reservedTypes = [...]string{
	"num",
	"byte",
	"char",
	"str",
	"bool",
	"any",
}

// IsReservedType 判断是否为保留类型
func IsReservedType(value string) bool {
	for _, item := range reservedTypes {
		if item == value {
			return true
		}
	}
	return false
}

type ChainType = uint8

const (
	ChainTypeDot      ChainType = 0b00001000
	ChainTypeComputed ChainType = 0b00000100
	ChainTypeCall     ChainType = 0b00000010
	ChainTypeStruct   ChainType = 0b00000001
)
