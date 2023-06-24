package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"strings"
)

func (m *Module) compileFile() {
	for _, stmt := range m.Ast.Body {
		m.compileStmt(stmt)
	}
}

func (m *Module) compileStmt(stmt *ast.Stmt) {
	switch (stmt.Node).(type) {
	case *ast.ImportDecl:
		m.compileImportDecl(stmt.Node.(*ast.ImportDecl), false)
	case *ast.FuncDecl:
		m.compileFuncDecl(stmt.Node.(*ast.FuncDecl), nil, false)
	case *ast.ImplDecl:
		m.compileImplDecl(stmt.Node.(*ast.ImplDecl))
	case *ast.VarDecl:
		m.compileVarDecl(stmt.Node.(*ast.VarDecl), false)
	case *ast.BlockStmt:
		m.compileBlockStmt(stmt.Node.(*ast.BlockStmt))
	case *ast.ReturnStmt:
		m.compileReturnStmt(stmt.Node.(*ast.ReturnStmt))
	case *ast.ExprStmt:
		m.compileExprStmt(stmt.Node.(*ast.ExprStmt))
	case *ast.IfStmt:
		m.compileIfStmt(stmt.Node.(*ast.IfStmt))
	case *ast.ForStmt:
		m.compileForStmt(stmt.Node.(*ast.ForStmt))
	case *ast.BreakStmt:
		m.compileBreakStmt(stmt.Node.(*ast.BreakStmt))
	case *ast.ContinueStmt:
		m.compileContinueStmt(stmt.Node.(*ast.ContinueStmt))
	case *ast.TAliasDecl:
		m.compileTAliasDecl(stmt.Node.(*ast.TAliasDecl), false)
	case *ast.TInterfaceDecl:
		m.compileTInterfaceDecl(stmt.Node.(*ast.TInterfaceDecl), false)
	case *ast.TStructDecl:
		m.compileTStructDecl(stmt.Node.(*ast.TStructDecl), false)
	case *ast.TEnumDecl:
		m.compileTEnumDecl(stmt.Node.(*ast.TEnumDecl), false)
	default:
		m.unexpectedPos(stmt.Start, "unexpected stmt")
	}
}

func (m *Module) compileImportDecl(node *ast.ImportDecl, isPrecompile bool) {
	builder := strings.Builder{}
	if node.Package != nil {
		builder.WriteString(node.Package.Name)
		builder.WriteString(":")
	}
	for i, item := range node.Paths {
		builder.WriteString(item.Name)
		if i < len(node.Paths)-1 {
			builder.WriteString(".")
		}
	}

	moduleId := builder.String()
	module, has := m.compiler.Modules.find(moduleId)

	if !has {
		_mod, err := NewModule(m.compiler).resolve(moduleId)
		if err != nil {
			m.unexpectedPos(node.Paths[0].Start, err.Error())
		} else {
			module = _mod
		}
	}

	local := node.Local
	if local == nil {
		local = node.Paths[len(node.Paths)-1]
	}

	if isPrecompile {
		value := &ModuleValue{
			Name:   local.Name,
			Module: module,
		}
		m.putValue(local, value, true)

		_, err := module.parse()
		if err != nil {
			m.unexpectedPos(node.Paths[0].Start, err.Error())
		}
		module.precompile()
	} else {
		module.compile()
	}
}

func (m *Module) compileFuncDecl(node *ast.FuncDecl, target *KindRef, isPrecompile bool) *FuncValue {
	name := node.Name
	var value *FuncValue

	if isPrecompile {
		value = &FuncValue{
			Name:    name.Name,
			KindRef: &KindRef{},
		}
		m.putValue(name, value, true)
		if node.Pub {
			m.exports.setValue(name.Name, value)
		}
		return value
	}

	// 实现 struct methods
	if target != nil {
		targetImpls := target.Ref.getImpl()
		if targetImpls.hasFunc(name.Name) {
			m.unexpectedPos(node.Name.Start, "Duplicate key: "+name.Name)
		}

		value = &FuncValue{
			Name:    name.Name,
			KindRef: &KindRef{},
		}
		targetImpls.addFunc(value)
	} else {
		value = m.findFunc(name, true)
	}

	value.KindRef.Ref = m.compileKindExpr(node.Kind).Ref

	// TODO ptr

	return value
}

