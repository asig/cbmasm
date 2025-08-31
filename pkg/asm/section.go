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
package asm

import (
	"github.com/asig/cbmasm/pkg/errors"
)

type Section struct {
	ignore    bool
	errorSink errors.Sink
	org       int
	bytes     []byte
}

func NewSection(org int, errorSink errors.Sink) *Section {
	return &Section{errorSink: errorSink, org: org, ignore: false}
}

func (section *Section) Emit(b byte) {
	section.bytes = append(section.bytes, b)
	section.ignore = false
}

func (section *Section) Org() int {
	if section == nil {
		return 0
	}
	return section.org
}

func (section *Section) Size() int {
	if section == nil {
		return 0
	}
	return len(section.bytes)
}

func (section *Section) PC() int {
	if section == nil {
		return 0
	}
	return section.org + len(section.bytes)
}

func (section *Section) applyPatch(p patch) {
	// TODO(asigner): Add warning for JMP ($xxFF)
	p.node.CheckRange(section.errorSink)
	val := p.node.Eval()
	if p.node.IsRelative() {
		val = val - (p.pc + 1)
		if val < -128 || val > 127 {
			section.errorSink.AddError(p.node.Pos(), "Branch target too far away.")
		}
	}
	size := p.node.ResultSize()
	pos := p.pc - section.org
	for size > 0 {
		section.bytes[pos] = byte(val & 0xff)
		val = val >> 8
		pos = pos + 1
		size = size - 1
	}
}
