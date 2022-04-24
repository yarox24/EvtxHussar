package output_manager

import (
	"github.com/Velocidex/ordereddict"
	"github.com/xuri/excelize/v2"
	"github.com/yarox24/EvtxHussar/common"
)

func (dc *OutputManager) CreateFileForExcelWriting() error {

	// Directory creation
	dir_error := common.EnsureDirectoryStructureIsCreated(dc.path)

	if dir_error != nil {
		return dir_error
	}

	dc.excel_file = excelize.NewFile()

	// NewSheet
	dc.excel_file.SetSheetName("Sheet1", dc.excel_sheet_name)

	return nil
}

func (dc *OutputManager) SetExcelSheetName(excel_sheet_name string) {
	dc.excel_sheet_name = excel_sheet_name
}

func (dc *OutputManager) AnalyzeColumnDataWidthAndEmptiness() { // wg_analyzer *sync.WaitGroup

	// Create analyzed structure
	dc.column_lengths = ordereddict.NewDict()
	//dc.column_emptiness = ordereddict.NewDict()

	// Set initial values based on length of headers

	for i := 0; i < len(dc.headers_list); i++ {
		val := dc.headers_list[i]
		dc.column_lengths.Set(val, int64(len(val)))

		// Assume column is empty
		//dc.column_emptiness.Set(val, true)
	}

	// Analyze all rows
	for i := 0; i < len(dc.rows_list_of_lists); i++ {
		row := dc.rows_list_of_lists[i]

		for col, val := range row {
			col_name := dc.headers_list[col]

			// Check length
			record_length, _ := dc.column_lengths.GetInt64(col_name)
			current_length := int64(len(val))

			if current_length > record_length && current_length <= 100 {
				dc.column_lengths.Update(col_name, int64(current_length))
			}

			// Check if is empty
			//if current_length > 0 {
			//	dc.column_emptiness.Update(col_name, false)
			//}
		}
	}

	//wg_analyzer.Done()
}

func ToInterfaceSliceWithStyle(ss []string, style_id int) []interface{} {
	iface := make([]interface{}, len(ss))

	for i := range ss {
		iface[i] = excelize.Cell{StyleID: style_id, Value: ss[i]}
	}

	return iface
}

func ToInterfaceSlice(ss []string) []interface{} {
	iface := make([]interface{}, len(ss))

	for i := range ss {
		iface[i] = ss[i]
	}

	return iface
}

func (dc *OutputManager) SaveAllDataToExcelFormatStreaming() error {

	// Ensure directory
	if err1 := dc.CreateFileForExcelWriting(); err1 != nil {
		return err1
	}

	stream_writer, err2 := dc.excel_file.NewStreamWriter(dc.excel_sheet_name)

	if err2 != nil {
		return err2
	}

	if len(dc.headers_list) == 0 {
		panic("Missing Excel headers!")
	}

	// Analyze data in extra goroutine
	//var wg_analyzer sync.WaitGroup
	//wg_analyzer.Add(1)
	dc.AnalyzeColumnDataWidthAndEmptiness() // &wg_analyzer

	// Before writing rows #Set auto column width - Analyzer dependent
	//wg_analyzer.Wait()

	for col_ind, col_name := range dc.column_lengths.Keys() {
		//column_letter, _ := excelize.ColumnNumberToName(col_ind + 1)
		column_len, _ := dc.column_lengths.GetInt64(col_name)
		stream_writer.SetColWidth(col_ind+1, col_ind+1, float64(column_len+2))
		//dc.excel_file.SetColWidth(dc.excel_sheet_name, column_letter, column_letter, float64(column_len+2))
	}

	// Bold style
	bold_style_id, _ := dc.excel_file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})

	// Write bold headers
	for i := 0; i < len(dc.headers_list); i++ {
		err_header := stream_writer.SetRow("A1", ToInterfaceSliceWithStyle(dc.headers_list, bold_style_id))

		if err_header != nil {
			return err_header
		}
	}

	// Write rows
	for i := 0; i < len(dc.rows_list_of_lists); i++ {
		row_interface := ToInterfaceSlice(dc.rows_list_of_lists[i])
		axis, _ := excelize.CoordinatesToCellName(1, i+2)

		err_header := stream_writer.SetRow(axis, row_interface)

		if err_header != nil {
			return err_header
		}
	}

	// Set auto-filter
	if len(dc.rows_list_of_lists) > 0 {
		lowest_right_cell, _ := excelize.CoordinatesToCellName(len(dc.headers_list), len(dc.rows_list_of_lists)+1)
		dc.excel_file.AutoFilter(dc.excel_sheet_name, "A1", lowest_right_cell, "")
	}

	// Flush
	stream_writer.Flush()

	if err2 := dc.excel_file.SaveAs(dc.path); err2 != nil {
		return err2
	}

	// End of writing part
	dc.excel_file.Close()
	dc.CloseAndFlushFileDescriptors()
	common.LogDebug("Closing streamed Excel file: " + dc.path)

	// Extra enhance part - Generated excel recovery attempt error when opening

	//common.LogDebug("Column analysis finished" + dc.path)
	//
	//var err4 error
	//dc.excel_file, err4 = excelize.OpenFile(dc.path)
	//x := dc.excel_file.GetSheetList()
	//
	//fmt.Println(x)
	//if err4 != nil {
	//	return err4
	//}

	// Hide empty columns - Analyzer dependent
	//if len(dc.rows_list_of_lists) > 0 {
	//	for col_ind, col_name := range dc.column_emptiness.Keys() {
	//		is_empty, _ := dc.column_emptiness.GetBool(col_name)
	//		col_letter, _ := excelize.ColumnNumberToName(col_ind + 1)
	//		if is_empty {
	//			dc.excel_file.SetColVisible(dc.excel_sheet_name, col_letter, false)
	//		}
	//	}
	//}

	//if err := dc.excel_file.Save(); err != nil {
	//	return err
	//}
	//if err := dc.excel_file.Close(); err != nil {
	//	return err
	//}

	return nil
}

