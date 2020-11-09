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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/expr"
	"github.com/asig/cbmasm/pkg/scanner"
	"github.com/asig/cbmasm/pkg/text"
)

// patch records nodes that can't be evaluated because of undefined nodes
type patch struct {
	pc   int       // Place to patch
	node expr.Node // Node that needs to be patched in
}

type param struct {
	rawText string // Only used in macro calls
	mode    AddressingMode
	val     expr.Node
}

type state int

const (
	stateAssemble state = iota
	stateRecordMacro
)

var conditionalTokens = map[scanner.TokenType]bool{
	scanner.Ifdef:  true,
	scanner.If:     true,
	scanner.Ifndef: true,
	scanner.Else:   true,
	scanner.Endif:  true,
}

var relOpToBinOp = map[scanner.TokenType]expr.BinaryOp{
	scanner.Eq: expr.Eq,
	scanner.Ne: expr.Ne,
	scanner.Lt: expr.Lt,
	scanner.Le: expr.Le,
	scanner.Gt: expr.Gt,
	scanner.Ge: expr.Ge,
}

type Assembler struct {
	includePaths  []string
	errorModifier errors.Modifier

	errors      []errors.Error
	warnings    []errors.Error
	scanner     *scanner.Scanner
	lookahead   scanner.Token
	tokenBuf    scanner.Token
	tokenBufSet bool

	assemblyEnabled stack
	state           state

	// current macro, only set when recording macros
	macro *macro

	// Code generation buffer
	section *Section

	// outstanding patches
	patchesPerLabel map[string][]patch

	// Symbol table
	symbols symbolTable

	// Macros
	macros map[string]*macro
}

func New(includePaths []string) *Assembler {
	a := &Assembler{
		includePaths: includePaths,
	}
	return a
}

func (a *Assembler) Assemble(t text.Text) {
	a.errors = nil
	a.warnings = nil
	a.section = nil
	a.patchesPerLabel = make(map[string][]patch)
	a.symbols = newSymbolTable()
	a.assemblyEnabled = stack{}
	a.assemblyEnabled.push(true)

	t = a.resolveIncludes(t)
	a.assembleText(t)

	ll := t.LastLine()
	p := text.Pos{Filename: ll.Filename, Line: ll.LineNumber, Col: 1}
	if a.state == stateRecordMacro {
		a.AddError(p, ".endm expected")
	}
	a.reportUnresolvedSymbols(p, func(string) bool { return true })
	if a.assemblyEnabled.len() > 1 {
		a.AddError(p, ".endif expected")
	}
}

func (a *Assembler) resolveIncludes(t text.Text) text.Text {
	res := text.Text{}
	for _, line := range t.Lines {
		a.beginLine(line)

		t, _, label := a.maybeLabel()
		if t.Type == scanner.Ident && t.StrVal == "include" {
			a.match(scanner.Ident)
			p := a.lookahead.Pos
			filename := a.lookahead.StrVal
			a.match(scanner.String)
			f := a.findIncludeFile(filename)
			if f == nil {
				a.AddError(p, "Can't find file %q in include paths.", filename)
				res.AppendLine(line)
				continue
			}
			content, err := ioutil.ReadFile(*f)
			if err != nil {
				a.AddError(p, "Can't read file %q: %s", *f, err)
			}

			if label != "" {
				l := a.scanner.Line()
				res.AppendLine(text.Line{l.Filename, l.LineNumber, []rune(label)})
			}
			included := a.resolveIncludes(text.Process(filename, string(content)))
			res.Append(included)
		} else {
			res.AppendLine(*a.scanner.Line())
		}
	}
	return res
}

func (a *Assembler) assembleText(t text.Text) {
	a.state = stateAssemble
	for _, line := range t.Lines {
		a.beginLine(line)
		a.processLine()
	}
}

func (a *Assembler) beginLine(line text.Line) {
	a.scanner = scanner.New(line, a)
	a.tokenBufSet = false
	a.lookahead = a.scanner.Scan()
}

func (a *Assembler) Origin() int {
	return a.section.Org()
}

