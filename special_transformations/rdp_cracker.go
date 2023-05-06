package special_transformations

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/common"
	"strings"
	"unicode/utf16"
)

// CURENTLY DISABLED!!!!
func build_hash_for_username(username string) string {
	usernameBytes := utf16.Encode([]rune(username))
	usernameBinary := make([]byte, len(usernameBytes)*2)
	for index, value := range usernameBytes {
		binary.LittleEndian.PutUint16(usernameBinary[index*2:], value)
	}
	hash := sha256.Sum256(usernameBinary)
	encodedHash := base64.StdEncoding.EncodeToString(hash[:])
	return encodedHash
}

// CURENTLY DISABLED!!!!
// Event ID 1029 - Microsoft-Windows-TerminalServices-RDPClient/Operational.evtx
func RDP1029DetermineUsername(ord_map *ordereddict.Dict, options map[string]string) {
	Input_field := options["input_field"]
	Output_field := options["output_field"]

	if !common.KeyExistsInOrderedDict(ord_map, Input_field) {
		panic("Wrong Yaml - field_extra_transformations - input_field")
	}

	if !common.KeyExistsInOrderedDict(ord_map, Output_field) {
		panic("Wrong Yaml - field_extra_transformations - output_field")
	}

	input_val, _ := ord_map.GetString(Input_field)

	if len(input_val) == 0 {
		return
	}

	// If different number of hyphens something is wrong
	if strings.Count(input_val, "-") != 1 {
		common.LogError("Wrong value in Base64_SHA256_UserName field")
		return
	}

	components := strings.Split(input_val, "-")

	domain_hash := ""
	username_hash := ""

	if len(components) == 2 {
		//USERNAME-
		if len(components[1]) == 0 {
			username_hash = components[0]
		} else {
			// DOMAIN-USERNAME
			domain_hash = components[0]
			username_hash = components[1]
		}
	} else {
		common.LogError("Split of Base64_SHA256_UserName field failed")
		return
	}

	//println(input_val)
	println(domain_hash)
	println(username_hash)

	//result := check_dictionary_for_username_presence(generateDictionary(5), username_hash)
	//
	//if len(result) > 0 {
	//	println("Cracked hash! => " + result)
	//}
	//
	//if len(domain_hash) > 0 {
	//	result2 := check_dictionary_for_username_presence(generateDictionary(5), domain_hash)
	//
	//	if len(result2) > 0 {
	//		println("Cracked domain hash! => " + result2)
	//	}
	//}

	//ord_map.Update(Output_field, base64hunter(input_val))

}

// CURENTLY DISABLED!!!!
func check_dictionary_for_username_presence(dictionary []string, expected_hash string) string {

	for _, word := range dictionary {
		calculated_hash := build_hash_for_username(word)
		if strings.Compare(calculated_hash, expected_hash) == 0 {
			return word
		}
	}
	return ""
}

// CURENTLY DISABLED!!!!
func generateDictionary(length int) []string {
	var strings []string
	for i := 'A'; i <= 'Z'; i++ {
		if length == 1 {
			strings = append(strings, string(i))
		} else {
			for _, s := range generateDictionary(length - 1) {
				strings = append(strings, fmt.Sprintf("%c%s", i, s))
			}
		}
	}
	return strings
}
