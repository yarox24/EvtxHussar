package engine

import (
	"bufio"
	"fmt"
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/common"
	"github.com/yarox24/EvtxHussar/eventmap"
	"github.com/yarox24/EvtxHussar/output_manager"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const CHANNEL_EVENT_CAPACITY = 300

type Layer2Info struct {
	Typ  string
	Name string
}

type Layer2Output struct {
	Typ                   string
	Category              string
	Subcategory           string
	Filename              string
	GlobalOutputDirectory string
}

type Layer2 struct {
	Info                        Layer2Info
	Output                      Layer2Output
	Aggregation_options         map[string]string
	Fields_remap                map[string]string
	Fields_remap_dict           *ordereddict.Dict
	Field_extra_transformations []common.Layer2FieldExtraTransformations
	Ordered_fields              []string
	Ordered_fields_enhanced     map[string]common.SingleField
	UsageCounter                int
}

type Layer2SingleLayer struct {
	l2_name             string
	latest_computer     string
	ch                  chan *eventmap.EventMap
	wg_to_close_channel *sync.WaitGroup
	parsing_info        *Layer2
	engine              *Engine
	workers_counter     int32
}

func NewLayer2SingleLayer(l2_name string, latest_computer string, engine *Engine) Layer2SingleLayer {
	return Layer2SingleLayer{
		l2_name:             l2_name,
		latest_computer:     latest_computer,
		ch:                  make(chan *eventmap.EventMap, CHANNEL_EVENT_CAPACITY),
		wg_to_close_channel: new(sync.WaitGroup),
		parsing_info:        engine.FindL2LayerByName(l2_name),
		engine:              engine,
		workers_counter:     0,
	}
}

func NewLayer2GlobalMemory() Layer2GlobalMemory {
	return Layer2GlobalMemory{
		Layers:    make([]Layer2SingleLayer, 0),
		Wg_l2_all: new(sync.WaitGroup),
	}
}

type Layer2GlobalMemory struct {
	Layers    []Layer2SingleLayer
	Wg_l2_all *sync.WaitGroup
}

func (l2mem *Layer2GlobalMemory) SetupChannels(engi *Engine, efi []common.EvtxFileInfo) {

	// For every file
	for i := 0; i < len(efi); i++ {
		latest_computer := efi[i].GetLatestComputer()
		channel := efi[i].GetChannel()

		// Find all layers 1 which support this channel
		for _, l1 := range engi.GetAllLayer1WhichSupportsChannel(channel) {
			sendto_layer2 := l1.Sendto_layer2
			l2mem.SetupLayer2SingleLayer(sendto_layer2, latest_computer, engi)

		}
	}
}

func (l2mem *Layer2GlobalMemory) FindLayer2SingleLayer(l2_name string, latest_computer string) *Layer2SingleLayer {
	for i := 0; i < len(l2mem.Layers); i++ {
		var l2_layer *Layer2SingleLayer = &l2mem.Layers[i]
		if strings.ToLower(l2_layer.l2_name) == strings.ToLower(l2_name) &&
			strings.ToLower(l2_layer.latest_computer) == strings.ToLower(latest_computer) {
			return l2_layer
		}
	}
	return nil
}

func (l2mem *Layer2GlobalMemory) SetupLayer2SingleLayer(l2_name string, latest_computer string, engine *Engine) {

	// Check if L2 layer already exists
	if l2mem.FindLayer2SingleLayer(l2_name, latest_computer) != nil {
		return
	}

	common.LogDebug(fmt.Sprintf("New L2: %s %s", l2_name, latest_computer))
	l2mem.Layers = append(l2mem.Layers, NewLayer2SingleLayer(l2_name, latest_computer, engine))
}

func RunL2WorkerFlat(l2s *Layer2SingleLayer) {

	// Intermediate form = JSONL
	ongoing_path := l2s.GetOngoingOutputPath()
	final_path := l2s.GetOutputPath()

	//New JSONL unsorted
	om1 := output_manager.NewOutputManager()
	om1.SetPath(ongoing_path)
	jsonl_err := om1.CreateFileForJSONLWriting(true)

	if jsonl_err != nil {
		common.LogCriticalErrorWithError("Cannot create ongoing file", jsonl_err)
	}

	for ev_map := range l2s.ch {

		// Common fields
		final_od := l2s.engine.ParseCommonFieldsOrderedDict(ev_map)

		//eid, _ := final_od.GetString("EID")
		//if eid == "91" {
		//	fmt.Println("Break")
		//}

		// Layer 2 fields
		l2_od := l2s.engine.ParseL2FieldsOrderedDict(l2s.l2_name, ev_map)

		// Combine both
		final_od.MergeFrom(l2_od)

		err := om1.AppendJSONLDataToFile(final_od)

		if err != nil {
			common.LogCriticalErrorWithError("Cannot write event to ongoing file", err)
		}
	}

	om1.CloseAndFlushFileDescriptors()
	common.LogDebug(fmt.Sprintf("Ongoing result saved to file : %s | %s | %d", l2s.l2_name, l2s.latest_computer, l2s.workers_counter))

	// Reading
	om2 := output_manager.NewOutputManager()
	om2.SetPath(ongoing_path)
	common.LogDebug(fmt.Sprintf("[Start] Loading ongoing file : %s", ongoing_path))
	om2.LoadPathDataAsJSONL(true)
	common.LogDebug(fmt.Sprintf("[End] Loading ongoing file : %s", ongoing_path))

	// Artifical headers (If not event is catched)
	m1 := l2s.engine.PrepareCommonFieldsEmptyOrderedDict()
	m2 := l2s.engine.PrepareLayer2FieldsEmptyOrderedDict(l2s.l2_name)
	m1.MergeFrom(m2)
	om2.SetHeadersIfNotSet(m1.Keys())

	om2.CloseAndFlushFileDescriptors()

	// Sorting and deduplications
	common.LogDebug(fmt.Sprintf("[Start] Sorting: %s", ongoing_path))

	om2.SortByAllColumnsAssumingFirstIsDate(true)
	//om2.SortByAllColumnsAssumingFirstIsDate(false)
	common.LogDebug(fmt.Sprintf("[End] Sorting %s", ongoing_path))
	om2.RemoveDuplicatesAssumingSorted()
	common.LogDebug(fmt.Sprintf("[End] Removing duplicates %s", ongoing_path))

	// Write to new file in desired output format
	om2.SetPath(final_path)
	defer om2.CloseAndFlushFileDescriptors()
	om2.SetExcelSheetName(l2s.latest_computer)
	common.LogDebug(fmt.Sprintf("[Start] Saving to output format %s", ongoing_path))
	err_is_saved := om2.SaveAllDataToProperFormat(l2s.GetOutputFormatName())

	if err_is_saved != nil {
		common.LogCriticalErrorWithError(fmt.Sprintf("Cannot properly saved in desired format: %s", final_path), err_is_saved)
	} else {
		os.Remove(ongoing_path)
	}

	common.LogInfo(fmt.Sprintf("Results saved to %s format - L2: %s | Hostname: %s | Nr of source files: %d", l2s.GetOutputFormatName(), l2s.l2_name, l2s.latest_computer, l2s.workers_counter))
}

func RunL2WorkerPowershellScriptblock(l2s *Layer2SingleLayer) {
	//log.Println("New layer 2 worker (Type: powershell_scriptblock)")

	// ScriptBlocks memory
	memory := make(map[string]common.PowerShellScriptBlockInfo, 1000)

	// Prepare fields
	matching_id_field := l2s.parsing_info.Aggregation_options["field_matching_id"]
	total_number_field := l2s.parsing_info.Aggregation_options["field_total_number"]
	current_number_field := l2s.parsing_info.Aggregation_options["field_current_number"]
	content_field := l2s.parsing_info.Aggregation_options["field_content"]
	filename_field := l2s.parsing_info.Aggregation_options["field_filename"]

	// Save into memory all scriptblocks
	for ev_map := range l2s.ch {

		// Layer 2 fields
		final_od := l2s.engine.ParseL2FieldsOrderedDict(l2s.l2_name, ev_map)

		matching_id, _ := final_od.GetString(matching_id_field)
		total_number, _ := final_od.GetString(total_number_field)
		total_number_int, _ := strconv.Atoi(total_number)
		current_number, _ := final_od.GetString(current_number_field)
		current_number_int, _ := strconv.Atoi(current_number)
		content, _ := final_od.GetString(content_field)
		filename, _ := final_od.GetString(filename_field)

		if len(matching_id) < 30 {
			common.LogError("Wrong matching id length " + matching_id)
		}

		// Check if empty structure exists
		if _, id_exists := memory[matching_id]; !id_exists {
			memory[matching_id] = common.PowerShellScriptBlockInfo{
				Total:    total_number_int,
				Segments: make(map[int]string, total_number_int),
				Path:     filename,
			}
		}

		// Add segment content
		memory[matching_id].Segments[current_number_int] = content
	}

	// Create output directory
	os.MkdirAll(l2s.GetScriptBlockOutputPath(""), 666)

	// Reconstruct scriptblocks - UTF-8 with no BOM
	for id, block_info := range memory {

		is_complete := "_incomplete"
		// Check completeness
		if len(block_info.Segments) == block_info.Total {
			is_complete = ""
		}

		// New filename
		part_filename := ""
		if len(block_info.Path) > 0 {
			extracted_filename := filepath.Base(block_info.Path)
			if len(extracted_filename) > 0 {
				part_filename = extracted_filename + "__"
			}
		}

		just_filename := fmt.Sprintf("%s%s%s.ps1", part_filename, id, is_complete)
		final_filename := l2s.GetScriptBlockOutputPath(just_filename)

		// Create new file
		psblock_file, err := os.Create(final_filename)
		psblock_datawriter := bufio.NewWriter(psblock_file)

		if err != nil {
			common.LogCriticalError("When writing ScriptBlock to : " + final_filename)
		}

		// Save segments to file
		for i := 1; i <= block_info.Total; i++ {
			if segment, seg_exists := block_info.Segments[i]; seg_exists {
				psblock_datawriter.WriteString(segment)
			} else {
				psblock_datawriter.WriteString("\r\n\r\n[Missing segment]\r\n\n")
			}
		}

		// Close file
		psblock_datawriter.Flush()
		psblock_file.Close()

		common.LogDebug(fmt.Sprintf("ScriptBlock saved %s", just_filename))
	}
	common.LogInfo(fmt.Sprintf("%d Scriptblocks saved | %s | %s", len(memory), l2s.l2_name, l2s.latest_computer))
}

func RunL2Worker(l2s *Layer2SingleLayer, Wg_l2_all *sync.WaitGroup) {
	worker_typ := l2s.parsing_info.Output.Typ
	defer Wg_l2_all.Done()

	common.LogDebug(fmt.Sprintf("Running L2 worker: %s | %s | Workers counter: %d", l2s.l2_name, l2s.latest_computer, l2s.workers_counter))

	if worker_typ == "flat" {
		RunL2WorkerFlat(l2s)
	} else if worker_typ == "powershell_scriptblock" {
		RunL2WorkerPowershellScriptblock(l2s)
	} else {
		panic("Unknown L2 worker typ!")
	}

	common.LogDebug(fmt.Sprintf("Ending L2Worker: %s | %s", l2s.l2_name, l2s.latest_computer))
}

func (l2s *Layer2SingleLayer) IncrementWorkerCounter() {
	//log.Printf("IncrementWorkerCounter: %s | %s \n", l2s.l2_name, l2s.latest_computer)
	l2s.wg_to_close_channel.Add(1)
	atomic.AddInt32(&l2s.workers_counter, 1)
}

func (l2mem *Layer2GlobalMemory) StartL2Workers(engi *Engine) {

	// Expect workers
	l2mem.Wg_l2_all.Add(len(l2mem.Layers))

	for i := 0; i < len(l2mem.Layers); i++ {
		var l2mem_layer = &l2mem.Layers[i]

		// Execute new parallel worker
		go RunL2Worker(l2mem_layer, l2mem.Wg_l2_all)
	}
}

func ChannelKiller(l2s *Layer2SingleLayer) {
	//log.Printf("Start ChannelKiller - pointer: %s \n", reflect.ValueOf(l2s.ch).Pointer())
	// Wait for sources sending events
	l2s.wg_to_close_channel.Wait()

	// Close channel
	//log.Printf("Closed channel for %s | %s \n", l2s.latest_computer, l2s.l2_name)
	close(l2s.ch)
}

func (l2mem *Layer2GlobalMemory) SetupChannelKillers() {

	for i := 0; i < len(l2mem.Layers); i++ {
		l2s := &l2mem.Layers[i]
		go ChannelKiller(l2s)
	}

}

func (l2s *Layer2SingleLayer) GetComputerDir() string {
	return l2s.parsing_info.Output.GlobalOutputDirectory + l2s.latest_computer + string(os.PathSeparator)
}

func (l2s *Layer2SingleLayer) GetOngoingOutputPath() string {
	return l2s.GetComputerDir() + l2s.parsing_info.Output.Category + string(os.PathSeparator) + l2s.parsing_info.Output.Filename + ".jsonl.gz.ongoing"
}

func (l2s *Layer2SingleLayer) GetOutputFormatName() string {
	return l2s.engine.OutputFormat
}

func (l2s *Layer2SingleLayer) GetOutputFormatExtension() string {
	switch l2s.engine.OutputFormat {
	case "csv":
		{
			return ".csv"
		}
	case "json":
		{
			return ".json"
		}
	case "jsonl":
		{
			return ".jsonl"
		}
	case "excel":
		{
			return ".xlsx"
		}
	default:
		panic("Unknown format type")
	}
}

func (l2s *Layer2SingleLayer) GetOutputPath() string {
	return l2s.GetComputerDir() + l2s.parsing_info.Output.Category + string(os.PathSeparator) + l2s.parsing_info.Output.Filename + l2s.GetOutputFormatExtension()
}

func (l2s *Layer2SingleLayer) GetScriptBlockOutputPath(filename string) string {
	return l2s.GetComputerDir() + l2s.parsing_info.Output.Category + string(os.PathSeparator) + l2s.parsing_info.Output.Subcategory + string(os.PathSeparator) + filename
}
