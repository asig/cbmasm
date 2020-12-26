/*
 * Copyright (c) 2020 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of cbmasm.
 *
 * cbmasm is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * cbmasm is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with cbmasm.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	input  *os.File
	output *os.File
)

func readUint32() (uint32, error) {
	buf := make([]byte, 4)
	read, err := input.Read(buf)
	if err != nil {
		return 0, err
	}
	if read != 4 {
		return 0, fmt.Errorf("Only read %d bytes instead of 4", read)
	}
	return uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2]) << 8 | uint32(buf[3]), nil
}

func readUint16() (uint16, error) {
	buf := make([]byte, 2)
	read, err := input.Read(buf)
	if err != nil {
		return 0, err
	}
	if read != 2 {
		return 0, fmt.Errorf("Only read %d bytes instead of 2", read)
	}
	return uint16(buf[0]) << 8 | uint16(buf[1]), nil
}

func readUint16LE() (uint16, error) {
	buf := make([]byte, 2)
	read, err := input.Read(buf)
	if err != nil {
		return 0, err
	}
	if read != 2 {
		return 0, fmt.Errorf("Only read %d bytes instead of 2", read)
	}
	return uint16(buf[1]) << 8 | uint16(buf[0]), nil
}

func readUint8() (uint8, error) {
	buf := make([]byte, 1)
	read, err := input.Read(buf)
	if err != nil {
		return 0, err
	}
	if read != 1 {
		return 0, fmt.Errorf("Only read %d bytes instead of 1", read)
	}
	return buf[0], nil
}

func readString(l int) (string, error) {
	buf := make([]byte, l)
	read, err := input.Read(buf)
	if err != nil {
		return "", err
	}
	if read != l {
		return "", fmt.Errorf("Only read %d bytes instead of %d", read, l)
	}
	for l > 0 && buf[l-1] == 0 {
		l = l - 1
	}
	return string(buf[0:l]), nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: sidconv input.sid output.asm\n")
	os.Exit(1)
}

func printErr(format string, vals ...interface{}) {
	fmt.Fprintf(os.Stderr, format, vals...)
	os.Exit(1)
}

func emit(format string, vals ...interface{}) {
	fmt.Fprintf(output, format, vals...)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	var err error
	input, err = os.Open(os.Args[1])
	if err != nil {
		printErr("Can't open input: %s", err)
	}
	defer input.Close()

	output, err = os.Create(os.Args[2])
	if err != nil {
		printErr("Can't open output: %s", err)
	}
	defer output.Close()


	magic, _ := readUint32()
	if magic != 0x50534944 {
		printErr("Wrong magic 0x%08x.\n", magic)
	}

	version, _ := readUint16()
	fmt.Printf("Version is %d.\n", version)

	dataOffset, _ := readUint16()
	fmt.Printf("Data offset is $%04x.\n", dataOffset)

	loadAddress, _ := readUint16()
	fmt.Printf("Load address is $%04x.\n", loadAddress)

	initAddress, _ := readUint16()
	fmt.Printf("Init address is $%04x.\n", initAddress)

	playAddress, _ := readUint16()
	fmt.Printf("Play address is $%04x.\n", playAddress)

	songs, _ := readUint16()
	fmt.Printf("%d songs.\n", songs)

	startSong, _ := readUint16()
	fmt.Printf("Start song is %d.\n", startSong)

	speed, _ := readUint32()
	fmt.Printf("Speed is %%%032b.\n", speed)

	name, _ := readString(32)
	fmt.Printf("Name is %q.\n", name)

	author, _ := readString(32)
	fmt.Printf("Author is %q.\n", author)

	released, _ := readString(32)
	fmt.Printf("Released is %q.\n", released)

	if version > 1 {
		flags, _ := readUint16()
		fmt.Printf("Flags: %%%016b.\n", flags)

		relocStartPage, _ := readUint8()
		fmt.Printf("RelocStartPage is $%02x.\n", relocStartPage)

		relocPages, _ := readUint8()
		fmt.Printf("RelocPages: %d.\n", relocPages)

		secondSIDAddress, _ := readUint8()
		fmt.Printf("2nd SID address is $%02x.\n", secondSIDAddress)

		thirdSIDAddress, _ := readUint8()
		fmt.Printf("3rd SID address is $%02x.\n", thirdSIDAddress)
	}

	if loadAddress == 0 {
		loadAddress, _ = readUint16LE()
		fmt.Printf("Load address from input is $%04x.\n", loadAddress)
	}
	cur, _ := input.Seek(0, io.SeekCurrent)
	end, _ := input.Seek(0, io.SeekEnd)
	input.Seek(cur, io.SeekStart)
	dataSize := uint16(end - cur)
	fmt.Printf("Remaining bytes: %d\n", dataSize)
	data := make([]byte, dataSize)
	input.Read(data)

	// For simplicity, we generate code that copies full pages. Compute # of pages and how much padding we need.
	padding := loadAddress % 256
	pages := (dataSize + padding + 255)/256

	// Generate file header
	emit("; %s\n", name)
	emit("; %s\n", author)
	emit("; %s\n", released)
	emit(";\n")
	emit("; generated with sidconv on %s\n",time.Now().Format(time.RFC822))
	emit("\n")

	// Generate the equs for init and play
	emit("songInit:\t.equ $%04x\n", initAddress)
	emit("songPlay:\t.equ $%04x\n", playAddress)

	// Generate the code that copies the song to the target location
	emit("songCopy:\n")
	emit("        lda #0\n")
	emit("        sta $fa\n")
	emit("        sta $fc\n")
	emit("        lda #>songdata\n")
	emit("        sta $fb\n")
	emit("        lda #$%02x\n", loadAddress/256)
	emit("        sta $fd\n")
	emit("        ldx #%d\n", pages)
	emit("_l1     ldy #0\n")
	emit("_l2     lda ($fa),y\n")
	emit("        sta ($fc),y\n")
	emit("        iny\n")
	emit("        bne _l2\n")
	emit("        inc $fb\n")
	emit("        inc $fd\n")
	emit("        dex\n")
	emit("        bne _l1\n")
	emit("        rts\n")


	// Generate song data, taking page boundaries and padding into account
	emit("        .align 256\n")
	emit("songdata:\n")
	if padding > 0 {
		emit("        .reserve %d\n", padding)
	}
	pos := uint16(0)
	for pos < dataSize {
		var strs []string
		max := pos + 16
		for pos < max && pos < dataSize {
			strs = append(strs, fmt.Sprintf("$%02x", data[pos]))
			pos = pos + 1
		}
		emit("        .byte %s\n", strings.Join(strs, ", "))
	}

	fmt.Printf("Code written to %s.\n", output.Name())
}
