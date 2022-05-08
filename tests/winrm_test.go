package tests

import (
	"github.com/yarox24/EvtxHussar/eventmap"
	"testing"
)

func TestWinRMEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 6
	easytesting6 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-WinRM_Operational_6.json", eng, "WinRMUniversal"))
	easytesting6.CheckDesiredValue("connection", "mydemohost/wsman?PSVersion=5.1.20348.320")
	easytesting6.CheckDesiredValue("connection (hostname)", "mydemohost")
	easytesting6.CheckDesiredValue("connection (powershell version)", "5.1.20348.320")

	// 7
	easytesting7 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-WinRM_Operational_7.json", eng, "WinRMUniversal"))
	easytesting7.CheckDesiredValue("errorCode", "2150859169")

	// 11
	easytesting11 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-WinRM_Operational_11.json", eng, "WinRMUniversal"))
	easytesting11.CheckDesiredValue("resourceUri", "http://schemas.microsoft.com/powershell/Microsoft.PowerShell")
	easytesting11.CheckDesiredValue("shellId", "26DAFDB8-CDB6-6417-9B88-83190B849771")

	// 41
	easytesting41 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-WinRM_Operational_41.json", eng, "WinRMUniversal"))
	easytesting41.CheckDesiredValue("applicationID", "ServerManager.exe")

	// 44
	easytesting44 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-WinRM_Operational_44.json", eng, "WinRMUniversal"))
	easytesting44.CheckDesiredValue("destination", "\u003clocal\u003e")

	// 47
	easytesting47 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-WinRM_Operational_47.json", eng, "WinRMUniversal"))
	easytesting47.CheckDesiredValue("operationType", "EnumerateInstances")
	easytesting47.CheckDesiredValue("namespaceName", "root\\microsoft\\windows\\storage")
	easytesting47.CheckDesiredValue("className", "*")

	// 91
	easytesting91 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-WinRM_Operational_91.json", eng, "WinRMUniversal"))
	easytesting91.CheckDesiredValue("ErrorCode", "15005")

	// 91 - Special UTF-16 Field EventPayload
	eventpayload := []uint8{104, 0, 116, 0, 116, 0, 112, 0, 58, 0, 47, 0, 47, 0, 115, 0, 99, 0, 104, 0, 101, 0, 109, 0, 97, 0, 115, 0, 46, 0, 109, 0, 105, 0, 99, 0, 114, 0, 111, 0, 115, 0, 111, 0, 102, 0, 116, 0, 46, 0, 99, 0, 111, 0, 109, 0, 47, 0, 112, 0, 111, 0, 119, 0, 101, 0, 114, 0, 115, 0, 104, 0, 101, 0, 108, 0, 108, 0, 47, 0, 77, 0, 105, 0, 99, 0, 114, 0, 111, 0, 115, 0, 111, 0, 102, 0, 116, 0, 46, 0, 80, 0, 111, 0, 119, 0, 101, 0, 114, 0, 83, 0, 104, 0, 101, 0, 108, 0, 108, 0, 0, 0}
	decoded_string := eventmap.ConvertAllTypesToString(eventpayload, "uint8slice_utf-16")
	if decoded_string != "http://schemas.microsoft.com/powershell/Microsoft.PowerShell" {
		t.Errorf("Error in decoding uint8slice_utf-16")
	}

	// 161
	easytesting161 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-WinRM_Operational_161.json", eng, "WinRMUniversal"))
	easytesting161.CheckDesiredValue("authFailureMessage", "WinRM cannot process the request. The following error with errorcode 0x80090350 occurred while using Negotiate authentication: An unknown security error occurred.  \r\n Possible causes are:\r\n  -The user name or password specified are invalid.\r\n  -Kerberos is used when no authentication method and no user name are specified.\r\n  -Kerberos accepts domain user names, but not local user names.\r\n  -The Service Principal Name (SPN) for the remote computer name and port does not exist.\r\n  -The client and remote computers are in different domains and there is no trust between the two domains.\r\n After checking for the above issues, try the following:\r\n  -Check the Event Viewer for events related to authentication.\r\n  -Change the authentication method; add the destination computer to the WinRM TrustedHosts configuration setting or use HTTPS transport.\r\n Note that computers in the TrustedHosts list might not be authenticated.\r\n   -For more information about WinRM configuration, run the following command: winrm help config.")
}
