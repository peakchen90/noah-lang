package lexer

type OpType uint8

const (
	OpNone OpType = 0b00000000

	_opbBase     OpType = 0b00000001
	OpBinaryLTR         = _opbBase << 0 // left to right
	OpBinaryRTL         = _opbBase << 1 // right to left
	OpBinaryType        = _opbBase << 2 // typeof value
	OpBinary            = OpBinaryLTR | OpBinaryRTL | OpBinaryType

	_opuBase       OpType = 0b00010000
	OpUnaryPrefix         = _opuBase << 0
	OpUnaryPostfix        = _opuBase << 1
	OpUnary               = OpUnaryPrefix | OpUnaryPostfix
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