func (a *Assembler) GetBytes() []byte {
	return a.section.bytes
}

func (a *Assembler) maybeLabel() (scanner.Token, text.Pos, string) {
	label := ""
	var labelPos text.Pos
	t := a.lookahead
	if t.Type == scanner.Ident {
		if t.Pos.Col == 1 {
			labelPos = t.Pos
			label = t.StrVal
			a.nextToken() // read over label
			if a.lookahead.Type == scanner.Colon {
				a.nextToken() // read over colon
			}
			t = a.lookahead
		} else {
			// Not at the beginning of the line, needs to be followed by a colon
			oldLookahead := a.lookahead
			a.nextToken()
			if a.lookahead.Type == scanner.Colon {
				labelPos = t.Pos
				label = t.StrVal
				a.nextToken() // read over colon
				t = a.lookahead
			} else {
				a.pushToken()
				a.lookahead = oldLookahead
			}
		}
	}
	return t, labelPos, label
}

func (a *Assembler) processLine() {
	t, labelPos, label := a.maybeLabel()

	errs := len(a.Errors())
	if _, found := conditionalTokens[t.Type]; found {
		a.maybeAddLabel(labelPos, label)
		switch t.Type {
		case scanner.Ifdef, scanner.Ifndef:
			negate := t.Type == scanner.Ifndef
			a.nextToken()
			s := a.lookahead.StrVal
			a.match(scanner.Ident)
			_, found := a.symbols.get(s)
			if negate {
				found = !found
			}
			a.assemblyEnabled.push(a.assemblyEnabled.top() && found)

		case scanner.If:
			a.nextToken()
			p := a.lookahead.Pos
			e := a.expr(2)
			if containsKey(relOpToBinOp, a.lookahead.Type) {
				binOp := relOpToBinOp[a.lookahead.Type]
				a.nextToken()
				e2 := a.expr(2)
				e = expr.NewBinaryOp(e, e2, binOp)
			}
			if !e.IsResolved() {
				a.AddError(p, "expression is not resolved")
				e = expr.NewConst(p, 1, 1)
			}
			a.assemblyEnabled.push(a.assemblyEnabled.top() && (e.Eval() != 0))

		case scanner.Else:
			a.nextToken()
			if a.assemblyEnabled.len() == 1 {
				a.AddError(t.Pos, ".else without .if/.ifdef/.ifndef")
				return
			}
			v := a.assemblyEnabled.pop()
			a.assemblyEnabled.push(a.assemblyEnabled.top() && !v)

		case scanner.Endif:
			a.nextToken()
			if a.assemblyEnabled.len() == 1 {
				a.AddError(t.Pos, ".endif without .if/.ifdef/.ifndef")
				return
			}
			a.assemblyEnabled.pop()
		}
	} else {
		if !a.assemblyEnabled.top() {
			// conditionally assembly is turned off, ignore this liune
			return
		}
		switch a.state {
		case stateAssemble:
			a.assembleLine(t, labelPos, label)
		case stateRecordMacro:
			a.recordMacro()
		}
	}
	if len(a.Errors()) <= errs {
		// Only match EOL if there were no errors reported.
		a.matchEol()
	}
}

func (a *Assembler) matchEol() {
	if a.lookahead.Type != scanner.Semicolon && a.lookahead.Type != scanner.Eol {
		a.AddError(a.lookahead.Pos, "';' or EOL expected")
	}
}

func (a *Assembler) maybeAddLabel(labelPos text.Pos, label string) {
	if label == "" {
		return
	}
	a.addLabel(labelPos, label)
}

