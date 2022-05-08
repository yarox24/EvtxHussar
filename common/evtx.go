package common

import (
	"encoding/binary"
	"github.com/Velocidex/ordereddict"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"www.velocidex.com/golang/evtx"
)

type EvtxFileInfo struct {
	path               string
	channel            string
	latest_computer    string
	is_empty           bool
	is_valid           bool
	will_be_processed  bool
	record_counter     int64
	header             evtx.EVTXHeader
	alternative_header EVTXHeaderAlternative
}

type EVTXHeaderAlternative struct {
	Magic           [8]byte
	OldestChunk     uint64
	CurrentChunkNum uint64
	NextRecordNum   uint64
	HeaderPart1Len  uint32
	MinorVersion    uint16
	MajorVersion    uint16
	HeaderSize      uint16
	ChunkCount      uint16
	_               [76]byte
	FileFlags       uint32
	CheckSum        uint32
}

func NewEvtxFileInfo(p string) EvtxFileInfo {
	return EvtxFileInfo{
		path:              p,
		channel:           "",
		latest_computer:   "",
		is_empty:          true,
		is_valid:          false,
		will_be_processed: false,
		record_counter:    0,
		header:            evtx.EVTXHeader{},
	}
}

func (efi *EvtxFileInfo) GetPath() string {
	return efi.path
}

func (efi *EvtxFileInfo) GetFilenameWithoutExtension() string {
	return strings.TrimSuffix(filepath.Base(efi.path), filepath.Ext(efi.path))
}

func (efi *EvtxFileInfo) GetChannel() string {
	return efi.channel
}

func (efi *EvtxFileInfo) GetLatestComputer() string {
	return efi.latest_computer
}

func (efi *EvtxFileInfo) IsEmpty() bool {
	return efi.is_empty
}

func (efi *EvtxFileInfo) GetNumberOfRecords() int64 {
	return efi.record_counter
}

func (efi *EvtxFileInfo) SetNumberOfRecords(record_counter int64) {
	efi.record_counter = record_counter
}

func (efi *EvtxFileInfo) IsValid() bool {
	return efi.is_valid
}

func (efi *EvtxFileInfo) WillBeProcessed() bool {
	return efi.will_be_processed
}

func (efi *EvtxFileInfo) EnableForProcessing() {
	efi.will_be_processed = true
}

func (efi *EvtxFileInfo) GetAlternativeHeader() *EVTXHeaderAlternative {
	return &efi.alternative_header
}

func readStructFromFile2(fd io.ReadSeeker, offset int64, obj interface{}) error {
	_, err := fd.Seek(offset, os.SEEK_SET)
	if err != nil {
		return errors.Wrap(err, "Seek")
	}

	err = binary.Read(fd, binary.LittleEndian, obj)
	if err != nil {
		return errors.Wrap(err, "Read")
	}

	return nil
}

//func GetSystemDate(rec *evtx.EventRecord) time.Time {
//	var ft FileTime
//	ft.LowDateTime = uint32(rec.Header.FileTime)
//	ft.HighDateTime = uint32(rec.Header.FileTime >> 32)
//	return ft.Time()
//}

func NewChunk(fd io.ReadSeeker, offset int64) (*evtx.Chunk, error) {
	self := &evtx.Chunk{Offset: offset, Fd: fd}
	_, err := fd.Seek(offset, os.SEEK_SET)
	if err != nil {
		return nil, errors.Wrap(err, "Seek")
	}

	err = binary.Read(fd, binary.LittleEndian, &self.Header)
	return self, errors.WithStack(err)
}

func (efi *EvtxFileInfo) DetermineParameters() {

	// Open file
	fd, err := os.OpenFile(efi.path, os.O_RDONLY, os.FileMode(0666))

	if err == nil {
		defer fd.Close()
	} else {
		LogCriticalErrorWithError("Error occured when opening evtx: "+efi.path, err)
		return
	}

	// Library specific
	header_block_size := int64(efi.header.HeaderBlockSize)

	// Grab last record from last chunk
	records_final := []*evtx.EventRecord{}
	last_chunk := &evtx.Chunk{}

	for offset := header_block_size + int64(efi.header.LastChunk)*evtx.EVTX_CHUNK_SIZE; offset >= header_block_size; offset -= evtx.EVTX_CHUNK_SIZE {
		chunk, err := NewChunk(fd, offset)

		if err != nil {
			continue
		}

		if string(chunk.Header.Magic[:]) == evtx.EVTX_CHUNK_HEADER_MAGIC {
			last_chunk = chunk
			break
		}
	}

	// Found anything?
	if string(last_chunk.Header.Magic[:]) != evtx.EVTX_CHUNK_HEADER_MAGIC {
		return
	}

	records, err := last_chunk.Parse(0)
	records_final = records
	if err != nil {
		return
	}

	// Check if file is empty
	if len(records_final) == 0 {
		return
	}

	var last_record *evtx.EventRecord = records_final[len(records_final)-1]

	event_map, event_ok := last_record.Event.(*ordereddict.Dict)

	if event_ok {
		event, ok := ordereddict.GetMap(event_map, "Event")

		if ok {

			// Latest computer
			latest_computer, ok_computer := ordereddict.GetString(event, "System.Computer")

			if ok_computer {
				efi.latest_computer = latest_computer
			}

			// Channel
			channel, ok_channel := ordereddict.GetString(event, "System.Channel")
			if ok_channel {
				efi.channel = channel
			}

			if !ok_channel || !ok_computer {
				LogError("Cannot determine channel and hostname for: " + efi.path)
				efi.is_valid = false
				return
			}
		} else {
			LogError("Cannot get EventMap for: " + efi.path)
			return
		}
	} else {
		LogError("Cannot convert Event to OrderedDict: " + efi.path)
		return
	}

	// File is not empty
	efi.is_empty = false

}

