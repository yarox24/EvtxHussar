package tests

import (
	"testing"
)

func TestCommonFields(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// System104
	easytesting104 := NewEasyTesting(t, UnmarshallAndParseEvent("System_104.json", eng, "AccountsUserRelatedOperations"))
	easytesting104.CheckDesiredValue("EventTime", "2022.05.19 13:30:57.3431922")
	easytesting104.CheckDesiredValue("EID", "104")
	easytesting104.CheckDesiredValue("Description", "The log file was cleared.")
	easytesting104.CheckDesiredValue("Computer", "DESKTOP-EvtxHussar")
	easytesting104.CheckDesiredValue("Channel", "System")
	easytesting104.CheckDesiredValue("Provider", "Microsoft-Windows-Eventlog")
	easytesting104.CheckDesiredValue("EventRecord ID", "555")
	easytesting104.CheckDesiredValue("System Process ID", "1232")
	easytesting104.CheckDesiredValue("Security User ID", "S-1-5-19")

}
