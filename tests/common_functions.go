package tests

import (
	"fmt"
	"github.com/Velocidex/ordereddict"
	"github.com/rs/zerolog"
	"github.com/yarox24/evtxhussar/common"
	"github.com/yarox24/evtxhussar/engine"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
)

// go test -args "F:\path_to\maps"
func GrabCmdlineMapPath() string {

	return "F:\\GoLangBase\\EvtxHussarProject\\maps"
	if len(os.Args) == 4 {
		return os.Args[3]
	} else {
		fmt.Println("You need to provide path to maps. Example: go test -args \"F:\\path_to\\maps\"")
	}

	return ""
}

func LoadEngine() *engine.Engine {
	maps_path := GrabCmdlineMapPath()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	maps_path2, _ := common.Determine_Maps_Path(maps_path)
	engi := engine.NewEngine("csv", maps_path2)
	engi.LoadLayer1()
	engi.LoadLayer2("F:\\GoLangBase\\EvtxHussarProject\\test_out_fake")
	engi.LoadParams()

	return &engi
}

func GetTestEventsPath(maps_path string) string {
	events_path_temp := strings.Split(maps_path, "\\")
	return strings.Join(events_path_temp[:len(events_path_temp)-2], "\\") + "\\tests\\events\\"
}

func GetTestFilesPath(maps_path string) string {
	events_files_temp := strings.Split(maps_path, "\\")
	return strings.Join(events_files_temp[:len(events_files_temp)-2], "\\") + "\\tests\\files\\"
}

func UnmarshallAndParseEvent(filename string, eng *engine.Engine, l2_name string) *ordereddict.Dict {
	test_event_path := GetTestEventsPath(eng.Maps_path) + filename

	o := ordereddict.NewDict()
	f, err := os.OpenFile(test_event_path, os.O_RDONLY, 0755)

	if err == nil {
		defer f.Close()

		b, err2 := ioutil.ReadAll(f)
		if err2 != nil {
			panic("When reading file")
		} else {
			err3 := o.UnmarshalJSON(b)

			if err3 == nil {
				return ParseEvent(eng, l2_name, o)
			} else {
				panic("When reading file")
			}
		}
	} else {
		panic("When reading file")
	}

	return nil
}

func ParseEvent(eng *engine.Engine, l2_name string, ev *ordereddict.Dict) *ordereddict.Dict {
	// Extract Event

	ev_map, ok_map := ordereddict.GetMap(ev, "Event")

	if !ok_map {
		panic("Test failed")
	}

	// Common fields
	final_od := eng.ParseCommonFieldsOrderedDict(ev_map)

	// Layer 2 fields
	l2_od := eng.ParseL2FieldsOrderedDict(l2_name, ev_map)

	// Combine both
	final_od.MergeFrom(l2_od)

	return final_od
}

type EasyTesting struct {
	t  *testing.T
	od *ordereddict.Dict
}

func NewEasyTesting(t *testing.T, od *ordereddict.Dict) *EasyTesting {
	return &EasyTesting{
		t:  t,
		od: od,
	}
}

func UniversalCheckDesiredValue(t *testing.T, current_value string, expected_value string) {
	if strings.ToLower(expected_value) != strings.ToLower(current_value) {
		t.Errorf("Value mismatch: %s <-> %s", expected_value, current_value)
	}
}

func UniversalCheckDesiredValueBoolean(t *testing.T, current_value bool, expected_value bool) {

	if expected_value != current_value {
		t.Errorf("Value mismatch (bool): %t <-> %t", expected_value, current_value)
	}
}

func UniversalCheckDesiredValueUint64(t *testing.T, current_value uint64, expected_value uint64) {

	if expected_value != current_value {
		t.Errorf("Value mismatch (bool): %d <-> %d", expected_value, current_value)
	}
}

func UniversalCheckDesiredValueint64(t *testing.T, current_value int64, expected_value int64) {

	if expected_value != current_value {
		t.Errorf("Value mismatch (bool): %d <-> %d", expected_value, current_value)
	}
}

func (et *EasyTesting) CheckDesiredValue(key string, expected_value string) {

	eid, _ := et.od.GetString("EID")
	if value, key_exists := et.od.GetString(key); key_exists {
		if strings.ToLower(value) != strings.ToLower(expected_value) {
			et.t.Errorf("[EID %s] Key %s is equal to: %s , but should be: %s\n", eid, key, value, expected_value)
		}
	} else {
		et.t.Errorf("[EID %s]Key %s doesn't exists in dictionary!\n", eid, key)
	}
}

func (et *EasyTesting) CheckDesiredLength(key string, expected_length int) {
	if value, key_exists := et.od.GetString(key); key_exists {
		if len(value) != expected_length {
			et.t.Errorf("Key %s has length of: %d , but should be: %d\n", key, len(value), expected_length)
		}
	} else {
		et.t.Errorf("Key %s doesn't exists in dictionary!\n", key)
	}
}

func FakeInspect_single_evtx(efi *common.EvtxFileInfo) {
	var wg sync.WaitGroup
	var inspector_channel chan struct{} = make(chan struct{}, 1)

	inspector_channel <- struct{}{}
	wg.Add(1)
	common.Inspect_single_evtx(&wg, efi, inspector_channel)

	wg.Wait()
	close(inspector_channel)
}

func FakeInspectEvtx(maps_path string, evtx_filename string) common.EvtxFileInfo {
	test_file_path := GetTestFilesPath(maps_path) + evtx_filename
	efi := common.NewEvtxFileInfo(test_file_path)

	FakeInspect_single_evtx(&efi)

	return efi
}

func TestL1Worker(maps_path string, evtx_filename string) int64 {
	test_file_path := GetTestFilesPath(maps_path) + evtx_filename
	efi := common.NewEvtxFileInfo(test_file_path)

	FakeInspect_single_evtx(&efi)

	var wg_all sync.WaitGroup
	var l2s_wg_to_close_channel_list []*sync.WaitGroup = make([]*sync.WaitGroup, 0)

	wg_all.Add(1)
	ch_limit_worker := make(chan struct{}, 1)
	ch_limit_worker <- struct{}{}
	HClist := engine.HostnameToChannels{
		Channels: nil,
	}
	var fake_counter uint64

	engine.RunL1Worker(&wg_all, &efi, HClist, l2s_wg_to_close_channel_list, ch_limit_worker, &fake_counter)
	wg_all.Wait()
	close(ch_limit_worker)

	return efi.GetNumberOfRecords()
}
