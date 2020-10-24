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

type SymbolRefNode struct {
	pos        text.Pos
	symbol     string
	maxSize    int
	val        int
	resolved   bool
	isRelative bool
}

func NewSymbolRef(pos text.Pos, symbol string, maxSize, val int) Node {
	return &SymbolRefNode{
		pos:        pos,
		symbol:     symbol,
		maxSize:    maxSize,
		val:        val,
		resolved:   true,
		isRelative: false,
	}
}

func NewUnresolvedSymbol(pos text.Pos, symbol string, maxSize int) Node {
	return &SymbolRefNode{
		pos:        pos,
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

func (n *SymbolRefNode) ForceSize(size int) bool {
	if !n.resolved {
		n.maxSize = size
		return true
	}

	if n.val < 1<<(size*8) {
		n.maxSize = size
		return true
	}
	return false
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
	return map[string]bool{n.symbol: true}
}

func (n *SymbolRefNode) MarkRelative() {
	n.isRelative = true
	n.maxSize = 1
}

func (n *SymbolRefNode) IsRelative() bool {
	return n.isRelative
}

func (n *SymbolRefNode) Pos() text.Pos {
	return n.pos
}
