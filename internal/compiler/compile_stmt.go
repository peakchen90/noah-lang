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
		m.compileFuncDecl(stmt.Node.(*ast.FuncDecl), nil)
	case *ast.ImplDecl:
		m.compileImplDecl(stmt.Node.(*ast.ImplDecl), false)
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
	case *ast.TTypeDecl:
		m.compileTTypeDecl(stmt.Node.(*ast.TTypeDecl), false)
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

	if local.Name == "self" {
		m.unexpectedPos(local.Start, "cannot use `self` as a local identifier")
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

func (m *Module) compileFuncSign(node *ast.FuncDecl, target *KindRef, isPrecompile bool) *FuncValue {
	name := node.Name
	var value *FuncValue

	if isPrecompile {
		if name.Name == "self" {
			m.unexpectedPos(name.Start, "identifier 'self' is not allowed")
		}

		value = &FuncValue{
			Name: name.Name,
			Kind: newKindRef(m, -1),
		}

		if target != nil {
			impls := target.current.getImpl()
			if impls.hasFunc(name.Name) {
				m.unexpectedPos(node.Name.Start, "duplicate key: "+name.Name)
			}
			impls.addFunc(value)
		} else {
			m.scopes.putValue(name, value, true)
			if node.Pub {
				m.exports.setValue(name.Name, value)
			}
		}

		return value
	}

	if target != nil {
		value = target.current.getImpl().getFunc(name.Name)
	} else {
		value = m.scopes.findFuncValue(name, true)
	}

	funcKindNode := node.Kind.Node.(*ast.TFuncKind)

	// compile func kind
	kind := m.compileKindExpr(node.Kind)
	value.Kind.current = kind.current
	funcKind := kind.current.(*TFunc)
	paramKinds := funcKind.Params
	// 校验 rest 参数类型
	if funcKind.RestParam && len(paramKinds) > 0 {
		restParamKind := paramKinds[len(paramKinds)-1]
		t, ok := restParamKind.current.(*TArray)
		if !ok || t.Len >= 0 {
			restKindNode := funcKindNode.Params[len(paramKinds)-1].Kind
			m.unexpectedPos(restKindNode.Start, "the rest parameter should be: []T")
		}
	}

	value.Ptr = 0 // TODO ptr

	return value
}

func (m *Module) compileFuncDecl(node *ast.FuncDecl, target *KindRef) {
	name := node.Name
	var value *FuncValue

	if target != nil {
		value = target.current.getImpl().getFunc(name.Name)
	} else {
		value = m.scopes.findFuncValue(name, true)
	}

	funcKindNode := node.Kind.Node.(*ast.TFuncKind)
	funcKind := value.Kind.current.(*TFunc)
	paramKinds := funcKind.Params

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
}

func (m *Module) compileImplDecl(node *ast.ImplDecl, onlyFuncSign bool) {
	target := m.compileKindExpr(node.Target)

	// 编译 impl 函数签名
	if onlyFuncSign {
		// push scope : 用于存放 self 指向
		m.scopes.push()
		m.scopes.putSelfKind(target)
		m.scopes.putSelfValue(&SelfValue{Kind: target})

		switch target.current.(type) {
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
			funcNode := stmt.Node.(*ast.FuncDecl)
			value := m.compileFuncSign(funcNode, target, true)
			value = m.compileFuncSign(funcNode, target, false)
			implValues[value.Name] = value
			implDecls[value.Name] = stmt
		}

		if node.Interface != nil {
			interfaceKind := m.compileKindExpr(node.Interface)
			t, ok := interfaceKind.current.(*TInterface)
			if ok {
				interfaceKind.refs = append(interfaceKind.refs, target)
				interfaceName := getKindExprString(node.Interface)
				for key, interfaceDeclKind := range t.Properties {
					if implValues[key] == nil {
						m.unexpectedPos(node.Body.Start, fmt.Sprintf("no implement method: %s.%s", interfaceName, key))
					}
					if !matchKind(interfaceDeclKind, implValues[key].Kind, true) {
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
	} else {
		// 编译 impl 函数
		for _, stmt := range node.Body.Node.(*ast.BlockStmt).Body {
			m.compileFuncDecl(stmt.Node.(*ast.FuncDecl), target)
		}
	}

}

func (m *Module) compileVarDecl(node *ast.VarDecl, isPrecompile bool) {
	name := node.Id
	if isPrecompile {
		if name.Name == "self" {
			m.unexpectedPos(name.Start, "identifier 'self' is not allowed")
		}
		scope := &VarValue{
			Name:  name.Name,
			Kind:  newKindRef(m, -1),
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
	value.Kind.current = kind.current
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

func (m *Module) processKindDecl(initKind *KindRef, name *ast.Identifier, pub bool, isPrecompile bool) *KindRef {
	if isPrecompile {
		if name.Name == "self" {
			m.unexpectedPos(name.Start, "identifier 'self' is not allowed")
		}
		m.scopes.putKind(name, initKind, true)
		if pub {
			m.exports.setKind(name.Name, initKind)
		}
		return initKind
	}
	return m.scopes.findIdentifierKind(name, true)
}

func (m *Module) compileTTypeDecl(node *ast.TTypeDecl, isPrecompile bool) {
	initKind := newKindRef(m, -1)
	initKind.current = &TCustom{}
	kind := m.processKindDecl(initKind, node.Name, node.Pub, isPrecompile)
	if isPrecompile {
		return
	}

	kind.current = &TCustom{
		Kind: m.compileKindExpr(node.Kind),
		Impl: newImpl(),
	}
}

func (m *Module) compileTInterfaceDecl(node *ast.TInterfaceDecl, isPrecompile bool) {
	initKind := newKindRef(m, helper.SmallCap)
	initKind.current = &TInterface{}
	kind := m.processKindDecl(initKind, node.Name, node.Pub, isPrecompile)
	if isPrecompile {
		return
	}

	_type := &TInterface{
		Properties: make(map[string]*KindRef),
	}
	kind.current = _type

	// push scope : 用于存放 self 指向
	m.scopes.push()
	m.scopes.putSelfKind(kind)

	for _, pair := range node.Properties {
		key := pair.Key.Name
		_, has := _type.Properties[key]
		if has {
			m.unexpectedPos(pair.Key.Start, "duplicate key: "+key)
		} else if key[0] == '_' {
			m.unexpectedPos(pair.Key.Start, "should not be private method: "+key)
		}
		_type.Properties[key] = m.compileKindExpr(pair.Kind)
	}

	m.scopes.pop()
}

func (m *Module) compileTStructDecl(node *ast.TStructDecl, isPrecompile bool) {
	initKind := newKindRef(m, helper.SmallCap)
	initKind.current = &TStruct{}
	kind := m.processKindDecl(initKind, node.Name, node.Pub, isPrecompile)
	if isPrecompile {
		return
	}

	result := m.compileStructKind(kind, node.Kind)
	kind.current = result.current
	kind.refs = result.refs
}

func (m *Module) compileTEnumDecl(node *ast.TEnumDecl, isPrecompile bool) {
	initKind := newKindRef(m, -1)
	initKind.current = &TEnum{}
	kind := m.processKindDecl(initKind, node.Name, node.Pub, isPrecompile)
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

	kind.current = &TEnum{
		Choices: choices,
	}
}
