        include "vic.i"

ptr1    .equ $fa ; zero page address used for memory indexing
ptr2    .equ $fc ; zero page address used for memory indexing

vram    .equ    $400

last_line .equ vram+24*40

        include "startup.i"

        jsr init_random
        jsr generate_line
        jsr scroll_up
        rts

        ; jsr generate_hard_scroll_up


generate_line:
        ; set addr to beginning of screen
        lda #<last_line
        sta ptr1
        lda #>last_line
        sta ptr1+1
        ldy #39 ; 1 line
_l      lda $D41B ; load random value
        and #1    ; limit to 0,1
        adc #77
_addr   sta (ptr1),y
        dey
        bpl _l
        rts

scroll_up:
        lda #0
        sta $d020
        ; set up dst address
        lda #<vram
        sta ptr1
        lda #>vram
        sta ptr1+1

        ; set up src address
        lda #<(vram+40)
        sta ptr2
        lda #>(vram+40)
        sta ptr2+1

        ldx #4      ; 4 loops, 240 bytes each
_lx     ldy #0
_ly     lda (ptr2),y
        sta (ptr1),y
        iny
        cpy #240
        bne _ly

        ; increment src and dst by 240
        clc
        lda ptr1
        adc #240
        sta ptr1
        bcc _l2
        inc ptr1+1
_l2:
        clc
        lda ptr2
        adc #240
        sta ptr2
        bcc _l3
        inc ptr2+1
_l3:

        dex
        bne _lx
        clc
        lda #1
        sta $d020
        rts

; From https://www.atarimagazines.com/compute/issue72/random_numbers.php

init_random:
    lda #$FF  ; maximum frequency value
    sta $D40E ; voice 3 frequency low byte
    sta $D40F ; voice 3 frequency high byte
    lda #$80  ; noise waveform, gate bit off
    sta $D412 ; voice 3 control register
    rts
