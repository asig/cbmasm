    .platform "c128"
DEBUG .equ 1

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
;

; This small intro is inspired by the famous maze generator
; 10 PRINT CHR$(205.5+RND(1));:GOTO 10

        .include "macros.i"
        .include "vic.i"

DEBUG_COL   .macro col
            .ifdef DEBUG
            lda col
            sta $d020
            .endif
            .endm

        .if PLATFORM = "c128"
ptr1    .equ $c8 ; abusing RS-232 input buffer
ptr2    .equ $ca ; abusing RS-232 output buffer
        .else
ptr1    .equ $f7 ; abusing RS-232 input buffer
ptr2    .equ $f9 ; abusing RS-232 output buffer
        .endif
ptr_rasterhandlers  .equ $fb

rasterline0 .equ 30             ; set vscroll
rasterline1 .equ 50             ; Switch to DGREY
rasterline2 .equ 60             ; Switch to GREY
rasterline3 .equ 65             ; Switch to LGREY
rasterline4 .equ 70             ; Switch to WHITE

rasterline5 .equ 80             ; Play song

rasterline6 .equ 50-4+20*8-20   ; Switch to LGREY
rasterline7 .equ 50-4+20*8-15   ; Switch to GREY
rasterline8 .equ 50-4+20*8-10   ; Switch to DGREY
rasterline9 .equ 50-4+20*8      ; Switch to BLACK
rasterline_last .equ 255      ; Handle sprites


vram            .equ $400
color_ram       .equ $d800

last_line       .equ vram+19*40
scroll_line     .equ vram+22*40
scroll_line_col .equ color_ram+22*40


sprite_ofs_top  .equ 54
sprite_ofs_left .equ 31

logo_top        .equ 10*8-21    ; centered on 20 lines of maze
logo_left       .equ 19*8-2*24    ; centered on 38 chars of maze


SPRITE_0_DATA   .equ vram+$3f8
SPRITE_1_DATA   .equ vram+$3f9
SPRITE_2_DATA   .equ vram+$3fa
SPRITE_3_DATA   .equ vram+$3fb
SPRITE_4_DATA   .equ vram+$3fc
SPRITE_5_DATA   .equ vram+$3fd
SPRITE_6_DATA   .equ vram+$3fe
SPRITE_7_DATA   .equ vram+$3ff


        .include "startup.i"
        jmp start

scrollx .byte 0
scrolly .byte 0
irqcnt  .byte 0
screenstart .byte 0
start:
        sei

        .if PLATFORM = "c128"
        lda #$3e        ; All RAM, except I/O range at $d000
        sta $ff00


        lda #$ff        ; Turn off BASIC7's raster irq handling...
         sta $d8         ; ... so that we can change the font

        lda $a04        ; Clear bit 0...
        and #%11111110  ; ...of $a04...
        sta $a04        ; ... to disable BASIC IRQ
        .endif

        ; Switch to new font
        lda VIC_MCR
        and #%11110001
        ora #%00001110 ; Bits 3-1 == 111 -> font starts at $3800
        sta VIC_MCR

        cli

        jsr songInit
        jsr init_sprites
        jsr init_random
        jsr clear_screen
        jsr draw_border
        jsr install_irq

        lda #COL_BLACK
        sta $d021
        sta $d020

loop:

        lda #1      ; wait 2 interrupts ...
        sta irqcnt
_w      lda irqcnt
        bne _w

;       
;       lda #1      ; wait another 1 interrupts ...
;       sta irqcnt
;w2     lda irqcnt
;;        bne _w2

        DEBUG_COL #COL_LBLUE
        jsr songPlay
        DEBUG_COL #COL_BLACK


        ; soft scroll the maze
        ldx scrolly
        dex
        bmi _hard_scroll
        stx scrolly
        jmp _end_scroll
_hard_scroll:
        DEBUG_COL #COL_YELLOW
        lda #7
        sta scrolly
        ;ora #%00010000
        ;sta $d011
        jsr scroll_up_fast
        jsr generate_line
        DEBUG_COL #COL_BLACK
