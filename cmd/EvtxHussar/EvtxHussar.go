package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/yarox24/EvtxHussar/common"
	"github.com/yarox24/EvtxHussar/engine"
	"strings"
)

var VERSION = "1.1"
var URL = "github.com/yarox24/EvtxHussar"

type Config struct {
	Recursive        bool     `arg:"-r,--recursive" help:"Recursive traversal for any input directories." default:"false""`
	Output_dir       string   `arg:"-o,--output_dir" help:"Reports will be saved in this directory (if doesn't exists it will be created)"`
	Output_format    string   `arg:"-f,--format" help:"Output data in one of the formats: Csv,JSON,JSONL,Excel" default:"Excel"`
	Input_evtx_paths []string `arg:"positional" help:"Path(s) to .evtx files or directories containing these files (can be mixed)"`
	WorkersLimit     int      `arg:"-w,--workers" help:"Max concurrent workers (.evtx opened)" default:"30"`
	MapsDirectory    string   `arg:"-m,--maps" help:"Custom directory with maps/ (Default: program directory)"`
	Debug            bool     `arg:"-d,--debug" help:"Be more verbose" default:"false""`
}

func (Config) Version() string {
	return "EvtxHussar " + VERSION
}

func (Config) Description() string {
	return "Initial triage of Windows Event logs"
}

// Profiling
//f, err := os.Create("profile.pb.gz")
//if err != nil {
//	log.Fatal(err)
//}
//err = pprof.StartCPUProfile(f)
//if err != nil {
//	log.Fatal(err)
//}
//defer pprof.StopCPUProfile()

func ValidateOutputFormat(p *arg.Parser, c *Config) {
	c.Output_format = strings.ToLower(c.Output_format)
	switch strings.ToLower(c.Output_format) {
	case "csv":
		{
			// OK
		}
	case "json":
		{
			// OK
		}
	case "jsonl":
		{
			// OK
		}
	case "excel":
		{
			// OK
		}
	default:
		p.Fail("When providing output format following values are valid: Csv,JSON,JSONL,Excel")
	}
}

func main() {

	// ASCI ART
	common.Hussar_art(VERSION, URL)

	// Config
	conf := Config{}
	p := arg.MustParse(&conf)
	ValidateOutputFormat(p, &conf)

	// Logging
	common.Logging_init(conf.Debug)
	common.LogDebugStructure("Parsed", conf, "Config")

	// Arguments check
	if len(conf.Input_evtx_paths) == 0 {
		common.LogCriticalError("You need to provide path(s) to directory(s) which contains .evtx or path to .evtx file(s) as argument(s)")
	}

	// Output dir creation
	common.Handle_output_directory(conf.Output_dir)

	// Load engine
	maps_path, err := common.Determine_Maps_Path(conf.MapsDirectory)

	if err != nil {
		common.LogCriticalError("Cannot find maps or maps\\{layer2,params} directory")
	}

	engi := engine.NewEngine(conf.Output_format, maps_path)
	engi.LoadLayer1()
	engi.LoadLayer2(conf.Output_dir)
	engi.LoadParams()

	// Generate list of .evtx files
	common.LogInfo("Generating list of .evtx files in provided paths...")
	var EfiList = common.Generate_list_of_files_to_process(conf.Input_evtx_paths, conf.Recursive)

	common.LogInfo(fmt.Sprintf("Inspecting %d found .evtx files", len(EfiList)))
	if len(EfiList) > 0 {
		// Inspect .evtx files
		EfiList = common.Inspect_evtx_paths(EfiList)
	} else {
		common.LogCriticalError("No .evtx files were found in given locations")
	}

	// Check for supported files
	for i := 0; i < len(EfiList); i++ {
		engi.IsEfiSupported(&EfiList[i])
	}

	// Dump info about files
	counters := struct {
		processing_counter int
		nr_of_invalid_evtx int
		nr_of_empty_evtx   int
	}{
		processing_counter: 0,
		nr_of_invalid_evtx: 0,
		nr_of_empty_evtx:   0,
	}

	for i := 0; i < len(EfiList); i++ {
		if EfiList[i].IsEmpty() {
			counters.nr_of_empty_evtx += 1
		}

		if !EfiList[i].IsValid() {
			counters.nr_of_invalid_evtx += 1
		}

		if EfiList[i].WillBeProcessed() {
			counters.processing_counter += 1
		}
	}

	common.LogInfo(fmt.Sprintf("Send to processing: %d files", counters.processing_counter))
	common.GrabUniversalLogger().Info().Int("nr_of_empty_evtx", counters.nr_of_empty_evtx).Int("nr_of_invalid_evtx", counters.nr_of_invalid_evtx).Msg("Summary")

	// Initalize Layer 1 Global struct
	var l1globmem = engine.NewLayer1GlobalMemory()

	// Initalize Layer 2 Global struct
	var l2globmem = engine.NewLayer2GlobalMemory()

	// Mixed layers setup steps
	l2globmem.SetupChannels(&engi, common.ReturnCopyOfSupportedEfiFileInfoElements(EfiList))
	l1globmem.SetupWorkers(&engi, common.ReturnCopyOfSupportedEfiFileInfoElements(EfiList), &l2globmem)

	l2globmem.SetupChannelKillers()

	l2globmem.StartL2Workers(&engi)
	l1globmem.StartL1Workers(&engi, common.ReturnCopyOfSupportedEfiFileInfoElements(EfiList), &l2globmem, conf.WorkersLimit)

	common.LogInfo("Start processing")
	l1globmem.Wg_l1_all.Wait()
	l2globmem.Wg_l2_all.Wait()
	common.LogInfo("End processing")

}
