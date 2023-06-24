package compiler

import (
	"fmt"
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/helper"
	"strings"
)

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
		panic("Internal Err")
	}
}

func (m *Module) compileImportDecl(node *ast.ImportDecl, isPrecompile bool) {
	builder := strings.Builder{}
	importPathStart := node.Paths[0].Start

	if node.Package != nil {
		importPathStart = node.Package.Start
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
		if !isPrecompile {
			panic("Internal Err")
		}

		_mod, err := NewModule(m.compiler).resolve(moduleId)
		if err != nil {
			m.unexpectedPos(importPathStart, err.Error())
		} else {
			module = _mod
		}
	}

	local := node.Local
	if local == nil {
		local = node.Paths[len(node.Paths)-1]
	}

	if isPrecompile {
		m.scopes.putModule(local, module, true)
		err := module.parse()
		if err != nil {
			m.unexpectedPos(importPathStart, err.Error())
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
			Name: name.Name,
			Kind: &KindRef{},
		}
		m.scopes.putValue(name, value, true)
		if node.Pub {
			m.exports.setValue(name.Name, value)
		}
		return value
	}

	// impl struct methods
	if target != nil {
		impls := target.Ref.getImpl()
		if impls.hasFunc(name.Name) {
			m.unexpectedPos(node.Name.Start, "duplicate key: "+name.Name)
		}

		value = &FuncValue{
			Name: name.Name,
			Kind: &KindRef{},
		}
		impls.addFunc(value)
	} else {
		value = m.scopes.findFuncValue(name, true)
	}

	funcKindNode := node.Kind.Node.(*ast.TFuncKind)

	// compile func kind
	kind := m.compileKindExpr(node.Kind)
	value.Kind.Ref = kind.Ref
	funcKind := kind.Ref.(*TFunc)
	paramKinds := funcKind.Params
	// 校验 rest 参数类型
	if funcKind.RestParam && len(paramKinds) > 0 {
		restParamKind := paramKinds[len(paramKinds)-1]
		t, ok := restParamKind.Ref.(*TArray)
		if !ok || t.Len >= 0 {
			restKindNode := funcKindNode.Params[len(paramKinds)-1].Kind
			m.unexpectedPos(restKindNode.Start, "the rest parameter should be: []T")
		}
	}

	// compile func params
	m.scopes.push()
	for i, param := range funcKindNode.Params {
		paramValue := &VarValue{
			Name:  name.Name,
			Kind:  paramKinds[i],
			Const: false,
			Ptr:   0, // TODO ptr
		}
		m.scopes.putValue(param.Name, paramValue, true)
	}

	// compile func body
	body := node.Body.Node.(*ast.BlockStmt)
	m.compileBlockStmt(body)
	m.scopes.pop()

	value.Ptr = 0 // TODO ptr

	return value
}

func (m *Module) compileImplDecl(node *ast.ImplDecl) {
	target := m.compileKindExpr(node.Target)

	// push scope : 用于存放 self 指向
	m.scopes.push()
	m.scopes.putSelfKind(target)
	m.scopes.putSelfValue(&SelfValue{Kind: target})

	switch target.Ref.(type) {
	case *TInterface:
		m.unexpectedPos(node.Target.Start, "cannot implements for `interface` type")
	case *TAny:
		m.unexpectedPos(node.Target.Start, "cannot implements for `any` type")
	case *TSelf:
		m.unexpectedPos(node.Target.Start, "cannot implements for `self` type")
	}

	implValues := make(map[string]*FuncValue)
	implDecls := make(map[string]*ast.Stmt)
	for _, stmt := range node.Body.Node.(*ast.BlockStmt).Body {
		val := m.compileFuncDecl(stmt.Node.(*ast.FuncDecl), target, false)
		implValues[val.Name] = val
		implDecls[val.Name] = stmt
	}

	if node.Interface != nil {
		t, ok := m.compileKindExpr(node.Interface).Ref.(*TInterface)
		if ok {
			t.Refs = append(t.Refs, target)
			interfaceName := getKindExprString(node.Interface)
			for key, interfaceDeclKind := range t.Properties {
				if implValues[key] == nil {
					m.unexpectedPos(node.Body.Start, fmt.Sprintf("no implement method: %s.%s", interfaceName, key))
				}
				if !matchKind(interfaceDeclKind, implValues[key].Kind) {
					funcNode := implDecls[key].Node.(*ast.FuncDecl)
					m.unexpectedPos(funcNode.Name.End, fmt.Sprintf("cannot match method signature: %s.%s", interfaceName, key))
				}
			}
		} else {
			if t == nil {
				m.unexpectedPos(node.Interface.Start, "cannot found: "+getKindExprString(node.Interface))
			}
			m.unexpectedPos(node.Interface.Start, "expect be an interface type")
		}
	}

	m.scopes.pop()
}

func (m *Module) compileVarDecl(node *ast.VarDecl, isPrecompile bool) {
	name := node.Id
	if isPrecompile {
		scope := &VarValue{
			Name:  name.Name,
			Kind:  &KindRef{},
			Const: node.Const,
		}
		m.scopes.putValue(name, scope, true)
		if node.Pub {
			m.exports.setValue(name.Name, scope)
		}
		return
	}

	value := m.scopes.findVarValue(name, true)

	// 变量类型
	var kind *KindRef
	if node.Kind != nil {
		kind = m.compileKindExpr(node.Kind)
	}
	if node.Init != nil {
		inferKind, err := m.inferKind(node.Init)
		if err != nil {
			m.unexpectedPos(node.Init.Start, err.Error())
		} else if kind == nil {
			kind = inferKind
		}

		// TODO maybe assign

		value.Ptr = 0 // TODO ptr
	}

	if kind == nil {
		m.unexpectedPos(node.Id.Start, "cannot infer variable type")
	}
	value.Kind.Ref = kind.Ref
}

func (m *Module) compileBlockStmt(node *ast.BlockStmt) {
	m.scopes.push()
	for _, stmt := range node.Body {
		m.compileStmt(stmt)
	}
	m.scopes.pop()
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
		m.scopes.putKind(name, kind, true)
		if pub {
			m.exports.setKind(name.Name, kind)
		}
		return kind
	}
	return m.scopes.findIdentifierKind(name, true)
}

func (m *Module) compileTAliasDecl(node *ast.TAliasDecl, isPrecompile bool) {
	kind := m.processKindDecl(node.Name, node.Pub, isPrecompile)
	if isPrecompile {
		return
	}

	kind.Ref = &TCustom{
		Kind: m.compileKindExpr(node.Kind),
		Impl: newImpl(),
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
	m.scopes.putSelfKind(kind)

	for _, pair := range node.Properties {
		key := pair.Key.Name
		_, has := refType.Properties[key]
		if has {
			m.unexpectedPos(pair.Key.Start, "duplicate key: "+key)
		} else if key[0] == '_' {
			m.unexpectedPos(pair.Key.Start, "should not be private method: "+key)
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
			m.unexpectedPos(item.Start, "duplicate item: "+name)
		}
		choices[name] = i
	}

	kind.Ref = &TEnum{
		Choices: choices,
	}
}
