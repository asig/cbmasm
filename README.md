# cbmasm

`cbmasm` is an assembler for the 6502 CPU family that supports local labels, macros, and conditional assembly.

## Usage
```bash
c128asm [inputfile] [outputfile] [-plain]
```
If `inputfile` and `outputfile` are not given, `cbmasm` reads from standard input and writes to standard output.

By default, the generated data starts with the load address, conforming to the "prg" format. If the `-plain` flag is 
set, the load address is suppressed.
