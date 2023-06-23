package lexer

type OpType uint8

const (
	OpNone OpType = 0b00000000

	OpBinary     OpType = 0b00001111
	OpBinaryLTR  OpType = 0b00000001 // left to right
	OpBinaryRTL  OpType = 0b00000010 // right to left
	OpBinaryType OpType = 0b00000100 // typeof value

	OpUnary        OpType = 0b11110000
	OpUnaryPrefix  OpType = 0b00010000
	OpUnaryPostfix OpType = 0b00100000
)

func (op OpType) IsOpNone() bool {
	return op == 0
}

func (op OpType) IsOpBinary() bool {
	return op&OpBinary > 0
}

func (op OpType) IsOpBinaryLTR() bool {
	return op&OpBinaryLTR > 0
}

func (op OpType) IsOpBinaryRTL() bool {
	return op&OpBinaryRTL > 0
}

func (op OpType) IsOpBinaryType() bool {
	return op&OpBinaryType > 0
}

func (op OpType) IsOpUnary() bool {
	return op&OpUnary > 0
}
func (op OpType) IsOpUnaryPrefix() bool {
	return op&OpUnaryPrefix > 0
}

func (op OpType) IsOpUnaryPostfix() bool {
	return op&OpUnaryPostfix > 0
}
