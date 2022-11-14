package eventmap

import (
	"fmt"
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/common"
	"github.com/yarox24/EvtxHussar/special_transformations"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type EventMap = ordereddict.Dict

func GetSystemTime(ev_map *ordereddict.Dict, highprecisioneventtime string) string {
	temp, _ := ordereddict.GetAny(ev_map, "System.TimeCreated.SystemTime")
	temp_float64, _ := temp.(float64)

	temp_time := common.ToTime(temp_float64)
	return common.SysTimeToString(temp_time, strings.ToLower(highprecisioneventtime) == "true")
}

func GetEID(ev_map *ordereddict.Dict) string {
	temp, _ := ordereddict.GetInt(ev_map, "System.EventID.Value")
	return strconv.Itoa(temp)
}

func GetCurrentComputer(ev_map *ordereddict.Dict) string {
	temp, _ := ordereddict.GetString(ev_map, "System.Computer")
	return temp
}

func GetChannel(ev_map *ordereddict.Dict) string {
	temp, _ := ordereddict.GetString(ev_map, "System.Channel")
	return temp
}
func GetProvider(ev_map *ordereddict.Dict) string {
	temp, _ := ordereddict.GetString(ev_map, "System.Provider.Name")
	return temp
}

func GetProviderGUID(ev_map *ordereddict.Dict) string {
	temp, is_succ := ordereddict.GetString(ev_map, "System.Provider.Guid")

	// Sometimes brackets are present, sometimes not
	if is_succ {
		if strings.Contains(temp, "{") {
			return temp
		} else {
			return "{" + temp + "}"
		}

	}

	return ""
}

func GetProviderName(ev_map *ordereddict.Dict) string {
	temp, _ := ordereddict.GetString(ev_map, "System.Provider.Name")
	return temp
}

func GetEventRecordID(ev_map *ordereddict.Dict) string {
	temp, _ := ordereddict.GetInt(ev_map, "System.EventRecordID")
	return strconv.Itoa(temp)
}

func GetKeywords(ev_map *ordereddict.Dict) string {

	// Apply only for Security logs
	if GetChannel(ev_map) == "Security" {
		kws, is_succ := ordereddict.GetAny(ev_map, "System.Keywords")

		if is_succ {
			// https://learn.microsoft.com/en-us/dotnet/api/system.diagnostics.eventing.reader.standardeventkeywords?view=netframework-4.8

			// Audit Success
			if kws.(uint64)&9007199254740992 > 0 {
				return "Audit Success"
				//Audit Failure
			} else if kws.(uint64)&4503599627370496 > 0 {
				return "Audit Failure"
			}
		}
	}

	return ""
}

func GetCorrelationActivityID(ev_map *ordereddict.Dict) string {
	temp, _ := ordereddict.GetString(ev_map, "System.Correlation.ActivityID")

	// brackets are not present!
	if temp == "" {
		return temp
	} else if strings.Contains(temp, "{") {
		return temp
	} else {
		return "{" + temp + "}"
	}
}

func GetEventRecordIDasNumber(ev_map *ordereddict.Dict) uint64 {
	temp, _ := ordereddict.GetInt(ev_map, "System.EventRecordID")
	return uint64(temp)
}

func GetSystemProcessID(ev_map *ordereddict.Dict) string {
	temp, _ := ordereddict.GetInt(ev_map, "System.Execution.ProcessID")
	return strconv.Itoa(temp)
}
func GetSecurityUserID(ev_map *ordereddict.Dict) string {
	temp, _ := ordereddict.GetString(ev_map, "System.Security.UserID")
	return temp
}

func ExtractAttribs(ev_map *ordereddict.Dict, attrib_extraction []common.ExtractedFunction, l1mode bool) *ordereddict.Dict {
	// Initial output dictionary
	var o = ordereddict.NewDict()
	o.SetCaseInsensitive()

	// Auto-determine content
	content_names := ev_map.Keys()

	if len(content_names) < 2 {
		return o
	}

	if len(content_names) != 2 {
		panic("Interesting case")
	}

	selected_content_name := ""
	for _, k := range content_names {
		if strings.ToLower(k) != "system" {
			selected_content_name = k
		}
	}

	if selected_content_name == "" {
		panic("Cannot auto-determine")
	}

	content, ok_content := ordereddict.GetMap(ev_map, selected_content_name)

	if !ok_content {
		return o
	}

	content.SetCaseInsensitive()

	// Set initial content
	o.MergeFrom(content)

	// Fix Data case (Single value [Name, Value])
	content_keys := o.Keys()
	if len(content_keys) >= 1 && strings.ToLower(content_keys[0]) == "data" {
		data_temp, _ := o.Get(content_keys[0])
		if data_temp_dict, ok := data_temp.(*ordereddict.Dict); ok {
			dict_keys := data_temp_dict.Keys()
			if strings.ToLower(dict_keys[0]) == "name" && strings.ToLower(dict_keys[1]) == "value" {
				key, _ := data_temp_dict.GetString(dict_keys[0])
				value, _ := data_temp_dict.Get(dict_keys[1])
				o.Set(key, value)
			}
		}
	}

	// Execute all functions sequentially
	for _, ef := range attrib_extraction {
		switch ef.Name {
		case "content_data_autonumbering":
			o = content_data_autonumbering(o, ef.Options)
		case "userdata_flatten_first_value":
			o = userdata_flatten_first_value(o, ef.Options)
		case "split_by_char_and_equal":
			o = split_by_char_and_equal(o, ef.Options)
		case "rename_field":
			if !l1mode {
				o = rename_field(o, ef.Options)
			}
		case "append_to_field":
			if !l1mode {
				o = append_to_field(o, ef.Options)
			}
		case "remove_key":
			o = remove_key(o, ef.Options)

		default:
			//fmt.Println("Fake panic")
			panic("Unsupported attrib_extraction function")
		}

	}

	return o
}

func ConvertAllTypesToString(val interface{}, display_as string) string {

	// Assert different types
	switch val.(type) {
	case string:
		return val.(string)
	case bool:
		return strconv.FormatBool(val.(bool))
	case uint8:
		if display_as == "hex" {
			return "0x" + strconv.FormatUint(uint64(val.(uint8)), 16)
		} else {
			return strconv.FormatUint(uint64(val.(uint8)), 10)
		}
	case uint16:
		if display_as == "hex" {
			return "0x" + strconv.FormatUint(uint64(val.(uint16)), 16)
		} else {
			return strconv.FormatUint(uint64(val.(uint16)), 10)
		}
	case uint64:
		if display_as == "hex" {
			return "0x" + strconv.FormatUint(val.(uint64), 16)
		} else {
			return strconv.FormatUint(val.(uint64), 10)
		}
	case uint32:
		if display_as == "hex" {
			return "0x" + strconv.FormatUint(uint64(val.(uint32)), 16)
		} else {
			return strconv.FormatUint(uint64(val.(uint32)), 10)
		}
	case float64:
		if display_as == "utctime" {
			temp_time := common.ToTime(val.(float64))
			return common.SysTimeToString(temp_time, true) + " UTC"
		} else {
			return fmt.Sprintf("%f", val.(float64))
		}
	case []uint8:
		if display_as == "uint8slice_utf-16" {
			bs_UTF16LE, _, err2 := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), val.([]uint8))

			if err2 == nil {
				return strings.TrimRight(string(bs_UTF16LE), "\x00")
			} else {
				common.LogErrorWithError("uint8slice_utf-16 decoding error", err2)
			}
		} else {
			common.LogError(fmt.Sprintf("Wrong format for value: %v", val))
			common.LogError(fmt.Sprintf("Wrong format for value: %s", val))
		}
	default:
		common.LogError(fmt.Sprintf("Wrong format for value: %v", val))
		common.LogError(fmt.Sprintf("Wrong format for value: %s", val))
	}

	return ""
}

