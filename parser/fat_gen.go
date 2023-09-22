
package parser

// Autogenerated code from fat_profile.json. Do not edit.

import (
    "encoding/binary"
    "fmt"
    "bytes"
    "io"
    "sort"
    "strings"
    "unicode/utf16"
    "unicode/utf8"
)

var (
   // Depending on autogenerated code we may use this. Add a reference
   // to shut the compiler up.
   _ = bytes.MinRead
   _ = fmt.Sprintf
   _ = utf16.Decode
   _ = binary.LittleEndian
   _ = utf8.RuneError
   _ = sort.Strings
   _ = strings.Join
   _ = io.Copy
)

func indent(text string) string {
    result := []string{}
    lines := strings.Split(text,"\n")
    for _, line := range lines {
         result = append(result, "  " + line)
    }
    return strings.Join(result, "\n")
}


type FATProfile struct {
    Off_DirectoryListing_Entries int64
    Off_FolderEntry_Name int64
    Off_FolderEntry_Attribute int64
    Off_FolderEntry__CreateTimeTenthSeconds int64
    Off_FolderEntry__CreateTime int64
    Off_FolderEntry__CreateDate int64
    Off_FolderEntry__LastAccessDate int64
    Off_FolderEntry__ClusterHigh int64
    Off_FolderEntry__LastModTime int64
    Off_FolderEntry__LastModDate int64
    Off_FolderEntry__ClusterLow int64
    Off_FolderEntry_FileSize int64
    Off_LFNEntry_Order int64
    Off_LFNEntry_Name1 int64
    Off_LFNEntry_Name2 int64
    Off_LFNEntry_Name3 int64
    Off_MBR_Oemname int64
    Off_MBR_Bytes_per_sector int64
    Off_MBR_Sectors_per_cluster int64
    Off_MBR_Reserved_sectors int64
    Off_MBR_Number_of_fats int64
    Off_MBR_Root_entries int64
    Off_MBR_Small_sectors int64
    Off_MBR_Sectors_per_fat int64
    Off_MBR_Large_sectors int64
    Off_MBR_Signature int64
    Off_MBR_Volume_serial_number int64
    Off_MBR_Volume_label int64
    Off_MBR_System_id int64
    Off_MBR_Magic int64
}

func NewFATProfile() *FATProfile {
    // Specific offsets can be tweaked to cater for slight version mismatches.
    self := &FATProfile{0,0,11,13,14,16,18,20,22,24,26,28,0,1,14,28,3,11,13,14,16,17,19,22,32,38,39,43,54,510}
    return self
}

func (self *FATProfile) DirectoryListing(reader io.ReaderAt, offset int64) *DirectoryListing {
    return &DirectoryListing{Reader: reader, Offset: offset, Profile: self}
}

func (self *FATProfile) FolderEntry(reader io.ReaderAt, offset int64) *FolderEntry {
    return &FolderEntry{Reader: reader, Offset: offset, Profile: self}
}

func (self *FATProfile) LFNEntry(reader io.ReaderAt, offset int64) *LFNEntry {
    return &LFNEntry{Reader: reader, Offset: offset, Profile: self}
}

func (self *FATProfile) MBR(reader io.ReaderAt, offset int64) *MBR {
    return &MBR{Reader: reader, Offset: offset, Profile: self}
}


type DirectoryListing struct {
    Reader io.ReaderAt
    Offset int64
    Profile *FATProfile
    
    _Entries []*FolderEntry
    _Entries_cached bool

}

func (self *DirectoryListing) Size() int {
    return 0
}

func (self *DirectoryListing) Entries() []*FolderEntry {
   return ParseArray_FolderEntry(self.Profile, self.Reader, self.Profile.Off_DirectoryListing_Entries + self.Offset, 512)
}
func (self *DirectoryListing) DebugString() string {
    result := fmt.Sprintf("struct DirectoryListing @ %#x:\n", self.Offset)
    return result
}

