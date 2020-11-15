;PLATFORM    .equ "C64"

PLATFORM_C128    .equ 1
;PLATFORM_C64    .equ 1

        .include "macros.i"
        .include "vic.i"


DEBUG_COL   .macro col
            ;.ifdef DEBUG
            lda col
            sta $d020
            ;.endif
            .endm

ptr1    .equ $fa ; zero page address used for memory indexing
ptr2    .equ $fc ; zero page address used for memory indexing


rasterline0 .equ 30 ; set vscroll
rasterline1 .equ 50 ; Switch to DGREY
rasterline2 .equ 60 ; Switch to GREY
rasterline3 .equ 65 ; Switch to LGREY
rasterline4 .equ 70 ; Switch to WHITE
rasterline5 .equ 50-4+20*8-20; Switch to LGREY
rasterline6 .equ 50-4+20*8-15 ; Switch to GREY
rasterline7 .equ 50-4+20*8-10 ; Switch to DGREY
rasterline8 .equ 50-4+20*8 ; Switch to BLACK and fixed scroll
rasterline9 .equ 50-4+21*8 ; Horizontal scroll

vram    .equ    $400
color_ram   .equ $d800

last_line   .equ vram+19*40
scroll_line  .equ vram+22*40
scroll_line_col .equ color_ram+22*40

        .include "startup.i"
        jmp start

scrollx .byte 0
scrolly .byte 0
wait    .byte 1
start:
        jsr init_random
        jsr clear_screen
        jsr draw_scroller
        jsr install_irq

        lda #0
        sta $d021
        sta $d020


loop:
        ldx #2      ; wait 3 times ...
_l1     lda #1      ; ... for raster line 250
_w1     cmp wait
        beq _w1
        sta wait

        DEBUG_COL #COL_YELLOW
        jsr do_scrolltext
        DEBUG_COL #COL_BLACK

        jsr scroll_colors
        dex
        bne _l1

        ; soft scroll the maze
        ldx scrolly
        dex
        bmi hard_scroll
        stx scrolly
        ;lda #%00010000
        ;ora scrolly
        ;sta $d011
        jmp end_scroll
hard_scroll:
        lda #7
        sta scrolly
        ;ora #%00010000
        ;sta $d011
        jsr scroll_up_fast
        jsr generate_line
end_scroll
        jmp loop


do_scrolltext:
        txa
        pha
        ; scroll text
        ldx scrollx
        dex
        bpl _no_scrolltext_update
        clc
        lda scroll_pos
        adc #1
        sta scroll_pos
        bcc _no_overflow
        inc scroll_pos+1
_no_overflow
        jsr draw_scroller
        ldx #7
_no_scrolltext_update
        stx scrollx

        pla
        tax
        rts

scroll_colors:
        txa
        pha
        ldx colofs
        dex
        txa
        and #15
        sta colofs
        clc
        adc #<scroll_cols
        sta ptr1
        lda #>scroll_cols
        sta ptr1+1
        SET16   ptr2, scroll_line_col
        ldy #39
_l2     lda (ptr1),y
        sta (ptr2),y
        dey
        bne _l2
        pla
        tax
        rts

draw_scroller:
        clc
        lda #<scroll_text
        adc scroll_pos+0
        sta ptr1
        lda #>scroll_text
        adc scroll_pos+1
        sta ptr1+1
        SET16   ptr2, scroll_line
        ldy #38
_l      lda (ptr1),y
        cmp #$ff
        bne _l2
        ldx #0
        stx scroll_pos
        stx scroll_pos+1
        jmp _l3
_l2     sta (ptr2),y
_l3     dey
        bpl _l
        rts


