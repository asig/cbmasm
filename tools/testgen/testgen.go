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
package main

import (
	"fmt"
	"strings"
)

var rawZ80 = `
ce 56		   ADC A,$56 
89			   ADC A,C 
8e			   ADC A,(HL) 
dd 8e 12		   ADC A,(IX+$12) 
fd 8e ee		   ADC A,(IY-$12) 
ed 7a		   ADC HL,SP 
c6 56		   ADD A,$56 
81			   ADD A,C 
86			   ADD A,(HL) 
dd 86 12		   ADD A,(IX+$12) 
fd 86 ee		   ADD A,(IY-$12) 
39			   ADD HL,SP 
dd 19		   ADD IX,DE 
fd 19		   ADD IY,DE 
e6 56		   AND $56 
a1			   AND C 
a6			   AND (HL) 
dd a6 12		   AND (IX+$12) 
fd a6 ee		   AND (IY-$12) 
cb 46		   BIT 0,(HL) 
dd cb 12 4e	   BIT 1,(IX+$12) 
fd cb ee 56	   BIT 2,(IY-$12) 
cb 59		   BIT 3,C 
cd 78 56		   CALL $5678 
c4 78 56		   CALL NZ,$5678 
3f			   CCF    
fe 56		   CP $56 
b9			   CP C 
ed a9		   CPD 
ed b9		   CPDR 
be			   CP (HL) 
ed a1		   CPI 
ed b1		   CPIR 
2f			   CPL 
27			   DAA 
0d			   DEC C 
1b			   DEC DE 
35			   DEC (HL) 
dd 2b		   DEC IX 
dd 35 12		   DEC (IX+$12) 
fd 2b		   DEC IY 
fd 35 ee		   DEC (IY-$12) 
f3			   DI 
fb			   EI 
08			   EX AF, AF' 
eb			   EX DE, HL 
e3			   EX (SP), HL 
dd e3		   EX (SP), IX 
fd e3		   EX (SP), IY 
d9			   EXX 
76			   HALT 
ed 46		   IM 0 
ed 56		   IM 1 
ed 5e		   IM 2 
db 78		   IN A,($78) 
ed 48		   IN C,(C) 
0c			   INC C 
13			   INC DE 
34			   INC (HL) 
dd 23		   INC IX 
dd 34 12		   INC (IX+$12) 
fd 23		   INC IY 
fd 34 ee		   INC (IY-$12) 
ed aa		   IND 
ed ba		   INDR 
ed a2		   INI 
ed b2		   INIR 
c3 78 56		   JP $5678 
e9			   JP (HL) 
dd e9		   JP (IX) 
fd e9		   JP (IY) 
d2 78 56		   JP NC,$5678 
32 78 56		   LD ($5678), A 
ed 43 78 56	   LD ($5678), BC 
22 78 56		   LD ($5678), HL 
dd 22 78 56	   LD ($5678), IX 
fd 22 78 56	   LD ($5678), IY 
3a 78 56		   LD A, ($5678) 
0a			   LD A, (BC) 
79			   LD A,C 
1a			   LD A, (DE) 
ed 57		   LD A, I 
41			   LD B,C 
01 78 56		   LD BC, $5678 
02			   LD (BC),A 
0e 56		   LD C, $56 
4e			   LD C, (HL) 
dd 4e 12		   LD C, (IX+$12) 
fd 4e ee		   LD C, (IY-$12) 
ed a8		   LDD 
ed 5b 78 56	   LD DE, ($5678) 
12			   LD (DE),A 
ed b8		   LDDR 
36 56		   LD (HL), $56 
2a 78 56		   LD HL, ($5678) 
71			   LD (HL),C 
ed a0		   LDI 
ed 47		   LD I, A 
ed b0		   LDIR 
dd 36 12 56	   LD (IX+$12), $56 
dd 71 12		   LD (IX+$12),C 
dd 2a 78 56	   LD IX, ($5678) 
dd 21 78 56	   LD IX, $5678 
fd 36 ee 56	   LD (IY-$12), $56 
fd 71 ee		   LD (IY-$12),C 
fd 2a 78 56	   LD IY, ($5678) 
fd 21 78 56	   LD IY, $5678 
ed 4f		   LD R, A 
f9			   LD SP, HL 
dd f9		   LD SP, IX 
fd f9		   LD SP, IY 
ed 44		   NEG 
00			   NOP 
f6 56		   OR $56 
b1			   OR C 
b6			   OR (HL) 
dd b6 12		   OR (IX+$12) 
dd b6 12		   OR (IX+$12) 
fd b6 ee		   OR (IY-$12) 
fd b6 ee		   OR (IY-$12) 
ed bb		   OTDR 
ed b3		   OTIR 
d3 17		   OUT (23), A 
ed 51		   OUT (C),D 
ed ab		   OUTD 
ed a3		   OUTI 
f1			   POP AF 
dd e1		   POP IX 
fd e1		   POP IY 
f5			   PUSH AF 
dd e5		   PUSH IX 
fd e5		   PUSH IY 
fd cb ee a6	   RES 4,(IY-$12) 
dd cb 12 ae	   RES 5,(IX+$12) 
cb b6		   RES 6,(HL) 
cb b9		   RES 7,C 
c9			   RET 
ed 4d		   RETI 
ed 45		   RETN 
e8			   RET PE 
17			   RLA 
cb 11		   RL C 
07			   RLCA 
cb 01		   RLC C 
cb 06		   RLC (HL) 
dd cb 12 06	   RLC (IX+$12) 
fd cb ee 06	   RLC (IY-$12) 
ed 6f		   RLD 
cb 16		   RL (HL) 
dd cb 12 16	   RL (IX+$12) 
fd cb ee 16	   RL (IY-$12) 
1f			   RRA 
cb 19		   RR C 
0f			   RRCA 
cb 09		   RRC C 
cb 0e		   RRC (HL) 
dd cb 12 0e	   RRC (IX+$12) 
fd cb ee 0e	   RRC (IY-$12) 
ed 67		   RRD 
cb 1e		   RR (HL) 
dd cb 12 1e	   RR (IX+$12) 
fd cb ee 1e	   RR (IY-$12) 
f7			   RST $30 
de 56		   SBC A,$56 
99			   SBC A,C 
9e			   SBC A,(HL) 
dd 9e 12		   SBC A,(IX+$12) 
fd 9e ee		   SBC A,(IY-$12) 
ed 52		   SBC HL,DE 
37			   SCF 
cb c1		   SET 0,C 
cb ce		   SET 1,(HL) 
dd cb 12 d6	   SET 2,(IX+$12) 
fd cb ee de	   SET 3,(IY-$12) 
cb 21		   SLA C 
cb 26		   SLA (HL) 
dd cb 12 26	   SLA (IX+$12) 
fd cb ee 26	   SLA (IY-$12) 
cb 29		   SRA C 
cb 2e		   SRA (HL) 
dd cb 12 2e	   SRA (IX+$12) 
fd cb ee 2e	   SRA (IY-$12) 
cb 39		   SRL C 
cb 3e		   SRL (HL) 
dd cb 12 3e	   SRL (IX+$12) 
fd cb ee 3e	   SRL (IY-$12) 
d6 56		   SUB $56 
91			   SUB C 
96			   SUB (HL) 
dd 96 12		   SUB (IX+$12) 
fd 96 ee		   SUB (IY-$12) 
ee 56		   XOR $56 
a9			   XOR C 
ae			   XOR (HL) 
dd ae 12		   XOR (IX+$12) 
fd ae ee		   XOR (IY-$12) 
`