//func (dc *OutputManager) SaveAllDataToExcelFormat() error {
//
//	// Ensure directory
//	if err1 := dc.CreateFileForExcelWriting(); err1 != nil {
//		return err1
//	}
//
//	if len(dc.headers_list) == 0 {
//		panic("Missing Excel headers!")
//	}
//
//	// Bold style
//	bold_style, _ := dc.excel_file.NewStyle(&excelize.Style{
//		Font: &excelize.Font{
//			Bold: true,
//		},
//	})
//
//	// Write bold headers
//	for i := 0; i < len(dc.headers_list); i++ {
//		axis, _ := excelize.CoordinatesToCellName(i+1, 1)
//		if err_cell := dc.excel_file.SetCellStr(dc.excel_sheet_name, axis, dc.headers_list[i]); err_cell != nil {
//			return err_cell
//		}
//
//		// Bold
//		dc.excel_file.SetCellStyle(dc.excel_sheet_name, axis, axis, bold_style)
//	}
//
//	// Analyze data in extra goroutine
//	var wg_analyzer sync.WaitGroup
//	wg_analyzer.Add(1)
//	go dc.AnalyzeColumnDataWidthAndEmptiness(&wg_analyzer)
//
//	// Write rows
//	for i := 0; i < len(dc.rows_list_of_lists); i++ {
//		row := dc.rows_list_of_lists[i]
//		axis, _ := excelize.CoordinatesToCellName(1, i+2)
//		dc.excel_file.SetSheetRow(dc.excel_sheet_name, axis, &row)
//	}
//
//	// Set auto-filter
//	if len(dc.rows_list_of_lists) > 0 {
//		lowest_right_cell, _ := excelize.CoordinatesToCellName(len(dc.headers_list), len(dc.rows_list_of_lists)+1)
//		dc.excel_file.AutoFilter(dc.excel_sheet_name, "A1", lowest_right_cell, "")
//	}
//
//	// Wait for analyzer
//	wg_analyzer.Wait()
//	// Set auto column width - Analyzer dependent
//	for col_ind, col_name := range dc.column_lengths.Keys() {
//		column_letter, _ := excelize.ColumnNumberToName(col_ind + 1)
//		column_len, _ := dc.column_lengths.GetInt64(col_name)
//		dc.excel_file.SetColWidth(dc.excel_sheet_name, column_letter, column_letter, float64(column_len+2))
//	}
//
//	// Hide empty columns - Analyzer dependent
//	if len(dc.rows_list_of_lists) > 0 {
//		for col_ind, col_name := range dc.column_emptiness.Keys() {
//			is_empty, _ := dc.column_emptiness.GetBool(col_name)
//			col_letter, _ := excelize.ColumnNumberToName(col_ind + 1)
//			if is_empty {
//				dc.excel_file.SetColVisible(dc.excel_sheet_name, col_letter, false)
//			}
//		}
//	}
//
//	if err2 := dc.excel_file.SaveAs(dc.path); err2 != nil {
//		return err2
//	}
//
//	return nil
//}
