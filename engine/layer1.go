package engine

import (
	"errors"
	"fmt"
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/common"
	"github.com/yarox24/EvtxHussar/eventmap"
	"golang.org/x/exp/slices"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"www.velocidex.com/golang/evtx"
)

type Layer1Events struct {
	Attrib_extraction []string
	Short_description string
	Provider_guid     string
}

type Layer1EventsEnhanced struct {
	Attrib_extraction []common.ExtractedFunction
	Short_description string
}

func NewLayer1EventsEnhanced(Attrib_extraction []string, Short_description string) Layer1EventsEnhanced {

	var l1enh = Layer1EventsEnhanced{
		Attrib_extraction: make([]common.ExtractedFunction, 0),
		Short_description: Short_description,
	}

	for _, attr := range Attrib_extraction {
		l1enh.Attrib_extraction = append(l1enh.Attrib_extraction, common.FunctionExtractor(attr))
	}

	return l1enh
}

type Layer1Info struct {
	Typ             string
	Source_comment  string
	Channel         string
	Content_data_in string
}

type Layer1 struct {
	Info           Layer1Info
	Sendto_layer2  string
	Ordered_fields []string
	Events         map[string]Layer1Events
	Options        map[string]string
}

func NewLayer1GlobalMemory() Layer1GlobalMemory {
	return Layer1GlobalMemory{
		Hclist:                 make(map[string]HostnameToChannels, 0),
		Wg_l1_all:              new(sync.WaitGroup),
		Atomic_Counter_Workers: 0,
	}
}

type Layer1GlobalMemory struct {
	Hclist                 map[string]HostnameToChannels
	Wg_l1_all              *sync.WaitGroup
	Atomic_Counter_Workers uint64
}

func (l1globmem *Layer1GlobalMemory) SetupWorkers(e *Engine, efi []common.EvtxFileInfo, l2globmem *Layer2GlobalMemory) {

	// For every supported file
	for i := 0; i < len(efi); i++ {
		latest_computer := efi[i].GetLatestComputer()
		channel := efi[i].GetChannel()
		l1list := e.GetAllLayer1WhichSupportsChannel(channel)

		for _, l1 := range l1list {
			l2s := l2globmem.FindLayer2SingleLayer(l1.Sendto_layer2, latest_computer)
			l2s.IncrementWorkerCounter()
			common.LogDebug(fmt.Sprintf("Preparing worker %s %s to send to L2: %s", latest_computer, channel, l2s.l2_name))

			for eid, _ := range l1.Events {
				l1globmem.HClistAddChan(latest_computer, l1.Info.Channel, eid, l1.Events[eid].Provider_guid, l2s.ch)
			}
		}
	}
}

func (l1globmem *Layer1GlobalMemory) HClistAddHostname(latest_computer string) {
	_, ok := l1globmem.Hclist[latest_computer]

	if !ok {
		l1globmem.Hclist[latest_computer] = NewHostnameToChannels()
	}
}

func (l1globmem *Layer1GlobalMemory) HClistAddChannel(latest_computer string, channel string) {

	// Hostname
	l1globmem.HClistAddHostname(latest_computer)

	_, ok := l1globmem.Hclist[latest_computer].Channels[channel]

	if !ok {
		l1globmem.Hclist[latest_computer].Channels[channel] = ChannelTOEID{
			Eid: make(map[string]EIDToChan),
		}
	}
}

func (l1globmem *Layer1GlobalMemory) HClistAddEID(latest_computer string, channel string, eid string, provider_guid string) {
	// Hostname + Channel
	l1globmem.HClistAddChannel(latest_computer, channel)

	_, ok := l1globmem.Hclist[latest_computer].Channels[channel].Eid[eid]

	if !ok {
		l1globmem.Hclist[latest_computer].Channels[channel].Eid[eid] = EIDToChan{
			Chans:         make([]chan *eventmap.EventMap, 0),
			Provider_guid: provider_guid,
		}
	}
}

func (l1globmem *Layer1GlobalMemory) HClistAddChan(latest_computer string, channel string, eid string, provider_guid string, ch chan *eventmap.EventMap) {

	// Hostname + Channel + EID
	l1globmem.HClistAddEID(latest_computer, channel, eid, provider_guid)

	// Add chan if not exists
	new_addr := reflect.ValueOf(ch).Pointer()

	for i := 0; i < len(l1globmem.Hclist[latest_computer].Channels[channel].Eid[eid].Chans); i++ {
		existing_addr := reflect.ValueOf(l1globmem.Hclist[latest_computer].Channels[channel].Eid[eid].Chans[i]).Pointer()
		if existing_addr == new_addr {
			common.LogDebug(fmt.Sprintf("Don't append Channel (Duplicate) | %s | %s | %s | %p", latest_computer, channel, eid, new_addr))
			return
		}
	}

	// Append to last element
	var etc = l1globmem.Hclist[latest_computer].Channels[channel].Eid[eid]
	etc.Chans = append(etc.Chans, ch)
	l1globmem.Hclist[latest_computer].Channels[channel].Eid[eid] = etc
}

