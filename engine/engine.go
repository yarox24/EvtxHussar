package engine

import (
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/common"
	"github.com/yarox24/EvtxHussar/eventmap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

type UnifiedChannelEvents struct {
	EIDs map[string]Layer1EventsEnhanced
}

func NewUnifiedChannelEvents() UnifiedChannelEvents {
	return UnifiedChannelEvents{
		EIDs: make(map[string]Layer1EventsEnhanced, 0),
	}
}

type Engine struct {
	Layer1               []Layer1
	Layer2               []Layer2
	DoubleQuotes         map[string]common.Params
	VariousMappers       map[string]common.Params
	UnifiedChannelEvents map[string]UnifiedChannelEvents
	Common               Layer1
	Maps_path            string
	OutputFormat         string
}

func NewEngine(output_format string, maps_path string) Engine {
	return Engine{
		Layer1:               make([]Layer1, 0),
		Layer2:               make([]Layer2, 0),
		DoubleQuotes:         make(map[string]common.Params, 0),
		VariousMappers:       make(map[string]common.Params, 0),
		UnifiedChannelEvents: make(map[string]UnifiedChannelEvents, 0),
		Common:               Layer1{},
		Maps_path:            maps_path,
		OutputFormat:         output_format,
	}
}

func (e *Engine) LoadLayer1() {

	l1_files, err := ioutil.ReadDir(e.Maps_path)

	if err != nil {
		common.LogCriticalErrorWithError("When reading maps directory", err)
	}
	common.LogDebugStructure("L1 files", l1_files, "l1_files")

	// Layer1 - maps/*.yaml
	for _, f := range l1_files {

		// Skip directories
		if f.IsDir() {
			continue
		}

		file_bytes, err := ioutil.ReadFile(e.Maps_path + f.Name())

		if err != nil {
			common.LogCriticalErrorWithError("When reading L1 map: "+f.Name(), err)
		}

		var l1 = new(Layer1)
		err = yaml.Unmarshal(file_bytes, l1)

		if err != nil {
			common.LogCriticalErrorWithError("When unmarshalling L1 map: "+f.Name(), err)
		}

		if l1.Info.Typ == "common" {
			e.Common = *l1
			common.LogDebug("Loaded common map (L1)")
		} else if l1.Info.Typ == "layer1" {
			e.Layer1 = append(e.Layer1, *l1)
			common.LogDebug("Loaded L1 (layer1) map: " + f.Name())

			// UnifiedChannelEvents
			e.AppendLayer1ToUnifiedChannelEvents(l1)

		} else {
			panic("YAML - LoadLayer1() - Unsupported Info.Typ")
		}
	}
}

func (e *Engine) LoadLayer2(Output_dir string) {

	l2_layer_dir := e.Maps_path + "layer2" + string(os.PathSeparator)
	l2_files, err := ioutil.ReadDir(l2_layer_dir)

	if err != nil {
		common.LogCriticalErrorWithError("When reading maps/layer2 directory", err)
	}
	common.LogDebugStructure("L2 files", l2_files, "l2_files")

	// Layer2 - maps/layer2/*.yaml
	for _, f := range l2_files {

		// Skip directories
		if f.IsDir() {
			continue
		}

		file_bytes, err := ioutil.ReadFile(l2_layer_dir + f.Name())

		if err != nil {
			common.LogCriticalErrorWithError("When reading L2 map: "+f.Name(), err)
		}

		var l2 = new(Layer2)
		err = yaml.Unmarshal(file_bytes, l2)

		if err != nil {
			common.LogCriticalErrorWithError("When unmarshalling L2 map: "+f.Name(), err)
		}

		if l2.Info.Typ == "layer2" {
			// Append OutputDir
			l2.Output.GlobalOutputDirectory = Output_dir + string(os.PathSeparator)

			// Case insensitive remmaping dict
			l2.Fields_remap_dict = ordereddict.NewDict()
			l2.Fields_remap_dict.SetCaseInsensitive()

			// Field extra transformations - Options append & Rename
			for i, trans := range l2.Field_extra_transformations {
				ef := common.FunctionExtractor(trans.Special_transform)
				l2.Field_extra_transformations[i].Special_transform = ef.Name
				l2.Field_extra_transformations[i].Options = ef.Options
			}

			// Enhance ordered fields
			l2.Ordered_fields_enhanced = make(map[string]common.SingleField, 0)

			for i := 0; i < len(l2.Ordered_fields); i++ {
				sf := e.SingleFieldExtractor(l2.Ordered_fields[i])
				l2.Ordered_fields_enhanced[strings.ToLower(sf.NiceName)] = sf
				l2.Ordered_fields[i] = sf.NiceName
			}

			for k, v := range l2.Fields_remap {
				l2.Fields_remap_dict.Set(k, v)
			}

			e.Layer2 = append(e.Layer2, *l2)
		} else {
			panic("YAML - LoadLayer2() - Unsupported Info.Typ")
		}
	}
}

func (e *Engine) IncreaseUsageCounterForLayer2(l2name string) {
	for i := 0; i < len(e.Layer2); i++ {
		if e.Layer2[i].Info.Name == l2name {
			e.Layer2[i].UsageCounter += 1
			return
		}
	}

	panic("Wrong name for sendto_layer2")
}

func (e *Engine) GetAllLayer1WhichSupportsChannel(channel string) []*Layer1 {
	var temp = make([]*Layer1, 0)

	for i := 0; i < len(e.Layer1); i++ {
		if strings.ToLower(e.Layer1[i].Info.Channel) == strings.ToLower(channel) {
			temp = append(temp, &e.Layer1[i])
		}
	}
	return temp
}

func (e *Engine) IsEfiSupported(efi *common.EvtxFileInfo) {

	// Is valid
	if !efi.IsValid() {
		return
	}

	// Is non-empty | This might not be necessary
	if efi.IsEmpty() {
		return
	}

	// Is channel supported?
	ch := strings.ToLower(efi.GetChannel())

	for _, llayer1 := range e.Layer1 {
		if ch == strings.ToLower(llayer1.Info.Channel) {
			efi.EnableForProcessing()
			e.IncreaseUsageCounterForLayer2(llayer1.Sendto_layer2)
		}
	}

	return
}

func (e *Engine) FindL2LayerByName(name string) *Layer2 {

	for i := 0; i < len(e.Layer2); i++ {
		var l2 = e.Layer2[i]
		if strings.ToLower(l2.Info.Name) == strings.ToLower(name) {
			return &l2
		}

	}

	return nil
}

func (e *Engine) PrepareCommonFieldsEmptyOrderedDict() *ordereddict.Dict {
	o := ordereddict.NewDict()
	o.SetCaseInsensitive()

	for _, v := range e.Common.Ordered_fields {
		o.Set(v, nil)
	}

	return o
}

func (e *Engine) PrepareLayer2FieldsEmptyOrderedDict(l2_name string) *ordereddict.Dict {
	o := ordereddict.NewDict()
	o.SetCaseInsensitive()

	active_l2 := e.FindL2LayerByName(l2_name)

	for _, k := range active_l2.Ordered_fields {
		o.Set(k, nil)
	}

	return o
}

func (e *Engine) PrepareCommonAndLayer2FieldsEmptyOrderedDict(l2_name string) *ordereddict.Dict {

	// Common
	ord_map := e.PrepareCommonFieldsEmptyOrderedDict()

	// Append Active Layer2
	ord_map.MergeFrom(e.PrepareLayer2FieldsEmptyOrderedDict(l2_name))

	return ord_map
}

func (e *Engine) GetCSVHeadersOrdered(l2_name string) []string {
	ord_map := e.PrepareCommonAndLayer2FieldsEmptyOrderedDict(l2_name)
	return common.OrderedDictToKeysOrderedStringList(ord_map)
}

func (e *Engine) ParseCommonFieldsOrderedDict(ev_map *eventmap.EventMap) *ordereddict.Dict {
	ord_map := e.PrepareCommonFieldsEmptyOrderedDict()

	//EventTime
	if _, eventtimeok := ord_map.Get("EventTime"); eventtimeok {
		ord_map.Update("EventTime", eventmap.GetSystemTime(ev_map, e.Common.Options["HighPrecisionEventTime"]))
	}

	//EID
	if _, eidok := ord_map.Get("EID"); eidok {
		ord_map.Update("EID", eventmap.GetEID(ev_map))
	}

	//EID [Description]
	if _, eidok := ord_map.Get("Description"); eidok {
		ord_map.Update("Description", e.GetEIDDescription(ev_map))
	}

	//Computer
	if _, current_computer_ok := ord_map.Get("Computer"); current_computer_ok {
		ord_map.Update("Computer", eventmap.GetCurrentComputer(ev_map))
	}

	//Channel
	if _, channel_ok := ord_map.Get("Channel"); channel_ok {
		ord_map.Update("Channel", eventmap.GetChannel(ev_map))
	}

	//Provider
	if _, provider_ok := ord_map.Get("Provider"); provider_ok {
		ord_map.Update("Provider", eventmap.GetProvider(ev_map))
	}

	//EventRecord ID
	if _, erid_ok := ord_map.Get("EventRecord ID"); erid_ok {
		ord_map.Update("EventRecord ID", eventmap.GetEventRecordID(ev_map))
	}

	//System Process ID
	if _, sysprocid_ok := ord_map.Get("System Process ID"); sysprocid_ok {
		ord_map.Update("System Process ID", eventmap.GetSystemProcessID(ev_map))
	}

	//Security User ID
	if _, secuid_ok := ord_map.Get("Security User ID"); secuid_ok {
		ord_map.Update("Security User ID", eventmap.GetSecurityUserID(ev_map))
	}

	return ord_map
}

func (e *Engine) ParseL2FieldsOrderedDict(l2_name string, ev_map *eventmap.EventMap) *ordereddict.Dict {

	channel := eventmap.GetChannel(ev_map)
	eid := eventmap.GetEID(ev_map)
	l2_current := e.FindL2LayerByName(l2_name)

	// Empty dict with fields
	ord_map := e.PrepareLayer2FieldsEmptyOrderedDict(l2_name)

	// Attrib extraction - RAW data types
	attrib_map := eventmap.ExtractAttribs(ev_map, e.UnifiedChannelEvents[channel].EIDs[eid].Attrib_extraction)

	// Convert to string type
	eventmap.MapAttribToOrderedMap(attrib_map, ord_map, l2_current.Fields_remap_dict, l2_current.Ordered_fields_enhanced)

	// Resolve - Mappers & Double Quotes (Optional)
	eventmap.ResolveMappersAndDoubleQuotesInPlace(ord_map, l2_current.Ordered_fields_enhanced, e.VariousMappers, e.GetDoubleQuotesForChannel(channel))

	// Special transformations
	if len(l2_current.Field_extra_transformations) > 0 {
		eventmap.ApplySpecialTransformations(ord_map, l2_current.Field_extra_transformations)
	}

	return ord_map
}

func (e *Engine) SingleFieldExtractor(function string) common.SingleField {
	var sf = common.SingleField{
		NiceName: "",
		Options:  make(map[string]string, 0),
	}

	temp1 := strings.SplitN(function, ":", 2)

	// Set name
	sf.NiceName = temp1[0]

	// Optional options
	if len(temp1) > 1 {
		remaining := temp1[1]
		// Options separated by ,
		temp2 := strings.Split(remaining, ",")

		for _, option := range temp2 {
			opt_split := strings.Split(option, "=")

			if len(opt_split) != 2 {
				panic("SingleFieldExtractor - wrong nr of fields after = split")
			}
			sf.Options[opt_split[0]] = opt_split[1]
		}
	}

	return sf
}

func (e *Engine) AppendLayer1ToUnifiedChannelEvents(l1 *Layer1) {
	channel := l1.Info.Channel

	// Create if not exists
	if _, exists := e.UnifiedChannelEvents[channel]; !exists {
		e.UnifiedChannelEvents[channel] = NewUnifiedChannelEvents()
	}

	// Iterate over current 1st layer events
	for eid, l1events := range l1.Events {
		e.UnifiedChannelEvents[channel].EIDs[eid] = NewLayer1EventsEnhanced(l1events.Attrib_extraction, l1events.Short_description)
	}

}

func (e *Engine) GetEIDDescription(ev_map *ordereddict.Dict) string {
	eid := eventmap.GetEID(ev_map)
	channel := eventmap.GetChannel(ev_map)

	return e.UnifiedChannelEvents[channel].EIDs[eid].Short_description
}

func (e *Engine) LoadParams() {

	params_dir := e.Maps_path + "params" + string(os.PathSeparator)
	params_files, err := ioutil.ReadDir(params_dir)

	if err != nil {
		common.LogCriticalErrorWithError("When reading params directory", err)
	}

	common.LogDebugStructure("Params files", params_files, "params_files")

	// Layer2 - maps/params/*.yaml
	for _, f := range params_files {

		// Skip directories
		if f.IsDir() {
			continue
		}

		file_bytes, err := ioutil.ReadFile(params_dir + f.Name())

		if err != nil {
			common.LogCriticalErrorWithError("When reading params: "+f.Name(), err)
		}

		var p = new(common.Params)
		err = yaml.Unmarshal(file_bytes, p)

		if err != nil {
			common.LogCriticalErrorWithError("When unmarshalling params: "+f.Name(), err)
		}

		if p.Info.Typ == "doublequotes" {
			channel := p.Info.Channel
			e.DoubleQuotes[channel] = *p
			common.LogDebug("Loaded params (doublequotes) map: " + f.Name())
		} else if p.Info.Typ == "mapper_number_to_string" || p.Info.Typ == "mapper_string_to_string" || p.Info.Typ == "mapper_bitwise_to_string" {
			name := p.Info.Name
			e.VariousMappers[name] = *p
			common.LogDebug("Loaded params (mapper) map: " + f.Name())
		} else {
			panic("YAML - LoadLayer2() - Unsupported Info.Typ")
		}
	}

}

func (e *Engine) GetDoubleQuotesForChannel(channel string) map[string]string {
	if p, exists := e.DoubleQuotes[channel]; exists {
		return p.Params
	}

	return nil
}
