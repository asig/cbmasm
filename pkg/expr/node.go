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

type Node interface {
	ResultSize() int
	ForceSize(size int) bool
	Eval() int

	Resolve(label string, val int)
	IsResolved() bool
	UnresolvedSymbols() map[string]bool

	MarkRelative()
	IsRelative() bool

	MarkSigned()
	IsSigned() bool

	SetRange(min, max int)
	Range() (Range, bool)
	CheckRange(sink errors.Sink)

	Pos() text.Pos
}

type Range struct {
	min, max int
}

type baseNode struct {
	signed bool
	r      *Range
}

func (n *baseNode) IsSigned() bool {
	return n.signed
}

func (n *baseNode) MarkSigned() {
	n.signed = true
}

func (n *baseNode) SetRange(min, max int) {
	n.r = &Range{min: min, max: max}
}

func (n *baseNode) Range() (Range, bool) {
	if n.r == nil {
		return Range{}, false
	}
	return *n.r, true
}

func checkRange(n Node, sink errors.Sink) {
	size := n.ResultSize()
	val := n.Eval()
	var min, max int
	if r, ok := n.Range(); ok {
		min = r.min
		max = r.max
	} else if n.IsSigned() {
		min = (-1) << (size*8 - 1)
		max = (1 << (size*8 - 1)) - 1
	} else {
		min = 0
		max = 1<<(size*8) - 1
	}
	if val < min || val > max {
		sink.AddError(n.Pos(), "Value out of range.")
	}
}