func (a *Assembler) assembleLine(t scanner.Token, labelPos text.Pos, label string) {
	if t.Type == scanner.Semicolon || t.Type == scanner.Eol {
		// Empty line. Add a label if necessary, and bail out.
		a.maybeAddLabel(labelPos, label)
		return
	}

	// Label checks
	switch t.Type {
	case scanner.Equ, scanner.Macro:
		// Label will be treated as name
		if label == "" {
			a.AddError(labelPos, "Label is necessary")
		}
	case scanner.Org:
		// Can't have a label
		if label != "" {
			a.AddError(labelPos, "Label is not allowed")
		}
	default:
		// In all other cases, add a label
		a.maybeAddLabel(labelPos, label)
	}

	switch t.Type {
	case scanner.Incbin:
		a.nextToken()
		p := a.lookahead.Pos
		filename := a.lookahead.StrVal
		a.match(scanner.String)
		f := a.findIncludeFile(filename)
		if f == nil {
			a.AddError(p, "Can't find file %q in include paths.", filename)
			break
		}
		data, err := ioutil.ReadFile(*f)
		if err != nil {
			a.AddError(p, "Can't read file %q: %s", *f, err)
		}
		for _, b := range data {
			a.emit(expr.NewConst(p, int(b), 1))
		}
	case scanner.Byte:
		a.nextToken()
		// handle byte consts
		nodes := a.dbOp()
		for a.lookahead.Type == scanner.Comma {
			a.nextToken()
			n2 := a.dbOp()
			nodes = append(nodes, n2...)
		}
		a.emit(nodes...)
	case scanner.Reserve:
		a.nextToken()
		// handle byte const
		pos := a.lookahead.Pos
		valNode := expr.NewConst(pos, 0, 1)
		sizeNode := a.expr(2)
		if !sizeNode.IsResolved() {
			a.AddError(pos, "Expression is unresolved")
			sizeNode = expr.NewConst(pos, 1, 2)
		}
		for a.lookahead.Type == scanner.Comma {
			a.nextToken()
			pos = a.lookahead.Pos
			vals := a.dbOp()
			if len(vals) > 1 {
				a.AddError(pos, "Strings not allowed.")
			}
			valNode = vals[0]
		}
		for i := 0; i < sizeNode.Eval(); i++ {
			a.emit(valNode)
		}
	case scanner.Word:
		a.nextToken()
		// handle wird const
		nodes := []expr.Node{a.expr(2)}
		for a.lookahead.Type == scanner.Comma {
			a.nextToken()
			n2 := a.expr(2)
			nodes = append(nodes, n2)
		}
		a.emit(nodes...)
	case scanner.Org:
		a.nextToken()
		// set origin
		orgNode := a.expr(2)
		org := 0
		if orgNode.IsResolved() {
			org = orgNode.Eval()
		} else {
			a.AddError(t.Pos, "Can't use forward declarations in .org")
			org = 0
		}
		if a.section != nil {
			max := a.section.PC()
			if org < max {
				a.AddError(t.Pos, "New origin %d is lower than current pc %d", org, max)
				org = max
			}
			toAdd := org - max
			for toAdd > 0 {
				a.section.Emit(0)
				toAdd = toAdd - 1
			}
		} else {
			a.section = NewSection(org, a)
		}
	case scanner.Equ:
		a.nextToken()
		// label is equ name!
		pos := t.Pos
		val := a.expr(2)
		err := a.addSymbol(label, symbolConst, val)
		if err != nil {
			a.AddError(pos, err.Error())
		}
	case scanner.Fail:
		a.nextToken()
		s := a.lookahead.StrVal
		a.match(scanner.String)
		a.AddError(t.Pos, s)
	case scanner.Macro:
		a.nextToken()
		// label is macroname!
		macroName := label
		a.macro = &macro{
			pos:  t.Pos,
			text: &text.Text{},
		}
		if _, found := Mnemonics[macroName]; found {
			a.AddError(labelPos, "Can't use mnemonic %q as macro name", macroName)
		}
		if err := a.symbols.add(symbol{name: macroName, kind: symbolMacro, m: a.macro}); err != nil {
			a.AddError(labelPos, "%q is already defined", macroName)
		}
		if a.lookahead.Type != scanner.Eol {
			a.macroParam()
			for a.lookahead.Type == scanner.Comma {
				a.nextToken()
				a.macroParam()
			}
		}
		a.state = stateRecordMacro
	case scanner.Endm:
		a.AddError(t.Pos, ".mend without .macro")
	case scanner.Ident:
		op := t.StrVal
		a.nextToken()
		if sym, found := a.symbols.get(op); found {
			if sym.kind != symbolMacro {
				a.AddError(t.Pos, "%q is not a macro", op)
				return
			}
			a.handleMacroInstantiation(sym.m, t.Pos)
		} else {
			// must be a mnemonic
			a.handleMnemonic(t)
		}
	default:
		a.AddError(t.Pos, "Identifier or directive expected")
		return
	}
}