var raw6502 = `
ca65 V2.19 - Git e95db43
Main file   : test6502.s
Current file: test6502.s

000000r 1  65 78                ADC $0078
000002r 1  75 78                ADC $0078,X
000004r 1  6D 34 12             ADC $1234
000007r 1  7D 34 12             ADC $1234,X
00000Ar 1  79 34 12             ADC $1234,Y
00000Dr 1  61 9A                ADC ($9A,X)
00000Fr 1  71 BC                ADC ($BC),Y
000011r 1  25 56                AND $0056
000013r 1  25 78                AND $0078
000015r 1  35 78                AND $0078,X
000017r 1  2D 34 12             AND $1234
00001Ar 1  3D 34 12             AND $1234,X
00001Dr 1  39 34 12             AND $1234,Y
000020r 1  21 9A                AND ($9A,X)
000022r 1  31 BC                AND ($BC),Y
000024r 1  06 78                ASL $0078
000026r 1  16 78                ASL $0078,X
000028r 1  0E 34 12             ASL $1234
00002Br 1  1E 34 12             ASL $1234,X
00002Er 1  0A                   ASL A
00002Fr 1  24 78                BIT $0078
000031r 1  2C 34 12             BIT $1234
000034r 1  00                   BRK
000035r 1  18                   CLC
000036r 1  D8                   CLD
000037r 1  58                   CLI
000038r 1  B8                   CLV
000039r 1  C5 56                CMP $0056
00003Br 1  C5 78                CMP $0078
00003Dr 1  D5 78                CMP $0078,X
00003Fr 1  CD 34 12             CMP $1234
000042r 1  DD 34 12             CMP $1234,X
000045r 1  D9 34 12             CMP $1234,Y
000048r 1  C1 9A                CMP ($9A,X)
00004Ar 1  D1 BC                CMP ($BC),Y
00004Cr 1  E4 56                CPX $0056
00004Er 1  E4 78                CPX $0078
000050r 1  EC 34 12             CPX $1234
000053r 1  C4 56                CPY $0056
000055r 1  C4 78                CPY $0078
000057r 1  CC 34 12             CPY $1234
00005Ar 1  C6 78                DEC $0078
00005Cr 1  D6 78                DEC $0078,X
00005Er 1  CE 34 12             DEC $1234
000061r 1  DE 34 12             DEC $1234,X
000064r 1  CA                   DEX
000065r 1  88                   DEY
000066r 1  45 56                EOR $0056
000068r 1  45 78                EOR $0078
00006Ar 1  55 78                EOR $0078,X
00006Cr 1  4D 34 12             EOR $1234
00006Fr 1  5D 34 12             EOR $1234,X
000072r 1  59 34 12             EOR $1234,Y
000075r 1  41 9A                EOR ($9A,X)
000077r 1  51 BC                EOR ($BC),Y
000079r 1  E6 78                INC $0078
00007Br 1  F6 78                INC $0078,X
00007Dr 1  EE 34 12             INC $1234
000080r 1  FE 34 12             INC $1234,X
000083r 1  E8                   INX
000084r 1  C8                   INY
000085r 1  4C 34 12             JMP $1234
000088r 1  6C CD AB             JMP ($ABCD)
00008Br 1  20 34 12             JSR $1234
00008Er 1  A5 56                LDA $0056
000090r 1  A5 78                LDA $0078
000092r 1  B5 78                LDA $0078,X
000094r 1  AD 34 12             LDA $1234
000097r 1  BD 34 12             LDA $1234,X
00009Ar 1  B9 34 12             LDA $1234,Y
00009Dr 1  A1 9A                LDA ($9A,X)
00009Fr 1  B1 BC                LDA ($BC),Y
0000A1r 1  A6 56                LDX $0056
0000A3r 1  A6 78                LDX $0078
0000A5r 1  B6 78                LDX $0078,Y
0000A7r 1  AE 34 12             LDX $1234
0000AAr 1  BE 34 12             LDX $1234,Y
0000ADr 1  A4 56                LDY $0056
0000AFr 1  A4 78                LDY $0078
0000B1r 1  B4 78                LDY $0078,X
0000B3r 1  AC 34 12             LDY $1234
0000B6r 1  BC 34 12             LDY $1234,X
0000B9r 1  46 78                LSR $0078
0000BBr 1  56 78                LSR $0078,X
0000BDr 1  4E 34 12             LSR $1234
0000C0r 1  5E 34 12             LSR $1234,X
0000C3r 1  4A                   LSR A
0000C4r 1  EA                   NOP
0000C5r 1  05 56                ORA $0056
0000C7r 1  05 78                ORA $0078
0000C9r 1  15 78                ORA $0078,X
0000CBr 1  0D 34 12             ORA $1234
0000CEr 1  1D 34 12             ORA $1234,X
0000D1r 1  19 34 12             ORA $1234,Y
0000D4r 1  01 9A                ORA ($9A,X)
0000D6r 1  11 BC                ORA ($BC),Y
0000D8r 1  48                   PHA
0000D9r 1  08                   PHP
0000DAr 1  68                   PLA
0000DBr 1  28                   PLP
0000DCr 1  26 78                ROL $0078
0000DEr 1  36 78                ROL $0078,X
0000E0r 1  2E 34 12             ROL $1234
0000E3r 1  3E 34 12             ROL $1234,X
0000E6r 1  2A                   ROL A
0000E7r 1  66 78                ROR $0078
0000E9r 1  76 78                ROR $0078,X
0000EBr 1  6E 34 12             ROR $1234
0000EEr 1  7E 34 12             ROR $1234,X
0000F1r 1  6A                   ROR A
0000F2r 1  40                   RTI
0000F3r 1  60                   RTS
0000F4r 1  E5 56                SBC $0056
0000F6r 1  E5 78                SBC $0078
0000F8r 1  F5 78                SBC $0078,X
0000FAr 1  ED 34 12             SBC $1234
0000FDr 1  FD 34 12             SBC $1234,X
000100r 1  F9 34 12             SBC $1234,Y
000103r 1  E1 9A                SBC ($9A,X)
000105r 1  F1 BC                SBC ($BC),Y
000107r 1  38                   SEC
000108r 1  F8                   SED
000109r 1  78                   SEI
00010Ar 1  85 78                STA $0078
00010Cr 1  95 78                STA $0078,X
00010Er 1  8D 34 12             STA $1234
000111r 1  9D 34 12             STA $1234,X
000114r 1  99 34 12             STA $1234,Y
000117r 1  81 9A                STA ($9A,X)
000119r 1  91 BC                STA ($BC),Y
00011Br 1  86 78                STX $0078
00011Dr 1  8E 34 12             STX $1234
000120r 1  96 78                STX $78,Y
000122r 1  84 78                STY $0078
000124r 1  8C 34 12             STY $1234
000127r 1  94 78                STY $78,X
000129r 1  AA                   TAX
00012Ar 1  A8                   TAY
00012Br 1  BA                   TSX
00012Cr 1  8A                   TXA
00012Dr 1  9A                   TXS
00012Er 1  98                   TYA
`

