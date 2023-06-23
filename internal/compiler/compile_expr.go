package compiler

import "github.com/peakchen90/noah-lang/internal/ast"

func (m *Module) compileExpr(expr *ast.Expr) {
	switch (expr.Node).(type) {
	case *ast.CallExpr:
	case *ast.MemberExpr:
	case *ast.BinaryExpr:
	case *ast.BinaryTypeExpr:
	case *ast.UnaryExpr:
	case *ast.FuncExpr:
	case *ast.StructExpr:
	case *ast.ArrayExpr:
	case *ast.IdentifierLiteral:
	case *ast.NumberLiteral:
	case *ast.BoolLiteral:
	case *ast.NullLiteral:
	case *ast.StringLiteral:
	case *ast.CharLiteral:
	}
}