func ReadSpecialOptionForKey(Ordered_fields_enhanced map[string]common.SingleField, key string, param string) string {

	if sf, key_exists := Ordered_fields_enhanced[strings.ToLower(key)]; key_exists {
		if param, param_exists := sf.Options[param]; param_exists {
			return param
		}
	}

	return ""
}

//func ReadSpecialOptionForKey(Ordered_fields_enhanced []common.SingleField, key string, param string) string {
//
//	// Special options
//	for _, ofe := range Ordered_fields_enhanced {
//		if strings.ToLower(ofe.NiceName) == strings.ToLower(key) {
//			if param, param_exists := ofe.Options[param]; param_exists {
//				return param
//			}
//		}
//	}
//	return ""
//}

func ExtraFixField(Ordered_fields_enhanced map[string]common.SingleField, key string, function string, value_to_fix string) string {

	switch function {
	case "fix_field":
		{
			fix_option := ReadSpecialOptionForKey(Ordered_fields_enhanced, key, function)

			// Fix field options
			if len(fix_option) > 0 {
				switch fix_option {
				case "cast_singlerune_to_number":
					{
						runes := []rune(value_to_fix)

						if len(runes) == 1 {
							nr := int(runes[0])
							return strconv.Itoa(nr)
						} else {
							//fmt.Println("Fake panic")
							panic("Runes invalid length")
						}
					}
				default:
					panic("Unknown option!")
				}
			}
		}
	default:
		panic("Unknown option!")
	}

	return value_to_fix
}

