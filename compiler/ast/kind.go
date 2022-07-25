package ast

type Kind uint8

const (
	KVoid      Kind = iota + 1 // `void`
	KNumber                    // `num`
	KString                    // `str`
	KBool                      // `bool`
	KChar                      // `char`
	KArray                     // `[]str`、`[5]num`
	KStruct                    // `type T {a: str, b: num}`
	KEnum                      // `type T {A, B}`
	KInterface                 // `interface I {}`
	KCustom                    // `type foo = num`
)

type KindMeta struct {
	Kind
	Data K
}

type K interface{ isKind() }

func (*KindString) isKind()    {}
func (*KindArray) isKind()     {}
func (*KindStruct) isKind()    {}
func (*KindEnum) isKind()      {}
func (*KindInterface) isKind() {}
func (*KindCustom) isKind()    {}
func (*KindFunction) isKind()  {}

type KindString struct {
	Len int
	Cap int
}

type KindArray struct {
	T      KindMeta
	Len    int
	Cap    int
	Vector bool
}

type KindStruct struct {
	Properties []KindProperty
}

type KindEnum struct {
	Items []string
}

type KindInterface struct {
	Properties []KindProperty
}

type KindCustom struct {
	T KindMeta
}

// helper kind

type KindProperty struct {
	Name string
	T    KindMeta
}

type KindFunction struct {
	Name      string
	Arguments []KindProperty
	T         KindMeta
}
