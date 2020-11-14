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
	// ResultSize returns the size of the result
	ResultSize() int

	// ForceSize forces a certain size, and returns false if the value is too big
	ForceSize(size int) bool

	// Eval evaluates the node and panics if the node is unresolved
	Eval() int

	// Resolve resolves symbols
	Resolve(label string, val int)

	// IsResolved returns whether the node is resolved
	IsResolved() bool

	// UnresolvedSymbols returns a list of symbols that are not yet resolved.
	UnresolvedSymbols() map[string]bool

	// MarkRelative marks the node as relative. When emitting such nodes,
	// instead of the absolute value, the difference to the PC is written out.
	MarkRelative()

	// IsRelative returns whether the node is relative.
	IsRelative() bool

	// MarkSigned marks the node as a sigend value
	MarkSigned()

	// IsSigned returns whether the node is signed
	IsSigned() bool

	// SetRange sets the valid range for this node
	SetRange(min, max int)

	// SetValidValues sets a list of valid values for this node.
	SetValidValues(v ...int)

	// Range returns the valid range, if any
	Range() (Range, bool)

	// Range returns the valid values, if any
	ValidValues() ([]int, bool)

	// CheckRange checks whether the value is in a valid range, and emits an error otherwise.
	// If a range or valid values are set, they are used to compute the validity. Otherwise, the size and whether it is
	// a signed value are used.
	CheckRange(sink errors.Sink)

	// Pos returns the position in the text of this node.
	Pos() text.Pos
}

type Range struct {
	min, max int
}

type baseNode struct {
	signed      bool
	validValues []int
	r           *Range
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

func (n *baseNode) SetValidValues(values ...int) {
	n.validValues = values
}

func (n *baseNode) Range() (Range, bool) {
	if n.r == nil {
		return Range{}, false
	}
	return *n.r, true
}

func (n *baseNode) ValidValues() ([]int, bool) {
	if n.validValues == nil {
		return nil, false
	}
	return n.validValues, true
}

func checkRange(n Node, sink errors.Sink) {
	size := n.ResultSize()
	val := n.Eval()

	if validVals, ok := n.ValidValues(); ok {
		for _, v := range validVals {
			if val == v {
				return
			}
		}
		sink.AddError(n.Pos(), "Value is not in list of supported values.")
		return
	}

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