_end_scroll
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
        bpl _l2
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
        bne _l2
        ; A contains 0 if we end up here
        sta scroll_pos
        sta scroll_pos+1
        jmp _l3
_l2     sta (ptr2),y
_l3     dey
        bpl _l
        rts


scroll_pos  .word 0
scroll_text .byte scr("                                      hello world, welcome to the awesome maze!!!!   and one and two and three and four and RESET!                                      "),0
colofs  .byte 7
scroll_cols
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_WHITE,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_WHITE,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_WHITE,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_WHITE,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE

        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_YELLOW,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_YELLOW,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_YELLOW,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE
        .byte COL_BLUE,COL_BLUE,COL_BLUE,COL_LBLUE,COL_LBLUE,COL_LBLUE,COL_CYAN,COL_CYAN,COL_LGREEN,COL_YELLOW,COL_LGREEN,COL_CYAN,COL_CYAN,COL_LBLUE,COL_LBLUE,COL_LBLUE

        .byte COL_BLACK,COL_BLACK,COL_BLACK,COL_DGREY,COL_DGREY,COL_DGREY,COL_GREY,COL_GREY,COL_LGREY,COL_WHITE,COL_LGREY,COL_GREY,COL_GREY,COL_DGREY,COL_DGREY,COL_DGREY
        .byte COL_BLACK,COL_BLACK,COL_BLACK,COL_DGREY,COL_DGREY,COL_DGREY,COL_GREY,COL_GREY,COL_LGREY,COL_WHITE,COL_LGREY,COL_GREY,COL_GREY,COL_DGREY,COL_DGREY,COL_DGREY
        .byte COL_BLACK,COL_BLACK,COL_BLACK,COL_DGREY,COL_DGREY,COL_DGREY,COL_GREY,COL_GREY,COL_LGREY,COL_WHITE,COL_LGREY,COL_GREY,COL_GREY,COL_DGREY,COL_DGREY,COL_DGREY
        .byte COL_BLACK,COL_BLACK,COL_BLACK,COL_DGREY,COL_DGREY,COL_DGREY,COL_GREY,COL_GREY,COL_LGREY,COL_WHITE,COL_LGREY,COL_GREY,COL_GREY,COL_DGREY,COL_DGREY,COL_DGREY


install_irq:
        sei

        ; set up rasterline interrupt for color changes
        lda #%01111111
        sta $dc0d       ;"Switch off" interrupts signals from CIA-1

        and $d011       ; Clear highest bit of Rasterline
        sta $d011

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


        ; Set up rasterhandlers ptr and first raster irq
        SET16 ptr_rasterhandlers, rasterhandlers
        lda #0
        sta curhandlerpos
        jsr set_next_rasterirq

        cli

        rts


set_next_rasterirq:
        ldy curhandlerpos
        lda (ptr_rasterhandlers), y
        sta $d012
        iny
        lda (ptr_rasterhandlers), y
        sta _irqdispatch+1
        iny
        lda (ptr_rasterhandlers), y
        sta _irqdispatch+2
        rts

irq:
        lda $d019
        sta $d019   ; Acknowledge interrupt

        ; dispatcher if-elses is way too slow. better use self-modifying code.
_irqdispatch
        jsr $ffff

        ; set up next handler
        lda curhandlerpos
        clc
        adc #3
        cmp #maxhandlerpos
        bne _cont
        lda #0
_cont   sta curhandlerpos
        jsr set_next_rasterirq
irq_chain:
        jmp $ffff   ; jump to original interrupt handler. Will be patched at runtime


rasterline_vscroll
        DEBUG_COL #2

        lda #1
        sta screenstart
        lda #7
        sta $d016

        lda #%00010000
        ora scrolly
        sta $d011

        rts

rasterline_dgrey
        DEBUG_COL #1
        lda #COL_DGREY
        sta $d021       ; Set up new color
        rts

rasterline_grey
        DEBUG_COL #2
        lda #COL_GREY
        sta $d021       ; Set up new color
        rts

