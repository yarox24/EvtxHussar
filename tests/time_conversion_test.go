package tests

import (
	"github.com/yarox24/EvtxHussar/common"
	"testing"
)

func CompareTime(t *testing.T, calculated_time_str string, expected_time_str string) {
	if calculated_time_str != expected_time_str {
		t.Errorf("Date mismatch: %s <-> %s", calculated_time_str, expected_time_str)
	}
}

func TestTimeConversion(t *testing.T) {

	// Cut down to precision of 5 digits (Nano)

	// Time 1
	t1_time := common.ToTime(float64(1650697323.305288))

	// Perfect time which won't be achieved: 		2022.04.23 07:02:03.3052876
	// Calculated imperfect time without cutting: 	2022.04.23 07:02:03.305288076
	CompareTime(t, common.SysTimeToString(t1_time, true), "2022.04.23 07:02:03.30528")

	// Time 2
	t2_time := common.ToTime(float64(1634892254.4073904))

	// Perfect time which won't be achieved: 		2021-10-22 08:44:14.4073911
	// Calculated imperfect time without cutting: 	2021.10.22 08:44:14.407390356
	CompareTime(t, common.SysTimeToString(t2_time, true), "2021.10.22 08:44:14.40739")

	// Time 3
	t3_time := common.ToTime(float64(1634892254.0))
	CompareTime(t, common.SysTimeToString(t3_time, true), "2021.10.22 08:44:14.0")

	// Time 4
	t4_time := common.ToTime(float64(1635503568.259472))

	// Perfect time which won't be achieved: 		2021-10-29T10:32:48.2594718
	// Calculated imperfect time without cutting: 	2021-10-29T10:32:48.259471893 Too long, but this time correct?
	CompareTime(t, common.SysTimeToString(t4_time, true), "2021.10.29 10:32:48.25947")
	CompareTime(t, common.SysTimeToString(t4_time, false), "2021.10.29 10:32:48")

}
