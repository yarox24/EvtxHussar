package engine

import (
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/common"
	"github.com/yarox24/EvtxHussar/eventmap"
	"strconv"
	"strings"
)

type LogicEngine struct {
	supported_eids ChannelTOEID
}

func (le *LogicEngine) SetSupportedEIDs(supported_eids ChannelTOEID) {
	le.supported_eids = supported_eids
}

func (le *LogicEngine) IsProviderTestPassed(ev_map *ordereddict.Dict, Provider_guid string, Provider_name string) bool {

	// Check Provider GUID if present
	if len(Provider_guid) > 0 {
		provider_guid := eventmap.GetProviderGUID(ev_map)
		if strings.ToLower(Provider_guid) != strings.ToLower(provider_guid) {
			return false
		}
	}

	// Check Provider Name if present
	if len(Provider_name) > 0 {
		provider_name := eventmap.GetProviderName(ev_map)
		if strings.ToLower(Provider_name) != strings.ToLower(provider_name) {
			return false
		}
	}

	return true
}

func (le *LogicEngine) SingleMatchDecimalEqual(event_value interface{}, ExpectedValue string) bool {

	switch v := event_value.(type) {
	case uint64:
		if ExpectedValueUint64, err := strconv.ParseUint(ExpectedValue, 10, 64); err == nil {
			if v == ExpectedValueUint64 {
				return true
			}
		}
	case uint32:
		if ExpectedValueUint32, err := strconv.ParseUint(ExpectedValue, 10, 32); err == nil {
			if v == uint32(ExpectedValueUint32) {
				return true
			}
		}

	default:
		common.LogError("SingleMatchDecimalEqual - unsupported input type")
	}

	return false
}

func (le *LogicEngine) SingleMatchSubstring(event_value interface{}, ExpectedValue string, CaseSensitive bool) bool {

	switch v := event_value.(type) {
	case string:
		if !CaseSensitive {
			v = strings.ToLower(v)
			ExpectedValue = strings.ToLower(ExpectedValue)
		}
		return strings.Contains(v, ExpectedValue)

	default:
		common.LogError("SingleMatchSubstring - unsupported input type")
	}

	return false
}

func (le *LogicEngine) CaseSensiveFlagToBool(CaseSensitive string) bool {
	if CaseSensitive == "0" || strings.ToLower(CaseSensitive) == "false" {
		return false
	}

	return true
}

func (le *LogicEngine) EvaluateRule(rule common.ExtractedLogic, attrib_map *ordereddict.Dict) bool {

	if rule.Method == "single_match" {
		FieldName := rule.Options["Field"]
		ExpectedValue := rule.Options["Value"]
		Function := rule.Options["Function"]

		FieldValue, FieldExists := attrib_map.Get(FieldName)

		if !FieldExists {
			return false
		}

		if Function == "DecimalEqual" {
			return le.SingleMatchDecimalEqual(FieldValue, ExpectedValue)
		} else if Function == "Substring" {

			return le.SingleMatchSubstring(FieldValue, ExpectedValue, le.CaseSensiveFlagToBool(rule.Options["CaseSensitive"]))
		} else {
			common.LogError("EvaluateRule - Unknown Function: " + Function)
		}
	} else {
		common.LogCriticalError("Unknown rule method: " + rule.Method)
	}

	return true
}

func (le *LogicEngine) EvaluateSingleOrContainer(container_rules []common.ExtractedLogic, attrib_map *ordereddict.Dict) bool {

	for _, rule := range container_rules {

		// If at least one rule is true, then everything is true, and calculation can be stop
		if le.EvaluateRule(rule, attrib_map) {
			return true
		}
	}

	return false
}

func (le *LogicEngine) EvaluateSingleAndContainer(container_rules []common.ExtractedLogic, attrib_map *ordereddict.Dict) bool {

	if len(container_rules) == 0 {
		return false
	}

	for _, rule := range container_rules {
		if !le.EvaluateRule(rule, attrib_map) {
			return false
		}
	}

	return true
}

