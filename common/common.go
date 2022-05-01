package common

import (
	"github.com/Velocidex/ordereddict"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	NOT_EXISTS = iota
	IS_FILE
	IS_DIR
)

func Determine_Maps_Path(MapsDirectoryConfig string) (string, error) {

	maps_dir := ""

	// Auto-detection based on binary
	if MapsDirectoryConfig == "" {
		LogDebug("Maps path auto-detection based on executable")
		path, _ := os.Executable()
		LogDebug("Executable path: " + path)
		maps_dir = filepath.Dir(path) + string(filepath.Separator) + "maps" + string(filepath.Separator)
	} else {
		LogDebug("Custom maps directory selected")
		maps_dir = strings.TrimSuffix(MapsDirectoryConfig, string(filepath.Separator))
		maps_dir += string(filepath.Separator)
	}
	LogDebugString("maps_dir", maps_dir)

	maps_layer2_dir := maps_dir + "layer2" + string(filepath.Separator)
	maps_params_dir := maps_dir + "params" + string(filepath.Separator)

	// Check maps is dir
	nr, err := determine_path(maps_dir)

	if nr != IS_DIR {
		LogError("maps_dir is not a directory")
		return "", err
	}

	// Check layer2 is dir
	nr2, err2 := determine_path(maps_layer2_dir)

	if nr2 != IS_DIR {
		LogError("maps/layer2 is not a directory")
		return "", err2
	}

	// Check params is dir
	nr3, err3 := determine_path(maps_params_dir)

	if nr3 != IS_DIR {
		LogError("maps/params is not a directory")
		return "", err3
	}

	return maps_dir, nil
}

func determine_path(p string) (int, error) {
	fi, err := os.Stat(p)

	if err != nil {
		return NOT_EXISTS, err
	}

	if fi.Mode().IsRegular() {
		return IS_FILE, nil
	} else if fi.IsDir() {
		return IS_DIR, nil
	}

	panic("Unsupported path")
}

func find_logs_in_directory(p string, recursive bool) []string {
	out := make([]string, 0, 10)

	if recursive {

		wf := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.Mode().IsDir() {
				return nil
			}

			// Filter to .evtx only
			matched, _ := filepath.Match("*.evtx", filepath.Base(strings.ToLower(path)))
			if matched {
				out = append(out, path)
			}

			return nil
		}

		filepath.Walk(p, wf)
	} else {
		files, err := ioutil.ReadDir(p)
		if err != nil {
			LogCriticalErrorWithError("When reading directory: "+p, err)
			return out
		}

		for _, f := range files {

			// Skip directories
			if f.IsDir() {
				continue
			}

			// Filter to .evtx only
			matched, _ := filepath.Match("*.evtx", strings.ToLower(f.Name()))
			if matched {
				out = append(out, path.Join(p, f.Name()))
			}
		}
	}

	return out
}

func UniqueElementsOfSliceInMemory(slice []string) []string {
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}

func OrderedDictToOrderedStringListValues(od *ordereddict.Dict) []string {
	l := make([]string, 0, od.Len())

	for _, key := range od.Keys() {
		val, _ := od.GetString(key)

		l = append(l, val)
	}

	return l
}

func HeadersAndRowListToOrderedDict(keys []string, values []string) *ordereddict.Dict {
	o := ordereddict.NewDict()

	if len(keys) != len(values) {
		panic("HeadersAndRowListToOrderedDict - lenght mismatch")
	}

	for i := 0; i < len(keys); i++ {
		o.Set(keys[i], values[i])
	}
	return o

}

func AppendNumberToPath(path_ string, nr int) string {
	dir, filename := path.Split(path_)
	var extension = filepath.Ext(filename)
	var filename_without_ext = strings.TrimRight(filename, extension)

	return path.Join(dir, filename_without_ext+"_"+strconv.Itoa(nr)+".xlsx")
}

func Min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
