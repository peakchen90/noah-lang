package compiler

import (
	"github.com/peakchen90/noah-lang/internal/helper"
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
