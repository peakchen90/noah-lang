package ast

type Kind struct {
	Type T
}

type T interface{ isKind() }

func (*TypeNumber) isKind()    {}
func (*TypeString) isKind()    {}
func (*TypeBool) isKind()      {}
func (*TypeChar) isKind()      {}
func (*TypeArray) isKind()     {}
func (*TypeStruct) isKind()    {}
func (*TypeEnum) isKind()      {}
func (*TypeInterface) isKind() {}
func (*TypeAlias) isKind()     {}

type (
	TypeNumber struct{}

	TypeString struct{}

	TypeBool struct{}

	TypeChar struct{}

	TypeArray struct {
		T      Kind
		Len    int
		Cap    int
		Vector bool
	}

	TypeStruct struct {
		Properties []KindProperty
		Extends    Kind
		Interface  Kind
	}

	TypeEnum struct {
		Items []string
	}

	TypeInterface struct {
		Properties []KindProperty
		Extends    Kind
	}

	TypeAlias struct {
		T Kind
	}
)

// helper kind

type KindProperty struct {
	Name string
	T    Kind
}

type KindFunction struct {
	Name      string
	Arguments []KindProperty
	T         Kind
}