rasterline_lgrey
        DEBUG_COL #3
        lda #COL_LGREY
        sta $d021       ; Set up new color
        rts

rasterline_white
        DEBUG_COL #4
        lda #COL_WHITE
        sta $d021       ; Set up new color
        rts

rasterline_signalirq
        DEBUG_COL #COL_CYAN
        lda screenstart
        beq _done
        DEBUG_COL #5
        dec irqcnt
        lda #0
        sta screenstart
_done
        rts


rasterline_black
        DEBUG_COL #9
        lda #COL_BLACK
        sta $d021       ; Set up new color

        lda #%00010111
        sta $d011

        jsr do_scrolltext
        lda scrollx
        sta $d016
        jsr scroll_colors

        rts

rasterline_sprites
        DEBUG_COL #10
        ; soft scroll the "border" char
        ldx fontdata+255*8+7
        ldy fontdata+255*8+6
        lda fontdata+255*8+5
        sta fontdata+255*8+7
        lda fontdata+255*8+4
        sta fontdata+255*8+6
        lda fontdata+255*8+3
        sta fontdata+255*8+5
        lda fontdata+255*8+2
        sta fontdata+255*8+4
        lda fontdata+255*8+1
        sta fontdata+255*8+3
        lda fontdata+255*8+0
        sta fontdata+255*8+2
        stx fontdata+255*8+1
        sty fontdata+255*8+0

        jsr move_sprites
        DEBUG_COL #0

        rts


init_sprites:
        lda #$ff
        sta SPRITES_VISIBLE ; All sprites visible
        lda #0
        sta SPRITES_DBL_W   ; All sprites double width
        sta SPRITES_DBL_H   ; All sprites double height

        lda #0
        sta SPRITES_PRIO    ; All sprite before background

        lda #sprite_1/64
        sta SPRITE_0_DATA

        lda #sprite_2/64
        sta SPRITE_1_DATA

        lda #sprite_3/64
        sta SPRITE_2_DATA

        lda #sprite_4/64
        sta SPRITE_3_DATA

        lda #sprite_box/64
        lda #sprite_filled/64
        sta SPRITE_0_DATA
        sta SPRITE_1_DATA
        sta SPRITE_2_DATA
        sta SPRITE_3_DATA
        sta SPRITE_4_DATA
        sta SPRITE_5_DATA
        sta SPRITE_6_DATA
        sta SPRITE_7_DATA

        ldx #1
        stx SPRITE_0_COL
        inx
        stx SPRITE_1_COL
        inx
        stx SPRITE_2_COL
        inx
        stx SPRITE_3_COL
        inx
        stx SPRITE_4_COL
        inx
        stx SPRITE_5_COL
        inx
        stx SPRITE_6_COL
        inx
        stx SPRITE_7_COL
        rts

sprite_0_pos_x  .byte  0+0*20
sprite_0_pos_y  .byte 64+0*20
sprite_1_pos_x  .byte  0+1*20
sprite_1_pos_y  .byte 64+1*20
sprite_2_pos_x  .byte  0+2*20
sprite_2_pos_y  .byte 64+2*20
sprite_3_pos_x  .byte  0+3*20
sprite_3_pos_y  .byte 64+3*20
sprite_4_pos_x  .byte  0+4*20
sprite_4_pos_y  .byte 64+4*20
sprite_5_pos_x  .byte  0+5*20
sprite_5_pos_y  .byte 64+5*20
sprite_6_pos_x  .byte  0+6*20
sprite_6_pos_y  .byte 64+6*20
sprite_7_pos_x  .byte  0+7*20
sprite_7_pos_y  .byte 64+7*20

sprite_inc_x  .byte 254
sprite_inc_y  .byte 255

MOVESPR   .macro xreg, yreg, posx, posy, incx, incy
        ldy posx
        lda (ptr1),y
        sta xreg
        ldy posy
        lda (ptr1),y
        sta yreg

        lda posx
        clc
        adc incx
        sta posx
        clc
        lda posy
        adc incy
        sta posy
        .endm

