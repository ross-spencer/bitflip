package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	app := cli.NewApp()
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "offset",
			Description: "flip a bit in a file at an offset",
			Usage:       "byte-offset filepath",
			ArgsUsage: `byte-offset filepath
			byte-offset: of the form <offset>@<bit>, where:
			- <offset> is a integer speciying in which byte the flip should occur
			- <bit> is an integer from 0 to 7 speciying which bit to flip (from LSB)

			filepath: is the file in which to do the bitflipping`,
			Action: func(cctx *cli.Context) error {
				if cctx.NArg() != 2 {
					return cli.ShowCommandHelp(cctx, "offset")
				}
				byteOffset, bitOffset, err := parseByteOffset(cctx.Args().Get(0))
				if err != nil {
					return err
				}
				filename := cctx.Args().Get(1)

				file, err := os.OpenFile(filename, os.O_RDWR, 0)
				if err != nil {
					return err
				}
				defer file.Close()

				log.Printf("flipping %dth bit of byte %d in file %q", bitOffset, byteOffset, filename)
				if err := flipBitAtOffset(file, int64(byteOffset), bitOffset); err != nil {
					return fmt.Errorf("flipping bit: %v", err)
				}

				return nil
			},
		},
		cli.Command{
			Name:        "random",
			Description: "flip a random bit in a file",
			Usage:       "filepath",
			ArgsUsage: `filepath
			filepath: is the file in which to do the random bitflipping`,
			Action: func(cctx *cli.Context) error {
				if cctx.NArg() != 1 {
					return cli.ShowCommandHelp(cctx, "random")
				}
				filename := cctx.Args().Get(0)

				file, err := os.OpenFile(filename, os.O_RDWR, 0)
				if err != nil {
					return err
				}
				defer file.Close()

				fi, err := file.Stat()
				if err != nil {
					return err
				}
				byteOffset := r.Int63n(fi.Size())
				bitOffset := uint8(r.Intn(8))

				log.Printf("flipping %dth bit of byte %d in file %q", bitOffset, byteOffset, filename)
				if err := flipBitAtOffset(file, int64(byteOffset), bitOffset); err != nil {
					return fmt.Errorf("flipping bit: %v", err)
				}

				return nil
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func parseByteOffset(offset string) (byteOffset uint64, bitOffset uint8, err error) {
	args := strings.Split(offset, "@")
	if len(args) > 2 {
		return byteOffset, bitOffset, errors.New("invalid offset string, too many '@' symbols found")
	}
	if len(args) < 2 {
		return byteOffset, bitOffset, errors.New("invalid offset string, no '@' symbols found")
	}
	iByteOffset, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return byteOffset, bitOffset, fmt.Errorf("invalid byte offset: %v", err)
	}
	iBitOffset, err := strconv.ParseUint(args[1], 10, 8)
	if err != nil {
		return byteOffset, bitOffset, fmt.Errorf("invalid bit offset: %v", err)
	}
	if iBitOffset > 7 {
		return byteOffset, bitOffset, fmt.Errorf("bit offset must be between 0 and 7, was %d", iBitOffset)
	}
	byteOffset = uint64(iByteOffset)
	bitOffset = uint8(iBitOffset)
	return byteOffset, bitOffset, nil
}

func flipBitAtOffset(file io.ReadWriteSeeker, byteOffset int64, bitOffset uint8) error {
	_, err := file.Seek(byteOffset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seeking to flip offset: %v", err)
	}
	oneByte := make([]byte, 1)
	if _, err := file.Read(oneByte); err != nil {
		return fmt.Errorf("reading data from flip offset: %v", err)
	}

	oneByte[0] = toggleNthBit(oneByte[0], bitOffset)

	_, err = file.Seek(-1, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("seeking back to flip offset: %v", err)
	}
	if _, err := file.Write(oneByte); err != nil {
		return fmt.Errorf("writing back the flipped bit: %v", err)
	}
	return nil
}

func toggleNthBit(b byte, n uint8) byte {
	singleBitMask := uint8(1) << (n)
	return b ^ singleBitMask
}
