package main

import (
	"fmt"
	"os"

	"github.com/Velocidex/go-fat/parser"
	kingpin "github.com/alecthomas/kingpin/v2"
	ntfs_parser "www.velocidex.com/golang/go-ntfs/parser"
)

var (
	stat_command = app.Command(
		"stat", "Stat a FAT file.")

	stat_command_file_arg = stat_command.Arg(
		"file", "The image file to inspect",
	).Required().OpenFile(os.O_RDONLY, os.FileMode(0666))

	stat_command_arg = stat_command.Arg(
		"path", "The first cluster to read from.",
	).Required().String()

	stat_command_image_offset = stat_command.Flag(
		"image_offset", "An offset into the file.",
	).Default("0").Int64()
)

func doStat() {
	reader, _ := ntfs_parser.NewPagedReader(
		getReader(*stat_command_file_arg), 1024, 10000)

	fat, err := parser.GetFATContext(reader)
	kingpin.FatalIfError(err, "Can not open filesystem")

	stat, err := fat.Stat(*stat_command_arg)
	kingpin.FatalIfError(err, "Can not stat file")

	Dump(stat)

	// Get the runlist
	stream, err := parser.NewFATReader(fat, stat.FirstCluster)
	kingpin.FatalIfError(err, "Can not open stream")

	bytes_per_sector := fat.Bytes_per_sector
	sectors_per_cluster := fat.Sectors_per_cluster

	fmt.Printf("\n\nSectors (%v)\n", int64(len(stream.Runs()))*sectors_per_cluster)
	for _, cluster := range stream.Runs() {
		// Each cluster has this many sectors
		for i := int64(0); i < sectors_per_cluster; i++ {
			fmt.Printf("%v ", cluster/bytes_per_sector+i)
		}
	}
	fmt.Println("")
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case "stat":
			doStat()
		default:
			return false
		}
		return true
	})
}
