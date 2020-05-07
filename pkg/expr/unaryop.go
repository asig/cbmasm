package expr

import "fmt"

type UnaryOp int

const (
	HiByte UnaryOp = iota
	LoByte
	Neg
	Not
)

type UnaryOpNode struct {
	node Node
	op   UnaryOp
}

func NewUnaryOp(node Node, op UnaryOp) Node {
	return &UnaryOpNode{
		node: node,
		op:   op,
	}
}

func (n *UnaryOpNode) ResultSize() int {
	switch n.op {
	case Neg:
		return n.node.ResultSize()
	case LoByte, HiByte:
		return 1
	}
	panic(fmt.Sprintf("Unimplemented UnaryOp %d", n.op))
}

func (n *UnaryOpNode) ForceSize(size int) bool {
	return n.node.ForceSize(size)
}

func (n *UnaryOpNode) Eval() int {
	if !n.IsResolved() {
		panic("Can't evaluate unresolved expr node")
	}
	v := n.node.Eval()
	switch n.op {
	case Neg:
		return -v
	case Not:
		return ^v
	case HiByte:
		return (v >> 8) & 0xff
	case LoByte:
		return v & 0xff
	}
	panic(fmt.Sprintf("Unimplemented UnaryOp %d", n.op))
}

func (n *UnaryOpNode) IsResolved() bool {
	return n.node.IsResolved()
}

func (n *UnaryOpNode) Resolve(label string, val int) {
	n.node.Resolve(label, val)
}

func (n *UnaryOpNode) UnresolvedSymbols() map[string]bool {
	return n.node.UnresolvedSymbols()
}

func (n *UnaryOpNode) MarkRelative() {
	n.node.MarkRelative()
}

func (n *UnaryOpNode) IsRelative() bool {
	return n.node.IsRelative()
}
