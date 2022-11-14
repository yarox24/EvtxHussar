package tests

import (
	"testing"
)

func TestCommonFields(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// System104
	easytesting104 := NewEasyTesting(t, UnmarshallAndParseEvent("System_104.json", eng, "AuditLogCleared"))
	easytesting104.CheckDesiredValue("EventTime", "2021.12.06 13:40:00.11111")
	easytesting104.CheckDesiredValue("EID", "104")
	easytesting104.CheckDesiredValue("Description", "The log file was cleared.")
	easytesting104.CheckDesiredValue("Computer", "DESKTOP-EvtxHussar")
	easytesting104.CheckDesiredValue("Channel", "System")
	easytesting104.CheckDesiredValue("Provider", "Microsoft-Windows-Eventlog")
	easytesting104.CheckDesiredValue("EventRecord ID", "555")
	easytesting104.CheckDesiredValue("System Process ID", "1232")
	easytesting104.CheckDesiredValue("Security User ID", "S-1-5-19")

	// Security 4674 (Keywords applied only to Security)
	easytesting4674 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4674.json", eng, "LogonsUniversal"))
	easytesting4674.CheckDesiredValue("Keywords", "Audit Success")
}