MOVESPR2  .macro xreg, yreg, posx, posy
        lda #posx
        sta xreg
        lda #posy
        sta yreg
        .endm

; Lissajous
;move_sprites:
;        SET16 ptr1, sintab
;
;        MOVESPR SPRITE_0_X, SPRITE_0_Y, sprite_0_pos_x, sprite_0_pos_y, sprite_inc_x, sprite_inc_y
;        MOVESPR SPRITE_1_X, SPRITE_1_Y, sprite_1_pos_x, sprite_1_pos_y, sprite_inc_x, sprite_inc_y
;        MOVESPR SPRITE_2_X, SPRITE_2_Y, sprite_2_pos_x, sprite_2_pos_y, sprite_inc_x, sprite_inc_y
;        MOVESPR SPRITE_3_X, SPRITE_3_Y, sprite_3_pos_x, sprite_3_pos_y, sprite_inc_x, sprite_inc_y
;        MOVESPR SPRITE_4_X, SPRITE_4_Y, sprite_4_pos_x, sprite_4_pos_y, sprite_inc_x, sprite_inc_y
;        MOVESPR SPRITE_5_X, SPRITE_5_Y, sprite_5_pos_x, sprite_5_pos_y, sprite_inc_x, sprite_inc_y
;        MOVESPR SPRITE_6_X, SPRITE_6_Y, sprite_6_pos_x, sprite_6_pos_y, sprite_inc_x, sprite_inc_y
;        MOVESPR SPRITE_7_X, SPRITE_7_Y, sprite_7_pos_x, sprite_7_pos_y, sprite_inc_x, sprite_inc_y
;        rts

move_sprites:
        ; Top left coords:  31, 54
        MOVESPR2 SPRITE_0_X, SPRITE_0_Y, sprite_ofs_left + logo_left + 0*24, sprite_ofs_top+logo_top+0*21
        MOVESPR2 SPRITE_1_X, SPRITE_1_Y, sprite_ofs_left + logo_left + 1*24, sprite_ofs_top+logo_top+0*21
        MOVESPR2 SPRITE_2_X, SPRITE_2_Y, sprite_ofs_left + logo_left + 2*24, sprite_ofs_top+logo_top+0*21
        MOVESPR2 SPRITE_3_X, SPRITE_3_Y, sprite_ofs_left + logo_left + 3*24, sprite_ofs_top+logo_top+0*21

        MOVESPR2 SPRITE_4_X, SPRITE_4_Y, sprite_ofs_left + logo_left + 0*24, sprite_ofs_top+logo_top+1*21
        MOVESPR2 SPRITE_5_X, SPRITE_5_Y, sprite_ofs_left + logo_left + 1*24, sprite_ofs_top+logo_top+1*21
        MOVESPR2 SPRITE_6_X, SPRITE_6_Y, sprite_ofs_left + logo_left + 2*24, sprite_ofs_top+logo_top+1*21
        MOVESPR2 SPRITE_7_X, SPRITE_7_Y, sprite_ofs_left + logo_left + 3*24, sprite_ofs_top+logo_top+1*21


        rts

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

draw_border
        lda #$ff
        sta vram+00*40
        sta vram+01*40
        sta vram+02*40
        sta vram+03*40
        sta vram+04*40
        sta vram+05*40
        sta vram+06*40
        sta vram+07*40
        sta vram+08*40
        sta vram+09*40
        sta vram+10*40
        sta vram+11*40
        sta vram+12*40
        sta vram+13*40
        sta vram+14*40
        sta vram+15*40
        sta vram+16*40
        sta vram+17*40
        sta vram+18*40
        sta vram+19*40

        sta vram+00*40+37
        sta vram+01*40+37
        sta vram+02*40+37
        sta vram+03*40+37
        sta vram+04*40+37
        sta vram+05*40+37
        sta vram+06*40+37
        sta vram+07*40+37
        sta vram+08*40+37
        sta vram+09*40+37
        sta vram+10*40+37
        sta vram+11*40+37
        sta vram+12*40+37
        sta vram+13*40+37
        sta vram+14*40+37
        sta vram+15*40+37
        sta vram+16*40+37
        sta vram+17*40+37
        sta vram+18*40+37
        sta vram+19*40+37

        rts