//func (efi *EvtxFileInfo) Print() {
//	log.Printf("[%s] Valid: %t | Channel: %s | Latest computer: %s | Empty: %t | Is supported? %t", efi.path, efi.IsValid(), efi.GetChannel(), efi.GetLatestComputer(), efi.IsEmpty(), efi.WillBeProcessed())
//}

func Inspect_evtx_paths(EfiList []EvtxFileInfo) []EvtxFileInfo {
	var wg sync.WaitGroup

	INSPECTORS_LIMIT := 80
	LogDebugString("INSPECTORS_LIMIT", INSPECTORS_LIMIT)
	var inspector_channel chan struct{} = make(chan struct{}, INSPECTORS_LIMIT) // Create channel for limiting concurrent open files

	// Pass efi elements
	for i, _ := range EfiList {
		inspector_channel <- struct{}{}
		wg.Add(1)
		go Inspect_single_evtx(&wg, &EfiList[i], inspector_channel)
	}

	wg.Wait()
	close(inspector_channel)

	// Statistics
	LogInfo("Finished inspecting")
	return EfiList
}

func Inspect_single_evtx(wg *sync.WaitGroup, efi *EvtxFileInfo, inspector_channel chan struct{}) {

	// Assume
	efi.is_empty = true
	efi.is_valid = false

	// Initial validation
	err := efi.Validate()

	if efi.is_valid {

		// Check if is empty
		if !(efi.header.Firstchunk == 0 && efi.header.LastChunk == 0 && efi.header.NextRecordID == 1) {
			efi.DetermineParameters()
		}

	} else {
		LogErrorWithError("Invalid .evtx file: ", err)
	}

	<-inspector_channel
	wg.Done()
}

func is_supported(minor, major uint16) bool {
	switch major {
	case 3:
		switch minor {
		case 0, 1, 2:
			return true
		}
	}
	return false
}

/* If the log is dirty, patch the log header with the values from the EOF record */
// ReactOS
//if (!LogFile->ReadOnly && IsLogDirty)
//{
//LogFile->Header.StartOffset = EofRec.BeginRecord;
//LogFile->Header.EndOffset   = EofRec.EndRecord;
//LogFile->Header.CurrentRecordNumber = EofRec.CurrentRecordNumber;
//LogFile->Header.OldestRecordNumber  = EofRec.OldestRecordNumber;
//}

func (efi *EvtxFileInfo) Validate() error {

	fstat, err := os.Stat(efi.path)

	if err != nil {
		return err
	}

	if fstat.Size() < 129 {
		return errors.New("File is smaller than 129 bytes")
	}

	//chunks, err := evtx.GetChunks(*parse_file)
	fd, err := os.OpenFile(efi.path, os.O_RDONLY, os.FileMode(0666))

	if err == nil {
		defer fd.Close()
	} else {
		return err
	}

	// Library specific validation
	err2 := readStructFromFile2(fd, 0, &efi.header)
	if err2 != nil {
		return errors.New("Invalid Evtx header")
	}
	err2 = readStructFromFile2(fd, 0, &efi.alternative_header)
	if err2 != nil {
		return errors.New("Invalid Evtx alternative header")
	}

	//MarshalIndent
	//hJSON, _ := json.MarshalIndent(efi.alternative_header, "", "  ")
	//LogDebug(fmt.Sprint("Analyzed file: %s:\n#v\n", efi.path, string(hJSON)))

	if string(efi.header.Magic[:]) != evtx.EVTX_HEADER_MAGIC {
		return errors.New("File is not an EVTX file (wrong magic).")
	}

	if !is_supported(efi.header.MinorVersion, efi.header.MajorVersion) {
		return errors.New("Unsupported EVTX version.")
	}

	efi.is_valid = true
	return nil
}

func ReturnCopyOfSupportedEfiFileInfoElements(efi []EvtxFileInfo) []EvtxFileInfo {
	var temp = make([]EvtxFileInfo, 0)

	for i := 0; i < len(efi); i++ {
		if efi[i].will_be_processed {
			temp = append(temp, efi[i])
		}
	}

	return temp
}

func Generate_list_of_files_to_process(Input_evtx_paths []string, Recursive bool) []EvtxFileInfo {
	var EfiList []EvtxFileInfo = make([]EvtxFileInfo, 0, 200)

	// Deduplicate by path
	Input_evtx_paths = UniqueElementsOfSliceInMemory(Input_evtx_paths)

	if len(Input_evtx_paths) > 0 {
		for _, p := range Input_evtx_paths {
			t, err := determine_path(p)

			if err != nil {
				LogCriticalErrorWithError("This file/dir will be ignored: "+p, err)
			} else {

				// Append file
				if t == IS_FILE {
					evtx_path, _ := filepath.Abs(p)
					EfiList = append(EfiList, NewEvtxFileInfo(evtx_path))
				} else if t == IS_DIR {
					for _, p_found := range find_logs_in_directory(p, Recursive) {
						evtx_path, _ := filepath.Abs(p_found)
						EfiList = append(EfiList, NewEvtxFileInfo(evtx_path))
					}
				}
			}
		}

	}

	return EfiList
}
