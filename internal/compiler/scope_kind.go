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
	_, has := i.methods[name]
	if has {
		// TODO
		panic("duplicate method " + name)
	}
	i.methods[name] = value
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

/* kind */

type Kind interface{ getImpl() *Impl }

func (t *TNumber) getImpl() *Impl    { return t.Impl }
func (t *TByte) getImpl() *Impl      { return t.Impl }
func (t *TChar) getImpl() *Impl      { return t.Impl }
func (t *TString) getImpl() *Impl    { return t.Impl }
func (t *TBool) getImpl() *Impl      { return t.Impl }
func (t *TAny) getImpl() *Impl       { return newImpl() }
func (t *TArray) getImpl() *Impl     { return t.Impl }
func (t *TFunc) getImpl() *Impl      { return t.Impl }
func (t *TInterface) getImpl() *Impl { return newImpl() }
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
		Impl *Impl
	}

	TArray struct {
		Kind Kind
		Len  int
		Impl *Impl
	}

	TFunc struct {
		Id           int
		Arguments    []Kind
		Return       Kind
		RestArgument bool
		Impl         *Impl
	}

	TStruct struct {
		Id         int
		Extends    []Kind
		Properties map[string]Kind
		Impl       *Impl
	}

	TInterface struct {
		Id         int
		Properties map[string]Kind
		Refers     []Kind
	}

	TEnum struct {
		Id      int
		Choices map[string]int
		Impl    *Impl
	}

	TCustom struct {
		Id   int
		Kind Kind
		Impl *Impl
	}
)

/* helpers */

var _typeId = 1 << 8

func getNextTypeId() int {
	_typeId++
	return _typeId
}

func compareKind(expected Kind, received Kind, isMatch bool) bool {
	_, ok := received.(*TAny)
	if ok {
		return true
	}

	switch expected.(type) {
	case *TNumber:
		_, ok = received.(*TNumber)
		return ok
	case *TByte:
		_, ok = received.(*TByte)
		return ok
	case *TChar:
		_, ok = received.(*TChar)
		return ok
	case *TString:
		_, ok = received.(*TString)
		return ok
	case *TBool:
		_, ok = received.(*TBool)
		return ok
	case *TAny:
		return true
	case *TArray:
		r, ok := received.(*TArray)
		if !ok {
			return false
		}

		e := expected.(*TArray)
		return e.Len == r.Len && compareKind(e.Kind, r.Kind, isMatch)
	case *TFunc:
		r, ok := received.(*TFunc)
		if !ok {
			return false
		}
		e := expected.(*TFunc)

		if isMatch {
			if e.RestArgument != r.RestArgument || len(e.Arguments) != len(r.Arguments) {
				return false
			}

			for i, arg := range e.Arguments {
				if !compareKind(arg, r.Arguments[i], isMatch) {
					return false
				}
			}
			return compareKind(e.Return, r.Return, isMatch)
		} else {
			return r.Id == e.Id
		}
	case *TStruct:
		r, ok := received.(*TStruct)
		if !ok {
			return false
		}
		e := expected.(*TStruct)

		if isMatch {
			// TODO think about extends

			if len(e.Properties) != len(r.Properties) {
				return false
			}
			for k, v := range e.Properties {
				if !compareKind(v, r.Properties[k], isMatch) {
					return false
				}
			}
			return true
		} else {
			return r.Id == e.Id
		}
	case *TInterface:
		r, ok := received.(*TInterface)
		if !ok {
			return false
		}
		e := expected.(*TInterface)

		if isMatch {
			if len(e.Properties) != len(r.Properties) {
				return false
			}
			for k, v := range e.Properties {
				if !compareKind(v, r.Properties[k], isMatch) {
					return false
				}
			}
			return true
		} else {
			return r.Id == e.Id
		}
	case *TEnum:
		r, ok := received.(*TEnum)
		if !ok {
			return false
		}
		e := expected.(*TEnum)

		if isMatch {
			if len(e.Choices) != len(r.Choices) {
				return false
			}
			for i, v := range e.Choices {
				if v != r.Choices[i] {
					return false
				}
			}
			return true
		} else {
			return r.Id == e.Id
		}
	case *TCustom:
		e := expected.(*TCustom)
		if isMatch {
			return compareKind(e.Kind, received, true)
		} else {
			r, ok := received.(*TCustom)
			return ok && r.Id == e.Id
		}
	}

	return false
}
