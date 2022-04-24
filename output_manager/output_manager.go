package output_manager

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"github.com/Velocidex/ordereddict"
	"github.com/xuri/excelize/v2"
	"github.com/yarox24/evtxhussar/common"
	"golang.org/x/exp/slices"
	"os"
	"sort"
	"strings"
	"sync"
)

type OutputManager struct {
	headers_list                 []string
	rows_list_of_lists           [][]string
	csv_writer                   *csv.Writer
	path                         string
	f                            *os.File
	f_buffered_writer            *bufio.Writer
	f_buffered_writer_gzip_outer *gzip.Writer
	excel_file                   *excelize.File
	excel_sheet_name             string
	column_lengths               *ordereddict.Dict
	column_emptiness             *ordereddict.Dict
	wg_analyzer                  *sync.WaitGroup
}

func NewOutputManager() OutputManager {

	// Wrap in CSV
	return OutputManager{
		headers_list:                 make([]string, 0),
		rows_list_of_lists:           make([][]string, 0),
		csv_writer:                   nil,
		path:                         "",
		f:                            nil,
		f_buffered_writer:            nil,
		f_buffered_writer_gzip_outer: nil,
		excel_file:                   nil,
		excel_sheet_name:             "Unnamed",
		column_lengths:               nil,
		column_emptiness:             nil,
		wg_analyzer:                  new(sync.WaitGroup),
	}
}

func (dc *OutputManager) SetPath(path string) {
	dc.path = path
}

func (dc *OutputManager) LoadPathFileDescriptor() error {
	// File open
	var err error
	dc.f, err = os.Open(dc.path)

	if err != nil {
		return err
	}

	return nil
}

func (dc *OutputManager) CreateEmptyFilesAndDirectories() error {

	// Directory creation
	dir_error := common.EnsureDirectoryStructureIsCreated(dc.path)

	if dir_error != nil {
		return dir_error
	}

	// File creation
	var err error
	dc.f, err = os.Create(dc.path)

	if err != nil {
		return err
	}

	return nil
}

func (dc *OutputManager) CloseAndFlushFileDescriptors() {
	// First close specific writers
	if dc.csv_writer != nil {
		dc.csv_writer.Flush()
	}

	// Gzip
	if dc.f_buffered_writer_gzip_outer != nil {
		dc.f_buffered_writer_gzip_outer.Close()
	}

	// Buffered
	if dc.f_buffered_writer != nil {
		dc.f_buffered_writer.Flush()
	}

	// At the end close file
	if dc.f != nil {
		dc.f.Close()
		dc.f = nil
	}
}

func (dc *OutputManager) SortByAllColumnsAssumingFirstIsDate(newest_first bool) {
	if len(dc.headers_list) > 0 && len(dc.rows_list_of_lists) > 2 {
		if newest_first {
			sort.Slice(dc.rows_list_of_lists, func(i, j int) bool {
				return strings.Join(dc.rows_list_of_lists[i], ".") > strings.Join(dc.rows_list_of_lists[j], ".")
			})
		} else {
			sort.Slice(dc.rows_list_of_lists, func(i, j int) bool {
				return strings.Join(dc.rows_list_of_lists[i], ".") < strings.Join(dc.rows_list_of_lists[j], ".")
			})
		}
	}
}

func equal(v1 []string, v2 []string) bool {
	return strings.Join(v1, "") == strings.Join(v2, "")
}

func (dc *OutputManager) RemoveDuplicatesAssumingSorted() {
	if len(dc.headers_list) > 0 && len(dc.rows_list_of_lists) >= 2 {
		dc.rows_list_of_lists = slices.CompactFunc(dc.rows_list_of_lists, equal)
	}
}

func (dc *OutputManager) SaveAllDataToProperFormat(format_name string) error {
	switch format_name {
	case "csv":
		{
			return dc.SaveAllDataToCSVFormat()
		}
	case "json":
		{
			return dc.SaveAllDataToJSONFormat()
		}
	case "jsonl":
		{
			return dc.SaveAllDataToJSONLFormat()
		}
	case "excel":
		{
			return dc.SaveAllDataToExcelFormatStreaming()
		}
	default:
		panic("Wrong format name")
	}

	return nil
}

func (dc *OutputManager) SetHeadersIfNotSet(list []string) {
	if len(dc.headers_list) == 0 {
		dc.headers_list = list
	}
}
