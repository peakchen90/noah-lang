package vm

// TODO

type ValueRef struct {
	Current Value
}

type Value interface{ isValue() }

func (*NumberValue) isValue()  {}
func (*ByteValue) isValue()    {}
func (*Uint32Value) isValue()  {}
func (*StringValue) isValue()  {}
func (*BoolValue) isValue()    {}
func (*ArrayValue) isValue()   {}
func (*StructValue) isValue()  {}
func (*PointerValue) isValue() {}

type NumberValue struct {
	Value float64
}

type ByteValue struct {
	Value uint8
}

type Uint32Value struct {
	Value uint32
}

type StringValue struct {
	Value string
}

type BoolValue struct {
	Value bool
}

type ArrayValue struct {
	Value []Value
	Len   int
}

type StructValue struct {
	Value map[string]Value
}

type PointerValue struct {
	Value Value
}
