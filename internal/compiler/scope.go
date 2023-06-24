package compiler

type Scope struct {
	module map[string]*Module
	value  map[string]Value
	kind   map[string]*KindRef
}

func newScope() *Scope {
	return &Scope{
		module: make(map[string]*Module),
		value:  make(map[string]Value),
		kind:   make(map[string]*KindRef),
	}
}

func (s *Scope) getModule(name string) *Module {
	return s.module[name]
}

func (s *Scope) getValue(name string) Value {
	return s.value[name]
}

func (s *Scope) getKind(name string) *KindRef {
	return s.kind[name]
}

func (s *Scope) setModule(name string, module *Module) {
	s.module[name] = module
}

func (s *Scope) setValue(name string, value Value) {
	s.value[name] = value
}

func (s *Scope) setKind(name string, kind *KindRef) {
	s.kind[name] = kind
}

func (s *Scope) has(name string) bool {
	return s.hasModule(name) || s.hasValue(name) || s.hasKind(name)
}

func (s *Scope) hasModule(name string) bool {
	_, has := s.module[name]
	return has
}

func (s *Scope) hasValue(name string) bool {
	_, has := s.value[name]
	return has
}

func (s *Scope) hasKind(name string) bool {
	_, has := s.kind[name]
	return has
}
