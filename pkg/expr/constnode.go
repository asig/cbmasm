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
	"github.com/asig/cbmasm/pkg/text"
)

type ConstNode struct {
	pos        text.Pos
	size       int
	val        int
	isRelative bool
}

func NewConst(pos text.Pos, val, size int) Node {
	return &ConstNode{
		pos:        pos,
		size:       size,
		val:        val,
		isRelative: false,
	}
}

func (n *ConstNode) ResultSize() int {
	return n.size
}

func (n *ConstNode) ForceSize(size int) bool {
	if n.val < 1<<(size*8) {
		n.size = size
		return true
	}
	return false
}

func (n *ConstNode) Eval() int {
	return n.val
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
