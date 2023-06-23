package compiler

import (
	"github.com/peakchen90/noah-lang/internal/helper"
	"strings"
)

/* 作用域栈 */

type ScopeStack struct {
	stack []*Scope
}

func newScopeStack() *ScopeStack {
	return &ScopeStack{
		stack: make([]*Scope, 0, helper.DefaultCap),
	}
}

func (s *ScopeStack) size() int {
	return len(s.stack)
}

func (s *ScopeStack) last() *Scope {
	size := s.size()
	if size > 0 {
		return s.stack[size-1]
	}
	return nil
}

func (s *ScopeStack) push() {
	s.stack = append(s.stack, &Scope{
		value: make(map[string]Value),
		kind:  make(map[string]Kind),
	})
}

func (s *ScopeStack) pop() {
	size := s.size()
	if size > 0 {
		s.stack = s.stack[:size-1]
	}
}

func (s *ScopeStack) putValue(name string, scope Value) bool {
	last := s.last()
	if last != nil {
		if last.has(name) {
			return false
		}

		last.setValue(name, scope)
		return true
	}
	return false
}

func (s *ScopeStack) putKind(name string, scope Kind) bool {
	last := s.last()
	if last != nil {
		if last.has(name) {
			return false
		}

		last.setKind(name, scope)
		return true
	}
	return false
}

func (s *ScopeStack) isExist(name string) bool {
	last := s.last()
	return last != nil && last.has(name)
}

func (s *ScopeStack) findValue(name string) Value {
	for i := s.size() - 1; i >= 0; i-- {
		scope := s.stack[i].getValue(name)
		if scope != nil {
			return scope
		}
	}
	return nil
}

func (s *ScopeStack) findKind(name string) Kind {
	for i := s.size() - 1; i >= 0; i-- {
		scope := s.stack[i].getKind(name)
		if scope != nil {
			return scope
		}
	}
	return nil
}

func (s *ScopeStack) findMemberKind(members []string) Kind {
	var outside *Module
	size := len(members)

	for i, id := range members {
		var maybeModule Value
		if outside != nil {
			maybeModule = outside.ScopeStack.findValue(id)
		} else {
			maybeModule = s.findValue(id)
		}

		if maybeModule != nil {
			module, ok := maybeModule.(*ModuleValue)
			if ok {
				outside = module.Module
			} else {
				// TODO panic
				panic("unexpected")
			}
		} else if i == size-1 {
			if outside != nil {
				return outside.PublicScope.getKind(id)
			} else {
				return s.findKind(id)
			}
		} else {
			// TODO not found
			panic("can not found: " + strings.Join(members, "."))
		}
	}

	return nil
}

func (s *ScopeStack) findFunc(name string) *FuncValue {
	scope := s.findValue(name)
	fnScope, ok := scope.(*FuncValue)
	if scope != nil && ok {
		return fnScope
	}
	return nil
}

func (s *ScopeStack) findVar(name string) *VarValue {
	scope := s.findValue(name)
	varScope, ok := scope.(*VarValue)
	if scope != nil && ok {
		return varScope
	}
	return nil
}
