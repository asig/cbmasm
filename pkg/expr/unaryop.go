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

import (
	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/text"
)

type UnaryOp struct {
	transformation func(int) int
	size           func(Node) int
}

var (
	HiByte = UnaryOp{
		transformation: func(v int) int { return (v >> 8) & 0xff },
		size:           func(_ Node) int { return 1 },
	}
	LoByte = UnaryOp{
		transformation: func(v int) int { return v & 0xff },
		size:           func(_ Node) int { return 1 },
	}
	Neg = UnaryOp{
		transformation: func(v int) int { return -v },
		size:           func(n Node) int { return n.ResultSize() },
	}
	Not = UnaryOp{
		transformation: func(v int) int { return ^v },
		size:           func(n Node) int { return n.ResultSize() },
	}
	ScreenCode = UnaryOp{
		transformation: func(v int) int { return int(petToScreen[v&0xff]) },
		size:           func(n Node) int { return 1 },
	}
	AsciiToPetscii = UnaryOp{
		transformation: func(v int) int { return int(ascToPet[v&0xff]) },
		size:           func(n Node) int { return 1 },
	}
)

type UnaryOpNode struct {
	baseNode

	pos  text.Pos
	node Node
	op   UnaryOp
}

func NewUnaryOp(pos text.Pos, node Node, op UnaryOp) Node {
	return &UnaryOpNode{
		pos:  pos,
		node: node,
		op:   op,
	}
}

func (n *UnaryOpNode) ResultSize() int {
	return n.op.size(n.node)
}

func (n *UnaryOpNode) ForceSize(size int) bool {
	return n.node.ForceSize(size)
}

func (n *UnaryOpNode) Eval() int {
	if !n.IsResolved() {
		panic("Can't evaluate unresolved expr node")
	}
	v := n.node.Eval()
	return n.op.transformation(v)
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

func (n *UnaryOpNode) Pos() text.Pos {
	return n.pos
}

func (n *UnaryOpNode) CheckRange(sink errors.Sink) {
	n.node.CheckRange(sink)
	checkRange(n, sink)
}
