package output_manager

import (
	"bufio"
	"encoding/json"
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/evtxhussar/common"
)

func (dc *OutputManager) CreateFileForJSONWriting() error {
	err := dc.CreateEmptyFilesAndDirectories()

	if err != nil {
		return err
	}

	// BUFFERING
	dc.f_buffered_writer = bufio.NewWriter(dc.f)

	return nil
}

func (dc *OutputManager) SaveAllDataToJSONFormat() error {

	// Initalize Buffered Writer
	if err1 := dc.CreateFileForJSONWriting(); err1 != nil {
		return err1
	}

	if len(dc.headers_list) == 0 {
		panic("Missing JSON headers!")
	}

	// Write rows
	mega_list := make([]*ordereddict.Dict, 0, len(dc.rows_list_of_lists))

	for i := 0; i < len(dc.rows_list_of_lists); i++ {
		row := dc.rows_list_of_lists[i]

		final_od := common.HeadersAndRowListToOrderedDict(dc.headers_list, row)
		mega_list = append(mega_list, final_od)
	}

	// Write at once
	var bytes_serialized []byte
	var err2 error = nil

	if bytes_serialized, err2 = json.MarshalIndent(mega_list, "", " "); err2 != nil {
		return err2
	}

	if _, err3 := dc.f_buffered_writer.Write(bytes_serialized); err3 != nil {
		return err3
	}

	// Flush
	dc.f_buffered_writer.Flush()

	return nil
}
