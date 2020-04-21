package asm

type Section struct {
	org int
	bytes []byte
}

func NewSection(org int) *Section {
	return &Section{org: org}
}

func (section *Section) Emit(b byte) {
	section.bytes = append(section.bytes, b)
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
	val := p.node.Eval()
	if p.node.IsRelative() {
		val = val - (p.pc+1)
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
