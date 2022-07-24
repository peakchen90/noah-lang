package ast

type NodeType uint8

const (
	NodeProgram NodeType = iota + 1
	NodeImportDeclaration
	NodeFunctionDeclaration
	NodeVariableDeclaration
	NodeBlockStatement
	NodeReturnStatement
	NodeExpressionStatement
	NodeIfStatement
	NodeLoopStatement
	NodeBreakStatement
	NodeContinueStatement

	NodeImportSpecifier
	NodeCallExpression
	NodeBinaryExpression
	NodeUnaryExpression
	NodeAssignmentExpression
	NodeIdentifier
	NodeNumberLiteral
	NodeBooleanLiteral
	NodeStringLiteral
)

type Position struct {
	Start int
	End   int
}
