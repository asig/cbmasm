package asm

import (
	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/expr"
	"github.com/asig/cbmasm/pkg/scanner"
	"github.com/asig/cbmasm/pkg/text"
	"io/ioutil"
	"os"
	"path/filepath"

	"fmt"
	"strings"
)

// patch records nodes that can't be evaluated because of undefined nodes
type patch struct {
	pc   int       // Place to patch
	node expr.Node // Node that needs to be patched in
}

type param struct {
	mode AddressingMode
	val  expr.Node
}

type Assembler struct {
	text         text.Text
	includePaths []string

	errors      []errors.Error
	warnings    []errors.Error
	scanner     *scanner.Scanner
	lookahead   scanner.Token
	tokenBuf    scanner.Token
	tokenBufSet bool

	// Code generation buffer
	section *Section

	// outstanding patches
	patchesPerLabel map[string][]patch

	// Symbol table
	symbols map[string]expr.Node
}

func New(t text.Text, includePaths []string) *Assembler {
	a := &Assembler{
		text:         t,
		includePaths: includePaths,
	}
	return a
}

func (a *Assembler) Assemble() {
	a.errors = nil
	a.warnings = nil
	a.section = nil
	a.patchesPerLabel = make(map[string][]patch)
	a.symbols = make(map[string]expr.Node)

	a.assembleText(a.text)

	p := text.Pos{Filename: a.text.Filename, Line: a.text.LastLine().LineNumber, Col: 1}
	a.reportUnresolvedLabels(p, func(string) bool { return true })
}

func (a *Assembler) assembleText(t text.Text) {
	for _, line := range t.Lines {
		a.scanner = scanner.New(t.Filename, line, a)
		a.tokenBufSet = false
		a.lookahead = a.scanner.Scan()

		a.assembleLine(line)
	}
}

func (a *Assembler) Origin() int {
	return a.section.Org()
}

func (a *Assembler) GetBytes() []byte {
	return a.section.bytes
}