func (a *Assembler) macroParam() {
	paramName := a.lookahead.StrVal
	paramPos := a.lookahead.Pos
	a.match(scanner.Ident)
	if err := a.macro.addParam(paramName); err != nil {
		a.AddError(paramPos, "Parameter %s is already used", paramName)
	}
}

type macroInvocation struct {
	callPos text.Pos
}

func (i *macroInvocation) Modify(err errors.Error) errors.Error {
	msg := err.Msg + fmt.Sprintf(" (called from %s, line %d)", i.callPos.Filename, i.callPos.Line)
	return errors.Error{err.Pos, msg}
}

func (a *Assembler) handleMacroInstantiation(m *macro, callPos text.Pos) {
	// Read actual params
	paramStart := a.lookahead.Pos
	var actParams []string
	if a.lookahead.Type != scanner.Semicolon && a.lookahead.Type != scanner.Eol {
		actParams = append(actParams, a.actMacroParam())
		for a.lookahead.Type == scanner.Comma {
			a.nextToken()
			actParams = append(actParams, a.actMacroParam())
		}
	}

	if len(actParams) != len(m.params) {
		a.AddError(paramStart, "Wrong number of arguments: %d expected, %d found", len(m.params), len(actParams))
		return
	}

	// Get copy of macro with parameters substituted
	t := text.Text{m.replaceParams(actParams)}

	savedErrorModifier := a.errorModifier
	a.errorModifier = &macroInvocation{callPos: callPos}
	a.assembleText(t)
	a.errorModifier = savedErrorModifier
}

func (a *Assembler) handleMnemonic(t scanner.Token) {
	op := strings.ToLower(t.StrVal)

	// must be a mnemonic
	opCodes, found := Mnemonics[op]
	if !found {
		a.AddError(t.Pos, fmt.Sprintf("%s is not a valid mnemonic", t.StrVal))
		return
	}
	param := a.param()
	opCode, found := opCodes[param.mode]
	if !found && param.mode == AM_Absolute {
		// Maybe it's a relative branch? let's check
		opCode, found = opCodes[AM_Relative]
		if found {
			// Yes, it is! Switch to relative addressing
			param.mode = AM_Relative
			param.val.MarkRelative()
		}
	}
	if !found {
		a.AddError(t.Pos, "Invalid parameter.")
	}

	// TODO(asigner): Add warning for JMP ($xxFF)
	a.emit(expr.NewConst(t.Pos, int(opCode), 1))
	if param.val != nil {
		a.emit(param.val)
	}
}

func (a *Assembler) recordMacro() {
	t, labelPos, label := a.maybeLabel()

	switch t.Type {
	case scanner.Macro:
		a.AddError(t.Pos, "Nested macros are not allowed")
	case scanner.Endm:
		// End of macro
		a.nextToken() // Read over ".endm"
		if label != "" {
			a.AddError(labelPos, "Labels not allowed for .endm")
		}
		a.state = stateAssemble
	default:
		// Just another macro line, add it to the current macro
		a.macro.text.Lines = append(a.macro.text.Lines, *a.scanner.Line())

		// Scan until we're at EOL to keep processLine() happy
		for a.lookahead.Type != scanner.Eol {
			a.nextToken()
		}
	}
}

func (a *Assembler) findIncludeFile(f string) *string {
	for _, path := range a.includePaths {
		fullFile := filepath.Join(path, f)
		if _, err := os.Stat(fullFile); err == nil {
			return &fullFile
		}
	}
	return nil
}

