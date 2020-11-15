
set16   .macro addr, val
        lda #<(val)
        sta addr
        lda #>(val)
        sta addr+1
        .endm

