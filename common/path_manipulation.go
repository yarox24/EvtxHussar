package common

import (
	"os"
	"path/filepath"
)

func Handle_output_directory(p string) {

	// Output directory must be provided
	if len(p) == 0 {
		LogCriticalError("Output directory not provided!")
	}

	p, _ = filepath.Abs(p)

	fi, err := os.Stat(p)

	// Directory doesn't exists
	if os.IsNotExist(err) {
		LogDebug("Create output report directory: " + p)
		err2 := os.MkdirAll(p, 0777)

		if err2 != nil {
			LogCriticalError("Cannot create output directory")
		}
	} else if err == nil {

		// Path exists but is not a directory
		if !fi.IsDir() {
			LogCriticalError(p + " exists and it's not a directory")
		}
	}
}

func EnsureDirectoryStructureIsCreated(full_path string) error {
	output_dirs := filepath.Dir(full_path)

	return os.MkdirAll(output_dirs, 0777)
}
