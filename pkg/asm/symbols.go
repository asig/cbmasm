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
	"fmt"
	"github.com/asig/cbmasm/pkg/expr"
	"strings"
)

type symbolKind int

const (
	symbolLabel symbolKind = iota
	symbolConst
	symbolMacro
)

type symbolType int

const (
	symbolTypeInt symbolType = iota
	symbolTypeString
)

type symbol struct {
	name string
	kind symbolKind
	typ  symbolType
	val  expr.Node // Only set for symbolKind in { symbolLabel, symbolConst }
	m    *macro    // only set for symbolKind in { symbolMacro }
}

type symbolTable struct {
	m map[string]*symbol
}

func newSymbolTable() symbolTable {
	return symbolTable{m: make(map[string]*symbol)}
}

func (t *symbolTable) add(sym symbol) error {
	n := strings.ToLower(sym.name)
	if _, found := t.m[n]; found {
		return fmt.Errorf("Symbol %q already defined", sym.name)
	}
	t.m[n] = &sym
	return nil
}

func (t *symbolTable) get(name string) (*symbol, bool) {
	n := strings.ToLower(name)
	s, found := t.m[n]
	if !found {
		return nil, false
	}
	return s, true
}

func (t *symbolTable) remove(name string) {
	n := strings.ToLower(name)
	delete(t.m, n)
}

func (t *symbolTable) removeMatching(predicate func(*symbol) bool) {
	for key, val := range t.m {
		if predicate(val) {
			delete(t.m, key)
		}
	}
}

func (t *symbolTable) symbols() []*symbol {
	var res []*symbol
	for _, s := range t.m {
		res = append(res, s)
	}
	return res
}
