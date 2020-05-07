package expr

import "fmt"

type BinaryOp int

const (
	Add BinaryOp = iota
	Sub
	Mul
	Mod
	Div
	And
	Or
	Xor
)

type BinaryOpNode struct {
	left, right Node
	op          BinaryOp
}

func NewBinaryOp(left, right Node, op BinaryOp) Node {
	return &BinaryOpNode{
		left:  left,
		right: right,
		op:    op,
	}
}

func max(i1, i2 int) int {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}

func (n *BinaryOpNode) ResultSize() int {
	return max(n.left.ResultSize(), n.right.ResultSize())
}

func (n *BinaryOpNode) ForceSize(size int) bool {
	b1 := n.left.ForceSize(size)
	b2 := n.right.ForceSize(size)
	return b1 && b2
}

func (n *BinaryOpNode) Eval() int {
	if !n.IsResolved() {
		panic("Can't evaluate non-const expr node")
	}
	l := n.left.Eval()
	r := n.right.Eval()
	switch n.op {
	case Add:
		return l + r
	case Sub:
		return l - r
	case Mul:
		return l * r
	case Mod:
		return l % r
	case Div:
		return l / r
	case And:
		return l & r
	case Or:
		return l | r
	case Xor:
		return l ^ r
	}
	panic(fmt.Sprintf("Unimplemented BinaryOp %d", n.op))
}

func (n *BinaryOpNode) IsResolved() bool {
	return n.left.IsResolved() && n.right.IsResolved()
}

func (n *BinaryOpNode) Resolve(label string, val int) {
	n.left.Resolve(label, val)
	n.right.Resolve(label, val)
}

func (n *BinaryOpNode) UnresolvedSymbols() map[string]bool {
	m := map[string]bool{}
	for s := range n.left.UnresolvedSymbols() {
		m[s] = true
	}
	for s := range n.right.UnresolvedSymbols() {
		m[s] = true
	}
	return m
}

func (n *BinaryOpNode) MarkRelative() {
	n.left.MarkRelative()
	n.right.MarkRelative()
}

func (n *BinaryOpNode) IsRelative() bool {
	return n.left.IsRelative() || n.right.IsRelative()
}
