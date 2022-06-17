package common

import "strings"

type ExtractedFunction struct {
	Name    string
	Options map[string]string
}

func FunctionExtractor(function string) ExtractedFunction {
	var ef = ExtractedFunction{
		Name:    "",
		Options: make(map[string]string, 0),
	}

	temp1 := strings.SplitN(function, ":", 2)

	// Set name
	ef.Name = temp1[0]

	// Optional options
	if len(temp1) > 1 {
		remaining := temp1[1]
		// Options separated by ,
		temp2 := strings.Split(remaining, ",")

		for _, option := range temp2 {
			opt_split := strings.Split(option, "=")

			if len(opt_split) != 2 {
				panic("FunctionExtractor - wrong nr of fields after = split")
			}
			ef.Options[opt_split[0]] = opt_split[1]
		}
	}

	return ef
}

type ExtractedLogic struct {
	Method  string
	Options map[string]string
}

func LogicExtractor(logic string) ExtractedLogic {
	var ef = ExtractedLogic{
		Method:  "",
		Options: make(map[string]string, 0),
	}

	temp1 := strings.SplitN(logic, ":", 2)

	// Set method
	ef.Method = temp1[0]

	// Optional options
	if len(temp1) > 1 {
		remaining := temp1[1]
		// Options separated by ,
		temp2 := strings.Split(remaining, ",")

		for _, option := range temp2 {
			opt_split := strings.Split(option, "=")

			if len(opt_split) != 2 {
				panic("ExtractedLogic - wrong nr of fields after = split")
			}
			ef.Options[opt_split[0]] = opt_split[1]
		}
	}

	return ef
}