func main() {
	lines := strings.Split(rawZ80, "\n")
	for _, l := range lines {
		l := strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		handleLineZ80(strings.TrimSpace(l))
	}

	fmt.Printf("-----------------------------------------------------------------\n")

	lines = strings.Split(raw6502, "\n")
	for _, l := range lines {
		l := strings.TrimSpace(l)
		if !strings.HasPrefix(l, "00") {
			continue
		}
		handleLine6502(strings.TrimSpace(l))
	}
}

func handleLineZ80(line string) {
	parts := strings.Split(line, "\t")
	bytes := parts[0]
	code := strings.TrimSpace(parts[len(parts)-1])

	var wantBytes []string
	for _, b := range strings.Split(bytes, " ") {
		wantBytes = append(wantBytes, fmt.Sprintf("0x%s", b))
	}

	fmt.Printf(
		`		{
			name: "Single instruction %s",
			text: "%s",
			want: []byte{%s},
		},
`, code, code, strings.Join(wantBytes, ", "))
}

func handleLine6502(line string) {
	bytes := strings.TrimSpace(line[11:32])
	code := strings.TrimSpace(line[32:])

	var wantBytes []string
	for _, b := range strings.Split(bytes, " ") {
		wantBytes = append(wantBytes, fmt.Sprintf("0x%s", b))
	}

	fmt.Printf(
		`		{
			name: "Single instruction %s",
			text: "%s",
			want: []byte{%s},
		},
`, code, code, strings.Join(wantBytes, ", "))
}