type FolderEntry struct {
    Reader io.ReaderAt
    Offset int64
    Profile *FATProfile
    
    __CreateTimeTenthSeconds uint8
    __CreateTimeTenthSeconds_cached bool

    __CreateTime uint16
    __CreateTime_cached bool

    __CreateDate uint16
    __CreateDate_cached bool

    __LastAccessDate uint16
    __LastAccessDate_cached bool

    __ClusterHigh uint16
    __ClusterHigh_cached bool

    __LastModTime uint16
    __LastModTime_cached bool

    __LastModDate uint16
    __LastModDate_cached bool

    __ClusterLow uint16
    __ClusterLow_cached bool

    _FileSize uint32
    _FileSize_cached bool

}

func (self *FolderEntry) Size() int {
    return 32
}


func (self *FolderEntry) Name() string {
  return ParseString(self.Reader, self.Profile.Off_FolderEntry_Name + self.Offset, 11)
}

func (self *FolderEntry) Attribute() *Flags {
   value := ParseUint8(self.Reader, self.Profile.Off_FolderEntry_Attribute + self.Offset)
   names := make(map[string]bool)


   if value & 4 != 0 {
      names["SYSTEM"] = true
   }

   if value & 8 != 0 {
      names["VOLUME_ID"] = true
   }

   if value & 16 != 0 {
      names["DIRECTORY"] = true
   }

   if value & 32 != 0 {
      names["ARCHIVE"] = true
   }

   if value & 15 != 0 {
      names["LFN"] = true
   }

   if value & 1 != 0 {
      names["READ_ONLY"] = true
   }

   if value & 2 != 0 {
      names["HIDDEN"] = true
   }

   return &Flags{Value: uint64(value), Names: names}
}


func (self *FolderEntry) _CreateTimeTenthSeconds() byte {
   if self.__CreateTimeTenthSeconds_cached {
       return self.__CreateTimeTenthSeconds
   }
   result := ParseUint8(self.Reader, self.Profile.Off_FolderEntry__CreateTimeTenthSeconds + self.Offset)
   self.__CreateTimeTenthSeconds = result
   self.__CreateTimeTenthSeconds_cached = true
   return result
}

func (self *FolderEntry) _CreateTime() uint16 {
   if self.__CreateTime_cached {
       return self.__CreateTime
   }
   result := ParseUint16(self.Reader, self.Profile.Off_FolderEntry__CreateTime + self.Offset)
   self.__CreateTime = result
   self.__CreateTime_cached = true
   return result
}

func (self *FolderEntry) _CreateDate() uint16 {
   if self.__CreateDate_cached {
       return self.__CreateDate
   }
   result := ParseUint16(self.Reader, self.Profile.Off_FolderEntry__CreateDate + self.Offset)
   self.__CreateDate = result
   self.__CreateDate_cached = true
   return result
}

func (self *FolderEntry) _LastAccessDate() uint16 {
   if self.__LastAccessDate_cached {
       return self.__LastAccessDate
   }
   result := ParseUint16(self.Reader, self.Profile.Off_FolderEntry__LastAccessDate + self.Offset)
   self.__LastAccessDate = result
   self.__LastAccessDate_cached = true
   return result
}

func (self *FolderEntry) _ClusterHigh() uint16 {
   if self.__ClusterHigh_cached {
       return self.__ClusterHigh
   }
   result := ParseUint16(self.Reader, self.Profile.Off_FolderEntry__ClusterHigh + self.Offset)
   self.__ClusterHigh = result
   self.__ClusterHigh_cached = true
   return result
}

func (self *FolderEntry) _LastModTime() uint16 {
   if self.__LastModTime_cached {
       return self.__LastModTime
   }
   result := ParseUint16(self.Reader, self.Profile.Off_FolderEntry__LastModTime + self.Offset)
   self.__LastModTime = result
   self.__LastModTime_cached = true
   return result
}

func (self *FolderEntry) _LastModDate() uint16 {
   if self.__LastModDate_cached {
       return self.__LastModDate
   }
   result := ParseUint16(self.Reader, self.Profile.Off_FolderEntry__LastModDate + self.Offset)
   self.__LastModDate = result
   self.__LastModDate_cached = true
   return result
}

func (self *FolderEntry) _ClusterLow() uint16 {
   if self.__ClusterLow_cached {
       return self.__ClusterLow
   }
   result := ParseUint16(self.Reader, self.Profile.Off_FolderEntry__ClusterLow + self.Offset)
   self.__ClusterLow = result
   self.__ClusterLow_cached = true
   return result
}

