package main

import (
	"io"
	"os"

	"github.com/Velocidex/go-fat/parser"
	kingpin "github.com/alecthomas/kingpin/v2"
	ntfs_parser "www.velocidex.com/golang/go-ntfs/parser"
)

var (
	icat_command = app.Command(
		"icat", "Cat a FAT Cluster.")

	icat_command_file_arg = icat_command.Arg(
		"file", "The image file to inspect",
	).Required().OpenFile(os.O_RDONLY, os.FileMode(0666))

	icat_command_arg = icat_command.Arg(
		"first_cluster", "The first cluster to read from.",
	).Required().Int32()

	icat_command_image_offset = icat_command.Flag(
		"image_offset", "An offset into the file.",
	).Default("0").Int64()
)

func doICAT() {
	reader, _ := ntfs_parser.NewPagedReader(
		getReader(*icat_command_file_arg), 1024, 10000)

	fat, err := parser.GetFATContext(reader)
	kingpin.FatalIfError(err, "Can not open filesystem")

	stream, err := fat.ICat(*icat_command_arg)
	kingpin.FatalIfError(err, "Can not read files")

	var fd io.WriteCloser = os.Stdout
	buf := make([]byte, 4096)
	var offset int64
	for {
		n, err := stream.ReadAt(buf, offset)
		if err == io.EOF {
			break
		}

		if err != nil {
			kingpin.FatalIfError(err, "Can not read files")
		}

		fd.Write(buf[:n])
		offset += int64(n)
	}
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case "icat":
			doICAT()
		default:
			return false
		}
		return true
	})
}
