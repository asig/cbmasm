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

wait_line .macro line
        lda #line
_l1     cmp $d012
        bne _l1
        .endm

        .include "vic.i"
        .include "startup.i"

        .if PLATFORM = "c128"
        lda #$3e  ; All RAM except I/O range at $d000
        sta $ff00
COL1    .equ COL_LBLUE
COL2    .equ COL_LGREEN
        .else
COL2    .equ COL_LBLUE
COL1    .equ COL_LGREEN
        .endif

        jsr songCopy
        jsr songInit

loop
        wait_line 70
        lda #COL1
        sta $d020
        jsr songPlay
        lda #COL2
        sta $d020

        jmp loop

        ; Use tools/sidconv to convert a sid file to asm source
        .include "Empty_512_bytes_4000.i"


