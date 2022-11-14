package tests

import (
	"testing"
)

func TestGeneralBootupRestartShutdownEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 16
	easytesting16 := NewEasyTesting(t, UnmarshallAndParseEvent("System_Microsoft-Windows-Kernel-Boot_16.json", eng, "General_BootupRestartShutdown"))
	easytesting16.CheckDesiredValue("ReasonCode", "3221225473")
	easytesting16.CheckDesiredValue("Reason", "A fatal error occurred processing the restoration data.\r\n")

	// 1001
	easytesting1001 := NewEasyTesting(t, UnmarshallAndParseEvent("System_Microsoft-Windows-WER-SystemErrorReporting_1001.json", eng, "General_BootupRestartShutdown"))
	easytesting1001.CheckDesiredValue("Bugcheck", "0x0000009f (0x0000000000000003, 0xffff940cc8dd1050, 0xfffff805597e0110, 0xffff940cf303d010)")
	easytesting1001.CheckDesiredValue("DumpPath", "C:\\Windows\\MEMORY.DMP")
	easytesting1001.CheckDesiredValue("ReportID", "3cab8006-21d7-30dc-5298-1a0b621c2383")

	// 1073
	easytesting1073 := NewEasyTesting(t, UnmarshallAndParseEvent("System_User32_1073.json", eng, "General_BootupRestartShutdown"))
	easytesting1073.CheckDesiredValue("SourceComputer", "DESKTOP-EvtxHussar")
	easytesting1073.CheckDesiredValue("SubjectUserName", "HUSS\\yar0")

	// 1074
	easytesting1074 := NewEasyTesting(t, UnmarshallAndParseEvent("System_User32_1074.json", eng, "General_BootupRestartShutdown"))
	easytesting1074.CheckDesiredValue("ProcessName", "C:\\Windows\\system32\\winlogon.exe (DESKTOP-EvtxHussar)")
	easytesting1074.CheckDesiredValue("SourceComputer", "WIN-OXY")
	easytesting1074.CheckDesiredValue("Reason", "Operating System: Upgrade (Planned)")
	easytesting1074.CheckDesiredValue("ReasonCode", "0x80020003")
	easytesting1074.CheckDesiredValue("Type", "restart")
	easytesting1074.CheckDesiredValue("Comment", "")
	easytesting1074.CheckDesiredValue("SubjectUserName", "NT AUTHORITY\\SYSTEM")

	// 4674
	easytesting4674 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4674.json", eng, "General_BootupRestartShutdown"))
	easytesting4674.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4674.CheckDesiredValue("SubjectUserName", "yar0")
	easytesting4674.CheckDesiredValue("SubjectDomainName", "HUSS")
	easytesting4674.CheckDesiredValue("SubjectLogonId", "0x104a6ff6a")
	easytesting4674.CheckDesiredValue("ObjectServer", "Win32 SystemShutdown module")
	easytesting4674.CheckDesiredValue("ObjectType", "-")
	easytesting4674.CheckDesiredValue("ObjectName", "-")
	easytesting4674.CheckDesiredValue("HandleId", "0x0")
	easytesting4674.CheckDesiredValue("AccessMask", "0")
	easytesting4674.CheckDesiredValue("PrivilegeList", "SeShutdownPrivilege")
	easytesting4674.CheckDesiredValue("ProcessId", "0x228")
	easytesting4674.CheckDesiredValue("ProcessName", "C:\\Windows\\System32\\wininit.exe")

	// 6008
	easytesting6008 := NewEasyTesting(t, UnmarshallAndParseEvent("System_6008.json", eng, "General_BootupRestartShutdown"))
	easytesting6008.CheckDesiredValue("StopTime", "11/24/2021 11:26:19 AM")

	// 6013
	easytesting6013 := NewEasyTesting(t, UnmarshallAndParseEvent("System_6013.json", eng, "General_BootupRestartShutdown"))
	easytesting6013.CheckDesiredValue("Uptime", "329198")
}
