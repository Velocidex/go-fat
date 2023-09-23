package parser

import (
	"io"
)

type FATReader struct {
	context           *FATContext
	offset_to_fat     int64
	total_fat_size    int64
	offset_to_data    int64 // Offset to cluster array start
	bytes_per_cluster int64
	runs              []int64

	// If present this is the directory we have
	Info *DirectoryEntry
}

// If this reader was opened as a result of directory walking we know
// the exact file size from the directory entry. Otherwise we just
// read all available clusters
func (self *FATReader) file_size() int64 {
	if self.Info == nil {
		return int64(len(self.runs)) * self.bytes_per_cluster
	}

	// Sometimes directories indicate that their size is 0 but this is
	// not true.
	if self.Info.IsDir && self.Info.Size == 0 {
		return int64(len(self.runs)) * self.bytes_per_cluster
	}

	return int64(self.Info.Size)
}

func (self *FATReader) Runs() []int64 {
	result := make([]int64, 0, len(self.runs))
	for _, r := range self.runs {
		result = append(result, self.offset_to_data+self.bytes_per_cluster*r)
	}
	return result
}

func (self *FATReader) ReadAt(buff []byte, off int64) (int, error) {
	current_cluster := int(off / self.bytes_per_cluster)
	// Read past the end of the runlist
	if current_cluster >= len(self.runs) || off >= self.file_size() {
		return 0, io.EOF
	}

	current_cluster_offset := off % self.bytes_per_cluster
	current_buff_offset := 0

	for current_buff_offset <= len(buff) && current_cluster < len(self.runs) {
		available_in_file := self.file_size() -
			(off + int64(current_buff_offset))
		available_in_cluster := self.bytes_per_cluster - current_cluster_offset
		available_in_buffer := int64(len(buff) - current_buff_offset)
		to_read := available_in_cluster
		if to_read > available_in_buffer {
			to_read = available_in_buffer
		}

		if to_read > available_in_file {
			to_read = available_in_file
		}

		if to_read == 0 {
			break
		}

		current_cluster_to_read := self.runs[current_cluster]
		n, err := self.context.DiskReader.ReadAt(
			buff[current_buff_offset:int(to_read)+current_buff_offset],
			self.offset_to_data+
				self.bytes_per_cluster*current_cluster_to_read+
				current_cluster_offset)
		if err != nil {
			return 0, err
		}

		// Prepare for the next cluster to read
		current_cluster++
		current_buff_offset += n
		current_cluster_offset = 0
	}

	if current_buff_offset == 0 {
		return 0, io.EOF
	}

	return current_buff_offset, nil
}

func (self *FATReader) readFAT12Entry(cluster int32) uint16 {
	offset := cluster + cluster>>1 // Multiply by 1.5
	value := ParseUint16(self.context.DiskReader,
		self.offset_to_fat+int64(offset))

	// Odd entries are treated different from even entries.
	if cluster&1 > 0 {
		return value >> 4
	}
	return value & 0xfff
}

// Try to parse the FAT based on FAT12. In FAT12 each block is 12 bits
// long. This means 3 bytes represent 2 entries
func (self *FATReader) parseFAT12(first_cluster int32) error {
	current_cluster := first_cluster
	end_of_fat := self.offset_to_fat + self.total_fat_size

	for {
		next_cluster := self.readFAT12Entry(current_cluster)

		// Last sector in the chain
		if next_cluster >= 0xFF8 {
			break
		}

		// This is a bad cluster - not sure what that means?
		if next_cluster == 0xFF7 {
			break
		}

		self.runs = append(self.runs, int64(next_cluster))
		current_cluster = int32(next_cluster)

		// The next cluster points outside the FAT!
		if int64(next_cluster+next_cluster>>1) > end_of_fat {
			break
		}
	}

	return nil
}

// Parse FAT16 structures
func (self *FATReader) parseFAT16(first_cluster int32) error {
	current_cluster := first_cluster
	end_of_fat := self.offset_to_fat + self.total_fat_size

	max_number_of_clusters := int64(4) * 1024 * 1024 * 1024 /
		self.context.Bytes_per_cluster

	for {
		next_cluster := ParseUint16(self.context.DiskReader,
			self.offset_to_fat+int64(current_cluster)*2)

		// Last sector in the chain
		if next_cluster >= 0xFFF8 {
			break
		}

		// This is a bad cluster - not sure what that means?
		if next_cluster == 0xFFF7 || next_cluster <= 2 {
			break
		}

		self.runs = append(self.runs, int64(next_cluster))
		current_cluster = int32(next_cluster)

		// The next cluster points outside the FAT!
		if int64(next_cluster)*2 > end_of_fat ||
			int64(len(self.runs)) > max_number_of_clusters {
			break
		}
	}

	return nil
}

// Parse FAT32 structures
func (self *FATReader) parseFAT32(first_cluster int32) error {
	current_cluster := first_cluster
	end_of_fat := self.offset_to_fat + self.total_fat_size

	max_number_of_clusters := int64(4) * 1024 * 1024 * 1024 /
		self.context.Bytes_per_cluster

	for {
		next_cluster := ParseUint32(self.context.DiskReader,
			self.offset_to_fat+int64(current_cluster)*4) & 0x0FFFFFFF

		// Last sector in the chain
		if next_cluster >= 0x0FFFFFF8 {
			break
		}

		// This is a bad cluster - not sure what that means?
		if next_cluster == 0x0FFFFF7 || next_cluster <= 2 {
			break
		}

		self.runs = append(self.runs, int64(next_cluster))
		current_cluster = int32(next_cluster)

		// The next cluster points outside the FAT!
		if int64(next_cluster)*2 > end_of_fat ||
			int64(len(self.runs)) > max_number_of_clusters {
			break
		}
	}

	return nil
}

func NewFATReader(context *FATContext, first_cluster int32) (*FATReader, error) {
	// Build the runlist is advance and store in memory.
	result := &FATReader{
		context:        context,
		offset_to_fat:  int64(context.reserved_sectors * context.Bytes_per_sector),
		total_fat_size: int64(context.sectors_per_fat * context.Bytes_per_sector),
		// Clusters start after the root directory.
		offset_to_data:    context.rootOffset() + int64(context.MBR.Root_entries()*32) - 2*context.Bytes_per_cluster,
		bytes_per_cluster: context.Bytes_per_cluster,
		runs:              []int64{int64(first_cluster)},
	}

	if context.FATType() == "FAT12" {
		return result, result.parseFAT12(first_cluster)
	}

	if context.FATType() == "FAT16" {
		return result, result.parseFAT16(first_cluster)
	}

	if context.FATType() == "FAT32" {
		return result, result.parseFAT32(first_cluster)
	}

	return result, nil
}
