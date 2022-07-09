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

        .if PLATFORM = "c128"

        ; BASIC header for Commodore 128
        .org $1c01
        .word _next          ; pointer to next line
        .word 1              ; line number (1)
        .byte $9e, "7181",0  ; SYS 7181
_next   .word 0              ; End of listing

        .else

        .if PLATFORM = "c64"

        ; BASIC header for Commodore 64
        .org $801
        .word _next          ; pointer to next line
        .word 1              ; line number (1)
        .byte $9e, "2061",0  ; SYS 2061
_next   .word 0              ; End of listing

        .else

        .if PLATFORM = "pet"

        ; BASIC header for PET
        .org $401
        .word _next          ; pointer to next line
        .word 1              ; line number (1)
        .byte $9e, "1037",0  ; SYS 1037
_next   .word 0              ; End of listing

        .else

        .fail "Unsupported platform."
        .endif
        .endif
        .endif

