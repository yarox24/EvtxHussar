package tests

import (
	"testing"
)

func TestWindowsFirewallEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 2002
	easytesting2002 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-Windows_Firewall_With_Advanced_Security_Firewall_2002.json", eng, "FirewallUniversal"))
	easytesting2002.CheckDesiredValue("SettingType", "WindowsFirewallCurrentProfile [Original value: 2]")
	easytesting2002.CheckDesiredValue("SettingValueText", "Offentlig")
	easytesting2002.CheckDesiredValue("Origin", "Local")
	easytesting2002.CheckDesiredValue("ModifyingUser", "S-1-5-18")
	easytesting2002.CheckDesiredValue("ModifyingApplication", "")

	// 2003
	easytesting2003 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-Windows_Firewall_With_Advanced_Security_Firewall_2003.json", eng, "FirewallUniversal"))
	easytesting2003.CheckDesiredValue("Profiles/NewProfile", "Domain")
	easytesting2003.CheckDesiredValue("SettingType", "WindowsFirewallActivating [Original value: 1]")
	easytesting2003.CheckDesiredValue("SettingValueText", "No")
	easytesting2003.CheckDesiredValue("Origin", "Local")
	easytesting2003.CheckDesiredValue("ModifyingUser", "S-1-5-18")
	easytesting2003.CheckDesiredValue("ModifyingApplication", "C:\\Windows\\SysWOW64\\netsh.exe")

	// 2004
	easytesting2004 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-Windows_Firewall_With_Advanced_Security_Firewall_2004.json", eng, "FirewallUniversal"))
	easytesting2004.CheckDesiredValue("RuleId", "{2E3B4B93-F13C-4155-8D2A-4BDC974ADB6F}")
	easytesting2004.CheckDesiredValue("RuleName", "WLpacSenseNdr")
	easytesting2004.CheckDesiredValue("Origin", "Local")
	easytesting2004.CheckDesiredValue("ApplicationPath", "")
	easytesting2004.CheckDesiredValue("ServiceName", "")
	easytesting2004.CheckDesiredValue("Direction", "In")
	easytesting2004.CheckDesiredValue("Protocol", "All")
	easytesting2004.CheckDesiredValue("LocalPorts/SourcePort", "")
	easytesting2004.CheckDesiredValue("RemotePorts/DestPort", "")
	easytesting2004.CheckDesiredValue("Action", "Allow")
	easytesting2004.CheckDesiredValue("Profiles/NewProfile", "All profiles")
	easytesting2004.CheckDesiredValue("Local/SourceAddresses", "*")
	easytesting2004.CheckDesiredValue("Remote/DestAddresses", "*")
	easytesting2004.CheckDesiredValue("RemoteMachineAuthorizationList", "")
	easytesting2004.CheckDesiredValue("RemoteUserAuthorizationList", "")
	easytesting2004.CheckDesiredValue("EmbeddedContext", "WLpacSenseNdr")
	easytesting2004.CheckDesiredValue("Flags", "1")
	easytesting2004.CheckDesiredValue("Active", "Yes")
	easytesting2004.CheckDesiredValue("EdgeTraversal", "None")
	easytesting2004.CheckDesiredValue("SecurityOptions", "None")
	easytesting2004.CheckDesiredValue("ModifyingUser", "S-1-5-18")
	easytesting2004.CheckDesiredValue("ModifyingApplication", "C:\\Windows\\System32\\svchost.exe")
	easytesting2004.CheckDesiredValue("SchemaVersion", "541")

	// 2006
	easytesting2006 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-Windows_Firewall_With_Advanced_Security_Firewall_2006.json", eng, "FirewallUniversal"))
	easytesting2006.CheckDesiredValue("RuleId", "{061FB391-2B2D-4CE5-A3B8-F3111A557E15}")
	easytesting2006.CheckDesiredValue("RuleName", "WLpacSenseNdr")
	easytesting2006.CheckDesiredValue("ModifyingUser", "S-1-5-18")
	easytesting2006.CheckDesiredValue("ModifyingApplication", "C:\\Windows\\System32\\svchost.exe")

	// 2010
	easytesting2010 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-Windows_Firewall_With_Advanced_Security_Firewall_2010.json", eng, "FirewallUniversal"))
	easytesting2010.CheckDesiredValue("InterfaceGuid", "BA416E75-AA93-4D27-84FA-BEDA9D551AF8")
	easytesting2010.CheckDesiredValue("InterfaceName", "ethernet_30701")
	easytesting2010.CheckDesiredValue("OldProfile", "Domain")
	easytesting2010.CheckDesiredValue("Profiles/NewProfile", "None")

	// 2011
	easytesting2011 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-Windows_Firewall_With_Advanced_Security_Firewall_2011.json", eng, "FirewallUniversal"))
	easytesting2011.CheckDesiredValue("ReasonCode", "The application is a system service")
	easytesting2011.CheckDesiredValue("ApplicationPath", "C:\\program files (x86)\\rarwin\\rarwin.clienthost.exe")
	easytesting2011.CheckDesiredValue("IPVersion", "IPv4")
	easytesting2011.CheckDesiredValue("Protocol", "UDP")
	easytesting2011.CheckDesiredValue("LocalPorts/SourcePort", "57127")
	easytesting2011.CheckDesiredValue("ProcessId", "3270")
	easytesting2011.CheckDesiredValue("ModifyingUser", "S-1-5-18")

	// Security

	// 4945
	easytesting4945 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4945.json", eng, "FirewallUniversal"))
	easytesting4945.CheckDesiredValue("Profiles/NewProfile", "Offentlig")
	easytesting4945.CheckDesiredValue("RuleId", "WiFiDirect-Driver-In-TCP")
	easytesting4945.CheckDesiredValue("RuleName", "WFD-driver (TCP-In)")

	// 4946
	easytesting4946 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4946.json", eng, "FirewallUniversal"))
	easytesting4946.CheckDesiredValue("Profiles/NewProfile", "All")
	easytesting4946.CheckDesiredValue("RuleId", "{8B012210-4B06-4476-9D00-C141388C16AA}")
	easytesting4946.CheckDesiredValue("RuleName", "Acrobat Reader")

	// 4953
	easytesting4953 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4953.json", eng, "FirewallUniversal"))
	easytesting4953.CheckDesiredValue("Profiles/NewProfile", "All")
	easytesting4953.CheckDesiredValue("RuleId", "WMPNetworkSvc-2")
	easytesting4953.CheckDesiredValue("RuleName", "-")

	// 5031
	easytesting5031 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_5031.json", eng, "FirewallUniversal"))
	easytesting5031.CheckDesiredValue("Profiles/NewProfile", "Public")
	easytesting5031.CheckDesiredValue("ApplicationPath", "C:\\program files (x86)\\wan\\wdclient\\issuser.exe")

	// 5152
	easytesting5152 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_5152.json", eng, "FirewallUniversal"))
	easytesting5152.CheckDesiredValue("ProcessId", "4")
	easytesting5152.CheckDesiredValue("ApplicationPath", "System")
	easytesting5152.CheckDesiredValue("Direction", "Outbound")
	easytesting5152.CheckDesiredValue("Local/SourceAddresses", "11.11.11.11")
	easytesting5152.CheckDesiredValue("LocalPorts/SourcePort", "137")
	easytesting5152.CheckDesiredValue("Remote/DestAddresses", "12.12.12.12")
	easytesting5152.CheckDesiredValue("RemotePorts/DestPort", "137")
	easytesting5152.CheckDesiredValue("Protocol", "UDP")
	easytesting5152.CheckDesiredValue("FilterRTID", "111537")
	easytesting5152.CheckDesiredValue("LayerName", "Connect")
	easytesting5152.CheckDesiredValue("LayerRTID", "48")

	// 5154
	easytesting5154 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_5154.json", eng, "FirewallUniversal"))
	easytesting5154.CheckDesiredValue("ProcessId", "4")
	easytesting5154.CheckDesiredValue("ApplicationPath", "System")
	easytesting5154.CheckDesiredValue("Local/SourceAddresses", "11.11.11.11")
	easytesting5154.CheckDesiredValue("LocalPorts/SourcePort", "139")
	easytesting5154.CheckDesiredValue("Protocol", "TCP")
	easytesting5154.CheckDesiredValue("FilterRTID", "0")
	easytesting5154.CheckDesiredValue("LayerName", "Listen")
	easytesting5154.CheckDesiredValue("LayerRTID", "40")

	// 5156
	easytesting5156 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_5156.json", eng, "FirewallUniversal"))
	easytesting5156.CheckDesiredValue("ProcessID", "652")
	easytesting5156.CheckDesiredValue("ApplicationPath", "\\device\\harddiskvolume1\\windows\\system32\\lsass.exe")
	easytesting5156.CheckDesiredValue("Direction", "Inbound")
	easytesting5156.CheckDesiredValue("Local/SourceAddresses", "11.11.11.11")
	easytesting5156.CheckDesiredValue("LocalPorts/SourcePort", "636")
	easytesting5156.CheckDesiredValue("Remote/DestAddresses", "12.12.12.12")
	easytesting5156.CheckDesiredValue("RemotePorts/DestPort", "44980")
	easytesting5156.CheckDesiredValue("Protocol", "TCP")
	easytesting5156.CheckDesiredValue("FilterRTID", "275032")
	easytesting5156.CheckDesiredValue("LayerName", "Receive/Accept")
	easytesting5156.CheckDesiredValue("LayerRTID", "44")
	easytesting5156.CheckDesiredValue("RemoteUserID", "S-1-0-0")
	easytesting5156.CheckDesiredValue("RemoteMachineID", "S-1-0-0")

	// 5157
	easytesting5157 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_5157.json", eng, "FirewallUniversal"))
	easytesting5157.CheckDesiredValue("ProcessID", "4")
	easytesting5157.CheckDesiredValue("ApplicationPath", "System")
	easytesting5157.CheckDesiredValue("Direction", "Outbound")
	easytesting5157.CheckDesiredValue("Local/SourceAddresses", "11.11.11.11")
	easytesting5157.CheckDesiredValue("LocalPorts/SourcePort", "137")
	easytesting5157.CheckDesiredValue("Remote/DestAddresses", "12.12.12.12")
	easytesting5157.CheckDesiredValue("RemotePorts/DestPort", "137")
	easytesting5157.CheckDesiredValue("Protocol", "UDP")
	easytesting5157.CheckDesiredValue("FilterRTID", "111537")
	easytesting5157.CheckDesiredValue("LayerName", "Connect")
	easytesting5157.CheckDesiredValue("LayerRTID", "48")
	easytesting5157.CheckDesiredValue("RemoteUserID", "S-1-0-0")
	easytesting5157.CheckDesiredValue("RemoteMachineID", "S-1-0-0")

	// 6406
	easytesting6406 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_6406.json", eng, "FirewallUniversal"))
	easytesting6406.CheckDesiredValue("ProductName", "Symantec Endpoint Protection")
	easytesting6406.CheckDesiredValue("Categories", "BootTimeRuleCategory, StealthRuleCategory, FirewallRuleCategory")

}
