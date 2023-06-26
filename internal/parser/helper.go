package parser

import (
	"github.com/peakchen90/noah-lang/internal/ast"
	"github.com/peakchen90/noah-lang/internal/lexer"
)

func newIdentifier(token *lexer.Token) *ast.Identifier {
	return &ast.Identifier{
		Name:     token.Value,
		Position: token.Position,
	}
}

func newKindIdentifier(token *lexer.Token) *ast.Identifier {
	return &ast.Identifier{
		Name:     token.Value,
		Position: token.Position,
	}
}

func newOperator(token *lexer.Token) *ast.Operator {
	return &ast.Operator{
		Value:    token.Value,
		Position: token.Position,
	}
}

func newIdentifierExpr(token *lexer.Token) *ast.Expr {
	return &ast.Expr{
		Node:     &ast.IdentifierLiteral{Name: newIdentifier(token)},
		Position: token.Position,
	}
}

func exprToKindExpr(expr *ast.Expr) *ast.KindExpr {
	if expr == nil {
		return nil
	}

	switch expr.Node.(type) {
	case *ast.IdentifierLiteral:
		return &ast.KindExpr{
			Node:     &ast.TIdentifier{Name: expr.Node.(*ast.IdentifierLiteral).Name},
			Position: expr.Position,
		}

	case *ast.MemberExpr:
		currentExpr := expr
		currentKind := &ast.TMemberKind{}
		result := &ast.KindExpr{
			Node:     currentKind,
			Position: currentExpr.Position,
		}

		for {
			switch currentExpr.Node.(type) {
			case *ast.MemberExpr:
				node := currentExpr.Node.(*ast.MemberExpr)
				if node.Computed {
					panic("Internal Err")
				}

				currentKind.Right = &ast.KindExpr{
					Node:     &ast.TIdentifier{Name: node.Property.Node.(*ast.IdentifierLiteral).Name},
					Position: node.Property.Position,
				}

				currentExpr = node.Object
				currentKind = &ast.TMemberKind{}
				currentKind.Left = &ast.KindExpr{
					Node:     currentKind,
					Position: currentExpr.Position,
				}
			case *ast.IdentifierLiteral:
				node := currentExpr.Node.(*ast.IdentifierLiteral)
				currentKind.Left = &ast.KindExpr{
					Node:     &ast.TIdentifier{Name: node.Name},
					Position: node.Name.Position,
				}
				return result
			}
		}
	default:
		panic("Internal Err")
	}
}
