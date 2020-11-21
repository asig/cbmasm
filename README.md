# cbmasm

`cbmasm` is an assembler for the Z80 and 6502 CPU family that supports local labels, macros, and conditional assembly.

## Usage
```bash
cbmasm [flags] [inputfile] [outputfile]
```
Supported flags are:
- `-D value`: defined symbols; can be repeated
- `-I value`: include paths; can be repeated
- `-dump_labels`: If true, the labels will be printed. (default true)
- `-listing`: If true, a listing is generated.
- `-plain`: If true, the load address is not added to the generated code.
- `-platform string`: Target platform. Supported values are: c128, c64 (default "c128")

If `inputfile` and `outputfile` are not given, `cbmasm` reads from standard input and writes to standard output.

The assumber starts in 6502 mode. By default, the generated data starts with the load address, conforming to Commodore's
"prg" format. If the `-plain` flag is set, the load address is suppressed.

For more details, read the [docs](Documentation.md)
## How to build
```bash
go build
```
