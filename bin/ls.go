package main

import (
	"os"

	"github.com/Velocidex/go-fat/parser"
	kingpin "github.com/alecthomas/kingpin/v2"
	ntfs_parser "www.velocidex.com/golang/go-ntfs/parser"
)

var (
	ls_command = app.Command(
		"ls", "List files.")

	ls_command_file_arg = ls_command.Arg(
		"file", "The image file to inspect",
	).Required().OpenFile(os.O_RDONLY, os.FileMode(0666))

	ls_command_arg = ls_command.Arg(
		"path", "The path to list separated by \\.",
	).Default("\\").String()

	ls_command_image_offset = ls_command.Flag(
		"image_offset", "An offset into the file.",
	).Default("0").Int64()
)

func doLS() {
	reader, _ := ntfs_parser.NewPagedReader(
		getReader(*ls_command_file_arg), 1024, 10000)

	fat, err := parser.GetFATContext(reader)
	kingpin.FatalIfError(err, "Can not open filesystem")

	entries, err := fat.ListDirectory(*ls_command_arg)
	kingpin.FatalIfError(err, "Can not list files")

	Dump(entries)

}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case "ls":
			doLS()
		default:
			return false
		}
		return true
	})
}
