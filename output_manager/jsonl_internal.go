package output_manager

import (
	"bufio"
	"compress/gzip"
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/common"
	"sync"
)

func (dc *OutputManager) CreateFileForJSONLWriting(compress bool) error {
	err := dc.CreateEmptyFilesAndDirectories()

	if err != nil {
		return err
	}

	// BUFFERING
	dc.f_buffered_writer = bufio.NewWriter(dc.f)

	// GZIP
	var err2 error
	if compress {
		dc.f_buffered_writer_gzip_outer, err2 = gzip.NewWriterLevel(dc.f_buffered_writer, gzip.BestSpeed)
	}

	if err2 != nil {
		return err2
	}

	return nil
}

func (dc *OutputManager) AppendJSONLDataToFile(od *ordereddict.Dict) error {

	buf, err := od.MarshalJSON()

	// Append \r\n
	buf = append(buf, 13, 10)

	if err != nil {
		return err
	}

	var err2 error
	// GZIP Case
	if dc.f_buffered_writer_gzip_outer != nil {
		_, err2 = dc.f_buffered_writer_gzip_outer.Write(buf)
	} else { // BUFFERED
		_, err2 = dc.f_buffered_writer.Write(buf)
	}

	if err2 != nil {
		return err2
	}

	return nil
}

func (dc *OutputManager) LoadPathDataAsJSONL(compressed bool) error {

	err1 := dc.LoadPathFileDescriptor()
	if err1 != nil {
		return err1
	}

	var gzip_reader *gzip.Reader = nil
	var err2 error
	var scanner *bufio.Scanner = nil

	if compressed {
		gzip_reader, err2 = gzip.NewReader(dc.f)

		if err2 != nil {
			return err2
		}
		defer gzip_reader.Close()
		scanner = bufio.NewScanner(gzip_reader)
	} else {
		scanner = bufio.NewScanner(dc.f)
	}
	defer dc.f.Close()

	var wg_scanner sync.WaitGroup
	var wg_unmarshal sync.WaitGroup
	ch := make(chan []byte, 10000)

	// Start scanner
	wg_scanner.Add(1)
	go func() {
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			text := scanner.Bytes()

			if len(text) > 0 {
				topass := make([]byte, len(text))
				copy(topass, text)
				ch <- topass
			}
		}
		wg_scanner.Done()
	}()

	// Start unmarshall
	wg_unmarshal.Add(1)
	go func() {
		for text := range ch {
			od := ordereddict.NewDict()
			err_json := od.UnmarshalJSON(text)

			if err_json != nil {
				panic("Error when unmarshalling ongoing file")
			}

			// First time assign headers
			if len(dc.headers_list) == 0 {
				dc.headers_list = od.Keys()
			}

			// Append values
			dc.rows_list_of_lists = append(dc.rows_list_of_lists, common.OrderedDictToOrderedStringListValues(od))
		}
		wg_unmarshal.Done()
	}()

	// End action
	wg_scanner.Wait()
	close(ch)
	wg_unmarshal.Wait()

	return nil
}

func (dc *OutputManager) SaveAllDataToJSONLFormat() error {

	// Initalize Buffered Writer
	if err1 := dc.CreateFileForJSONLWriting(false); err1 != nil {
		return err1
	}

	if len(dc.headers_list) == 0 {
		panic("Missing JSONL headers!")
	}

	// Write rows
	for i := 0; i < len(dc.rows_list_of_lists); i++ {
		row := dc.rows_list_of_lists[i]

		final_od := common.HeadersAndRowListToOrderedDict(dc.headers_list, row)
		if err2 := dc.AppendJSONLDataToFile(final_od); err2 != nil {
			return err2
		}
	}

	// Flush
	dc.f_buffered_writer.Flush()

	return nil
}
