; Example that switches from 8502 to z80
; Shamelessly borrowed from https://www.reddit.com/r/c128/comments/9nxhcc/8502_z80_switchover_a_simple_example/,
; slightly adapted for cbmasm.
; See also the article in Transactor Magazine, Volume 7, Issue 3, page 48 ff.
; (https://archive.org/details/transactor-magazines-v7-i03/page/n51/mode/2up)

PLATFORM_C128   .equ 1

    .include "startup.i"

    sei         ; disable the system interrupt
    lda $ff00   ; the mmu configuration register mirror, seen by all configurations
    pha         ; save the mmu cr into the stack, for safe-keeping

    lda #$c3    ; our first injectable z80 opcode C3 .. which is JP .. the jump instruction
    sta $ffee   ; stash z80 JUMP into $ffee
    lda #<z80code ; lowbyte of address which z80 JUMP will use
    sta $ffef   ; stash lowbyte into $ffef
    lda #>z80code ; highbyte of address which z80 JUMP will use
    sta $fff0   ; stash highbyte into $fff0
    ; at this point we have preloaded $ffee to $fff0 with JP z80code

    lda #$3e    ; load up the #$3e byte, for mmu cr
    sta $ff00   ; set the mmu configuration register mirror with #$3e

    lda $d505   ; load up the mode configuration register
    pha         ; save the mode configuration register to the stack, for safe-keeping
    lda #$b0    ; load up the mode configuration register for z80 action
    sta $d505   ; trigger z80 active .. with #$b0 byte .. 8502 is suspended
    ; z80 comes alive at $ffee and first thing it sees is JP z80code
    nop         ; make sure to add at least one nop here for 8502 restart, after z80 jumps to the bootlink routine
    pla         ; 8502's first active instruction upon restart .. pull mode configuration register from stack
    sta $d505   ; restore mode configuration register
    pla         ; pull mmu configuration register from stack
    sta $ff00   ; restore mmu configuration register
    cli         ; restore system interrupt
    rts         ; return from subroutine


z80code:
    .cpu "z80"
    LD A, $3F       ; load up the #$3f byte, for mmu cr
    LD ($FF00),A    ; set the mmu configuration register mirror with #$3f
    LD A, $53       ; dummy fill-byte seed.

    ; Running the following loop directly in VRAM won't work. It will just fill the memory
    ; with garbage. I have *no* idea why this is happening. Maybe the VIC and Z80 are competing
    ; for RAM access? If you know, tell me!
    LD ($C000),A    ; stash fill-byte at $0400, the location of standard 40x25 VIC-II text screen
    LD HL, $C000    ; load HL with address value #$c000, prepwork for LDIR
    LD DE, $C001    ; load DE with address value #$c001, prepwork for LDIR
    LD BC, $03FF    ; load BC with address value #$03FF (1023 decimal), prepwork for LDIR
    LDIR            ; repeat HL to DE, #$03FF times (re: fill the text screen with #$03FF bytes)

    LD HL, $C000    ; load HL with address value #$c000, prepwork for LDIR
    LD DE, $0400    ; load DE with address value #$0400, prepwork for LDIR
    LD BC, $0400    ; load BC with address value #$0400 (1024 decimal), prepwork for LDIR
    LDIR            ; repeat HL to DE, #$03FF times (re: fill the text screen with #$03FF bytes)


    JP $FFE0        ; jump to the bootlink routine in the Z-80 ROM, 8502 is switched on there.
    NOP             ; add one z80 NOP for safety.