generate_line:
        ; set addr to beginning of screen
        SET16 ptr1, last_line
        ldy #39     ; 39 characters
_l      jsr random  ; load random value
        and #1      ; limit to 0,1
        clc
        adc #205
_addr   sta (ptr1),y
        dey
        bpl _l
        lda #255
        sta last_line+0
        sta last_line+37
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
        ; char 0 is border, no need to copy        
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
        ; char 37 is border, no need to copy        
        ; chars 38 and 39 are not needed because we're in 38-col-mode
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

; from https://gist.github.com/bhickey/0de228c02cc60b5965582d2d946d8c38,
; based on http://www.retroprogramming.com/2017/07/xorshift-pseudorandom-numbers-in-z80.html
init_random:
        lda $d012   ; Initialize rnd with current rasterline
        bne _l      ; but make sure we don't use 0.
        lda #1
_l      sta rndval

random:
        lda rndval
        asl a
        eor rndval
        sta rndval
        lsr a
        eor rndval
        sta rndval
        asl a
        asl a
        eor rndval
        sta rndval
        rts
rndval: .reserve 1

RASTER_HANDLER  .macro line, handler
        .byte line, <handler, >handler
        .endm

curhandlerpos
        .byte 0
rasterhandlers:
        RASTER_HANDLER 30, rasterline_vscroll
        RASTER_HANDLER 50, rasterline_dgrey
        RASTER_HANDLER 60, rasterline_grey
        RASTER_HANDLER 65, rasterline_lgrey
        RASTER_HANDLER 70, rasterline_white

        RASTER_HANDLER  75, rasterline_signalirq

        RASTER_HANDLER 50-4+20*8-20, rasterline_lgrey
        RASTER_HANDLER 50-4+20*8-15, rasterline_grey
        RASTER_HANDLER 50-4+20*8-10, rasterline_dgrey
        RASTER_HANDLER 50-4+20*8, rasterline_black
        RASTER_HANDLER 255, rasterline_sprites
maxhandlerpos   .equ *-rasterhandlers

sintab:
        .byte 128,131,134,137,140,143,146,149,152,156,159,162,165,168,171,174
        .byte 176,179,182,185,188,191,193,196,199,201,204,206,209,211,213,216
        .byte 218,220,222,224,226,228,230,232,234,236,237,239,240,242,243,245
        .byte 246,247,248,249,250,251,252,252,253,254,254,255,255,255,255,255
        .byte 255,255,255,255,255,255,254,254,253,252,252,251,250,249,248,247
        .byte 246,245,243,242,240,239,237,236,234,232,230,228,226,224,222,220
        .byte 218,216,213,211,209,206,204,201,199,196,193,191,188,185,182,179
        .byte 176,174,171,168,165,162,159,156,152,149,146,143,140,137,134,131
        .byte 128,124,121,118,115,112,109,106,103, 99, 96, 93, 90, 87, 84, 81
        .byte  79, 76, 73, 70, 67, 64, 62, 59, 56, 54, 51, 49, 46, 44, 42, 39
        .byte  37, 35, 33, 31, 29, 27, 25, 23, 21, 19, 18, 16, 15, 13, 12, 10
        .byte   9,  8,  7,  6,  5,  4,  3,  3,  2,  1,  1,  0,  0,  0,  0,  0
        .byte   0,  0,  0,  0,  0,  0,  1,  1,  2,  3,  3,  4,  5,  6,  7,  8
        .byte   9, 10, 12, 13, 15, 16, 18, 19, 21, 23, 25, 27, 29, 31, 33, 35
        .byte  37, 39, 42, 44, 46, 49, 51, 54, 56, 59, 62, 64, 67, 70, 73, 76
        .byte  79, 81, 84, 87, 90, 93, 96, 99,103,106,109,112,115,118,121,124



        .align 256 ; Ensure page boundary so that indexing is easy
