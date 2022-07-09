        include "vic.i"

        include "startup.i"

border_col .equ $d020
background_col .equ $d021
target .equ loop

        ldx #0
        ldy #1
loop    stx border_col
        stx background_col
        inx
        dey
        jmp target

        rts
