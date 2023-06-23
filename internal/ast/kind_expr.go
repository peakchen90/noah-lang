package ast

import "github.com/peakchen90/noah-lang/internal/helper"

type KE interface{ isKindExpr() }

func (*TypeNumber) isKindExpr()     {}
func (*TypeByte) isKindExpr()       {}
func (*TypeChar) isKindExpr()       {}
func (*TypeString) isKindExpr()     {}
func (*TypeBool) isKindExpr()       {}
func (*TypeAny) isKindExpr()        {}
func (*TypeArray) isKindExpr()      {}
func (*TypeIdentifier) isKindExpr() {}
func (*TypeMemberKind) isKindExpr() {}
func (*TypeFuncKind) isKindExpr()   {}
func (*TypeStructKind) isKindExpr() {}

type (
	TypeNumber struct{}

	TypeByte struct{}

	TypeChar struct{}

	TypeString struct{}

	TypeBool struct{}

	TypeAny struct{}

	TypeArray struct {
		Kind *KindExpr
		Len  *Expr
	}

	TypeIdentifier struct {
		Name *KindIdentifier
	}

	TypeMemberKind struct {
		Left  *KindExpr
		Right *KindExpr
	}

	TypeFuncKind struct {
		Arguments []*Argument
		Return    *KindExpr
	}

	TypeStructKind struct {
		Extends    []*KindExpr
		Properties []*KindProperty
	}
)

func (t *TypeMemberKind) ToMemberIds() []string {
	members := make([]string, 0, helper.SmallCap)

	left, ok := t.Left.Node.(*TypeMemberKind)
	if ok {
		for _, s := range left.ToMemberIds() {
			members = append(members, s)
		}

	} else {
		members = append(members, t.Left.Node.(*TypeIdentifier).Name.Name)
	}

	members = append(members, t.Right.Node.(*TypeIdentifier).Name.Name)

	return members
}