sprite_cols:
        .byte COL_RED
        .byte COL_RED
        .byte COL_PINK
        .byte COL_RED
        .byte COL_PINK
        .byte COL_PINK
        .byte COL_YELLOW
        .byte COL_PINK
        .byte COL_YELLOW
        .byte COL_YELLOW
        .byte COL_WHITE
        .byte COL_YELLOW
        .byte COL_WHITE
        .byte COL_WHITE
        .byte COL_YELLOW
        .byte COL_WHITE
        .byte COL_YELLOW
        .byte COL_YELLOW
        .byte COL_PINK
        .byte COL_YELLOW
        .byte COL_PINK
        .byte COL_PINK
        .byte COL_RED
        .byte COL_PINK
        .byte COL_RED
        .byte COL_RED
sprite_cols_count       .equ *-sprite_cols
        .byte COL_RED
        .byte COL_RED
        .byte COL_PINK
        .byte COL_RED
        .byte COL_PINK
        .byte COL_PINK
        .byte COL_YELLOW
        .byte COL_PINK
        .byte COL_YELLOW
        .byte COL_YELLOW
        .byte COL_WHITE
        .byte COL_YELLOW
        .byte COL_WHITE
        .byte COL_WHITE
        .byte COL_YELLOW
        .byte COL_WHITE
        .byte COL_YELLOW
        .byte COL_YELLOW
        .byte COL_PINK
        .byte COL_YELLOW
        .byte COL_PINK
        .byte COL_PINK
        .byte COL_RED
        .byte COL_PINK
        .byte COL_RED
        .byte COL_RED









; sprite 0 / singlecolor / color: $01
        .align 64
sprite_0:
        .byte %11111111,%11111111,%11000000
        .byte %11111111,%11111111,%11000000
        .byte %00000111,%11111000,%00000000
        .byte %00000111,%11111000,%00000000
        .byte %00000111,%11111000,%00000000
        .byte %00000111,%11111000,%00000000
        .byte %00000111,%11111000,%00000000
        .byte %00000111,%11111000,%00000000
        .byte %00000111,%11111000,%00000000
        .byte %00000111,%11111000,%00000000
        .byte %00000111,%11111000,%00000000
        .byte %00000111,%11111000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000

; sprite 1 / singlecolor / color: $01
        .align 64
sprite_1:
        .byte %11110000,%00000000,%00011111
        .byte %11111000,%00000000,%00111111
        .byte %11111100,%00000000,%01111111
        .byte %11111110,%00000000,%11111111
        .byte %11111111,%00000001,%11111111
        .byte %11111111,%10000011,%11111111
        .byte %11111111,%11000111,%11111111
        .byte %11111111,%11101111,%11111111
        .byte %11100111,%11111110,%01111111
        .byte %11100011,%11111100,%01111111
        .byte %11100001,%11111000,%01111111
        .byte %11100000,%11110000,%01111111
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000

; sprite 2 / singlecolor / color: $01
        .align 64
sprite_2:
        .byte %00000000,%00111100,%00000000
        .byte %00000000,%01111110,%00000000
        .byte %00000000,%11111111,%00000000
        .byte %00000001,%11111111,%10000000
        .byte %00000011,%11111111,%11000000
        .byte %00000111,%11111111,%11100000
        .byte %00001111,%00111111,%11110000
        .byte %00011110,%00011111,%11111000
        .byte %00111100,%00001111,%11111100
        .byte %01111000,%00000111,%11111110
        .byte %11111111,%11110011,%11111111
        .byte %11111111,%11110001,%11111111
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000

; sprite 3 / singlecolor / color: $01
        .align 64
sprite_3:
        .byte %11111110,%00001111,%11100000
        .byte %11111110,%00001111,%11100000
        .byte %11111110,%00001111,%11100000
        .byte %11111110,%00001111,%11100000
        .byte %11111111,%11111111,%11100000
        .byte %11111111,%11111111,%11100000
        .byte %11111111,%11111111,%11100000
        .byte %11111110,%00001111,%11100000
        .byte %11111110,%00001111,%11100000
        .byte %11111110,%00001111,%11100000
        .byte %11111110,%00001111,%11100000
        .byte %11111110,%00001111,%11100000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000