func (self *FolderEntry) FileSize() uint32 {
   if self._FileSize_cached {
      return self._FileSize
   }
   result := ParseUint32(self.Reader, self.Profile.Off_FolderEntry_FileSize + self.Offset)
   self._FileSize = result
   self._FileSize_cached = true
   return result
}
func (self *FolderEntry) DebugString() string {
    result := fmt.Sprintf("struct FolderEntry @ %#x:\n", self.Offset)
    result += fmt.Sprintf("  Name: %v\n", string(self.Name()))
    result += fmt.Sprintf("  Attribute: %v\n", self.Attribute().DebugString())
    result += fmt.Sprintf("  _CreateTimeTenthSeconds: %#0x\n", self._CreateTimeTenthSeconds())
    result += fmt.Sprintf("  _CreateTime: %#0x\n", self._CreateTime())
    result += fmt.Sprintf("  _CreateDate: %#0x\n", self._CreateDate())
    result += fmt.Sprintf("  _LastAccessDate: %#0x\n", self._LastAccessDate())
    result += fmt.Sprintf("  _ClusterHigh: %#0x\n", self._ClusterHigh())
    result += fmt.Sprintf("  _LastModTime: %#0x\n", self._LastModTime())
    result += fmt.Sprintf("  _LastModDate: %#0x\n", self._LastModDate())
    result += fmt.Sprintf("  _ClusterLow: %#0x\n", self._ClusterLow())
    result += fmt.Sprintf("  FileSize: %#0x\n", self.FileSize())
    return result
}

type LFNEntry struct {
    Reader io.ReaderAt
    Offset int64
    Profile *FATProfile
    
    _Order int8
    _Order_cached bool

}

func (self *LFNEntry) Size() int {
    return 32
}

func (self *LFNEntry) Order() int8 {
   if self._Order_cached {
       return self._Order
   }
   result := ParseInt8(self.Reader, self.Profile.Off_LFNEntry_Order + self.Offset)
   self._Order = result
   self._Order_cached = true
   return result
}


func (self *LFNEntry) Name1() string {
  return ParseUTF16String(self.Reader, self.Profile.Off_LFNEntry_Name1 + self.Offset, 10)
}


func (self *LFNEntry) Name2() string {
  return ParseUTF16String(self.Reader, self.Profile.Off_LFNEntry_Name2 + self.Offset, 12)
}


func (self *LFNEntry) Name3() string {
  return ParseUTF16String(self.Reader, self.Profile.Off_LFNEntry_Name3 + self.Offset, 4)
}
func (self *LFNEntry) DebugString() string {
    result := fmt.Sprintf("struct LFNEntry @ %#x:\n", self.Offset)
    result += fmt.Sprintf("  Order: %#0x\n", self.Order())
    result += fmt.Sprintf("  Name1: %v\n", string(self.Name1()))
    result += fmt.Sprintf("  Name2: %v\n", string(self.Name2()))
    result += fmt.Sprintf("  Name3: %v\n", string(self.Name3()))
    return result
}

type MBR struct {
    Reader io.ReaderAt
    Offset int64
    Profile *FATProfile
    
    _Bytes_per_sector uint16
    _Bytes_per_sector_cached bool

    _Sectors_per_cluster uint8
    _Sectors_per_cluster_cached bool

    _Reserved_sectors uint16
    _Reserved_sectors_cached bool

    _Number_of_fats uint8
    _Number_of_fats_cached bool

    _Root_entries uint16
    _Root_entries_cached bool

    _Small_sectors uint16
    _Small_sectors_cached bool

    _Sectors_per_fat uint16
    _Sectors_per_fat_cached bool

    _Large_sectors uint32
    _Large_sectors_cached bool

    _Signature uint8
    _Signature_cached bool

    _Volume_serial_number uint32
    _Volume_serial_number_cached bool

    _Magic uint16
    _Magic_cached bool

}

func (self *MBR) Size() int {
    return 512
}


func (self *MBR) Oemname() string {
  return ParseString(self.Reader, self.Profile.Off_MBR_Oemname + self.Offset, 8)
}

func (self *MBR) Bytes_per_sector() uint16 {
   if self._Bytes_per_sector_cached {
       return self._Bytes_per_sector
   }
   result := ParseUint16(self.Reader, self.Profile.Off_MBR_Bytes_per_sector + self.Offset)
   self._Bytes_per_sector = result
   self._Bytes_per_sector_cached = true
   return result
}

