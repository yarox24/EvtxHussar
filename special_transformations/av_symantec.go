package special_transformations

import (
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/common"
	"regexp"
)

func AVSymantecExtract(ord_map *ordereddict.Dict, options map[string]string) {
	Input_field := options["input_field"]
	Output_field := options["output_field"]
	Scope := options["scope"]

	if !common.KeyExistsInOrderedDict(ord_map, Input_field) {
		panic("Wrong Yaml - field_extra_transformations - input_field")
	}

	if !common.KeyExistsInOrderedDict(ord_map, Output_field) {
		panic("Wrong Yaml - field_extra_transformations - output_field")
	}

	input_val, _ := ord_map.GetString(Input_field)
	if len(input_val) > 0 {
		if Scope == "description_path" {
			ord_map.Update(Output_field, description_path_parser(input_val))
		} else {
			panic("Wrong Yaml - field_extra_transformations - scope")
		}
	}
}

func description_path_parser(value string) string {
	paramsMap := make(map[string]string)

	// Regex 1
	var re = regexp.MustCompile(`(?i)(?:application|Application path): (?P<path>.*)`)
	match := re.FindStringSubmatch(value)

	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	// Regex 2
	var re2 = regexp.MustCompile(`(?i)'(?P<path>[^']+)'`)
	match2 := re2.FindStringSubmatch(value)

	for i, name := range re2.SubexpNames() {
		if i > 0 && i <= len(match2) {
			paramsMap[name] = match2[i]
		}
	}

	if val, exists := paramsMap["path"]; exists {
		return val
	} else {
		return ""
	}

}
