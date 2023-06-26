package compiler

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/bytecode"
)

func (m *Module) compileExpr(expr *ast.Expr) *bytecode.NValue {
	switch (expr.Node).(type) {
	case *ast.CallExpr:
		return m.compileCallExpr(expr.Node.(*ast.CallExpr))
	case *ast.MemberExpr:
		return m.compileMemberExpr(expr.Node.(*ast.MemberExpr))
	case *ast.BinaryExpr:
		return m.compileBinaryExpr(expr.Node.(*ast.BinaryExpr))
	case *ast.BinaryTypeExpr:
		return m.compileBinaryTypeExpr(expr.Node.(*ast.BinaryTypeExpr))
	case *ast.UnaryExpr:
		return m.compileUnaryExpr(expr.Node.(*ast.UnaryExpr))
	case *ast.FuncExpr:
		return m.compileFuncExpr(expr.Node.(*ast.FuncExpr))
	case *ast.StructExpr:
		return m.compileStructExpr(expr.Node.(*ast.StructExpr))
	case *ast.ArrayExpr:
		return m.compileArrayExpr(expr.Node.(*ast.ArrayExpr))
	case *ast.IdentifierLiteral:
		return m.compileIdentifierLiteral(expr.Node.(*ast.IdentifierLiteral))
	case *ast.NumberLiteral:
		return m.compileNumberLiteral(expr.Node.(*ast.NumberLiteral))
	case *ast.BoolLiteral:
		return m.compileBoolLiteral(expr.Node.(*ast.BoolLiteral))
	case *ast.NullLiteral:
		return m.compileNullLiteral(expr.Node.(*ast.NullLiteral))
	case *ast.StringLiteral:
		return m.compileStringLiteral(expr.Node.(*ast.StringLiteral))
	case *ast.CharLiteral:
		return m.compileCharLiteral(expr.Node.(*ast.CharLiteral))
	default:
		panic("Internal Err")
	}
}

func (m *Module) compileCallExpr(expr *ast.CallExpr) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileMemberExpr(expr *ast.MemberExpr) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileBinaryExpr(expr *ast.BinaryExpr) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileBinaryTypeExpr(expr *ast.BinaryTypeExpr) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileUnaryExpr(expr *ast.UnaryExpr) *bytecode.NValue {
	compileValue := bytecode.NewNValue()

	kind, _ := m.inferKind(expr.Argument)

	switch expr.Operator.Value {
	// number op
	case "+", "-", "++", "--":
		kind.current = typeNumber

	// logic
	case "!":
		kind.current = typeBool

	// bit op
	case "~":
		kind.current = typeNumber

	default:
		panic("Internal Err")
	}

	return compileValue
}

func (m *Module) compileFuncExpr(expr *ast.FuncExpr) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileStructExpr(expr *ast.StructExpr) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileArrayExpr(expr *ast.ArrayExpr) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileIdentifierLiteral(expr *ast.IdentifierLiteral) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileNumberLiteral(expr *ast.NumberLiteral) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileBoolLiteral(expr *ast.BoolLiteral) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileNullLiteral(expr *ast.NullLiteral) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileStringLiteral(expr *ast.StringLiteral) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}

func (m *Module) compileCharLiteral(expr *ast.CharLiteral) *bytecode.NValue {
	compileValue := bytecode.NewNValue()
	return compileValue
}
