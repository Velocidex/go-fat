package parser

import (
	"io"
	"sort"
	"strings"
	"time"
)

type DirectoryEntry struct {
	Name         string
	ShortName    string
	Size         uint64
	Mtime        time.Time
	Atime        time.Time
	Ctime        time.Time
	Attribute    string
	IsDir        bool
	IsDeleted    bool
	FirstCluster int32
}

func (self *FATContext) Open(path string) (*FATReader, error) {
	stat, err := self.Stat(path)
	if err != nil {
		return nil, err
	}

	// Now try to open the stream
	stream, err := NewFATReader(self, stat.FirstCluster)
	if err != nil {
		return nil, err
	}
	stream.Info = stat

	return stream, nil
}

func (self *FATContext) OpenComponents(components []string) (*FATReader, error) {
	stat, err := self.StatComponents(components)
	if err != nil {
		return nil, err
	}

	// Now try to open the stream
	stream, err := NewFATReader(self, stat.FirstCluster)
	if err != nil {
		return nil, err
	}
	stream.Info = stat
	return stream, nil
}

func (self *FATContext) Stat(path string) (*DirectoryEntry, error) {
	components := GetComponents(path)
	return self.StatComponents(components)
}

func (self *FATContext) StatComponents(
	components []string) (*DirectoryEntry, error) {
	// Start at the root directory
	stream := self.root_directory
	directory_size := 512

component_search:
	for idx, component := range components {
		component = strings.ToLower(component)

		entries := self.listDirectory(stream, directory_size)
		if len(entries) == 0 {
			return nil, notFoundError
		}

		for _, e := range entries {
			if strings.ToLower(e.ShortName) == component ||
				strings.ToLower(e.Name) == component {

				// Report the final component
				if idx == len(components)-1 {
					return e, nil
				}

				// Recurse into directories
				if e.IsDir {
					fat_stream, err := NewFATReader(self, e.FirstCluster)
					if err != nil {
						return nil, err
					}
					// Number of entries in this directory
					directory_size = int(fat_stream.file_size()) / 32
					stream = fat_stream

					continue component_search
				}
			}
		}

		// None of the entries matches the component.
		return nil, notFoundError
	}

	// Return an entry for the root directory
	return &DirectoryEntry{
		IsDir:        true,
		FirstCluster: 0,
	}, nil
}

func (self *FATContext) listDirectory(
	stream io.ReaderAt, number_of_entries int) []*DirectoryEntry {
	results := []*DirectoryEntry{}

	entries := ParseArray_FolderEntry(self.Profile, stream,
		0, number_of_entries)

	type name_fragment struct {
		order int
		idx   int
		name  string
	}

	var current_name []name_fragment
	for idx, entry := range entries {
		// fmt.Println(entry.DebugString())

		if entry.Attribute().IsSet("LFN") {
			lfn_entry := self.Profile.LFNEntry(stream, entry.Offset)
			current_name = append(current_name, name_fragment{
				order: int(lfn_entry.Order()),
				idx:   idx,
				name:  lfn_entry.Name1() + lfn_entry.Name2() + lfn_entry.Name3(),
			})
			continue
		}

		name := ""
		sort.Slice(current_name, func(i, j int) bool {
			// If the order is the same in all parts then use the
			// index as a tie breaker.
			if current_name[i].order == current_name[j].order {
				return current_name[i].idx > current_name[j].idx
			}

			return current_name[i].order < current_name[j].order
		})

		for _, fragment := range current_name {
			name += fragment.name
		}

		name = strings.SplitN(name, "\x00", 2)[0]
		short_name := entry.Name()
		extension := ""

		if len(short_name) == 11 {
			extension = strings.TrimRight(short_name[8:], " \x00")
			short_name = strings.TrimRight(short_name[:8], " \x00")
			if extension != "" {
				short_name += "." + extension
			}
		}

		if len(short_name) == 0 {
			break
		}

		if name == "" {
			name = short_name
		}

		deleted := false
		if len(short_name) > 0 && short_name[0] == 0xe5 {
			deleted = true
			short_name = "_" + short_name[1:]
		}

		results = append(results, &DirectoryEntry{
			Name:      name,
			ShortName: short_name,
			Size:      uint64(entry.FileSize()),
			Ctime: time.Date(
				1980+int(entry._CreateDate()>>9),      // Year
				time.Month(entry._CreateDate()>>5&15), // Month
				int(entry._CreateDate()&31),           // Day
				int(entry._CreateTime()>>11),          // Hour
				int(entry._CreateTime()>>5&63),        // Minutes
				int(entry._CreateTime()&31)*2,         // Seconds
				int(entry._CreateTimeTenthSeconds())*10000000,
				time.UTC,
			),

			Mtime: time.Date(
				1980+int(entry._LastModDate()>>9),      // Year
				time.Month(entry._LastModDate()>>5&15), // Month
				int(entry._LastModDate()&31),           // Day
				int(entry._LastModTime()>>11),          // Hour
				int(entry._LastModTime()>>5&63),        // Minutes
				int(entry._LastModTime()&31)*2,         // Seconds
				0,
				time.UTC,
			),

			Atime: time.Date(
				1980+int(entry._LastAccessDate()>>9),      // Year
				time.Month(entry._LastAccessDate()>>5&15), // Month
				int(entry._LastAccessDate()&31),           // Day
				0,                                         // Hour
				0,                                         // Minutes
				0,                                         // Seconds
				0,
				time.UTC,
			),

			Attribute: getAttributes(entry),
			IsDir:     entry.Attribute().IsSet("DIRECTORY"),
			IsDeleted: deleted,
			FirstCluster: int32(entry._ClusterHigh())<<16 |
				int32(entry._ClusterLow()),
		})
		current_name = nil
	}

	return results
}

func getAttributes(entry *FolderEntry) string {
	result := ""
	attr := entry.Attribute()
	if attr.IsSet("READ_ONLY") {
		result += "R"
	}

	if attr.IsSet("HIDDEN") {
		result += "H"
	}
	if attr.IsSet("SYSTEM") {
		result += "S"
	}
	if attr.IsSet("DIRECTORY") {
		result += "D"
	}
	if attr.IsSet("ARCHIVE") {
		result += "A"
	}

	return result
}
