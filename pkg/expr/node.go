package expr

type Node interface {
	ResultSize() int
	ForceSize(size int)
	Eval() int

	Resolve(label string, val int)
	IsResolved() bool
	UnresolvedSymbols() map[string]bool

	MarkRelative()
	IsRelative() bool
}