func MapAttribToOrderedMap(attrib_map *ordereddict.Dict, ord_map *ordereddict.Dict, Fields_remap *ordereddict.Dict, Ordered_fields_enhanced map[string]common.SingleField) {

	// For every attrib param, try to match it to final_map
	for _, key := range attrib_map.Keys() {
		val, _ := attrib_map.Get(key)

		// DIRECTLY
		if _, exists := ord_map.Get(key); exists {
			proper_key_case := common.ProperOrderedMapKeyCase(ord_map, key)
			always_string := ConvertAllTypesToString(val, ReadSpecialOptionForKey(Ordered_fields_enhanced, key, "display_as"))

			// Extra fix
			always_string = ExtraFixField(Ordered_fields_enhanced, key, "fix_field", always_string)

			copy_to_field := ReadSpecialOptionForKey(Ordered_fields_enhanced, key, "copy_raw_value_to_output_field")

			if len(copy_to_field) > 0 {
				ord_map.Update(copy_to_field, always_string)
			}

			// Extra copy
			ord_map.Update(proper_key_case, always_string)
		} else {

			// Remap exists in YAML
			if remmaped_key, remap_ok := Fields_remap.GetString(key); remap_ok {
				always_string := ConvertAllTypesToString(val, ReadSpecialOptionForKey(Ordered_fields_enhanced, remmaped_key, "display_as"))
				if _, ord_key_exists := ord_map.Get(remmaped_key); ord_key_exists {
					proper_remmaped_key_case := common.ProperOrderedMapKeyCase(ord_map, remmaped_key)
					ord_map.Update(proper_remmaped_key_case, always_string)
				}
			}
		}
	}

	// Change every remaining nil to empty string
	for _, key := range ord_map.Keys() {
		val, _ := ord_map.Get(key)
		if val == nil {
			ord_map.Update(key, "")
		}
	}

}

func GetOriginalDisplayValueForMapperStringToString(current_params common.Params, value string) string {
	return fmt.Sprintf(" [Original value: %s]", value)
}

func GetOriginalDisplayValueForMapperNumberToString(current_params common.Params, value string) string {
	nr, err_nr := strconv.Atoi(value)

	if err_nr == nil {
		switch current_params.Info.Display_original {
		case "dec":
			return fmt.Sprintf(" [Original value: %d]", nr)
		case "hex":
			return fmt.Sprintf(" [Original value: 0x%X]", nr)
		case "none":
			return ""
		default:
			panic("Wrong display_original type")
		}
	}

	return ""
}

