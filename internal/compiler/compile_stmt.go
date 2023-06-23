package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
)

func (m *Module) preCompile() {
	stack := m.ScopeStack
	stack.push()

	for _, stmt := range m.Ast.Body {
		switch stmt.Node.(type) {
		case *ast.UseModuleStmt:
			m.compileUseModuleStmt(stmt.Node.(*ast.UseModuleStmt))
		case *ast.FuncDecl:
			m.compileFuncDecl(stmt.Node.(*ast.FuncDecl), nil, true)
		case *ast.ImplDecl:
			m.compileImplDecl(stmt.Node.(*ast.ImplDecl), true)
		case *ast.VarDecl:
			m.compileVarDecl(stmt.Node.(*ast.VarDecl), true)
		case *ast.TypeAliasDecl:
			m.compileTypeAliasDecl(stmt.Node.(*ast.TypeAliasDecl))
		case *ast.TypeInterfaceDecl:
			m.compileTypeInterfaceDecl(stmt.Node.(*ast.TypeInterfaceDecl))
		case *ast.TypeStructDecl:
			m.compileTypeStructDecl(stmt.Node.(*ast.TypeStructDecl))
		case *ast.TypeEnumDecl:
			m.compileTypeEnumDecl(stmt.Node.(*ast.TypeEnumDecl))
		}
	}
}

func (m *Module) compileFile() {
	for _, stmt := range m.Ast.Body {
		m.compileStmt(stmt)
	}
}

func (m *Module) compileStmt(stmt *ast.Stmt) {
	switch (stmt.Node).(type) {
	case *ast.UseModuleStmt:
	case *ast.FuncDecl:
	case *ast.ImplDecl:
	case *ast.VarDecl:
	case *ast.BlockStmt:
	case *ast.ReturnStmt:
	case *ast.ExprStmt:
	case *ast.IfStmt:
	case *ast.ForStmt:
	case *ast.BreakStmt:
	case *ast.ContinueStmt:
	}
}

func (m *Module) compileUseModuleStmt(node *ast.UseModuleStmt) {
	// TODO
}

func (m *Module) compileFuncDecl(node *ast.FuncDecl, target Kind, isPrecompile bool) *FuncValue {
	name := node.Name.Name
	var value *FuncValue

	if isPrecompile {
		value = &FuncValue{
			Name: name,
			Kind: m.compileKindExpr(node.Kind),
		}

		if target != nil {
			if target.getImpl().getFunc(name) != nil {
				// TODO
				panic("duplicate method " + name)
			}
			target.getImpl().addFunc(value)
		} else {
			m.validateValueScope(node.Name)
			m.ScopeStack.putValue(name, value)
			if node.Pub {
				m.PublicScope.setValue(name, value)
			}
		}
	} else {
		if target != nil {
			value = target.getImpl().getFunc(name)
		} else {
			value = m.ScopeStack.findValue(name).(*FuncValue)
		}

		// TODO ptr
	}

	return value
}

func (m *Module) compileImplDecl(node *ast.ImplDecl, isPrecompile bool) {
	target := m.compileKindExpr(node.Target)
	switch target.(type) {
	case *TInterface:
		// TODO
		panic("can not impl interface")
	case *TAny:
		// TODO
		panic("can not impl any")
	}

	impls := make(map[string]*FuncValue)
	for _, stmt := range node.Body {
		val := m.compileFuncDecl(stmt.Node.(*ast.FuncDecl), target, isPrecompile)
		if isPrecompile {
			impls[val.Name] = val
		}
	}

	if isPrecompile {
		if node.Interface != nil {
			t, ok := m.compileKindExpr(node.Interface).(*TInterface)
			if ok {
				t.Refers = append(t.Refers, target)
				for k, v := range t.Properties {
					if impls[k] == nil {
						// TODO
						panic("missing func: " + k)
					}
					if !compareKind(v, impls[k].Kind, true) {
						// TODO
						panic("can math func sign: " + k)
					}
				}
			} else {
				// TODO
				panic("unexpected kind")
			}
		}
	}
}

func (m *Module) compileVarDecl(node *ast.VarDecl, isPrecompile bool) {
	name := node.Id.Name

	if isPrecompile {
		m.validateValueScope(node.Id)
		scope := &VarValue{
			Name:  name,
			Kind:  m.compileKindExpr(node.Kind),
			Const: node.Const,
		}
		m.ScopeStack.putValue(name, scope)
		if node.Pub {
			m.PublicScope.setValue(name, scope)
		}
	} else {
		// TODO assignment
		// TODO infer kind
	}
}

func (m *Module) compileTypeAliasDecl(node *ast.TypeAliasDecl) {
	m.validateKindScope(node.Name)

	name := node.Name.Name
	kind := &TCustom{
		Id:   getNextTypeId(),
		Kind: m.compileKindExpr(node.Kind),
	}

	m.ScopeStack.putKind(name, kind)
	if node.Pub {
		m.PublicScope.setKind(name, kind)
	}
}

func (m *Module) compileTypeInterfaceDecl(node *ast.TypeInterfaceDecl) {
	m.validateKindScope(node.Name)

	name := node.Name.Name
	props := make(map[string]Kind)

	for _, item := range node.Properties {
		key := item.Key.Name
		_, has := props[key]
		if has {
			// TODO duplicate
			panic("duplicate " + key)

		}
		props[key] = m.compileKindExpr(item.Kind)
	}

	kind := &TInterface{
		Id:         getNextTypeId(),
		Properties: props,
		Refers:     make([]Kind, 0, helper.DefaultCap),
	}

	m.ScopeStack.putKind(name, kind)
	if node.Pub {
		m.PublicScope.setKind(name, kind)
	}
}

func (m *Module) compileTypeStructDecl(node *ast.TypeStructDecl) {
	m.validateKindScope(node.Name)

	name := node.Name.Name
	kind := m.compileStructKind(node.Kind.Node.(*ast.TypeStructKind), true)

	m.ScopeStack.putKind(name, kind)
	if node.Pub {
		m.PublicScope.setKind(name, kind)
	}
}

func (m *Module) compileTypeEnumDecl(node *ast.TypeEnumDecl) {
	m.validateKindScope(node.Name)

	name := node.Name.Name
	choices := make(map[string]int)

	for i, item := range node.Choices {
		choiceName := item.Name
		_, has := choices[choiceName]
		if has {
			// TODO
			panic("duplicate " + choiceName)
		}
		choices[choiceName] = i
	}

	kind := &TEnum{
		Id:      getNextTypeId(),
		Choices: choices,
	}

	m.ScopeStack.putKind(name, kind)
	if node.Pub {
		m.PublicScope.setKind(name, kind)
	}
}