func (self *MBR) Sectors_per_cluster() byte {
   if self._Sectors_per_cluster_cached {
       return self._Sectors_per_cluster
   }
   result := ParseUint8(self.Reader, self.Profile.Off_MBR_Sectors_per_cluster + self.Offset)
   self._Sectors_per_cluster = result
   self._Sectors_per_cluster_cached = true
   return result
}

func (self *MBR) Reserved_sectors() uint16 {
   if self._Reserved_sectors_cached {
       return self._Reserved_sectors
   }
   result := ParseUint16(self.Reader, self.Profile.Off_MBR_Reserved_sectors + self.Offset)
   self._Reserved_sectors = result
   self._Reserved_sectors_cached = true
   return result
}

func (self *MBR) Number_of_fats() byte {
   if self._Number_of_fats_cached {
       return self._Number_of_fats
   }
   result := ParseUint8(self.Reader, self.Profile.Off_MBR_Number_of_fats + self.Offset)
   self._Number_of_fats = result
   self._Number_of_fats_cached = true
   return result
}

func (self *MBR) Root_entries() uint16 {
   if self._Root_entries_cached {
       return self._Root_entries
   }
   result := ParseUint16(self.Reader, self.Profile.Off_MBR_Root_entries + self.Offset)
   self._Root_entries = result
   self._Root_entries_cached = true
   return result
}

func (self *MBR) Small_sectors() uint16 {
   if self._Small_sectors_cached {
       return self._Small_sectors
   }
   result := ParseUint16(self.Reader, self.Profile.Off_MBR_Small_sectors + self.Offset)
   self._Small_sectors = result
   self._Small_sectors_cached = true
   return result
}

func (self *MBR) Sectors_per_fat() uint16 {
   if self._Sectors_per_fat_cached {
       return self._Sectors_per_fat
   }
   result := ParseUint16(self.Reader, self.Profile.Off_MBR_Sectors_per_fat + self.Offset)
   self._Sectors_per_fat = result
   self._Sectors_per_fat_cached = true
   return result
}

func (self *MBR) Large_sectors() uint32 {
   if self._Large_sectors_cached {
      return self._Large_sectors
   }
   result := ParseUint32(self.Reader, self.Profile.Off_MBR_Large_sectors + self.Offset)
   self._Large_sectors = result
   self._Large_sectors_cached = true
   return result
}

func (self *MBR) Signature() byte {
   if self._Signature_cached {
       return self._Signature
   }
   result := ParseUint8(self.Reader, self.Profile.Off_MBR_Signature + self.Offset)
   self._Signature = result
   self._Signature_cached = true
   return result
}

func (self *MBR) Volume_serial_number() uint32 {
   if self._Volume_serial_number_cached {
      return self._Volume_serial_number
   }
   result := ParseUint32(self.Reader, self.Profile.Off_MBR_Volume_serial_number + self.Offset)
   self._Volume_serial_number = result
   self._Volume_serial_number_cached = true
   return result
}


func (self *MBR) Volume_label() string {
  return ParseString(self.Reader, self.Profile.Off_MBR_Volume_label + self.Offset, 11)
}


func (self *MBR) System_id() string {
  return ParseString(self.Reader, self.Profile.Off_MBR_System_id + self.Offset, 8)
}