; sprite 4 / singlecolor / color: $01
        .align 64
sprite_4:
        .byte %11111111,%11111111,%11000000
        .byte %11111111,%11111111,%11000000
        .byte %00000000,%11111111,%10000000
        .byte %00000001,%11111111,%00000000
        .byte %00000011,%11111110,%00000000
        .byte %00000111,%11111100,%00000000
        .byte %00001111,%11111000,%00000000
        .byte %00011111,%11110000,%00000000
        .byte %00111111,%11100000,%00000000
        .byte %01111111,%11000000,%00000000
        .byte %11111111,%11111111,%11000000
        .byte %11111111,%11111111,%11000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000

; sprite 5 / singlecolor / color: $01
        .align 64
sprite_5:
        .byte %00011111,%11111111,%11100000
        .byte %01111111,%11111111,%11100000
        .byte %01111111,%10000000,%00000000
        .byte %11111111,%00000000,%00000000
        .byte %11111111,%01111111,%11100000
        .byte %11111111,%01111111,%11100000
        .byte %11111111,%01111111,%11100000
        .byte %11111111,%00000000,%00000000
        .byte %11111111,%00000000,%00000000
        .byte %01111111,%10000000,%00000000
        .byte %01111111,%11111111,%11100000
        .byte %00011111,%11111111,%11100000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000

; sprite 6 / singlecolor / color: $01
        .align 64
sprite_6:
        .byte %00011111,%11111111,%00000000
        .byte %00111111,%11111111,%10000000
        .byte %01111110,%00001111,%11000000
        .byte %01111100,%00000111,%11000000
        .byte %11111100,%00100111,%11100000
        .byte %11111100,%01100111,%11100000
        .byte %11111100,%11000111,%11100000
        .byte %11111100,%10000111,%11100000
        .byte %01111100,%00000111,%11000000
        .byte %01111110,%00001111,%11000000
        .byte %00111111,%11111111,%10000000
        .byte %00011111,%11111111,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000

; sprite 7 / singlecolor / color: $01
        .align 64
sprite_7:
        .byte %11111111,%11111110,%00000000
        .byte %11111111,%11111111,%00000000
        .byte %00000000,%00111111,%10000000
        .byte %00000000,%00011111,%10000000
        .byte %00000000,%00111111,%10000000
        .byte %01111111,%11111111,%00000000
        .byte %11111111,%11111110,%00000000
        .byte %11111110,%00000000,%00000000
        .byte %11111110,%00000000,%00000000
        .byte %11111110,%00000000,%00000000
        .byte %11111111,%11111111,%10000000
        .byte %11111111,%11111111,%10000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000
        .byte %00000000,%00000000,%00000000


; sprite 7 / singlecolor / color: $01
        .align 64
sprite_box:
        .byte %11111111,%11111111,%11111111
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %10000000,%00000000,%00000001
        .byte %11111111,%11111111,%11111111

        .align 64
sprite_filled:
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111
        .byte %11111111,%11111111,%11111111

        .org $3800
fontdata:
        .incbin "font.bin"

        ; Use tools/sidconv to convert a sid file to asm source
songInit:	.equ $41ed
songPlay:	.equ $4004
        .org $4000
