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

type ConstNode struct {
	baseNode

	pos        text.Pos
	size       int
	val        int
	strval     string
	floatval   float64
	typ        NodeType
	isRelative bool
}

func NewConst(pos text.Pos, val, size int) Node {
	return &ConstNode{
		pos:        pos,
		size:       size,
		val:        val,
		typ:        NodeType_Int,
		isRelative: false,
	}
}

func NewStrConst(pos text.Pos, val string) Node {
	return &ConstNode{
		pos:        pos,
		size:       len(val),
		strval:     val,
		typ:        NodeType_String,
		isRelative: false,
	}
}

func NewFloatConst(pos text.Pos, val float64) Node {
	return &ConstNode{
		pos:        pos,
		size:       5,
		floatval:   val,
		typ:        NodeType_Float,
		isRelative: false,
	}
}

func (n *ConstNode) Type() NodeType {
	return n.typ
}

func (n *ConstNode) ResultSize() int {
	return n.size
}

func (n *ConstNode) ForceSize(size int) bool {
	var min, max int
	if n.Type() == NodeType_Float {
		min = 5
		max = 5
	} else if n.IsSigned() {
		min = (-1) << (size*8 - 1)
		max = (1 << (size*8 - 1)) - 1
	} else {
		min = 0
		max = 1<<(size*8) - 1
	}
	n.size = size
	if min <= n.val && n.val <= max {
		return true
	}
	return false
}

func (n *ConstNode) Eval() int {
	if n.typ != NodeType_Int {
		panic("type is not int")
	}
	return n.val
}

func (n *ConstNode) EvalFloat() float64 {
	if n.typ != NodeType_Float {
		panic("type is not float")
	}
	return n.floatval
}

func (n *ConstNode) EvalStr() string {
	if n.typ != NodeType_String {
		panic("type is not string")
	}
	return n.strval
}

func (n *ConstNode) IsResolved() bool {
	return true
}

func (n *ConstNode) Resolve(_ string, _ int) {
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

func (n *ConstNode) Pos() text.Pos {
	return n.pos
}

func (n *ConstNode) CheckRange(sink errors.Sink) {
	checkRange(n, sink)
}
