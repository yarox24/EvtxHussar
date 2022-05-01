package common

import (
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
)

//DebugLevel Level
//InfoLevel
//WarnLevel
//ErrorLevel
//FatalLevel
//PanicLevel
//NoLevel
//Disabled

func Logging_init(Debug bool) {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if runtime.GOOS == "windows" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out: colorable.NewColorableStdout(),
		})
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out: os.Stdout,
		})
	}

}

func LogInfo(msg string) {
	log.Info().Msg(msg)
}

func LogWarn(msg string) {
	log.Warn().Msg(msg)
}

func LogImprove(msg string) {
	log.Debug().Msg("[Improve] " + msg)
}

func GrabUniversalLogger() *zerolog.Logger {
	return &log.Logger
}

func LogInfoStructure(msg string, s interface{}, s_name string) {
	log.Info().Interface(s_name, s).Msg(msg)
}

func LogDebug(msg string) {
	log.Debug().Msg(msg)
}

func LogDebugString(key string, val interface{}) {
	switch val.(type) {
	case string:
		log.Debug().Str(key, val.(string)).Msg("Variable")
	case int:
		log.Debug().Int(key, val.(int)).Msg("Variable")
	default:
		panic("LogDebugString - Unsupported type")
	}

}

func LogCriticalError(msg string) {
	log.Error().Msg(msg)
	os.Exit(-1)
}

func LogCriticalErrorWithError(msg string, e error) {
	log.Error().Err(e).Msg(msg)
	os.Exit(5)
}

func LogError(msg string) {
	log.Error().Msg(msg)
}

func LogErrorWithError(msg string, e error) {
	log.Error().Err(e).Msg(msg)
}

func LogDebugStructure(msg string, s interface{}, s_name string) {
	log.Debug().Interface(s_name, s).Msg(msg)
}