func (a *Assembler) assembleLine(line text.Line) {
	var labelPos text.Pos
	label := ""

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

	if t.Type == scanner.Semicolon || t.Type == scanner.Eol {
		// Empty line. Add a label if necessary, and bail out.
		if label != "" {
			a.addLabel(labelPos, label)
		}
		return
	}

	if t.Type != scanner.Ident {
		a.AddError(t.Pos, fmt.Sprintf("expected %s, got %s", scanner.Ident, t.Type))
		return
	}
	op := strings.ToLower(t.StrVal)
	a.nextToken()

	// Label checks
	switch op {
	case ".equ", ".macro":
		// Label will be treated as name
		if label == "" {
			a.AddError(labelPos, fmt.Sprintf("Label is necessary"))
		}
	case ".org":
		// must not have a label
		if label != "" {
			a.AddError(labelPos, fmt.Sprintf("Labels not allowed for .org"))
		}
	default:
		// In all other cases, add a label
		if label != "" {
			a.addLabel(labelPos, label)
		}
	}

	switch op {
	case "incbin":
		p := a.lookahead.Pos
		filename := a.lookahead.StrVal
		a.match(scanner.String)
		f := a.findIncludeFile(filename)
		if f == nil {
			a.AddError(p, fmt.Sprintf("Can't find file %q in include paths.", filename))
			break
		}
		data, err := ioutil.ReadFile(*f)
		if err != nil {
			a.AddError(p, fmt.Sprintf("Can't read file %q: %s", *f, err))
		}
		for _, b := range data {
			a.emit(expr.NewConst(int(b), 1))
		}
	case "include":
		p := a.lookahead.Pos
		filename := a.lookahead.StrVal
		a.match(scanner.String)
		f := a.findIncludeFile(filename)
		if f == nil {
			a.AddError(p, fmt.Sprintf("Can't find file %q in include paths.", filename))
			break
		}
		content, err := ioutil.ReadFile(*f)
		if err != nil {
			a.AddError(p, fmt.Sprintf("Can't read file %q: %s", *f, err))
		}
		a.assembleText(text.Process(filename, string(content)))
	case ".byte":
		// handle byte consts
		nodes := a.dbOp(1)
		for a.lookahead.Type == scanner.Comma {
			a.nextToken()
			n2 := a.dbOp(1)
			nodes = append(nodes, n2...)
		}
		a.emit(nodes...)
	case ".reserve":
		// handle byte const
		pos := a.lookahead.Pos
		valNode := expr.NewConst(0, 1)
		sizeNode := a.expr(2)
		if !sizeNode.IsResolved() {
			a.AddError(pos, "Expression is unresolved")
			sizeNode = expr.NewConst(1, 2)
		}
		for a.lookahead.Type == scanner.Comma {
			a.nextToken()
			pos = a.lookahead.Pos
			vals := a.dbOp(1)
			if len(vals) > 1 {
				a.AddError(pos, "Strings not allowed.")
			}
			valNode = vals[0]
		}
		for i := 0; i < sizeNode.Eval(); i++ {
			a.emit(valNode)
		}
	case ".word":
		// handle wird const
		nodes := a.dbOp(2)
		for a.lookahead.Type == scanner.Comma {
			a.nextToken()
			n2 := a.dbOp(2)
			nodes = append(nodes, n2...)
		}
		a.emit(nodes...)
	case ".org":
		// set origin
		orgNode := a.expr(2)
		org := 0
		if orgNode.IsResolved() {
			org = orgNode.Eval()
		} else {
			a.AddError(t.Pos, fmt.Sprintf("Can't use forward declarations in .org"))
			org = 0
		}
		if a.section != nil {
			max := a.section.PC()
			if org < max {
				a.AddError(t.Pos, fmt.Sprintf("New origin %d is lower than current pc %d", org, max))
				org = max
			}
			toAdd := org - max
			for toAdd > 0 {
				a.section.Emit(0)
				toAdd = toAdd - 1
			}
		} else {
			a.section = NewSection(org)
		}
	case ".equ":
		// label is equ name!
		if _, found := a.symbols[label]; found {
			a.AddError(t.Pos, fmt.Sprintf("Symbol %s already exists.", label))
			return
		}
		val := a.expr(2)
		a.addSymbol(label, val)
	case ".macro":
		// label is macroname!
		// begin macro
	case ".mend":
		// end maccro
	default:
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
		a.emit(expr.NewConst(int(opCode), 1))
		if param.val != nil {
			a.emit(param.val)
		}
	}

	if a.lookahead.Type != scanner.Semicolon && a.lookahead.Type != scanner.Eol {
		a.AddError(a.lookahead.Pos, "';' or EOL expected")
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

	if a.lookahead.Type == scanner.Semicolon || a.lookahead.Type == scanner.Eol {
		// No param, implied addressing mode'
		return param{mode: AM_Implied}
	}

	switch a.lookahead.Type {
	case scanner.Hash:
		am := AM_Immediate
		var node expr.Node
		a.nextToken()
		switch a.lookahead.Type {
		case scanner.Lt:
			a.nextToken()
			node = expr.NewUnaryOp(a.expr(2), expr.LoByte)
		case scanner.Gt:
			a.nextToken()
			node = expr.NewUnaryOp(a.expr(2), expr.HiByte)
		default:
			node = a.expr(1)
		}
		return param{mode: am, val: node}
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
					a.AddError(a.lookahead.Pos, fmt.Sprintf("Address $%x is too large, only 8 bits allowed", node.Eval()))
				}
			} else {
				// We can't, so complain
				a.AddError(a.lookahead.Pos, fmt.Sprintf("Address $%x is too large, only 8 bits allowed", node.Eval()))
			}
			reg := a.lookahead.StrVal
			pos := a.lookahead.Pos
			a.match(scanner.Ident)
			if strings.ToLower(reg) != "x" {
				a.AddError(pos, fmt.Sprintf("Register X expected, found %s.", reg))
			}
			am = AM_IndexedIndirect
			a.match(scanner.RParen)
			return param{mode: am, val: node}
		} else {
			a.match(scanner.RParen)
			if a.lookahead.Type == scanner.Comma {
				// AM_IndirectIndexed  // ($aa),Y
				a.nextToken()
				if node.ResultSize() > 1 {
					// Let see if we can enforce size
					if !node.ForceSize(1) {
						a.AddError(a.lookahead.Pos, fmt.Sprintf("Address $%x is too large, only 8 bits allowed", node.Eval()))
					}
				} else {
					// We can't, so complain
					a.AddError(a.lookahead.Pos, fmt.Sprintf("Address $%x is too large, only 8 bits allowed", node.Eval()))
				}
				reg := a.lookahead.StrVal
				pos := a.lookahead.Pos
				a.match(scanner.Ident)
				if strings.ToLower(reg) != "y" {
					a.AddError(pos, fmt.Sprintf("Register Y expected, found %s.", reg))
				}
				am = AM_IndirectIndexed
			}
			return param{mode: am, val: node}
		}

	default:
		if a.lookahead.Type == scanner.Ident && strings.ToLower(a.lookahead.StrVal) == "a" {
			a.nextToken()
			return param{mode: AM_Accumulator, val: nil}
		}
		am := AM_Absolute
		node := a.expr(2)
		am = am.withSize(node.ResultSize())
		if a.lookahead.Type == scanner.Comma {
			a.nextToken()
			s := a.lookahead.StrVal
			pos := a.lookahead.Pos
			a.match(scanner.String)
			if strings.ToLower(s) != "x" && strings.ToLower(s) != "y" {
				a.AddError(pos, fmt.Sprintf("Expected 'X' or 'Y', but got %s.", s))
				s = "x"
			}
			am = am.withIndex(s)
		}
		return param{mode: am, val: node}
	}
}