func (a *Assembler) actMacroParam() string {
	// actmacroparam := ["#" ["<"|">"]] expr .

	startPos := a.lookahead.Pos
	if a.lookahead.Type == scanner.Hash {
		a.nextToken()
		if a.lookahead.Type == scanner.Lt || a.lookahead.Type == scanner.Gt {
			a.nextToken()
		}
	}
	a.expr(2)
	endPos := a.lookahead.Pos
	return strings.TrimSpace(a.scanner.Line().Extract(startPos, endPos))
}

func (a *Assembler) param() param {
	// param := "#" ["<"|">"] expr
	//       | expr
	//       | expr "," "X"
	//       | expr "," "Y"
	//       | "(" expr ")"
	//       | "(" expr "," "X" ")"
	//       | "(" expr "," "Y" ")"
	//       | "(" expr ") ""," "X"
	//       | "(" expr ")" "," "Y"

	startPos := a.lookahead.Pos

	if a.lookahead.Type == scanner.Semicolon || a.lookahead.Type == scanner.Eol {
		// No param, implied addressing mode'
		return param{rawText: "", mode: AM_Implied}
	}

	switch a.lookahead.Type {
	case scanner.Hash:
		am := AM_Immediate
		var node expr.Node
		a.nextToken()
		p := a.lookahead.Pos
		switch a.lookahead.Type {
		case scanner.Lt:
			a.nextToken()
			node = expr.NewUnaryOp(p, a.expr(2), expr.LoByte)
		case scanner.Gt:
			a.nextToken()
			node = expr.NewUnaryOp(p, a.expr(2), expr.HiByte)
		default:
			node = a.expr(1)
		}
		endPos := a.lookahead.Pos
		return param{rawText: a.scanner.Line().Extract(startPos, endPos), mode: am, val: node}
	case scanner.LParen:
		// AM_AbsoluteIndirect // ($aaaa)
		// AM_IndexedIndirect  // ($aa,X)
		// AM_IndirectIndexed  // ($aa),Y
		a.nextToken()
		node := a.expr(2)
		am := AM_AbsoluteIndirect

		if a.lookahead.Type == scanner.Comma {
			// AM_IndexedIndirect  // ($aa,X)
			a.nextToken()
			if node.ResultSize() > 1 {
				// Let see if we can enforce size
				if !node.ForceSize(1) {
					a.AddError(a.lookahead.Pos, "Address $%x is too large, only 8 bits allowed", node.Eval())
				}
			} else {
				// We can't, so complain
				a.AddError(a.lookahead.Pos, "Address $%x is too large, only 8 bits allowed", node.Eval())
			}
			reg := a.lookahead.StrVal
			pos := a.lookahead.Pos
			a.match(scanner.Ident)
			if strings.ToLower(reg) != "x" {
				a.AddError(pos, "Register X expected, found %s.", reg)
			}
			am = AM_IndexedIndirect
			a.match(scanner.RParen)
			endPos := a.lookahead.Pos
			return param{rawText: a.scanner.Line().Extract(startPos, endPos), mode: am, val: node}
		} else {
			a.match(scanner.RParen)
			if a.lookahead.Type == scanner.Comma {
				// AM_IndirectIndexed  // ($aa),Y
				a.nextToken()
				if node.ResultSize() > 1 {
					// Let see if we can enforce size
					if !node.ForceSize(1) {
						a.AddError(a.lookahead.Pos, "Address $%x is too large, only 8 bits allowed", node.Eval())
					}
				} else {
					// We can't, so complain
					a.AddError(a.lookahead.Pos, "Address $%x is too large, only 8 bits allowed", node.Eval())
				}
				reg := a.lookahead.StrVal
				pos := a.lookahead.Pos
				a.match(scanner.Ident)
				if strings.ToLower(reg) != "y" {
					a.AddError(pos, "Register Y expected, found %s.", reg)
				}
				am = AM_IndirectIndexed
			}
			endPos := a.lookahead.Pos
			return param{rawText: a.scanner.Line().Extract(startPos, endPos), mode: am, val: node}
		}

	default:
		if a.lookahead.Type == scanner.Ident && strings.ToLower(a.lookahead.StrVal) == "a" {
			a.nextToken()
			endPos := a.lookahead.Pos
			return param{rawText: a.scanner.Line().Extract(startPos, endPos), mode: AM_Accumulator, val: nil}
		}
		am := AM_Absolute
		node := a.expr(2)
		am = am.withSize(node.ResultSize())
		if a.lookahead.Type == scanner.Comma {
			a.nextToken()
			s := a.lookahead.StrVal
			pos := a.lookahead.Pos
			a.match(scanner.Ident)
			if strings.ToLower(s) != "x" && strings.ToLower(s) != "y" {
				a.AddError(pos, "Expected 'X' or 'Y', but got %s.", s)
				s = "x"
			}
			am = am.withIndex(s)
		}
		endPos := a.lookahead.Pos
		return param{rawText: a.scanner.Line().Extract(startPos, endPos), mode: am, val: node}
	}
}

