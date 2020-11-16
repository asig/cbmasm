# cbmasm

`cbmasm` is an assembler for the Z80 and 6502 CPU family that supports local labels, macros, and conditional assembly.

## Usage
```bash
cbmasm [inputfile] [outputfile] [-plain] [-I includedir] [-D sym]
```
If `inputfile` and `outputfile` are not given, `cbmasm` reads from standard input and writes to standard output.

The assumber starts in 6502 mode. By default, the generated data starts with the load address, conforming to Commodore's
"prg" format. If the `-plain` flag is set, the load address is suppressed.

For more details, read the [docs](Documentation.md)
## How to build
```bash
go build
```
