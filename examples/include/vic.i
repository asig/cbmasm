; Copyright (c) 2021 Andreas Signer <asigner@gmail.com>
;
; This file is part of cbmasm.
;
; cbmasm is free software: you can redistribute it and/or
; modify it under the terms of the GNU General Public License as
; published by the Free Software Foundation, either version 3 of the
; License, or (at your option) any later version.
;
; cbmasm is distributed in the hope that it will be useful,
; but WITHOUT ANY WARRANTY; without even the implied warranty of
; MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
; GNU General Public License for more details.
;
; You should have received a copy of the GNU General Public License
; along with cbmasm.  If not, see <http://www.gnu.org/licenses/>.

; Colors
COL_BLACK	.equ $00
COL_WHITE	.equ $01
COL_RED     .equ $02
COL_CYAN	.equ $03
COL_PURPLE	.equ $04
COL_GREEN	.equ $05
COL_BLUE	.equ $06
COL_YELLOW	.equ $07
COL_ORANGE	.equ $08
COL_BROWN	.equ $09
COL_PINK	.equ $0A
COL_DGREY	.equ $0B
COL_GREY	.equ $0C
COL_LGREEN	.equ $0D
COL_LBLUE	.equ $0E
COL_LGREY	.equ $0F

; VIC registers
VIC_BASE .equ $d000
SPRITE_0_X .equ VIC_BASE+0
SPRITE_0_Y .equ VIC_BASE+1
SPRITE_1_X .equ VIC_BASE+2
SPRITE_1_Y .equ VIC_BASE+3
SPRITE_2_X .equ VIC_BASE+4
SPRITE_2_Y .equ VIC_BASE+5
SPRITE_3_X .equ VIC_BASE+6
SPRITE_3_Y .equ VIC_BASE+7
SPRITE_4_X .equ VIC_BASE+8
SPRITE_4_Y .equ VIC_BASE+9
SPRITE_5_X .equ VIC_BASE+10
SPRITE_5_Y .equ VIC_BASE+11
SPRITE_6_X .equ VIC_BASE+12
SPRITE_6_Y .equ VIC_BASE+13
SPRITE_7_X .equ VIC_BASE+14
SPRITE_7_Y .equ VIC_BASE+15
SPRITES_MAX_X .equ VIC_BASE+16
SPRITES_VISIBLE .equ    VIC_BASE+21
SPRITES_DBL_H   .equ    VIC_BASE+23
VIC_MCR         .equ    VIC_BASE+24     ; VIC Memory Control Register
SPRITES_PRIO    .equ    VIC_BASE+27
SPRITES_MULTICOL    .equ    VIC_BASE+28
SPRITES_DBL_W   .equ    VIC_BASE+29
SPRITES_SPR_COLL    .equ    VIC_BASE+30
SPRITES_BKG_COLL    .equ    VIC_BASE+31
SPRITE_MULTICOL_0 .equ VIC_BASE+37
SPRITE_MULTICOL_1 .equ VIC_BASE+38
SPRITE_0_COL .equ VIC_BASE+39
SPRITE_1_COL .equ VIC_BASE+40
SPRITE_2_COL .equ VIC_BASE+41
SPRITE_3_COL .equ VIC_BASE+42
SPRITE_4_COL .equ VIC_BASE+43
SPRITE_5_COL .equ VIC_BASE+44
SPRITE_6_COL .equ VIC_BASE+45
SPRITE_7_COL .equ VIC_BASE+46