func (a *Assembler) dbOp() []expr.Node {
	p := a.lookahead.Pos
	switch {
	case a.lookahead.Type == scanner.Lt:
		a.nextToken()
		n := a.expr(1)
		return []expr.Node{expr.NewUnaryOp(p, n, expr.LoByte)}
	case a.lookahead.Type == scanner.Gt:
		a.nextToken()
		n := a.expr(1)
		return []expr.Node{expr.NewUnaryOp(p, n, expr.HiByte)}
	case a.lookahead.Type == scanner.Ident && strings.ToLower(a.lookahead.StrVal) == "scr":
		// "screen" "(" basicDbOp { "," basicDbOp } ")"
		a.nextToken()
		a.match(scanner.LParen)
		nodes := wrapWithUnaryOp(a.basicDbOp(), expr.ScreenCode)
		for a.lookahead.Type == scanner.Comma {
			a.nextToken()
			n2 := wrapWithUnaryOp(a.basicDbOp(), expr.ScreenCode)
			nodes = append(nodes, n2...)
		}
		a.match(scanner.RParen)
		return nodes
	default:
		return a.basicDbOp()
	}
}

func wrapWithUnaryOp(nodes []expr.Node, op expr.UnaryOp) []expr.Node {
	var newNodes []expr.Node
	for _, n := range nodes {
		newNodes = append(newNodes, expr.NewUnaryOp(n.Pos(), n, op))
	}
	return newNodes
}

func (a *Assembler) basicDbOp() []expr.Node {
	p := a.lookahead.Pos
	switch a.lookahead.Type {
	case scanner.String:
		str := a.lookahead.StrVal
		a.nextToken()
		var res []expr.Node
		for _, c := range str {
			n := expr.NewConst(p, int(c), 1)
			res = append(res, expr.NewUnaryOp(p, n, expr.AsciiToPetscii))
		}
		return res
	default:
		return []expr.Node{a.expr(1)}
	}
}

func containsKey(m map[scanner.TokenType]expr.BinaryOp, key scanner.TokenType) bool {
	_, found := m[key]
	return found
}

func (a *Assembler) expr(size int) expr.Node {
	// expr := ["-"] term { "+"|"-"|"|" term } .
	neg := false
	var negPos text.Pos
	if a.lookahead.Type == scanner.Minus {
		neg = true
		negPos = a.lookahead.Pos
		a.nextToken()
	}
	node := a.term(size)
	if neg {
		node = expr.NewUnaryOp(negPos, node, expr.Neg)
	}

	ops := map[scanner.TokenType]expr.BinaryOp{
		scanner.Plus:  expr.Add,
		scanner.Minus: expr.Sub,
		scanner.Bar:   expr.Or,
	}

	for containsKey(ops, a.lookahead.Type) {
		op := ops[a.lookahead.Type]
		a.nextToken()
		n2 := a.term(size)
		node = expr.NewBinaryOp(node, n2, op)
	}
	return node
}

