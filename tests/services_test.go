package tests

import (
	"testing"
)

func TestServicesSecurityEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 4697 Version 1
	easytesting4697 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4697.json", eng, "ServicesUniversal"))
	easytesting4697.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4697.CheckDesiredValue("SubjectUserName", "DESKTOP-EvtxHussar$")
	easytesting4697.CheckDesiredValue("SubjectDomainName", "HUSSAR")
	easytesting4697.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4697.CheckDesiredValue("ServiceName", "DepSvc_2158")
	easytesting4697.CheckDesiredValue("ImagePath/ServiceFileName", "C:\\WINDOWS\\system32\\svchost.exe -k DepSvcGroup -p")
	easytesting4697.CheckDesiredValue("ServiceType", "User share process (instance) [Original value: 0xE0]")
	easytesting4697.CheckDesiredValue("ServiceStartType", "Demand start [Original value: 3]")
	easytesting4697.CheckDesiredValue("ServiceAccount/AccountName", "LocalSystem")
	easytesting4697.CheckDesiredValue("ClientProcessStartKey", "1207372833353327")
	easytesting4697.CheckDesiredValue("ClientProcessId", "2452")
	easytesting4697.CheckDesiredValue("ParentProcessId", "888")
}

func TestServicesSystemEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 7024
	easytesting7024 := NewEasyTesting(t, UnmarshallAndParseEvent("System_7024.json", eng, "ServicesUniversal"))
	easytesting7024.CheckDesiredValue("ServiceName", "Hipsotek Evolution")
	easytesting7024.CheckDesiredValue("Error", "%%3")

	// 7034
	easytesting7034 := NewEasyTesting(t, UnmarshallAndParseEvent("System_7034.json", eng, "ServicesUniversal"))
	easytesting7034.CheckDesiredValue("ServiceName", "Extra protection service")

	// 7040
	easytesting7040 := NewEasyTesting(t, UnmarshallAndParseEvent("System_7040.json", eng, "ServicesUniversal"))
	easytesting7040.CheckDesiredValue("ServiceName", "Background Intelligent Transfer Service")
	easytesting7040.CheckDesiredValue("ServiceStartTypeOld", "auto start")
	easytesting7040.CheckDesiredValue("ServiceStartTypeNew", "demand start")
	easytesting7040.CheckDesiredValue("ExtraServiceName", "BITS")

	// 7045
	easytesting7045 := NewEasyTesting(t, UnmarshallAndParseEvent("System_7045.json", eng, "ServicesUniversal"))
	easytesting7045.CheckDesiredValue("ServiceName", "Adobe Light Filter")
	easytesting7045.CheckDesiredValue("ImagePath/ServiceFileName", "\\SystemRoot\\system32\\DRIVERS\\alf.sys")
	easytesting7045.CheckDesiredValue("ServiceStartType", "system start")
	easytesting7045.CheckDesiredValue("ServiceAccount/AccountName", "")

	// 7046
	easytesting7046 := NewEasyTesting(t, UnmarshallAndParseEvent("System_7046.json", eng, "ServicesUniversal"))
	easytesting7046.CheckDesiredValue("ServiceName", "Windows Firewall")
}
