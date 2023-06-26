package compiler

import (
	"errors"
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"strings"
)

/* 作用域栈 */

type ScopeStack struct {
	module *Module
	stack  []*Scope
}

func newScopeStack(module *Module) *ScopeStack {
	return &ScopeStack{
		module: module,
		stack:  make([]*Scope, 0, helper.DefaultCap),
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
	s.stack = append(s.stack, newScope())
}

func (s *ScopeStack) pop() {
	size := s.size()
	if size > 0 {
		s.stack = s.stack[:size-1]
	}
}

func (s *ScopeStack) putModule(name *ast.Identifier, module *Module, isPanic bool) {
	last := s.last()
	if last != nil {
		if last.has(name.Name) {
			if isPanic {
				s.module.unexpectedPos(name.Start, "identifier has already been declared: "+name.Name)
			}
		}

		last.setModule(name.Name, module)
	}
}

func (s *ScopeStack) putValue(name *ast.Identifier, value Value, isPanic bool) {
	last := s.last()
	if last != nil {
		if last.has(name.Name) {
			if isPanic {
				s.module.unexpectedPos(name.Start, "identifier has already been declared: "+name.Name)
			}
		}

		last.setValue(name.Name, value)
	}
}

func (s *ScopeStack) putSelfValue(value Value) {
	last := s.last()
	if last != nil {
		last.setValue("self", value)
	}
}

func (s *ScopeStack) putKind(name *ast.Identifier, kind *KindRef, isPanic bool) {
	last := s.last()
	if last != nil {
		if last.has(name.Name) {
			if isPanic {
				s.module.unexpectedPos(name.Start, "identifier has already been declared: "+name.Name)
			}
		}

		last.setKind(name.Name, kind)
	}
}

func (s *ScopeStack) putSelfKind(kind *KindRef) {
	last := s.last()
	if last != nil {
		last.setKind("self", kind)
	}
}

func (s *ScopeStack) findModule(name *ast.Identifier, isPanic bool) *Module {
	for i := s.size() - 1; i >= 0; i-- {
		module := s.stack[i].getModule(name.Name)
		if module != nil {
			return module
		}
	}

	if isPanic {
		s.module.unexpectedPos(name.Start, name.Name+" is not defined")
	}

	return nil
}

func (s *ScopeStack) findValue(name *ast.Identifier, isPanic bool) Value {
	for i := s.size() - 1; i >= 0; i-- {
		value := s.stack[i].getValue(name.Name)
		if value != nil {
			return value
		}
	}

	if isPanic {
		s.module.unexpectedPos(name.Start, name.Name+" is not defined")
	}

	return nil
}

func (s *ScopeStack) findFuncValue(name *ast.Identifier, isPanic bool) *FuncValue {
	scope := s.findValue(name, isPanic)
	value, ok := scope.(*FuncValue)
	if scope != nil && ok {
		return value
	}

	if isPanic {
		s.module.unexpectedPos(name.Start, name.Name+" is not found")
	}

	return nil
}

func (s *ScopeStack) findVarValue(name *ast.Identifier, isPanic bool) *VarValue {
	scope := s.findValue(name, isPanic)
	value, ok := scope.(*VarValue)
	if scope != nil && ok {
		return value
	}

	if isPanic {
		s.module.unexpectedPos(name.Start, name.Name+" is not found")
	}

	return nil
}

func (s *ScopeStack) findSelfValue(name *ast.Identifier, isPanic bool) *SelfValue {
	scope := s.findValue(name, isPanic)
	value, ok := scope.(*SelfValue)
	if scope != nil && ok {
		return value
	}

	if isPanic {
		s.module.unexpectedPos(name.Start, name.Name+" is not found")
	}

	return nil
}

// TODO
func (s *ScopeStack) findMemberValue(expr *ast.Expr, isPanic bool) Value {
	//	memberIdStack := make([]*ast.Expr, 0, helper.DefaultCap)
	//	current := expr
	//
	//outer:
	//	for {
	//		switch current.Node.(type) {
	//		case *ast.MemberExpr:
	//			node := current.Node.(*ast.MemberExpr)
	//			if node.Computed {
	//				// TODO
	//				s.module.compileExpr(node.Property)
	//			} else {
	//				memberIdStack = append(memberIdStack, node.Property)
	//			}
	//			current = node.Object
	//		case *ast.IdentifierLiteral:
	//			memberIdStack = append(memberIdStack, current)
	//			break outer
	//		default:
	//			panic("Internal Err")
	//		}
	//	}
	//
	//	var kind *KindRef
	//	module := s.module
	//
	//	handleValue := func(value Value) {
	//		switch value.(type) {
	//		case *FuncValue:
	//		case *VarValue:
	//		case *SelfValue:
	//		default:
	//			panic("Internal Err")
	//		}
	//	}
	//
	//	handleKind := func(kind *KindRef) {
	//
	//	}
	//
	//outer:
	//	for i := len(memberIdStack) - 1; i >= 0; i-- {
	//		item := memberIdStack[i]
	//
	//		switch item.Node.(type) {
	//		case *ast.IdentifierLiteral:
	//			node := item.Node.(*ast.IdentifierLiteral)
	//			if m := module.scopes.findModule(node.Name, false); m != nil {
	//				module = m
	//				continue outer
	//			} else if v := module.scopes.findValue(node.Name, false); v != nil {
	//				handleValue(v)
	//			} else if k := module.scopes.findIdentifierKind(node.Name, false); k != nil {
	//				handleKind(k)
	//			}
	//
	//		}
	//		node := item.Node.(*ast.IdentifierLiteral)
	//
	//		if i > 0 {
	//			found := module.scopes.findModule(node.Name, false)
	//			if found == nil {
	//				break
	//			}
	//			module = found
	//		} else {
	//			kind = module.scopes.findIdentifierKind(node.Name, false)
	//		}
	//	}
	//
	//	if kind == nil {
	//		if isPanic {
	//			s.module.unexpectedPos(expr.Start, " is not found")
	//		}
	//	}

	return nil
}

func (s *ScopeStack) findKind(name string) (*KindRef, error) {
	for i := s.size() - 1; i >= 0; i-- {
		scope := s.stack[i].getKind(name)
		if scope != nil {
			return scope, nil
		}
	}
	return nil, errors.New(name + " is not found")
}

func (s *ScopeStack) findIdentifierKind(name *ast.Identifier, isPanic bool) *KindRef {
	kind, err := s.findKind(name.Name)

	if err != nil && isPanic {
		s.module.unexpectedPos(name.Start, name.Name+" is not found")
	}

	return kind
}

func (s *ScopeStack) findMemberKind(kindExpr *ast.KindExpr, isPanic bool) *KindRef {
	memberIdStack := make([]*ast.KindExpr, 0, helper.SmallCap)
	current := kindExpr

outer:
	for {
		switch current.Node.(type) {
		case *ast.TMemberKind:
			node := current.Node.(*ast.TMemberKind)
			memberIdStack = append(memberIdStack, node.Right)
			current = node.Left
		case *ast.TIdentifier:
			memberIdStack = append(memberIdStack, current)
			break outer
		default:
			panic("Internal Err")
		}
	}

	var kind *KindRef
	module := s.module
	builder := strings.Builder{}

	for i := len(memberIdStack) - 1; i >= 0; i-- {
		item := memberIdStack[i]
		node := item.Node.(*ast.TIdentifier)
		builder.WriteString(node.Name.Name)

		if i > 0 {
			found := module.scopes.findModule(node.Name, false)
			if found == nil {
				break
			}
			module = found
			builder.WriteByte('.')
		} else {
			kind = module.scopes.findIdentifierKind(node.Name, false)
		}
	}

	if kind == nil {
		if isPanic {
			s.module.unexpectedPos(kindExpr.Start, builder.String()+" is not found")
		}
	}

	return kind
}

func (s *ScopeStack) findSelfKind(kindExpr *ast.KindExpr, isPanic bool) *KindRef {
	_, ok := kindExpr.Node.(*ast.TSelf)
	if !ok {
		panic("Internal Err")
	}

	kind, err := s.findKind("self")

	if err != nil && isPanic {
		s.module.unexpectedPos(kindExpr.Start, "`self` is not allowed here")
	}

	return kind
}
