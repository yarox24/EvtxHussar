package tests

import (
	"testing"
)

func TestRDPEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine()

	// Microsoft-Windows-TerminalServices-LocalSessionManager_Operational

	// 17
	easytesting17 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-TerminalServices-LocalSessionManager_Operational_17.json", eng, "RDPUniversal"))
	easytesting17.CheckDesiredValue("Status Code", "0x8007045b")

	// 21
	easytesting21 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-TerminalServices-LocalSessionManager_Operational_21.json", eng, "RDPUniversal"))
	easytesting21.CheckDesiredValue("SourceIP", "11.11.11.11")
	easytesting21.CheckDesiredValue("SessionID", "13")
	easytesting21.CheckDesiredValue("User", "HUSS\\buz")

	// 36
	easytesting36 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-TerminalServices-LocalSessionManager_Operational_36.json", eng, "RDPUniversal"))
	easytesting36.CheckDesiredValue("SessionId", "4294967295")
	easytesting36.CheckDesiredValue("State", "0")
	easytesting36.CheckDesiredValue("StateName", "Initialized")
	easytesting36.CheckDesiredValue("Event", "1")
	easytesting36.CheckDesiredValue("EventName", "EvCreated")
	easytesting36.CheckDesiredValue("ErrorCode", "0xd00002fe")

	// Microsoft-Windows-RemoteDesktopServices-RdpCoreTS

	// 65
	easytesting65 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-RemoteDesktopServices-RdpCoreTS_Operational_65.json", eng, "RDPUniversal"))
	easytesting65.CheckDesiredValue("ConnectionName", "RDP-Tcp#0")

	// 131
	easytesting131 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-RemoteDesktopServices-RdpCoreTS_Operational_131.json", eng, "RDPUniversal"))
	easytesting131.CheckDesiredValue("ConnType", "TCP")
	easytesting131.CheckDesiredValue("SourceIP", "11.11.11.11:5555")

	// 168
	easytesting168 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-RemoteDesktopServices-RdpCoreTS_Operational_168.json", eng, "RDPUniversal"))
	easytesting168.CheckDesiredValue("MonitorWidth", "2246")
	easytesting168.CheckDesiredValue("MonitorHeight", "1297")
	easytesting168.CheckDesiredValue("ServerName", "HUSS")

	// Microsoft-Windows-TerminalServices-RDPClient_Operational

	// 226
	easytesting226 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-TerminalServices-RDPClient_Operational_226.json", eng, "RDPUniversal"))
	easytesting226.CheckDesiredValue("Event", "8")
	easytesting226.CheckDesiredValue("EventName", "TsSslEventHandshakeContinueFailed")
	easytesting226.CheckDesiredValue("ErrorCode", "0x80004005")

	// 1024
	easytesting1024 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-TerminalServices-RDPClient_Operational_1024.json", eng, "RDPUniversal"))
	easytesting1024.CheckDesiredValue("TargetIP", "huss.huss.pl")

	// 1027
	easytesting1027 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-TerminalServices-RDPClient_Operational_1027.json", eng, "RDPUniversal"))
	easytesting1027.CheckDesiredValue("DomainName", "HUSS")
	easytesting1027.CheckDesiredValue("SessionId", "4")

	// Microsoft-Windows-TerminalServices-RemoteConnectionManager_Operational
	// 1149
	easytesting1149 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-TerminalServices-RemoteConnectionManager_Operational_1149.json", eng, "RDPUniversal"))
	easytesting1149.CheckDesiredValue("User", "adminadmin")
	easytesting1149.CheckDesiredValue("DomainName", "huss")
	easytesting1149.CheckDesiredValue("SourceIP", "11.11.11.11")

}
