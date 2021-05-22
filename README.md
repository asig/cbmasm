# cbmasm

`cbmasm` is an assembler targeted primarily to the Commodore 128 that allows the programmer to switch between
8510 and Z80 assembly code in a single source file. 

It supports Commodore's `prg` format natively, but can also generate code for pretty much anything that uses a MOS6502
or Z80 CPU. Besides that, `cbmasm` comes with all the features that you expect from a decent assembler: local labels,
macros, conditional assembly, and many more.

It does not support illegal 6510 operations.

## Usage
```bash
cbmasm [flags] [inputfile] [outputfile]
```
Supported flags are:
- `-D sym`: defines a symbol; can be repeated
- `-I value`: include paths; can be repeated
- `-cpu string`: CPU to assemble code for (if not specified otherwise in the source code).   
  Supported values are `6502`, `z80`; default is `6502`
- `-dump_labels`: If true, the labels will be printed.
- `-listing`: If true, a listing is generated.
- `-plain`: If true, the load address is not added to the generated code.
- `-platform string`: Target platform.  
  Supported values are `c128`, `c64`; default is `c128`

If `inputfile` and `outputfile` are not given, `cbmasm` reads from standard input and writes to standard output.

The assumber starts in 6502 mode. By default, the generated data starts with the load address, conforming to Commodore's
"prg" format. If the `-plain` flag is set, the load address is suppressed.

For more details, read the [docs](Documentation.md).

## Building
To build `cbmasm`, just run the following command in the projects rood directory:
```bash
go build           
```                                                                              

You can also run all the tests to verify the binary is behaving as expected:

```bash
go test ./...

```

## Project layout

TODO describe directories

## License
Copyright (c) 2020 Andreas Signer.  
Licensed under [GPLv3](https://www.gnu.org/licenses/gpl-3.0).