scroll_pos  .word 0
scroll_text .byte scr("                                      hello world, welcome to the awesome maze!!!!   and one and two and three and four and RESET!                                      "),$ff
colofs  .byte 7
scroll_cols
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_YELLOW,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_YELLOW,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_YELLOW,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_YELLOW,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE

        .byte COL_BLACK,COL_BLACK,COL_BLACK,COL_DGREY,COL_DGREY,COL_DGREY,COL_GREY,COL_GREY,COL_LGREY,COL_WHITE,COL_LGREY,COL_GREY,COL_GREY,COL_DGREY,COL_DGREY,COL_DGREY
        .byte COL_BLACK,COL_BLACK,COL_BLACK,COL_DGREY,COL_DGREY,COL_DGREY,COL_GREY,COL_GREY,COL_LGREY,COL_WHITE,COL_LGREY,COL_GREY,COL_GREY,COL_DGREY,COL_DGREY,COL_DGREY
        .byte COL_BLACK,COL_BLACK,COL_BLACK,COL_DGREY,COL_DGREY,COL_DGREY,COL_GREY,COL_GREY,COL_LGREY,COL_WHITE,COL_LGREY,COL_GREY,COL_GREY,COL_DGREY,COL_DGREY,COL_DGREY
        .byte COL_BLACK,COL_BLACK,COL_BLACK,COL_DGREY,COL_DGREY,COL_DGREY,COL_GREY,COL_GREY,COL_LGREY,COL_WHITE,COL_LGREY,COL_GREY,COL_GREY,COL_DGREY,COL_DGREY,COL_DGREY

        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_WHITE,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE


install_irq:
        sei

        ; set up rasterline interrupt for color changes
        lda #%01111111
        sta $dc0d       ;"Switch off" interrupts signals from CIA-1

        and $d011       ; Clear highest bit of Rasterline
        sta $d011
        lda #rasterline0
        sta $d012       ; first color change on line 50

        lda  #%00000001
        sta  $d01a      ; Enable raster interrupt signals from VIC

        ; store the old irq vector in the proper jmp
        lda $314
        sta irq_chain+1
        lda $315
        sta irq_chain+2

        ; set the new irq handler
        SET16 $314, irq

        lda $d019   ; Clear any pending...
        sta $d019   ; ... VIC interrupt

        cli

        rts


irq:
        lda $d019
        sta $d019   ; Acknowledge interrupt

        ; dispatcher if-elses is way too slow. better use self-modifying code.

_irqdispatch
        jmp _rasterline0

_rasterline0
        lda #8
        sta $d016

        lda #%00010000
        ora scrolly
        sta $d011
        lda #rasterline1
        sta $d012
        SET16 _irqdispatch+1, _rasterline1
        jmp _endirq

_rasterline1
        lda #COL_DGREY
        sta $d021       ; Set up new color
        lda #rasterline2
        sta $d012
        SET16 _irqdispatch+1, _rasterline2
        jmp _endirq

_rasterline2
        lda #COL_GREY
        sta $d021       ; Set up new color
        lda #rasterline3
        sta $d012
        SET16 _irqdispatch+1, _rasterline3
        jmp _endirq

_rasterline3
        lda #COL_LGREY
        sta $d021       ; Set up new color
        lda #rasterline4
        sta $d012
        SET16 _irqdispatch+1, _rasterline4
        jmp _endirq

_rasterline4
        lda #COL_WHITE
        sta $d021       ; Set up new color
        lda #rasterline5
        sta $d012
        SET16 _irqdispatch+1, _rasterline5
        jmp _endirq

_rasterline5
        lda #COL_LGREY
        sta $d021       ; Set up new color
        lda #rasterline6
        sta $d012
        SET16 _irqdispatch+1, _rasterline6
        jmp _endirq

_rasterline6
        lda #COL_GREY
        sta $d021       ; Set up new color
        lda #rasterline7
        sta $d012
        SET16 _irqdispatch+1, _rasterline7
        jmp _endirq

_rasterline7
        lda #COL_DGREY
        sta $d021       ; Set up new color
        lda #rasterline8
        sta $d012
        SET16 _irqdispatch+1, _rasterline8
        jmp _endirq

_rasterline8
        lda #COL_BLACK
        sta $d021       ; Set up new color

        lda #%00010111
        sta $d011

        lda #rasterline9
        sta $d012
        SET16 _irqdispatch+1, _rasterline9
        jmp _endirq