func (le *LogicEngine) CanWeFinishNow(global_logic string, all_containers_result []bool) (bool, bool) {
	can_we_finish := false
	result := false

	if global_logic == "or" {
		for _, single_res := range all_containers_result {
			if single_res {
				result = true
				can_we_finish = true
				break
			}
		}
	}

	if global_logic == "and" {
		if len(all_containers_result) > 0 {
			result = true
		}

		for _, single_res := range all_containers_result {
			if single_res == false {
				result = false
				can_we_finish = true
				break
			}
		}
	}

	return can_we_finish, result
}

func (le *LogicEngine) AreLogicalTestsPassed(ev_map *ordereddict.Dict, MR MatchingRulesT, attrib_extraction []common.ExtractedFunction) bool {

	// Results for all of the containers, but not all of them will be always evaluated
	var all_containers_result = make([]bool, 0)

	// Attrib extraction - RAW data types
	attrib_map := eventmap.ExtractAttribs(ev_map, attrib_extraction, true)

	// OR  containers first
	for _, or_container_rules := range MR.Container_OrEnhanced {
		all_containers_result = append(all_containers_result, le.EvaluateSingleOrContainer(or_container_rules, attrib_map))
		can_we_finish, result := le.CanWeFinishNow(MR.Global_Logic, all_containers_result)
		if can_we_finish {
			return result
		}
	}

	// AND containers next
	for _, and_container_rules := range MR.Container_AndEnhanced {
		all_containers_result = append(all_containers_result, le.EvaluateSingleAndContainer(and_container_rules, attrib_map))
		can_we_finish, result := le.CanWeFinishNow(MR.Global_Logic, all_containers_result)
		if can_we_finish {
			return result
		}
	}

	// Full check, no shortcuts
	_, result := le.CanWeFinishNow(MR.Global_Logic, all_containers_result)

	return result
}

func (le *LogicEngine) ReturnMatchingChanFullInfo(ev_map *ordereddict.Dict) []ChanFullInfo {
	var cfi_list = make([]ChanFullInfo, 0)

	// Validate Event ID
	eid := eventmap.GetEID(ev_map)

	eid_struct, eid_supported := le.supported_eids.Eid[eid]

	if eid_supported {

		// Check every chan if is suitable
		for i := range eid_struct.Chans {
			cfi_curr := eid_struct.Chans[i]

			// Provided GUID or Name check
			if !le.IsProviderTestPassed(ev_map, cfi_curr.Provider_guid, cfi_curr.Provider_name) {
				continue
			}

			// Perform logic filtering
			mr := cfi_curr.Matching_Rules
			if len(mr.Global_Logic) > 0 && (len(mr.Container_OrEnhanced) > 0 || len(mr.Container_AndEnhanced) > 0) {
				if !le.AreLogicalTestsPassed(ev_map, mr, cfi_curr.Attrib_extraction) {
					continue
				}

				// All logical tests passed
				cfi_list = append(cfi_list, cfi_curr)
			} else {
				// No logic rules so we passed all tests
				cfi_list = append(cfi_list, cfi_curr)
			}
		}
	}

	return cfi_list
}

func NewLogicEngine() LogicEngine {
	return LogicEngine{
		supported_eids: ChannelTOEID{},
	}
}

func (MRT *MatchingRulesT) EnhanceRulesInPlace() {

	for _, group := range MRT.Container_Or {
		var new_group = make([]common.ExtractedLogic, 0)
		for _, single_logic_string := range group {
			new_group = append(new_group, common.LogicExtractor(single_logic_string))
		}
		MRT.Container_OrEnhanced = append(MRT.Container_OrEnhanced, new_group)
	}

	// Enhance logic AND
	for _, group := range MRT.Container_And {
		var new_group = make([]common.ExtractedLogic, 0)
		for _, single_logic_string := range group {
			new_group = append(new_group, common.LogicExtractor(single_logic_string))
		}
		MRT.Container_AndEnhanced = append(MRT.Container_AndEnhanced, new_group)
	}
}
