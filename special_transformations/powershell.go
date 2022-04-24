package special_transformations

import (
	"encoding/base64"
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/common"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"strings"
)

func Base64powershellhunter(ord_map *ordereddict.Dict, options map[string]string) {
	Input_field := options["input_field"]
	Output_field := options["output_field"]

	if !common.KeyExistsInOrderedDict(ord_map, Input_field) {
		panic("Wrong Yaml - field_extra_transformations - input_field")
	}

	if !common.KeyExistsInOrderedDict(ord_map, Output_field) {
		panic("Wrong Yaml - field_extra_transformations - output_field")
	}

	input_val, _ := ord_map.GetString(Input_field)
	ord_map.Update(Output_field, base64hunter(input_val))

}

func base64hunter(s string) string {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz{}"

	// Remove all caret characters
	temp := strings.ReplaceAll(s, "^", "")

	// Remove false positives
	temp = strings.ReplaceAll(s, "-ExecutionPolicy", "")

	// Potentially encoded base64
	if strings.Contains(strings.ToLower(temp), strings.ToLower("-e")) {

		sp := strings.Split(s, " ")

		// Sort by longest??

		for _, v := range sp {
			b, err := base64.StdEncoding.DecodeString(v)
			if err == nil {

				bs_UTF16LE, _, err2 := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), b)
				if err2 == nil {
					out := string(bs_UTF16LE)

					// Quick check if we don't have garbage on out
					if strings.ContainsAny(out, alphabet) {
						return out
					}
				}
			}
		}
	}

	return ""
}