_rasterline9
        lda scrollx
        sta $d016

        lda #rasterline0
        sta $d012
        lda #0
        sta wait
        SET16 _irqdispatch+1, _rasterline0
_endirq
irq_chain:
        jmp $ffff   ; jump to original interrupt handler. Will be patched at runtime


clear_screen
        ; empty video ram
        lda #<vram
        sta ptr1
        lda #>vram
        sta ptr1+1

        ldx #25     ; 20 lines, 40 bytes each
_lx     ldy #0
        lda #160    ; inverted space
_ly     sta (ptr1),y
        iny
        cpy #40
        bne _ly

        ; increment ptr by 40
        clc
        lda ptr1
        adc #40
        sta ptr1
        bcc _l2
        inc ptr1+1
_l2:
        dex
        bne _lx
        clc

clear_color: ; empty color ram
        lda #<color_ram
        sta ptr1
        lda #>color_ram
        sta ptr1+1

        ldx #25      ; 20 loops, 40 bytes each
_lx     ldy #0
        lda #COL_BLACK
_ly     sta (ptr1),y
        iny
        cpy #40
        bne _ly

        ; increment dst by 40
        clc
        lda ptr1
        adc #40
        sta ptr1
        bcc _l2
        inc ptr1+1
_l2:
        dex
        bne _lx
        clc
        rts

generate_line:
        ; set addr to beginning of screen
        SET16 ptr1, last_line
        ldy #39 ; 40 characters
        clc     ; make sure carry is not set
_l      lda $D41B ; load random value
        and #1    ; limit to 0,1
        adc #205
_addr   sta (ptr1),y
        dey
        bpl _l
        rts

scroll_up:
        ; set up dst address
        SET16 ptr1, vram

        ; set up src address
        SET16 ptr2, vram+40

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
        rts

copy_char   .macro adr
        lda adr
        sta adr-40
        .endm

copy_line   .macro line
        copy_char vram+line*40+0
        copy_char vram+line*40+1
        copy_char vram+line*40+2
        copy_char vram+line*40+3
        copy_char vram+line*40+4
        copy_char vram+line*40+5
        copy_char vram+line*40+6
        copy_char vram+line*40+7
        copy_char vram+line*40+8
        copy_char vram+line*40+9
        copy_char vram+line*40+10
        copy_char vram+line*40+11
        copy_char vram+line*40+12
        copy_char vram+line*40+13
        copy_char vram+line*40+14
        copy_char vram+line*40+15
        copy_char vram+line*40+16
        copy_char vram+line*40+17
        copy_char vram+line*40+18
        copy_char vram+line*40+19
        copy_char vram+line*40+20
        copy_char vram+line*40+21
        copy_char vram+line*40+22
        copy_char vram+line*40+23
        copy_char vram+line*40+24
        copy_char vram+line*40+25
        copy_char vram+line*40+26
        copy_char vram+line*40+27
        copy_char vram+line*40+28
        copy_char vram+line*40+29
        copy_char vram+line*40+30
        copy_char vram+line*40+31
        copy_char vram+line*40+32
        copy_char vram+line*40+33
        copy_char vram+line*40+34
        copy_char vram+line*40+35
        copy_char vram+line*40+36
        copy_char vram+line*40+37
        copy_char vram+line*40+38
        copy_char vram+line*40+39
        .endm

scroll_up_fast:
        copy_line 1
        copy_line 2
        copy_line 3
        copy_line 4
        copy_line 5
        copy_line 6
        copy_line 7
        copy_line 8
        copy_line 9
        copy_line 10
        copy_line 11
        copy_line 12
        copy_line 13
        copy_line 14
        copy_line 15
        copy_line 16
        copy_line 17
        copy_line 18
        copy_line 19
        rts

; From https://www.atarimagazines.com/compute/issue72/random_numbers.php

init_random:
    lda #$FF  ; maximum frequency value
    sta $D40E ; voice 3 frequency low byte
    sta $D40F ; voice 3 frequency high byte
    lda #$80  ; noise waveform, gate bit off
    sta $D412 ; voice 3 control register
    rts