func (a *Assembler) term(size int) expr.Node {
	// term := factor { "*"|"/"|"%"|"&"|"^" factor } .
	ops := map[scanner.TokenType]expr.BinaryOp{
		scanner.Asterisk:  expr.Mul,
		scanner.Slash:     expr.Div,
		scanner.Percent:   expr.Mod,
		scanner.Ampersand: expr.And,
		scanner.Caret:     expr.Xor,
	}

	node := a.factor(size)
	for containsKey(ops, a.lookahead.Type) {
		op := ops[a.lookahead.Type]
		a.nextToken()
		n2 := a.factor(size)
		node = expr.NewBinaryOp(node, n2, op)
	}
	return node
}

func (a *Assembler) factor(size int) expr.Node {
	// factor := "~" factor | number | ident | '*'.
	var node expr.Node
	switch a.lookahead.Type {
	case scanner.Tilde:
		p := a.lookahead.Pos
		a.nextToken()
		node = a.factor(size)
		node = expr.NewUnaryOp(p, node, expr.Not)
	case scanner.Number:
		p := a.lookahead.Pos
		val := a.lookahead.IntVal
		if !checkSize(size, int(val)) {
			a.AddError(p, "Constant $%x (decimal %d) is wider than %d bits", val, val, size*8)
			break
		}
		node = expr.NewConst(p, int(val), size)
		a.nextToken()
	case scanner.Ident:
		p := a.lookahead.Pos
		sym := a.lookahead.StrVal
		node = nil
		if s, found := a.symbols.get(sym); found {
			if s.val.IsResolved() {
				node = expr.NewSymbolRef(p, sym, size, s.val.Eval())
			}
		}
		if node == nil {
			node = expr.NewUnresolvedSymbol(p, sym, size)
		}
		a.nextToken()
	case scanner.LParen:
		a.nextToken()
		node = a.expr(size)
		a.match(scanner.RParen)
	case scanner.Asterisk:
		p := a.lookahead.Pos
		if size < 2 {
			a.AddError(a.lookahead.Pos, "Current PC is 16 bits wide, expected is a %d bit wide value", size*8)
			break
		}
		node = expr.NewConst(p, a.section.PC(), size)
		a.nextToken()
	default:
		a.AddError(a.lookahead.Pos, "'~', '*', number or identifier expected, found %s", a.lookahead.Type)
		node = expr.NewConst(a.lookahead.Pos, 0, 1)
	}
	return node
}

func checkSize(maxSize int, val int) bool {
	uv := uint64(val)
	uv = uv >> (maxSize * 8)
	return uv == 0
}

func (a *Assembler) emit(nodes ...expr.Node) {
	for _, n := range nodes {
		a.emitNode(n)
	}
}

func (a *Assembler) emitNode(n expr.Node) {
	if a.section == nil {
		a.AddError(a.scanner.LineStart(), "No .org specified")
		a.section = NewSection(0, a)
	}
	var val, size int
	if !n.IsResolved() {
		// register a patch, and emit 0 bytes
		a.registerPatch(a.section.PC(), n)
		val = 0
		size = n.ResultSize()
	} else {
		val = n.Eval()
		size = n.ResultSize()
		if n.IsRelative() {
			val = val - (a.section.PC() + 1)
			if val < -128 || val > 127 {
				a.AddError(n.Pos(), "Branch target too far away.")
			}
		}
	}
	for size > 0 {
		a.section.Emit(byte(val & 0xff))
		val = val >> 8
		size = size - 1
	}
}

func (a *Assembler) registerPatch(pc int, n expr.Node) {
	for label := range n.UnresolvedSymbols() {
		patches := a.patchesPerLabel[label]
		patches = append(patches, patch{pc: pc, node: n})
		a.patchesPerLabel[label] = patches
	}
}

func (a *Assembler) addLabel(pos text.Pos, label string) {
	pc := a.section.PC()
	err := a.addSymbol(label, symbolLabel, expr.NewConst(pos, pc, 2))
	if err != nil {
		a.AddError(pos, err.Error())
		return
	}

	if !isLocalLabel(label) {
		a.reportUnresolvedSymbols(pos, isLocalLabel)
		a.clearLocalLabels()
	}
}

