package tests

import (
	"testing"
)

func TestAuditLogClearedSystemEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 104
	easytesting104 := NewEasyTesting(t, UnmarshallAndParseEvent("System_104.json", eng, "AuditLogCleared"))
	easytesting104.CheckDesiredValue("SubjectUserName", "SYSTEM")
	easytesting104.CheckDesiredValue("SubjectDomainName", "NT AUTHORITY")
	easytesting104.CheckDesiredValue("Channel", "System")
	easytesting104.CheckDesiredValue("BackupPath", "sample_path")

}

func TestAuditLogClearedSecurityEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 1102
	easytesting1102 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_1102.json", eng, "AuditLogCleared"))
	easytesting1102.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting1102.CheckDesiredValue("SubjectUserName", "Administrator")
	easytesting1102.CheckDesiredValue("SubjectDomainName", "DESKTOP-EvtxHussar")
	easytesting1102.CheckDesiredValue("SubjectLogonId", "0x7c5")

}
