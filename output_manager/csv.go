package output_manager

import (
	"encoding/csv"
	"strings"
)

//func (dc *OutputManager) LoadPathDataAsCSV() error {
//
//	err1 := dc.LoadPathFileDescriptor()
//	if err1 != nil {
//		return err1
//	}
//
//	// Wrap in CSV
//	var err2 error
//	csv_reader := csv.NewReader(dc.f)
//	dc.rows_list_of_lists, err2 = csv_reader.ReadAll()
//
//	if err2 != nil {
//		return err2
//	}
//
//	if len(dc.rows_list_of_lists) >= 2 {
//		// Headers split
//		dc.headers_list = dc.rows_list_of_lists[0]
//
//		// Remove header from rows
//		dc.rows_list_of_lists = dc.rows_list_of_lists[1:]
//	}
//
//	return nil
//}

func (dc *OutputManager) CreateFileForCSVWriting() error {
	err := dc.CreateEmptyFilesAndDirectories()

	if err != nil {
		return err
	}

	dc.csv_writer = csv.NewWriter(dc.f)
	return nil
}

//func (dc *OutputManager) LoadPathDataAsCSV() error {
//
//	err1 := dc.LoadPathFileDescriptor()
//	if err1 != nil {
//		return err1
//	}
//
//	// Wrap in CSV
//	var err2 error
//	csv_reader := csv.NewReader(dc.f)
//	dc.rows_list_of_lists, err2 = csv_reader.ReadAll()
//
//	if err2 != nil {
//		return err2
//	}
//
//	if len(dc.rows_list_of_lists) >= 2 {
//		// Headers split
//		dc.headers_list = dc.rows_list_of_lists[0]
//
//		// Remove header from rows
//		dc.rows_list_of_lists = dc.rows_list_of_lists[1:]
//	}
//
//	return nil
//}

func disarm_newlines(s string) string {
	s1 := strings.ReplaceAll(s, "\r", "\\r")
	return strings.ReplaceAll(s1, "\n", "\\n")
}

func apply_disarm_newlines_on_list(l []string) []string {
	for i := 0; i < len(l); i++ {
		l[i] = disarm_newlines(l[i])
	}

	return l
}

func (dc *OutputManager) SaveAllDataToCSVFormat() error {

	// Initalize CSV Writer
	if err1 := dc.CreateFileForCSVWriting(); err1 != nil {
		return err1
	}

	if len(dc.headers_list) == 0 {
		panic("Missing CSV headers!")
	}

	// Write headers
	if err2 := dc.csv_writer.Write(apply_disarm_newlines_on_list(dc.headers_list)); err2 != nil {
		return err2
	}

	// Write rows
	for i := 0; i < len(dc.rows_list_of_lists); i++ {

		if err3 := dc.csv_writer.Write(apply_disarm_newlines_on_list(dc.rows_list_of_lists[i])); err3 != nil {
			return err3
		}
	}

	// Flush
	dc.csv_writer.Flush()

	return nil
}
