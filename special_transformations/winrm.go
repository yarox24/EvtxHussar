package special_transformations

import (
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/common"
	"regexp"
	"strings"
)

func WinRMStringExtract(ord_map *ordereddict.Dict, options map[string]string) {
	Input_field := options["input_field"]
	Output_field := options["output_field"]
	Extract_part := options["extract_part"]

	if !common.KeyExistsInOrderedDict(ord_map, Input_field) {
		panic("Wrong Yaml - field_extra_transformations - input_field")
	}

	if !common.KeyExistsInOrderedDict(ord_map, Output_field) {
		panic("Wrong Yaml - field_extra_transformations - output_field")
	}

	input_val, _ := ord_map.GetString(Input_field)
	if len(input_val) > 0 {
		ord_map.Update(Output_field, connection_string_parser(input_val, Extract_part))
	}
}

func connection_string_parser(value string, extract_part string) string {
	var re = regexp.MustCompile(`(?mi)(?P<hostname>.*?)/wsman\??(?P<params>.*)`)
	paramsMap := make(map[string]string)
	match := re.FindStringSubmatch(value)

	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	if len(paramsMap["params"]) > 0 {

		// Split by;
		psplit := strings.Split(paramsMap["params"], ";")

		for _, pslipt_param := range psplit {
			equal_split := strings.Split(pslipt_param, "=")

			if len(equal_split) == 2 {
				paramsMap[strings.ToLower(strings.TrimSpace(equal_split[0]))] = strings.TrimSpace(equal_split[1])
			} else {
				common.LogError("Parsing Wsman string error - connection_string_parser")
			}
		}

	}
	if val, exists := paramsMap[extract_part]; exists {
		return val
	} else {
		return ""
	}

}