func (self *MBR) Magic() uint16 {
   if self._Magic_cached {
       return self._Magic
   }
   result := ParseUint16(self.Reader, self.Profile.Off_MBR_Magic + self.Offset)
   self._Magic = result
   self._Magic_cached = true
   return result
}
func (self *MBR) DebugString() string {
    result := fmt.Sprintf("struct MBR @ %#x:\n", self.Offset)
    result += fmt.Sprintf("  Oemname: %v\n", string(self.Oemname()))
    result += fmt.Sprintf("  Bytes_per_sector: %#0x\n", self.Bytes_per_sector())
    result += fmt.Sprintf("  Sectors_per_cluster: %#0x\n", self.Sectors_per_cluster())
    result += fmt.Sprintf("  Reserved_sectors: %#0x\n", self.Reserved_sectors())
    result += fmt.Sprintf("  Number_of_fats: %#0x\n", self.Number_of_fats())
    result += fmt.Sprintf("  Root_entries: %#0x\n", self.Root_entries())
    result += fmt.Sprintf("  Small_sectors: %#0x\n", self.Small_sectors())
    result += fmt.Sprintf("  Sectors_per_fat: %#0x\n", self.Sectors_per_fat())
    result += fmt.Sprintf("  Large_sectors: %#0x\n", self.Large_sectors())
    result += fmt.Sprintf("  Signature: %#0x\n", self.Signature())
    result += fmt.Sprintf("  Volume_serial_number: %#0x\n", self.Volume_serial_number())
    result += fmt.Sprintf("  Volume_label: %v\n", string(self.Volume_label()))
    result += fmt.Sprintf("  System_id: %v\n", string(self.System_id()))
    result += fmt.Sprintf("  Magic: %#0x\n", self.Magic())
    return result
}

type Flags struct {
    Value uint64
    Names  map[string]bool
}

func (self Flags) DebugString() string {
    names := []string{}
    for k, _ := range self.Names {
      names = append(names, k)
    }

    sort.Strings(names)

    return fmt.Sprintf("%d (%s)", self.Value, strings.Join(names, ","))
}

func (self Flags) IsSet(flag string) bool {
    result, _ := self.Names[flag]
    return result
}

func (self Flags) Values() []string {
    result := make([]string, 0, len(self.Names))
    for k, _ := range self.Names {
       result = append(result, k)
    }
    return result
}


func ParseArray_FolderEntry(profile *FATProfile, reader io.ReaderAt, offset int64, count int) []*FolderEntry {
    result := make([]*FolderEntry, 0, count)
    for i:=0; i<count; i++ {
      value := profile.FolderEntry(reader, offset)
      result = append(result, value)
      offset += int64(value.Size())
    }
    return result
}

func ParseInt8(reader io.ReaderAt, offset int64) int8 {
    result := make([]byte, 1)
    _, err := reader.ReadAt(result, offset)
    if err != nil {
       return 0
    }
    return int8(result[0])
}

func ParseUint16(reader io.ReaderAt, offset int64) uint16 {
    data := make([]byte, 2)
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return binary.LittleEndian.Uint16(data)
}

func ParseUint32(reader io.ReaderAt, offset int64) uint32 {
    data := make([]byte, 4)
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return binary.LittleEndian.Uint32(data)
}

func ParseUint8(reader io.ReaderAt, offset int64) byte {
    result := make([]byte, 1)
    _, err := reader.ReadAt(result, offset)
    if err != nil {
       return 0
    }
    return result[0]
}

func ParseTerminatedString(reader io.ReaderAt, offset int64) string {
   data := make([]byte, 1024)
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
     return ""
   }
   idx := bytes.Index(data[:n], []byte{0})
   if idx < 0 {
      idx = n
   }
   return string(data[0:idx])
}

func ParseString(reader io.ReaderAt, offset int64, length int64) string {
   data := make([]byte, length)
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
      return ""
   }
   return string(data[:n])
}


func ParseTerminatedUTF16String(reader io.ReaderAt, offset int64) string {
   data := make([]byte, 1024)
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
     return ""
   }

   idx := bytes.Index(data[:n], []byte{0, 0})
   if idx < 0 {
      idx = n-1
   }
   if idx%2 != 0 {
      idx += 1
   }
   return UTF16BytesToUTF8(data[0:idx], binary.LittleEndian)
}

func ParseUTF16String(reader io.ReaderAt, offset int64, length int64) string {
   data := make([]byte, length)
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
     return ""
   }
   return UTF16BytesToUTF8(data[:n], binary.LittleEndian)
}

func UTF16BytesToUTF8(b []byte, o binary.ByteOrder) string {
	if len(b) < 2 {
		return ""
	}

	if b[0] == 0xff && b[1] == 0xfe {
		o = binary.BigEndian
		b = b[2:]
	} else if b[0] == 0xfe && b[1] == 0xff {
		o = binary.LittleEndian
		b = b[2:]
	}

	utf := make([]uint16, (len(b)+(2-1))/2)

	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = o.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}

	return string(utf16.Decode(utf))
}

