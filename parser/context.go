// https://wiki.osdev.org/FAT

package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	notFoundError = errors.New("Not found")
)

type FATContext struct {
	DiskReader        io.ReaderAt
	Profile           *FATProfile
	MBR               *MBR
	bytes_per_cluster int64

	bytes_per_sector int64
	total_sectors    int64
	root_dir_sectors int64
	data_sectors     int64
	total_clusters   int64

	root_directory io.ReaderAt
}

func (self *FATContext) FATType() string {
	//https://wiki.osdev.org/FAT
	if self.MBR.Bytes_per_sector() == 0 {
		return "ExFAT"
	}

	total_clusters := self.total_sectors * self.bytes_per_sector /
		self.bytes_per_cluster

	if total_clusters < 4085 {
		return "FAT12"
	}
	if total_clusters < 65525 {
		return "FAT16"
	}
	return "FAT32"
}

func (self *FATContext) DebugString() string {
	result := self.MBR.DebugString()
	result += fmt.Sprintf("\nOffset of root directory %v\n", self.rootOffset())
	result += fmt.Sprintf("bytes_per_sector %v\n", self.bytes_per_sector)
	result += fmt.Sprintf("total_sectors %v\n", self.total_sectors)
	result += fmt.Sprintf("bytes_per_cluster %v\n", self.bytes_per_cluster)
	result += fmt.Sprintf("FAT Type %v\n", self.FATType())

	return result
}

func (self *FATContext) rootOffset() int64 {
	return int64(
		(self.MBR.Reserved_sectors() +
			uint16(self.MBR.Number_of_fats())*self.MBR.Sectors_per_fat()) * self.MBR.Bytes_per_sector())
}

func (self *FATContext) ListDirectory(path string) ([]*DirectoryEntry, error) {
	components := GetComponents(path)
	if len(components) == 0 {
		return self.listDirectory(self.root_directory, 512), nil
	}

	stream, err := self.Open(path)
	if err != nil {
		return nil, err
	}

	if !stream.Info.IsDir {
		return nil, errors.New("Not a directory")
	}

	return self.listDirectory(stream, 512), nil
}

func (self *FATContext) ICat(first_cluster int32) (io.ReaderAt, error) {
	return NewFATReader(self, first_cluster)
}

func GetFATContext(reader io.ReaderAt) (*FATContext, error) {
	result := &FATContext{
		DiskReader: reader,
		Profile:    NewFATProfile(),
	}

	result.MBR = result.Profile.MBR(reader, 0)
	result.bytes_per_cluster = int64(result.MBR.Bytes_per_sector()) *
		int64(result.MBR.Sectors_per_cluster())

	result.bytes_per_sector = int64(result.MBR.Bytes_per_sector())
	if result.bytes_per_sector == 0 {
		return nil, errors.New("invalid bytes_per_sector")
	}

	result.total_sectors = int64(result.MBR.Small_sectors())
	if result.total_sectors == 0 {
		result.total_sectors = int64(result.MBR.Large_sectors())
	}

	root_directory := make([]byte, 512*32)
	_, err := reader.ReadAt(root_directory, result.rootOffset())
	if err != nil {
		return nil, err
	}
	result.root_directory = bytes.NewReader(root_directory)

	return result, nil
}
