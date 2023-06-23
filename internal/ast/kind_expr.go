package ast

type KE interface{ isKindExpr() }

func (*TNumber) isKindExpr()     {}
func (*TByte) isKindExpr()       {}
func (*TChar) isKindExpr()       {}
func (*TString) isKindExpr()     {}
func (*TBool) isKindExpr()       {}
func (*TAny) isKindExpr()        {}
func (*TArray) isKindExpr()      {}
func (*TIdentifier) isKindExpr() {}
func (*TMemberKind) isKindExpr() {}
func (*TFuncKind) isKindExpr()   {}
func (*TStructKind) isKindExpr() {}

type (
	TNumber struct{}

	TByte struct{}

	TChar struct{}

	TString struct{}

	TBool struct{}

	TAny struct{}

	TArray struct {
		Kind *KindExpr
		Len  *Expr
	}

	TIdentifier struct {
		Name *Identifier
	}

	TMemberKind struct {
		Left  *KindExpr
		Right *KindExpr
	}

	TFuncKind struct {
		Arguments []*Argument
		Return    *KindExpr
	}

	TStructKind struct {
		Extends    []*KindExpr
		Properties []*KindProperty
	}
)
