; Copyright (c) 2023 Andreas Signer <asigner@gmail.com>
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
;

; CP/M program that prints "Hello, world!"

BDOS    .equ    $0005

PRINT_STRING   .equ    $09

    .cpu "z80"
    .platform "c128"
    .encoding "ascii"
    .output "plain"

    .org $100

    LD  DE, msg
    LD  C, PRINT_STRING
    CALL    BDOS
    RET

msg .byte "Hello, world!\n\r",'$'
