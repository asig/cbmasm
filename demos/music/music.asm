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
        .include "songdata.i"

