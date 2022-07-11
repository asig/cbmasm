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
	"fmt"

	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/text"
)

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
	Eq
	Ne
	Lt
	Le
	Gt
	Ge
)

type BinaryOpNode struct {
	baseNode

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

func (n *BinaryOpNode) boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (n *BinaryOpNode) Eval() int {
	if !n.IsResolved() {
		panic("Can't evaluate non-const expr node")
	}
	if n.Type() != NodeType_Int {
		panic("Can't Eval() a string or float node")
	}

	if n.left.Type() == NodeType_Int {
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
		case Eq:
			return n.boolToInt(l == r)
		case Ne:
			return n.boolToInt(l != r)
		case Lt:
			return n.boolToInt(l < r)
		case Le:
			return n.boolToInt(l <= r)
		case Gt:
			return n.boolToInt(l > r)
		case Ge:
			return n.boolToInt(l >= r)
		}
		panic(fmt.Sprintf("Unimplemented BinaryOp %d", n.op))
	}
	// Must be string
	l := n.left.EvalStr()
	r := n.right.EvalStr()
	switch n.op {
	case Eq:
		return n.boolToInt(l == r)
	case Ne:
		return n.boolToInt(l != r)
	case Lt:
		return n.boolToInt(l < r)
	case Le:
		return n.boolToInt(l <= r)
	case Gt:
		return n.boolToInt(l > r)
	case Ge:
		return n.boolToInt(l >= r)
	}
	panic(fmt.Sprintf("BinaryOp %d not supported for strings", n.op))
}

func (n *BinaryOpNode) EvalFloat() float64 {
	if !n.IsResolved() {
		panic("Can't evaluate non-const expr node")
	}
	if n.Type() != NodeType_Float {
		panic("Can't EvalFloat() a non-float node")
	}

	var l, r float64

	switch n.left.Type() {
	case NodeType_Int:
		l = float64(n.left.Eval())
	case NodeType_Float:
		l = n.left.EvalFloat()
	default:
		panic("Left side is neither int nor float")
	}
	switch n.right.Type() {
	case NodeType_Int:
		r = float64(n.right.Eval())
	case NodeType_Float:
		r = n.right.EvalFloat()
	default:
		panic("Right side is neither int nor float")
	}
	switch n.op {
	case Add:
		return l + r
	case Sub:
		return l - r
	case Mul:
		return l * r
	default:
		panic(fmt.Sprintf("Unsupported operation %d", n.op))
	}
}

func (n *BinaryOpNode) EvalStr() string {
	panic("Can't evaluate BinaryOp node as string")
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

func (n *BinaryOpNode) Pos() text.Pos {
	return n.left.Pos()
}

func (n *BinaryOpNode) CheckRange(sink errors.Sink) {
	checkRange(n, sink)
}

func (n *BinaryOpNode) Type() NodeType {
	if n.left.Type() == NodeType_Float || n.right.Type() == NodeType_Float {
		return NodeType_Float
	}
	return NodeType_Int
}
