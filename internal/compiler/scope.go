package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
)

type Scope struct {
	value map[string]Value
	kind  map[string]Kind
}

func (s *Scope) getValue(name string) Value {
	return s.value[name]
}

func (s *Scope) getKind(name string) Kind {
	return s.kind[name]
}

func (s *Scope) setValue(name string, scope Value) {
	s.value[name] = scope
}

func (s *Scope) setKind(name string, scope Kind) {
	s.kind[name] = scope
}

func (s *Scope) has(name string) bool {
	return s.hasValue(name) || s.hasKind(name)
}

func (s *Scope) hasValue(name string) bool {
	return s.getValue(name) != nil
}

func (s *Scope) hasKind(name string) bool {
	return s.getKind(name) != nil
}

/* compile module helpers  */

func (m *Module) validateValueScope(id *ast.Identifier) {
	if m.ScopeStack.isExist(id.Name) {
		// TODO panic
		panic("duplicate" + id.Name)
	}
}

func (m *Module) validateKindScope(kindId *ast.KindIdentifier) {
	if m.ScopeStack.isExist(kindId.Name) {
		// TODO panic
		panic("duplicate" + kindId.Name)
	}
}