func (a *Assembler) dbOp(size int) []expr.Node {
	switch a.lookahead.Type {
	case scanner.Lt:
		a.nextToken()
		n := a.expr(size)
		return []expr.Node{expr.NewUnaryOp(n, expr.LoByte)}
	case scanner.Gt:
		a.nextToken()
		n := a.expr(size)
		return []expr.Node{expr.NewUnaryOp(n, expr.HiByte)}
	case scanner.String:
		str := a.lookahead.StrVal
		a.nextToken()
		var res []expr.Node
		for _, c := range str {
			res = append(res, expr.NewConst(int(c), 1))
		}
		return res
	default:
		return []expr.Node{a.expr(size)}
	}
}

func containsKey(m map[scanner.TokenType]expr.BinaryOp, key scanner.TokenType) bool {
	_, found := m[key]
	return found
}

func (a *Assembler) expr(size int) expr.Node {
	// expr := ["-"] term { "+"|"-"|"|" term } .
	neg := false
	if a.lookahead.Type == scanner.Minus {
		neg = true
		a.nextToken()
	}
	node := a.term(size)
	if neg {
		node = expr.NewUnaryOp(node, expr.Neg)
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
	node := expr.NewConst(0, size)
	switch a.lookahead.Type {
	case scanner.Tilde:
		a.nextToken()
		node = a.factor(size)
		node = expr.NewUnaryOp(node, expr.Not)
	case scanner.Number:
		val := a.lookahead.IntVal
		if !checkSize(size, int(val)) {
			a.AddError(a.lookahead.Pos, fmt.Sprintf("Constant $%x (decimal %d) is wider than %d bits", val, val, size*8))
			break
		}
		node = expr.NewConst(int(val), size)
		a.nextToken()
	case scanner.Ident:
		sym := a.lookahead.StrVal
		node = nil
		if val, found := a.symbols[sym]; found {
			if val.IsResolved() {
				node = expr.NewSymbolRef(sym, size, val.Eval())
			}
		}
		if node == nil {
			node = expr.NewUnresolvedSymbol(sym, size)
		}
		a.nextToken()
	case scanner.LParen:
		a.nextToken()
		node = a.expr(size)
		a.match(scanner.RParen)
	case scanner.Asterisk:
		if size < 2 {
			a.AddError(a.lookahead.Pos, fmt.Sprintf("Current PC is 16 bits wide, expected is a %d bit wide value", size*8))
			break
		}
		node = expr.NewConst(a.section.PC(), size)
		a.nextToken()
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
		a.section = NewSection(0)
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
	if _, found := a.symbols[label]; found {
		a.AddError(pos, fmt.Sprintf("Symbol %q already defined.", label))
		return
	}

	pc := a.section.PC()
	a.addSymbol(label, expr.NewConst(pc, 2))

	if !isLocalLabel(label) {
		a.reportUnresolvedLabels(pos, isLocalLabel)
		a.clearLocalLabels()
	}
}

func (a *Assembler) clearLocalLabels() {
	for symbol, _ := range a.symbols {
		if isLocalLabel(symbol) {
			delete(a.symbols, symbol)
		}
	}
	for label := range a.patchesPerLabel {
		if isLocalLabel(label) {
			delete(a.patchesPerLabel, label)
		}
	}
}

func (a *Assembler) reportUnresolvedLabels(errorPos text.Pos, filterFunc func(string) bool) {
	p := text.Pos{Filename: a.text.Filename, Line: a.text.LastLine().LineNumber, Col: 1}
	seen := make(map[string]bool)
	for symbol, node := range a.symbols {
		if !filterFunc(symbol) {
			continue
		}
		if !node.IsResolved() {
			syms := node.UnresolvedSymbols()
			if len(syms) > 0 {
				var symnames []string
				for s := range syms {
					symnames = append(symnames, s)
				}
				a.AddError(p, fmt.Sprintf("Undefined symbols in definition of %s: %s", symbol, strings.Join(symnames, ", ")))
				seen[symbol] = true
			} else {
				a.AddError(p, fmt.Sprintf("Undefined label %q", symbol))
				seen[symbol] = true
			}
		}
	}
	for label := range a.patchesPerLabel {
		if !filterFunc(label) {
			continue
		}
		if !seen[label] {
			a.AddError(p, fmt.Sprintf("Undefined label %q", label))
			seen[label] = true
		}
	}
}

func isLocalLabel(label string) bool {
	return strings.HasPrefix(label, "_")
}

func (a *Assembler) addSymbol(symbol string, val expr.Node) {
	a.symbols[symbol] = val
	if !val.IsResolved() {
		return
	}
	a.resolveDependencies(symbol, val)
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
	//resolved := map[string]expr.Node {}
	for name, node := range a.symbols {
		if node.IsResolved() {
			continue
		}
		node.Resolve(symbol, val.Eval())
		if node.IsResolved() {
			a.resolveDependencies(name, node)
		}
	}
}

func (a *Assembler) AddError(pos text.Pos, message string) {
	a.errors = append(a.errors, errors.Error{pos, message})
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

func (a *Assembler) match(t scanner.TokenType) {
	if a.lookahead.Type != t {
		a.AddError(a.lookahead.Pos, fmt.Sprintf("Expected %s, but found %s", t, a.lookahead.Type))
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