func l1close_wait_groups_in_loop(l2s_wg_to_close_channel_list []*sync.WaitGroup) {
	for _, l2s_wg := range l2s_wg_to_close_channel_list {
		l2s_wg.Done()
	}
}

func serialize_event(ev *ordereddict.Dict) {
	ev_map, _ := ordereddict.GetMap(ev, "Event")

	// Filename
	filename := fmt.Sprintf("%s_%s.json", eventmap.GetChannel(ev_map), eventmap.GetEID(ev_map))
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, " ", "_")
	out_path := "X:\\saved_test_events\\" + filename

	// Tamper time
	system_ev, is_ok1 := ev_map.Get("System")
	if !is_ok1 {
		panic("is_ok1")
	}

	system_ev_dict, _ := system_ev.(*ordereddict.Dict)
	time_created, is_ok2 := system_ev_dict.Get("TimeCreated")

	if !is_ok2 {
		panic("is_ok2")
	}
	time_created_dict, _ := time_created.(*ordereddict.Dict)
	time_created_dict.Update("SystemTime", 1638798000.1111111)

	// Tamper hostname
	system_ev_dict.Update("Computer", "DESKTOP-EvtxHussar")

	// Tamper Event Record ID
	system_ev_dict.Update("EventRecordID", 111)

	// Tamper Security UserID
	security, is_ok3 := system_ev_dict.Get("Security")

	if is_ok3 {
		security_dict, is_ok_security := security.(*ordereddict.Dict)

		if is_ok_security && slices.Contains(security_dict.Keys(), "UserID") {
			security_dict.Update("UserID", "S-1-5-18")
		}

	}

	// Check if file already exists
	if _, err := os.Stat(out_path); errors.Is(err, os.ErrNotExist) {
		b, err_serial := ev.MarshalJSON()

		if err_serial == nil {
			f_out, _ := os.Create(out_path)

			f_out.Write(b)
			f_out.Close()
		}
	}

}

func RunL1Worker(Wg_l1_all *sync.WaitGroup, efi *common.EvtxFileInfo, Hclist HostnameToChannels, l2s_wg_to_close_channel_list []*sync.WaitGroup, ch_limit_worker chan struct{}, Atomic_Counter_Workers *uint64) {
	// Run only when not exceeding limit
	<-ch_limit_worker

	//log.Printf("Run L1Worker: %s | %s \n", efi.GetLatestComputer(), efi.GetChannel())
	defer Wg_l1_all.Done()
	defer l1close_wait_groups_in_loop(l2s_wg_to_close_channel_list)
	defer atomic.AddUint64(Atomic_Counter_Workers, ^uint64(0))

	supported_eids := Hclist.Channels[efi.GetChannel()]

	// Open evtx file
	fd, err := os.OpenFile(efi.GetPath(), os.O_RDONLY, os.FileMode(0666))

	if err == nil {
		defer fd.Close()
	} else {
		common.LogCriticalErrorWithError("Error occured when opening evtx: "+efi.GetPath(), err)
		return
	}

	chunks, err_chunks := evtx.GetChunks(fd)

	// Flags is dirty
	const IS_DIRTY = 0x1

	if efi.GetAlternativeHeader().FileFlags == IS_DIRTY {
		common.LogInfo(fmt.Sprintf("Dirty file detected: %s", efi.GetPath()))
		// => Parsing all found chunks
	} else {
		// => Cut off chunks over header number
		header_chunks_counts := int(efi.GetAlternativeHeader().ChunkCount)
		found_chunks_count := len(chunks)

		if header_chunks_counts < found_chunks_count {
			chunks = chunks[0:header_chunks_counts]
		}
	}

	if err_chunks != nil {
		common.LogErrorWithError("Evtx chunks error: "+efi.GetPath(), err_chunks)
		return
	}

	var record_counter int64 = 0

	for _, chunk := range chunks {
		records, err_chunk := chunk.Parse(0)

		if err_chunk != nil {
			common.LogError("Chunk parsing error: " + efi.GetPath())
			continue
		}

		for _, i := range records {
			ev, ok := i.Event.(*ordereddict.Dict)

			if ok {
				ev_map, ok_map := ordereddict.GetMap(ev, "Event")

				// Now count
				record_counter += 1

				if !ok_map {
					common.LogError("Event parsing error: " + efi.GetPath())
					continue
				}

				eid := eventmap.GetEID(ev_map)

				eid_struct, eid_supported := supported_eids.Eid[eid]

				// Check if EID should be further send
				if eid_supported {

					// Check Provider GUID if present
					if len(eid_struct.Provider_guid) > 0 {
						provider_guid := eventmap.GetProviderGUID(ev_map)
						if strings.ToLower(eid_struct.Provider_guid) != strings.ToLower(provider_guid) {
							continue
						}
					}

					//serialize_event(ev)

					// Send to all related channels matched event
					for _, cha := range eid_struct.Chans {
						cha <- ev_map
					}
				}
			}
		}
	}
	efi.SetNumberOfRecords(record_counter)

	// PowerShell records validation
	//out, err := exec.Command("powershell", "-NoProfile", fmt.Sprintf("Get-Winevent -path \"%s\" | Measure-Object | Select-Object Count| ConvertTo-Json", efi.GetPath())).CombinedOutput()
	//if err != nil {
	//	panic("Powershell error")
	//}
	//// JSON convert
	//type Measure struct {
	//	Count int
	//}
	//var m Measure
	//
	//err2 := json.Unmarshal(out, &m)
	//if err2 != nil {
	//	panic(err2)
	//}
	//fmt.Println(m.Count)
	//
	//// Ultimate check
	//if int(record_counter) != m.Count {
	//	fmt.Println("evtx mismatch error found!")
	//}

	common.LogDebug(fmt.Sprintf("Finished L1Worker: %s | %s | Records: %d", efi.GetLatestComputer(), efi.GetChannel(), record_counter))
}

