package expr

type ConstNode struct {
	size int
	val int
	isRelative bool
}

func NewConst(val, size int) Node {
	return &ConstNode{
		size:     size,
		val:      val,
		isRelative: false,
	}
}

func (n *ConstNode) ResultSize() int {
	return n.size
}

func (n *ConstNode) ForceSize(size int) {
	if n.val < 1 << (size*8) {
		n.size = size
	}
}

func (n *ConstNode) Eval() int {
	return n.val
}

func (n *ConstNode) IsResolved() bool {
	return true
}

func (n *ConstNode) Resolve(label string, val int) {
}

func (n *ConstNode) UnresolvedSymbols() map[string]bool {
	return nil
}

func (n *ConstNode) MarkRelative() {
	n.isRelative = true
}

func (n *ConstNode) IsRelative() bool {
	return n.isRelative
}
