package compiler

/* impls */

type Impl struct {
	methods map[string]*FuncValue
}

func newImpl() *Impl {
	return &Impl{
		methods: make(map[string]*FuncValue),
	}
}

func (i *Impl) addFunc(value *FuncValue) {
	name := value.Name
	i.methods[name] = value
}

func (i *Impl) hasFunc(name string) bool {
	_, has := i.methods[name]
	return has
}

func (i *Impl) getFunc(name string) *FuncValue {
	return i.methods[name]
}

func (i *Impl) getPubFunc(name string) *FuncValue {
	value := i.getFunc(name)
	if value != nil && value.Name[0] != '_' {
		return value
	}
	return nil
}

/* kind ref */

type KindRef struct {
	Ref Kind
}

/* kind */

type Kind interface{ getImpl() *Impl }

func (t *TNumber) getImpl() *Impl    { return t.Impl }
func (t *TByte) getImpl() *Impl      { return t.Impl }
func (t *TChar) getImpl() *Impl      { return t.Impl }
func (t *TString) getImpl() *Impl    { return t.Impl }
func (t *TBool) getImpl() *Impl      { return t.Impl }
func (t *TAny) getImpl() *Impl       { return nil }
func (t *TSelf) getImpl() *Impl      { return nil }
func (t *TArray) getImpl() *Impl     { return t.Impl }
func (t *TFunc) getImpl() *Impl      { return t.Impl }
func (t *TInterface) getImpl() *Impl { return nil }
func (t *TStruct) getImpl() *Impl    { return t.Impl }
func (t *TEnum) getImpl() *Impl      { return t.Impl }
func (t *TCustom) getImpl() *Impl    { return t.Impl }

type (
	TNumber struct {
		Impl *Impl
	}

	TByte struct {
		Impl *Impl
	}

	TChar struct {
		Impl *Impl
	}

	TString struct {
		Impl *Impl
	}

	TBool struct {
		Impl *Impl
	}

	TAny struct {
	}

	TSelf struct {
		KindRef *KindRef
	}

	TArray struct {
		KindRef *KindRef
		Len     int // -1 means vector array
		Impl    *Impl
	}

	TFunc struct {
		Params    []*KindRef
		Return    *KindRef
		RestParam bool
		Impl      *Impl
	}

	TStruct struct {
		Extends    []*KindRef
		Properties map[string]*KindRef
		Impl       *Impl
	}

	TInterface struct {
		Properties map[string]*KindRef
		Refs       []*KindRef
	}

	TEnum struct {
		Choices map[string]int
		Impl    *Impl
	}

	TCustom struct {
		KindRef *KindRef
		Impl    *Impl
	}
)

/* 类型常量 */

var (
	typeNumber = &TNumber{Impl: newImpl()}
	typeByte   = &TByte{Impl: newImpl()}
	typeChar   = &TChar{Impl: newImpl()}
	typeString = &TString{Impl: newImpl()}
	typeBool   = &TBool{Impl: newImpl()}
	typeAny    = &TAny{}
)
