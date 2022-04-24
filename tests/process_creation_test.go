package tests

import (
	"testing"
)

func TestProcessCreationSecurityEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 4688 Version 2
	easytesting4688 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4688.json", eng, "ProcessCreation"))
	easytesting4688.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4688.CheckDesiredValue("SubjectUserName", "DESKTOP-EvtxHussar$")
	easytesting4688.CheckDesiredValue("SubjectDomainName", "HUSSAR")
	easytesting4688.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4688.CheckDesiredValue("ProcessId", "0x794")
	easytesting4688.CheckDesiredValue("ProcessName", "C:\\Windows\\SysWOW64\\cmd.exe")
	easytesting4688.CheckDesiredValue("TokenElevationType", "TokenElevationTypeDefault (1)")
	easytesting4688.CheckDesiredValue("ParentProcessId", "0x7e1")
	easytesting4688.CheckDesiredValue("CommandLine", "C:\\Windows\\SysWOW64\\cmd.exe dosomething")
	easytesting4688.CheckDesiredValue("TargetUserSid", "S-1-0-0")
	easytesting4688.CheckDesiredValue("TargetUserName", "-")
	easytesting4688.CheckDesiredValue("TargetDomainName", "-")
	easytesting4688.CheckDesiredValue("TargetLogonId", "0x0")
	easytesting4688.CheckDesiredValue("ParentProcessName", "C:\\Program Files\\VNC\\vnc.exe")
	easytesting4688.CheckDesiredValue("MandatoryLabel", "System integrity [Original value: S-1-16-16384]")

	// 4689
	easytesting4689 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4689.json", eng, "ProcessCreation"))
	easytesting4689.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4689.CheckDesiredValue("SubjectUserName", "DESKTOP-EvtxHussar$")
	easytesting4689.CheckDesiredValue("SubjectDomainName", "HUSSAR")
	easytesting4689.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4689.CheckDesiredValue("Status", "0")
	easytesting4689.CheckDesiredValue("ProcessId", "0x32c8")
	easytesting4689.CheckDesiredValue("ProcessName", "C:\\Windows\\System32\\wbem\\WMIC.exe")

}

func TestProcessCreationSysmonEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 1
	easytesting1 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-Sysmon_Operational_1.json", eng, "ProcessCreation"))
	easytesting1.CheckDesiredValue("RuleName", "-")
	easytesting1.CheckDesiredValue("ProcessGuid", "9FA3FE52-C6ED-6257-A33A-000000005F00")
	easytesting1.CheckDesiredValue("ProcessId", "0x1ba8")
	easytesting1.CheckDesiredValue("ProcessName", "C:\\Windows\\Sysmon64.exe")
	easytesting1.CheckDesiredValue("FileVersion", "13.33")
	easytesting1.CheckDesiredValue("Description", "System activity monitor")
	easytesting1.CheckDesiredValue("Product", "Sysinternals Sysmon")
	easytesting1.CheckDesiredValue("Company", "Sysinternals - www.sysinternals.com")
	easytesting1.CheckDesiredValue("OriginalFileName", "-")
	easytesting1.CheckDesiredValue("CommandLine", "C:\\WINDOWS\\Sysmon64.exe")
	easytesting1.CheckDesiredValue("CurrentDirectory", "C:\\WINDOWS\\system32\\")
	easytesting1.CheckDesiredValue("TargetUserName", "ZARZĄDZANIE NT\\SYSTEM")
	easytesting1.CheckDesiredValue("LogonGuid", "9FA3FE52-1ED0-6254-E703-000000000000")
	easytesting1.CheckDesiredValue("TargetLogonId", "0x3e7")
	easytesting1.CheckDesiredValue("TerminalSessionId", "0")
	easytesting1.CheckDesiredValue("IntegrityLevel", "System")
	easytesting1.CheckDesiredValue("Hashes", "SHA256=7ABCC0EDB4C0F9F47A1BAC0D06401504E1E91C4B1A4E679F01AD539F364D6881")
	easytesting1.CheckDesiredValue("ParentProcessGuid", "9FA3FE52-1ED0-6254-0B00-000000005F00")
	easytesting1.CheckDesiredValue("ParentProcessId", "0x798")
	easytesting1.CheckDesiredValue("ParentProcessName", "C:\\Windows\\System32\\services.exe")
	easytesting1.CheckDesiredValue("ParentCommandLine", "C:\\WINDOWS\\system32\\services.exe")
	easytesting1.CheckDesiredValue("SubjectUserName", "ZARZĄDZANIE NT\\SYSTEM")

	// 5
	easytesting5 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-Sysmon_Operational_5.json", eng, "ProcessCreation"))
	easytesting5.CheckDesiredValue("RuleName", "-")
	easytesting5.CheckDesiredValue("ProcessGuid", "9FA3FE52-C6EB-6257-9D3A-000000005F00")
	easytesting5.CheckDesiredValue("ProcessId", "0x67d4")
	easytesting5.CheckDesiredValue("ProcessName", "F:\\Sysmon\\Sysmon64.exe")
	easytesting5.CheckDesiredValue("TargetUserName", "DESKTOP-EvtxHussar\\Hussar")

}