func GetOriginalDisplayValueForMapperBitwiseToString(current_params common.Params, nr int64) string {

	switch current_params.Info.Display_original {
	case "dec":
		return fmt.Sprintf(" [Original value: %d]", nr)
	case "hex":
		return fmt.Sprintf(" [Original value: 0x%X]", nr)
	case "none":
		return ""
	default:
		panic("Wrong display_original type")
	}

	return ""
}

func ResolveForMapperNumberToString(VariousMappers map[string]common.Params, map_name string, value string, sf_name string) string {

	// Get correct map
	current_params, found_map := VariousMappers[map_name]

	if !found_map {
		panic("Yaml map error")
	}

	// Convert string to int
	nr, err_nr := strconv.Atoi(value)

	if err_nr == nil {
		if nice_name, nice_name_found := current_params.Params_number[nr]; nice_name_found {
			// switch options: hex, dec, none

			return nice_name + GetOriginalDisplayValueForMapperNumberToString(current_params, value)
		} else {
			common.LogImprove(fmt.Sprintf("Append more values to ResolveForMapperNumberToString: %s %s [%s]", map_name, value, sf_name))
		}
	}

	return value
}

func ResolveForMapperStringToString(VariousMappers map[string]common.Params, map_name string, value string, sf_name string) string {

	// Get correct map
	current_params, found_map := VariousMappers[map_name]

	if !found_map {
		panic("Yaml map error")
	}

	if nice_name, nice_name_found := current_params.Params[value]; nice_name_found {
		return nice_name + GetOriginalDisplayValueForMapperStringToString(current_params, value)
	} else {
		common.LogImprove(fmt.Sprintf("Append more values to ResolveForMapperStringToString: %s %s [%s]", map_name, value, sf_name))
	}

	return value
}

func ResolveForMapperBitwiseToString(VariousMappers map[string]common.Params, map_name string, value string, sf_name string) string {

	// Get correct map
	current_params, found_map := VariousMappers[map_name]

	if !found_map {
		panic("Yaml map error")
	}

	var int64_value int64
	var err error

	lower_string := strings.ToLower(value)
	if strings.HasPrefix(lower_string, "0x") {
		int64_value, err = strconv.ParseInt(strings.TrimPrefix(lower_string, "0x"), 16, 0)
	} else {
		int64_value, err = strconv.ParseInt(value, 10, 0)
	}

	if err != nil {
		return value
	}

	active_flags := make([]string, 0)

	// From lowest values
	map_keys := make([]int, 0)
	for k, _ := range current_params.Params_number {
		map_keys = append(map_keys, k)
	}

	sort.Slice(map_keys, func(i, j int) bool { return map_keys[i] < map_keys[j] })

	// Bitwise check
	for _, int_flag := range map_keys {
		description := current_params.Params_number[int_flag]

		if int64(int_flag)&int64_value > 0 {
			active_flags = append(active_flags, description)
		}
	}

	if len(active_flags) > 0 {
		return fmt.Sprintf("%s%s", strings.Join(active_flags, " | "), GetOriginalDisplayValueForMapperBitwiseToString(current_params, int64_value))
	} else {
		if int64_value != 0 {
			common.LogImprove(fmt.Sprintf("Not found value for: %s [%s]", value, sf_name))
		}
	}

	return value
}

