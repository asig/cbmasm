        .ifdef PLATFORM_C128

        ; BASIC header for Commodore 128
        .org $1c01
        .word _next          ; pointer to next line
        .word 65535          ; line number (65535)
        .byte $9e, "7181",0  ; SYS 7181
_next   .word 0              ; End of listing

        .else

        .ifdef PLATFORM_C64

        ; BASIC header for Commodore 64
        .org $801
        .word _next          ; pointer to next line
        .word 65535          ; line number (65535)
        .byte $9e, "2061",0  ; SYS 2061
_next   .word 0              ; End of listing

        .else
        .fail "Neither PLATFORM_C64 nor PLATFORM_C128 defined."
        .endif

        .endif