func ConcurrencyUnlockNewWorkers(ch_limit_worker chan struct{}, nr_of_files int) {

	// Keep new workers at limit
	for i := 0; i < nr_of_files; i++ {
		ch_limit_worker <- struct{}{}
	}

	close(ch_limit_worker)
}

func PrintRemaingWorkers(Atomic_Counter_Workers *uint64) {
	i := 0

	for {
		// First report after one minute
		time.Sleep(time.Minute * time.Duration(i+1))
		remaining_workers := atomic.LoadUint64(Atomic_Counter_Workers)

		if remaining_workers == 0 {
			common.LogWarn("All .evtx workers finished. Post-processing tasks (sorting, deduplication) are running ...")
			break
		} else if remaining_workers == 1 {
			common.LogWarn(fmt.Sprintf("%d .evtx worker is still running ...", remaining_workers))
		} else {
			common.LogWarn(fmt.Sprintf("%d .evtx workers are still running ...", remaining_workers))
		}

		if i < 5 {
			i += 1
		}
	}
}

func (l1globmem *Layer1GlobalMemory) StartL1Workers(e *Engine, efi []common.EvtxFileInfo, l2globmem *Layer2GlobalMemory, WorkersLimit int) {

	// Expect
	l1globmem.Wg_l1_all.Add(len(efi))
	atomic.AddUint64(&l1globmem.Atomic_Counter_Workers, uint64(len(efi)))
	// Limit conncurent workers
	ch_limit_worker := make(chan struct{}, WorkersLimit)

	go ConcurrencyUnlockNewWorkers(ch_limit_worker, len(efi))
	go PrintRemaingWorkers(&l1globmem.Atomic_Counter_Workers)

	// For every supported file
	for i := 0; i < len(efi); i++ {
		latest_computer := efi[i].GetLatestComputer()
		channel := efi[i].GetChannel()
		l1list := e.GetAllLayer1WhichSupportsChannel(channel)

		var l2s_wg_to_close_channel_list []*sync.WaitGroup = make([]*sync.WaitGroup, 0)

		for _, l1 := range l1list {
			l2s := l2globmem.FindLayer2SingleLayer(l1.Sendto_layer2, latest_computer)
			l2s_wg_to_close_channel_list = append(l2s_wg_to_close_channel_list, l2s.wg_to_close_channel)
		}

		go RunL1Worker(l1globmem.Wg_l1_all, &efi[i], l1globmem.Hclist[latest_computer], l2s_wg_to_close_channel_list, ch_limit_worker, &l1globmem.Atomic_Counter_Workers)
	}
}

type HostnameToChannels struct {
	Channels map[string]ChannelTOEID
}

type EIDToChan struct {
	Chans         []chan *eventmap.EventMap
	Provider_guid string
}

type ChannelTOEID struct {
	Eid map[string]EIDToChan
}

func NewHostnameToChannels() HostnameToChannels {
	return HostnameToChannels{
		Channels: make(map[string]ChannelTOEID, 0),
	}
}
