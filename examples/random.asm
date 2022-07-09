; From https://www.atarimagazines.com/compute/issue72/random_numbers.php

init_random:
    lda #$FF  ; maximum frequency value
    sta $D40E ; voice 3 frequency low byte
    sta $D40F ; voice 3 frequency high byte
    lda #$80  ; noise waveform, gate bit off
    sta $D412 ; voice 3 control register
    rts
