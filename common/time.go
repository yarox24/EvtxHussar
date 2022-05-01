package common

import (
	"math"
	"strconv"
	"strings"
	"time"
)

//func Float64ToByte(f float64) []byte {
//	var buf bytes.Buffer
//	err := binary.Write(&buf, binary.LittleEndian, f)
//	if err != nil {
//		LogErrorWithError("binary.Write failed:", err)
//	}
//	return buf.Bytes()
//}

func ToTime(f float64) time.Time {
	sec := int64(0)
	dec := int64(0)

	sec_f, dec_f := math.Modf(f)
	sec = int64(sec_f)
	dec = int64(dec_f * 1e9)
	return time.Unix(int64(sec), int64(dec))
}

func SysTimeToString(t time.Time, highprecisioneventtime bool) string {
	nanosecond := ""

	if highprecisioneventtime {
		nanosecond += "."
		nano := t.Nanosecond()
		nano_str := strconv.Itoa(nano)
		if len(nano_str) > 5 {
			nano_str = nano_str[:5]
		}
		if len(nano_str) < 5 {
			nano_str += strings.Repeat("0", 5-len(nano_str))
		}
		nanosecond += nano_str
	}

	return t.UTC().Format("2006.01.02 15:04:05") + nanosecond
}
