package common

import "github.com/Velocidex/ordereddict"

type SingleField struct {
	NiceName string
	Options  map[string]string
}

type Layer2FieldExtraTransformations struct {
	Input_field       string
	Output_field      string
	Special_transform string
	Options           map[string]string
}

type Params2Info struct {
	Typ              string
	Channel          string
	Name             string
	Display_original string
}

type Params struct {
	Info          Params2Info
	Params        map[string]string
	Params_number map[int]string
}

func KeyExistsInOrderedDict(od *ordereddict.Dict, key string) bool {
	_, key_exists := od.Get(key)
	return key_exists
}

func OrderedDictToKeysOrderedStringList(ord_map *ordereddict.Dict) []string {
	final_list := make([]string, 0, ord_map.Len())

	for _, k := range ord_map.Keys() {
		final_list = append(final_list, k)
	}
	return final_list
}

type PowerShellScriptBlockInfo struct {
	Total    int
	Segments map[int]string
	Path     string
}
