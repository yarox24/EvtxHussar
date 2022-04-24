package eventmap

import (
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/evtxhussar/common"
	"strconv"
	"strings"
)

func content_data_autonumbering(o *ordereddict.Dict, opt map[string]string) *ordereddict.Dict {

	data_strings, is_slice := o.GetStrings("Data")

	if is_slice {
		for i, s := range data_strings {
			o.Set("autonumbered"+strconv.Itoa(i), s)
		}
	}

	return o
}

func content_data_autonaming(o *ordereddict.Dict, opt map[string]string) *ordereddict.Dict {

	for _, k := range o.Keys() {
		val, _ := o.Get(k)

		if val_str, is_str := val.(string); is_str {

			//o.Set(k, strings.TrimSpace(val_str))
			o.Set(k, val_str)
		} else if dict, is_dict := val.(*ordereddict.Dict); is_dict {
			if len(dict.Keys()) == 2 {
				name, is_name := dict.GetString("NiceName")
				value, is_value := dict.Get("Value")

				if is_name && is_value {
					o.Set(name, value)
				} else {
					common.LogError("content_data_autonaming - is_name && is_value")
				}
			} else {
				common.LogError("content_data_autonaming - len(dict.Keys())")
			}
		} else {
			o.Set(k, val)
		}

	}

	return o
}

func userdata_flatten_first_value(o *ordereddict.Dict, options map[string]string) *ordereddict.Dict {

	// Find first key
	keys := o.Keys()

	if len(keys) != 1 {
		panic("userdata_flatten_first_value - wrong number of keys")
	}

	if len(keys) > 0 {
		first_key, _ := o.Get(keys[0])

		switch v := first_key.(type) {
		case *ordereddict.Dict:
			o.MergeFrom(v)
		default:
			panic("userdata_flatten_first_value - Wrong type")
		}
	}

	return o
}

func split_by_char_and_equal(o *ordereddict.Dict, opt map[string]string) *ordereddict.Dict {
	input_field, found_key := o.GetString(opt["input_field"])

	split_char := strings.ReplaceAll(opt["split_char"], "\"", "")

	if found_key {
		tab_split := strings.Split(input_field, split_char)

		for _, ss := range tab_split {
			if len(strings.TrimSpace(ss)) == 0 {
				continue
			}

			key_value_split := strings.SplitN(ss, "=", 2)
			if len(key_value_split) != 2 {
				common.LogDebug("split_by_tab_and_equal - Error")
				continue
				//return o
			}

			key := strings.TrimSpace(key_value_split[0])
			value := strings.TrimSpace(key_value_split[1])

			o.Set(key, value)
		}
	}

	return o
}

func rename_field(o *ordereddict.Dict, opt map[string]string) *ordereddict.Dict {
	// Previous value
	val, ok := o.Get(opt["input_field"])

	if ok {
		o.Set(opt["output_field"], val)
		o.Delete(opt["input_field"])
	}

	return o
}

func remove_key(o *ordereddict.Dict, opt map[string]string) *ordereddict.Dict {
	o.Delete(opt["input_field"])
	return o
}
