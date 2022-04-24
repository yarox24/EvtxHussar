package common

import (
	"bytes"
	"encoding/binary"
	"time"
)

func Float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		LogErrorWithError("binary.Write failed:", err)
	}
	return buf.Bytes()
}

type FiletimeCrossPlatform struct {
	LowDateTime  uint32
	HighDateTime uint32
}

// Nanoseconds returns Filetime ft in nanoseconds
// since Epoch (00:00:00 UTC, January 1, 1970).
func (ft *FiletimeCrossPlatform) Nanoseconds() int64 {
	// 100-nanosecond intervals since January 1, 1601
	nsec := int64(ft.HighDateTime)<<32 + int64(ft.LowDateTime)
	// change starting time to the Epoch (00:00:00 UTC, January 1, 1970)
	nsec -= 116444736000000000
	// convert into nanoseconds
	nsec *= 100
	return nsec
}

// toTime converts an 8-byte Windows Filetime to time.Time.
func ToTime(t []byte) time.Time {

	ft := FiletimeCrossPlatform{
		LowDateTime:  binary.LittleEndian.Uint32(t[:4]),
		HighDateTime: binary.LittleEndian.Uint32(t[4:]),
	}

	return time.Unix(0, ft.Nanoseconds()).UTC()
}

// Time return a golang Time type from the FileTime
//func (ft FileTime) Time() time.Time {
//	ns := (ft.MSEpoch() - 116444736000000000) * 100
//	return time.Unix(0, int64(ns)).UTC()
//}
//
//// MSEpoch returns the FileTime as a Microsoft epoch, the number of 100 nano second periods elapsed from January 1, 1601 UTC.
//func (ft FileTime) MSEpoch() int64 {
//	return (int64(ft.HighDateTime) << 32) + int64(ft.LowDateTime)
//}
