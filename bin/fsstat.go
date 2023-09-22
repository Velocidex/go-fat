package main

import (
	"fmt"

	"github.com/Velocidex/go-fat/parser"
	"github.com/alecthomas/kingpin/v2"
	ntfs_parser "www.velocidex.com/golang/go-ntfs/parser"
)

var (
	fs_stat_command = app.Command(
		"fsstat", "inspect the MFT record.")

	fs_stat_command_file_arg = fs_stat_command.Arg(
		"file", "The image file to inspect",
	).Required().File()

	fs_stat_command_image_offset = fs_stat_command.Flag(
		"image_offset", "The offset in the image to use.",
	).Int64()
)

func doFSSTAT() {
	reader, _ := ntfs_parser.NewPagedReader(
		getReader(*fs_stat_command_file_arg), 1024, 10000)

	fat, err := parser.GetFATContext(reader)
	kingpin.FatalIfError(err, "Can not open filesystem")

	fmt.Println(fat.DebugString())
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case "fsstat":
			doFSSTAT()

		default:
			return false
		}
		return true
	})
}
