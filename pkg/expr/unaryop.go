/*
 * Copyright (c) 2020 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of cbmasm.
 *
 * cbmasm is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * cbmasm is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with cbmasm.  If not, see <http://www.gnu.org/licenses/>.
 */
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