func (m *Module) compileImplDecl(node *ast.ImplDecl) {
	target := m.compileKindExpr(node.Target)

	// push scope : 用于存放 self 指向
	m.scopes.push()
	m.putSelfKind(target)
	m.putSelfValue(&SelfValue{KindRef: target})

	switch target.Ref.(type) {
	case *TInterface:
		m.unexpectedPos(node.Target.Start, "Cannot implements for interface type")
	case *TAny:
		m.unexpectedPos(node.Target.Start, "Cannot implements for any type")
	case *TSelf:
		m.unexpectedPos(node.Target.Start, "Cannot implements for self type")
	}

	implValues := make(map[string]*FuncValue)
	implDecls := make(map[string]*ast.Stmt)
	for _, stmt := range node.Body {
		val := m.compileFuncDecl(stmt.Node.(*ast.FuncDecl), target, false)
		implValues[val.Name] = val
		implDecls[val.Name] = stmt
	}

	if node.Interface != nil {
		t, ok := m.compileKindExpr(node.Interface).Ref.(*TInterface)
		if ok {
			t.Refs = append(t.Refs, target)
			for key, kind := range t.Properties {
				if implValues[key] == nil {
					m.unexpectedPos(node.Target.Start, "No implement method: "+key)
				}
				if !compareKind(kind, implValues[key].KindRef, true) {
					// TODO panic
					m.unexpectedPos(implDecls[key].Start, "Unable to match interface method signature: "+key)
				}
			}
		} else {
			if t == nil {
				m.unexpectedPos(node.Interface.Start, "Cannot found: "+getKindExprId(node.Interface))
			}
			m.unexpectedPos(node.Interface.Start, "Expect be an interface type")
		}
	}

	m.scopes.pop()
}

func (m *Module) compileVarDecl(node *ast.VarDecl, isPrecompile bool) {
	name := node.Id
	if isPrecompile {
		scope := &VarValue{
			Name:    name.Name,
			KindRef: &KindRef{},
			Const:   node.Const,
		}
		m.putValue(name, scope, true)
		if node.Pub {
			m.exports.setValue(name.Name, scope)
		}
		return
	}

	value := m.findVar(name, true)
	value.KindRef.Ref = m.compileKindExpr(node.Kind).Ref

	// TODO assignment
	// TODO infer kind
}

func (m *Module) compileBlockStmt(node *ast.BlockStmt) {
}

func (m *Module) compileReturnStmt(node *ast.ReturnStmt) {
}

func (m *Module) compileExprStmt(node *ast.ExprStmt) {
}

func (m *Module) compileIfStmt(node *ast.IfStmt) {
}

func (m *Module) compileForStmt(node *ast.ForStmt) {
}

func (m *Module) compileBreakStmt(node *ast.BreakStmt) {
}

func (m *Module) compileContinueStmt(node *ast.ContinueStmt) {
}

/* type decl */

func (m *Module) processKindDecl(name *ast.Identifier, pub bool, isPrecompile bool) *KindRef {
	if isPrecompile {
		kind := &KindRef{}
		m.putKind(name, kind, true)
		if pub {
			m.exports.setKind(name.Name, kind)
		}
		return kind
	}
	return m.findIdentifierKind(name, true)
}

func (m *Module) compileTAliasDecl(node *ast.TAliasDecl, isPrecompile bool) {
	kind := m.processKindDecl(node.Name, node.Pub, isPrecompile)
	if isPrecompile {
		return
	}

	kind.Ref = &TCustom{
		KindRef: m.compileKindExpr(node.Kind),
		Impl:    newImpl(),
	}
}

func (m *Module) compileTInterfaceDecl(node *ast.TInterfaceDecl, isPrecompile bool) {
	kind := m.processKindDecl(node.Name, node.Pub, isPrecompile)
	if isPrecompile {
		return
	}

	refType := &TInterface{
		Properties: make(map[string]*KindRef),
		Refs:       make([]*KindRef, 0, helper.DefaultCap),
	}
	kind.Ref = refType

	// push scope : 用于存放 self 指向
	m.scopes.push()
	m.putSelfKind(kind)

	for _, pair := range node.Properties {
		key := pair.Key.Name
		_, has := refType.Properties[key]
		if has {
			m.unexpectedPos(pair.Key.Start, "Duplicate key: "+key)
		} else if key[0] == '_' {
			m.unexpectedPos(pair.Key.Start, "Should not be private method: "+key)
		}
		refType.Properties[key] = m.compileKindExpr(pair.Kind)
	}

	m.scopes.pop()
}

func (m *Module) compileTStructDecl(node *ast.TStructDecl, isPrecompile bool) {
	kind := m.processKindDecl(node.Name, node.Pub, isPrecompile)
	if isPrecompile {
		return
	}

	kind.Ref = m.compileStructKind(node.Kind.Node.(*ast.TStructKind)).Ref
}

func (m *Module) compileTEnumDecl(node *ast.TEnumDecl, isPrecompile bool) {
	kind := m.processKindDecl(node.Name, node.Pub, isPrecompile)
	if isPrecompile {
		return
	}

	choices := make(map[string]int)

	for i, item := range node.Choices {
		name := item.Name
		_, has := choices[name]
		if has {
			m.unexpectedPos(item.Start, "Duplicate item: "+name)
		}
		choices[name] = i
	}

	kind.Ref = &TEnum{
		Choices: choices,
	}
}