songdata:
        .reserve 1
        .byte $4c, $07, $40, $4c, $b7, $40, $a2, $00, $a9, $01, $85, $80, $8a, $a8, $a5, $80
        .byte $99, $ab, $48, $0a, $85, $80, $bd, $30, $41, $99, $2b, $49, $2a, $9d, $30, $41
        .byte $90, $06, $a5, $80, $69, $00, $85, $80, $18, $98, $69, $0c, $a8, $10, $df, $e8
        .byte $e0, $0c, $d0, $d4, $a0, $0e, $8c, $18, $d4, $a2, $02, $bd, $ea, $41, $99, $2f
        .byte $48, $a9, $0e, $99, $1a, $48, $98, $e9, $07, $a8, $ca, $d0, $ee, $60, $bc, $15
        .byte $48, $bd, $02, $48, $79, $74, $41, $dd, $02, $48, $9d, $02, $48, $b0, $03, $fe
        .byte $03, $48, $98, $0a, $0a, $7d, $16, $48, $fd, $15, $48, $a8, $b9, $51, $41, $9d
        .byte $04, $48, $b9, $3c, $41, $c9, $f0, $90, $12, $29, $0f, $7d, $18, $48, $7d, $2c
        .byte $48, $a8, $b9, $2a, $49, $9d, $00, $48, $b9, $aa, $48, $9d, $01, $48, $fe, $16
        .byte $48, $bd, $16, $48, $29, $03, $d0, $0a, $bc, $15, $48, $b9, $6d, $41, $4a, $4a
        .byte $4a, $4a, $9d, $16, $48, $a0, $07, $bd, $00, $48, $9d, $00, $d4, $e8, $88, $d0
        .byte $f6, $e0, $15, $d0, $03, $60, $a2, $00, $de, $2e, $48, $10, $91, $a9, $06, $9d
        .byte $2e, $48, $de, $2a, $48, $10, $87, $fe, $19, $48, $bc, $19, $48, $bd, $1b, $48
        .byte $9d, $2a, $48, $b9, $7b, $41, $f0, $cd, $30, $1e, $9d, $18, $48, $a9, $00, $9d
        .byte $16, $48, $bc, $15, $48, $b9, $6d, $41, $9d, $03, $48, $b9, $66, $41, $9d, $06
        .byte $48, $a9, $09, $9d, $04, $48, $10, $96, $c9, $ff, $f0, $10, $c9, $df, $29, $0f
        .byte $90, $05, $9d, $15, $48, $10, $c0, $9d, $1b, $48, $10, $bb, $fe, $1a, $48, $bc
        .byte $1a, $48, $b9, $bd, $41, $30, $05, $9d, $19, $48, $10, $ae, $c9, $ff, $f0, $07
        .byte $29, $0f, $9d, $2c, $48, $10, $e5, $bd, $2f, $48, $9d, $1a, $48, $10, $e0, $0c
        .byte $1c, $2d, $3e, $51, $66, $7b, $91, $a9, $c3, $dd, $fa, $f3, $f7, $f0, $f4, $f7
        .byte $f0, $10, $af, $06, $f0, $f0, $f0, $14, $0c, $e0, $fc, $ef, $f0, $fc, $fc, $f0
        .byte $41, $41, $41, $41, $41, $41, $41, $81, $40, $41, $41, $40, $41, $41, $80, $11
        .byte $81, $40, $11, $21, $40, $6f, $6f, $95, $ec, $a9, $79, $6e, $15, $18, $38, $30
        .byte $38, $3b, $36, $1c, $0c, $00, $15, $00, $74, $a5, $00, $ff, $81, $e3, $11, $1d
        .byte $82, $e4, $55, $82, $e3, $11, $81, $1d, $e4, $55, $e3, $18, $ff, $e2, $55, $e5
        .byte $30, $35, $3a, $41, $35, $3a, $29, $ff, $8f, $e0, $35, $e1, $33, $31, $00, $e0
        .byte $35, $e1, $33, $e0, $2e, $00, $ff, $e6, $82, $38, $37, $87, $ff, $33, $80, $2e
        .byte $2c, $ff, $31, $8f, $00, $81, $00, $ff, $8f, $2e, $00, $ff, $f0, $02, $02, $f5
        .byte $02, $02, $f8, $02, $fa, $02, $f1, $02, $02, $ff, $13, $ff, $1e, $1e, $2d, $33
        .byte $2d, $33, $2d, $38, $2d, $33, $2d, $33, $3e, $2d, $f7, $33, $f0, $2d, $f9, $33
        .byte $f0, $2d, $38, $2d, $33, $2d, $33, $3e, $ff, $00, $0e, $10, $a2, $00, $8a, $9d
        .byte $00, $48, $e8, $d0, $fa, $4c, $01, $40
