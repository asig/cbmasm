package expr

type SymbolRefNode struct {
	symbol     string
	maxSize    int
	val        int
	resolved   bool
	isRelative bool
}

func NewSymbolRef(symbol string, maxSize, val int) Node {
	return &SymbolRefNode{
		symbol:     symbol,
		maxSize:    maxSize,
		val:        val,
		resolved:   true,
		isRelative: false,
	}
}

func NewUnresolvedSymbol(symbol string, maxSize int) Node {
	return &SymbolRefNode{
		symbol:     symbol,
		maxSize:    maxSize,
		val:        0,
		resolved:   false,
		isRelative: false,
	}
}

func (n *SymbolRefNode) ResultSize() int {
	return n.maxSize
}

func (n *SymbolRefNode) ForceSize(size int) {
	if !n.resolved {
		n.maxSize = size
	} else {
		if n.val < 1<<(size*8) {
			n.maxSize = size
		}
	}
}

func (n *SymbolRefNode) Eval() int {
	if !n.IsResolved() {
		panic("Can't evaluate unresolved symbol node")
	}
	return n.val
}

func (n *SymbolRefNode) IsResolved() bool {
	return n.resolved
}

func (n *SymbolRefNode) Resolve(symbol string, val int) {
	if symbol == n.symbol {
		n.val = val
		n.resolved = true
	}
}

func (n *SymbolRefNode) UnresolvedSymbols() map[string]bool {
	if n.resolved {
		return nil
	}
	return map[string]bool{ n.symbol: true}
}

func (n *SymbolRefNode) MarkRelative() {
	n.isRelative = true
	n.maxSize = 1
}

func (n *SymbolRefNode) IsRelative() bool {
	return n.isRelative
}
