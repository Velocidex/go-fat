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
	FAT32MBR          *FAT32MBR
	Bytes_per_cluster int64

	Bytes_per_sector    int64
	Sectors_per_cluster int64
	total_sectors       int64
	root_dir_sectors    int64
	data_sectors        int64
	total_clusters      int64

	number_of_fats   int64
	sectors_per_fat  int64
	reserved_sectors int64

	root_directory io.ReaderAt
}

func (self *FATContext) FATType() string {
	//https://wiki.osdev.org/FAT
	if self.MBR.Bytes_per_sector() == 0 {
		return "ExFAT"
	}

	if self.FAT32MBR != nil {
		return "FAT32"
	}

	total_clusters := self.total_sectors * self.Bytes_per_sector /
		self.Bytes_per_cluster

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

	if self.FAT32MBR != nil {
		result += "\n" + self.FAT32MBR.DebugString()
	}

	result += fmt.Sprintf("\nOffset of root directory %v\n", self.rootOffset())
	result += fmt.Sprintf("bytes_per_sector %v\n", self.Bytes_per_sector)
	result += fmt.Sprintf("total_sectors %v\n", self.total_sectors)
	result += fmt.Sprintf("bytes_per_cluster %v\n", self.Bytes_per_cluster)
	result += fmt.Sprintf("FAT Type %v\n", self.FATType())

	return result
}

func (self *FATContext) rootOffset() int64 {
	return (self.reserved_sectors +
		self.number_of_fats*self.sectors_per_fat) * self.Bytes_per_sector
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

	// The number of entries is related to the total size of the
	// stream.
	directory_size := int(stream.file_size()) / 32
	return self.listDirectory(stream, directory_size), nil
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

	// Check signature - if it is invalid then try to open as FAT32
	sig := result.MBR.Signature()
	if sig != 0x28 && sig != 0x29 {
		// Maybe this is a fat32
		result.FAT32MBR = result.Profile.FAT32MBR(reader, 0)

		sig = result.FAT32MBR.Signature()
		if sig != 0x28 && sig != 0x29 {
			return nil, errors.New("Invalid signature")
		}
		result.sectors_per_fat = int64(result.FAT32MBR.SectorsPerFat())

	} else {
		result.sectors_per_fat = int64(result.MBR.Sectors_per_fat())

	}

	result.Bytes_per_sector = int64(result.MBR.Bytes_per_sector())
	if result.Bytes_per_sector == 0 {
		return nil, errors.New("invalid bytes_per_sector")
	}

	result.Sectors_per_cluster = int64(result.MBR.Sectors_per_cluster())
	result.reserved_sectors = int64(result.MBR.Reserved_sectors())
	result.number_of_fats = int64(result.MBR.Number_of_fats())
	result.Bytes_per_cluster = result.Bytes_per_sector *
		result.Sectors_per_cluster

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
