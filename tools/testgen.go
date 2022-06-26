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
0000        ADC $0056       6D 56 00
0003        ADC $0078       6D 78 00
0006        ADC $0078,X     7D 78 00
0009        ADC $1234       6D 34 12
000C        ADC $1234,X     7D 34 12
000F        ADC $1234,Y     79 34 12
0012        ADC ($9A,X)     61 9A
0014        ADC ($BC),Y     71 BC
0016        AND $0056       2D 56 00
0019        AND $0078       2D 78 00
001C        AND $0078,X     3D 78 00
001F        AND $1234       2D 34 12
0022        AND $1234,X     3D 34 12
0025        AND $1234,Y     39 34 12
0028        AND ($9A,X)     21 9A
002A        AND ($BC),Y     31 BC
002C        ASL $0078       0E 78 00
002F        ASL $0078,X     1E 78 00
0032        ASL $1234       0E 34 12
0035        ASL $1234,X     1E 34 12
0038        ASL A           0A
0039        BIT $0078       2C 78 00
003C        BIT $1234       2C 34 12
003F        BRK             00
0040        CLC             18
0041        CLD             D8
0042        CLI             58
0043        CLV             B8
0044        CMP $0056       CD 56 00
0047        CMP $0078       CD 78 00
004A        CMP $0078,X     DD 78 00
004D        CMP $1234       CD 34 12
0050        CMP $1234,X     DD 34 12
0053        CMP $1234,Y     D9 34 12
0056        CMP ($9A,X)     C1 9A
0058        CMP ($BC),Y     D1 BC
005A        CPX $0056       EC 56 00
005D        CPX $0078       EC 78 00
0060        CPX $1234       EC 34 12
0063        CPY $0056       CC 56 00
0066        CPY $0078       CC 78 00
0069        CPY $1234       CC 34 12
006C        DEC $0078       CE 78 00
006F        DEC $0078,X     DE 78 00
0072        DEC $1234       CE 34 12
0075        DEC $1234,X     DE 34 12
0078        DEX             CA
0079        DEY             88
007A        EOR $0056       4D 56 00
007D        EOR $0078       4D 78 00
0080        EOR $0078,X     5D 78 00
0083        EOR $1234       4D 34 12
0086        EOR $1234,X     5D 34 12
0089        EOR $1234,Y     59 34 12
008C        EOR ($9A,X)     41 9A
008E        EOR ($BC),Y     51 BC
0090        INC $0078       EE 78 00
0093        INC $0078,X     FE 78 00
0096        INC $1234       EE 34 12
0099        INC $1234,X     FE 34 12
009C        INX             E8
009D        INY             C8
009E        JMP $1234       4C 34 12
00A1        JMP ($ABCD)     6C CD AB
00A4        JSR $1234       20 34 12
00A7        LDA $0056       AD 56 00
00AA        LDA $0078       AD 78 00
00AD        LDA $0078,X     BD 78 00
00B0        LDA $1234       AD 34 12
00B3        LDA $1234,X     BD 34 12
00B6        LDA $1234,Y     B9 34 12
00B9        LDA ($9A,X)     A1 9A
00BB        LDA ($BC),Y     B1 BC
00BD        LDX $0056       AE 56 00
00C0        LDX $0078       AE 78 00
00C3        LDX $0078,Y     BE 78 00
00C6        LDX $1234       AE 34 12
00C9        LDX $1234,Y     BE 34 12
00CC        LDY $0056       AC 56 00
00CF        LDY $0078       AC 78 00
00D2        LDY $0078,X     BC 78 00
00D5        LDY $1234       AC 34 12
00D8        LDY $1234,X     BC 34 12
00DB        LSR $0078       4E 78 00
00DE        LSR $0078,X     5E 78 00
00E1        LSR $1234       4E 34 12
00E4        LSR $1234,X     5E 34 12
00E7        LSR A           4A
00E8        NOP             EA
00E9        ORA $0056       0D 56 00
00EC        ORA $0078       0D 78 00
00EF        ORA $0078,X     1D 78 00
00F2        ORA $1234       0D 34 12
00F5        ORA $1234,X     1D 34 12
00F8        ORA $1234,Y     19 34 12
00FB        ORA ($9A,X)     01 9A
00FD        ORA ($BC),Y     11 BC
00FF        PHA             48
0100        PHP             08
0101        PLA             68
0102        PLP             28
0103        ROL $0078       2E 78 00
0106        ROL $0078,X     3E 78 00
0109        ROL $1234       2E 34 12
010C        ROL $1234,X     3E 34 12
010F        ROL A           2A
0110        ROR $0078       6E 78 00
0113        ROR $0078,X     7E 78 00
0116        ROR $1234       6E 34 12
0119        ROR $1234,X     7E 34 12
011C        ROR A           6A
011D        RTI             40
011E        RTS             60
011F        SBC $0056       ED 56 00
0122        SBC $0078       ED 78 00
0125        SBC $0078,X     FD 78 00
0128        SBC $1234       ED 34 12
012B        SBC $1234,X     FD 34 12
012E        SBC $1234,Y     F9 34 12
0131        SBC ($9A,X)     E1 9A
0133        SBC ($BC),Y     F1 BC
0135        SEC             38
0136        SED             F8
0137        SEI             78
0138        STA $0078       8D 78 00
013B        STA $0078,X     9D 78 00
013E        STA $1234       8D 34 12
0141        STA $1234,X     9D 34 12
0144        STA $1234,Y     99 34 12
0147        STA ($9A,X)     81 9A
0149        STA ($BC),Y     91 BC
014B        STX $0078       8E 78 00
014E        STX $1234       8E 34 12
0151        STX $78,Y       96 78
0153        STY $0078       8C 78 00
0156        STY $1234       8C 34 12
0159        STY $78,X       94 78
015B        TAX             AA
015C        TAY             A8
015D        TSX             BA
015E        TXA             8A
015F        TXS             9A
0160        TYA             98
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
		if len(l) == 0 {
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
	for _,b := range strings.Split(bytes, " ") {
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
	code := strings.TrimSpace(line[12:28])
	bytes := strings.TrimSpace(line[28:])

	var wantBytes []string
	for _,b := range strings.Split(bytes, " ") {
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