func (a *Assembler) clearLocalLabels() {
	a.symbols.remove(func(sym *symbol) bool { return isLocalLabel(sym.name) })
	for label := range a.patchesPerLabel {
		if isLocalLabel(label) {
			delete(a.patchesPerLabel, label)
		}
	}
}

func (a *Assembler) reportUnresolvedSymbols(errorPos text.Pos, filterFunc func(string) bool) {
	seen := make(map[string]bool)
	for _, symbol := range a.symbols.symbols() {
		if !filterFunc(symbol.name) {
			continue
		}
		if symbol.kind == symbolMacro {
			continue
		}
		if !symbol.val.IsResolved() {
			syms := symbol.val.UnresolvedSymbols()
			if len(syms) > 0 {
				var symnames []string
				for s := range syms {
					symnames = append(symnames, s)
				}
				a.AddError(errorPos, "Undefined symbols in definition of %s: %s", symbol.name, strings.Join(symnames, ", "))
				seen[symbol.name] = true
			} else {
				a.AddError(errorPos, "Undefined label %q", symbol.name)
				seen[symbol.name] = true
			}
		}
	}
	for label := range a.patchesPerLabel {
		if !filterFunc(label) {
			continue
		}
		if !seen[label] {
			a.AddError(errorPos, "Undefined label %q", label)
			seen[label] = true
		}
	}
}

func isLocalLabel(label string) bool {
	return strings.HasPrefix(label, "_")
}

func (a *Assembler) addSymbol(name string, kind symbolKind, val expr.Node) error {
	err := a.symbols.add(symbol{name: name, val: val, kind: kind})
	if err != nil {
		return err
	}
	if val.IsResolved() {
		a.resolveDependencies(name, val)
	}
	return nil
}

func (a *Assembler) resolveDependencies(symbol string, val expr.Node) {
	// Try to resolve as many patches as we can
	patches := a.patchesPerLabel[symbol]
	adjustedPatches := []patch{}
	for _, p := range patches {
		p.node.Resolve(symbol, val.Eval())
		if p.node.IsResolved() {
			a.section.applyPatch(p)
		} else {
			adjustedPatches = append(adjustedPatches, p)
		}
	}
	if len(adjustedPatches) == 0 {
		delete(a.patchesPerLabel, symbol)
	} else {
		a.patchesPerLabel[symbol] = adjustedPatches
	}

	// Now, resolve any symbols
	for _, sym := range a.symbols.symbols() {
		if sym.kind == symbolMacro {
			continue
		}
		if sym.val.IsResolved() {
			continue
		}
		sym.val.Resolve(symbol, val.Eval())
		if sym.val.IsResolved() {
			a.resolveDependencies(sym.name, sym.val)
		}
	}
}

func (a *Assembler) AddError(pos text.Pos, message string, args ...interface{}) {
	err := errors.Error{pos, fmt.Sprintf(message, args...)}
	if a.errorModifier != nil {
		err = a.errorModifier.Modify(err)
	}
	a.errors = append(a.errors, err)
}

func (a *Assembler) Errors() []errors.Error {
	return a.errors
}

func (a *Assembler) AddWarning(pos text.Pos, message string) {
	a.warnings = append(a.errors, errors.Error{pos, message})
}

func (a *Assembler) Warnings() []errors.Error {
	return a.warnings
}

func (a *Assembler) Labels() map[string]int {
	res := make(map[string]int)
	for _, sym := range a.symbols.symbols() {
		if sym.kind == symbolLabel {
			res[sym.name] = sym.val.Eval()
		}
	}
	return res
}

func (a *Assembler) match(t scanner.TokenType) {
	if a.lookahead.Type != t {
		a.AddError(a.lookahead.Pos, "Expected %s, but found %s", t, a.lookahead.Type)
	}
	a.nextToken()
}

func (a *Assembler) nextToken() {
	if a.tokenBufSet {
		a.lookahead = a.tokenBuf
		a.tokenBufSet = false
		return
	}
	a.lookahead = a.scanner.Scan()
}

func (a *Assembler) pushToken() {
	a.tokenBuf = a.lookahead
	a.tokenBufSet = true
}