func ResolveMappersAndDoubleQuotesInPlace(ord_map *ordereddict.Dict, Ordered_fields_enhanced map[string]common.SingleField, VariousMappers map[string]common.Params, doublequotes map[string]string, SIDList map[string]string) {

	for sf_name, sf := range Ordered_fields_enhanced {
		if len(sf.Options) > 0 {

			for opt_k, opt_v := range sf.Options {
				switch opt_k {
				case "mapper_number_to_string":
					current_val, found_val := ord_map.GetString(sf_name)
					if found_val && len(current_val) > 0 {
						ord_map.Update(sf.NiceName, ResolveForMapperNumberToString(VariousMappers, opt_v, current_val, sf_name))
					}
				case "mapper_string_to_string":
					current_val, found_val := ord_map.GetString(sf_name)
					if found_val && len(current_val) > 0 {
						ord_map.Update(sf.NiceName, ResolveForMapperStringToString(VariousMappers, opt_v, current_val, sf_name))
					}
				case "mapper_bitwise_to_string":
					current_val, found_val := ord_map.GetString(sf_name)
					if found_val && len(current_val) > 0 {
						ord_map.Update(sf.NiceName, ResolveForMapperBitwiseToString(VariousMappers, opt_v, current_val, sf_name))
					}
				case "resolve":
					current_val, found_val := ord_map.GetString(sf_name)

					if found_val && len(current_val) > 0 {
						ord_map.Update(sf.NiceName, ResolveDoubleQuotesInPlace(doublequotes, SIDList, opt_v, current_val, sf_name))
					}
				}
			}
		}

	}
}

func ResolveDoubleQuotesInPlace(double_quotes map[string]string, SIDList map[string]string, opt_v string, current_val string, sf_name string) string {

	if len(double_quotes) == 0 {
		return current_val
	}

	switch opt_v {
	case "doublequotes":
		{
			if strings.Contains(current_val, "%%") {
				var re = regexp.MustCompile(`(?m)(%%\d+)\b`)

				values_slice := common.UniqueElementsOfSliceInMemory(re.FindAllString(current_val, -1))

				temp := current_val
				for _, existing_double_percent := range values_slice {
					if nice_value, nice_value_exists := double_quotes[existing_double_percent]; nice_value_exists {
						temp = strings.ReplaceAll(temp, existing_double_percent, nice_value)
					} else {
						common.LogImprove(fmt.Sprintf("Find value for %s [%s]", current_val, sf_name))
					}
				}
				return temp
			}
		}
	case "doublesids":
		{
			if strings.Contains(current_val, "%{S-") {
				var re = regexp.MustCompile(`(?mi)(%{(?P<sid>S-[^}]*?)})`)

				values_slice := common.UniqueElementsOfSliceInMemory(re.FindAllString(current_val, -1))

				temp := strings.TrimLeft(current_val, "\r\n")
				for _, existing_double_sid := range values_slice {

					extracted_sid := strings.Trim(existing_double_sid, "%{}")
					if nice_value, nice_value_exists := SIDList[extracted_sid]; nice_value_exists {
						temp = strings.ReplaceAll(temp, existing_double_sid, nice_value)
					} else {
						if len(extracted_sid) < 40 {
							common.LogImprove(fmt.Sprintf("Find SID for %s [%s]", existing_double_sid, sf_name))
						}

					}
				}
				return temp
			}
		}
	default:
		panic("Wrong resolve parameter")
	}

	return current_val
}

func ApplySpecialTransformations(ord_map *ordereddict.Dict, Field_extra_transformations []common.Layer2FieldExtraTransformations) {
	opt := make(map[string]string, 0)

	for _, st := range Field_extra_transformations {
		opt["input_field"] = st.Input_field
		opt["output_field"] = st.Output_field

		for k, v := range st.Options {
			opt[k] = v
		}

		if strings.ToLower(st.Special_transform) == "base64powershellhunter" {
			special_transformations.Base64powershellhunter(ord_map, opt)
		} else if strings.ToLower(st.Special_transform) == "xml_scheduled_task" {
			special_transformations.XMLScheduledTask(ord_map, opt)
		} else if strings.ToLower(st.Special_transform) == "winrm_string_extract" {
			special_transformations.WinRMStringExtract(ord_map, opt)
		} else if strings.ToLower(st.Special_transform) == "av_symantec" {
			special_transformations.AVSymantecExtract(ord_map, opt)
		}
	}

}
