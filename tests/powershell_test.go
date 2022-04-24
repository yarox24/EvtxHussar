package tests

import (
	"testing"
)

func TestWindowsPowerShellEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 400
	easytesting400 := NewEasyTesting(t, UnmarshallAndParseEvent("Windows_PowerShell_400.json", eng, "PowerShellUniversal"))
	easytesting400.CheckDesiredValue("NewEngineState", "Available")
	easytesting400.CheckDesiredValue("PreviousEngineState", "None")
	easytesting400.CheckDesiredValue("SequenceNumber", "13")
	easytesting400.CheckDesiredValue("HostName", "ConsoleHost")
	easytesting400.CheckDesiredValue("HostVersion", "5.1.14409.1018")
	easytesting400.CheckDesiredValue("HostId", "bf81775a-bf81-473d-8f8e-bf81bf81bec2")
	easytesting400.CheckDesiredValue("HostApplication", "c:\\windows\\sysnative\\WindowsPowerShell\\v1.0\\powershell.exe -ExecutionPolicy Bypass -enc RwBlAHQALQBEAGEAdABlAA==")
	easytesting400.CheckDesiredValue("HostApplication (Base64 decoded)", "Get-Date")
	easytesting400.CheckDesiredValue("EngineVersion", "5.1.14409.1018")
	easytesting400.CheckDesiredValue("RunspaceId", "12e27fd9-281d-43bb-8135-geb0f27a38b5")
	easytesting400.CheckDesiredValue("PipelineId", "")

	//600
	easytesting600 := NewEasyTesting(t, UnmarshallAndParseEvent("Windows_PowerShell_600.json", eng, "PowerShellUniversal"))
	easytesting600.CheckDesiredValue("ProviderName", "Registry")
	easytesting600.CheckDesiredValue("NewProviderState", "Started")
	easytesting600.CheckDesiredValue("SequenceNumber", "1")
	easytesting600.CheckDesiredValue("HostName", "Default Host")
	easytesting600.CheckDesiredValue("HostVersion", "5.1.19041.610")
	easytesting600.CheckDesiredValue("HostId", "n32376fa-3500-4d44-8e09-792119fa644f")
	easytesting600.CheckDesiredValue("HostApplication", "C:\\WINDOWS\\System32\\remote.exe Disable")
	easytesting600.CheckDesiredValue("EngineVersion", "")
	easytesting600.CheckDesiredValue("RunspaceId", "")

	//800
	easytesting800 := NewEasyTesting(t, UnmarshallAndParseEvent("Windows_PowerShell_800.json", eng, "PowerShellUniversal"))
	easytesting800.CheckDesiredValue("DetailSequence", "1")
	easytesting800.CheckDesiredValue("DetailTotal", "1")
	easytesting800.CheckDesiredValue("SequenceNumber", "35")
	easytesting800.CheckDesiredValue("UserId", "NT AUTHORITY\\LOCAL SERVICE")
	easytesting800.CheckDesiredValue("HostName", "ConsoleHost")
	easytesting800.CheckDesiredValue("HostVersion", "5.1.19041.1237")
	easytesting800.CheckDesiredValue("HostId", "27a58bfa-e839-5bcd-d91f-5689d47d32c3")
	easytesting800.CheckDesiredValue("HostApplication", "C:\\Windows\\system32\\WindowsPowerShell\\v1.0\\powershell.exe -ExecutionPolicy AllSigned -NoProfile -NonInteractive -Command & {$OutputEncoding = [Console]::OutputEncoding =[System.Text.Encoding]::UTF8;$scriptFileStream = [System.IO.File]::Open('C:\\ProgramData\\Microsoft\\Windows Defender Advanced Threat Protection\\DataCollection\\xxx.ps1', [System.IO.FileMode]::Open, [System.IO.FileAccess]::Read, [System.IO.FileAccess]::Read);$calculatedHash = Get-FileHash 'C:\\ProgramData\\Microsoft\\Windows Defender Advanced Threat Protection\\DataCollection\\xxx.ps1' -Algorithm SHA256;if (!($calculatedHash.Hash -eq 'b545de56a106513269ea74f81e2268beae0fef0cff87fee4a98a66da7a643cc8')) { exit 324;}; . 'C:\\ProgramData\\Microsoft\\Windows Defender Advanced Threat Protection\\DataCollection\\xxx.ps1' }")
	easytesting800.CheckDesiredValue("EngineVersion", "5.1.19041.1237")
	easytesting800.CheckDesiredValue("RunspaceId", "hb18d2se-363e-54b4-13fd-abb094570d0e")
	easytesting800.CheckDesiredValue("PipelineId", "1")
	easytesting800.CheckDesiredValue("ScriptName", "C:\\ProgramData\\Microsoft\\Windows Defender Advanced Threat Protection\\DataCollection\\xxx.ps1")
	easytesting800.CheckDesiredValue("CommandLine", "Add-Type -TypeDefinition $Source -Language CSharp -IgnoreWarnings")
	easytesting800.CheckDesiredLength("CommandInvocation/ParameterBinding", 5829)
}

func TestMicrosoftWindowsPowerShellOperationalEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 4100
	easytesting4100 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-PowerShell_Operational_4100.json", eng, "PowerShellUniversal"))
	easytesting4100.CheckDesiredValue("Severity", "Warning")
	easytesting4100.CheckDesiredValue("HostName", "OpsMgr PowerShell Host")
	easytesting4100.CheckDesiredValue("HostVersion", "7.0.5000.0")
	easytesting4100.CheckDesiredValue("HostID", "4f80b870-db49-4aff-bf8e-ec910d0b6978")
	easytesting4100.CheckDesiredValue("HostApplication", "C:\\Program Files\\Microsoft Monitoring Agent\\Agent\\MonitoringHost.exe -Embedding")
	easytesting4100.CheckDesiredValue("EngineVersion", "4.0")
	easytesting4100.CheckDesiredValue("RunspaceID", "a132d323-abff-1ff6-b46f-2ed01b71fc4d")
	easytesting4100.CheckDesiredValue("PipelineID", "2")
	easytesting4100.CheckDesiredValue("CommandName", "")
	easytesting4100.CheckDesiredValue("SequenceNumber", "111696")
	easytesting4100.CheckDesiredValue("User", "NTNT\\SYSTEM")
	easytesting4100.CheckDesiredValue("Shell ID", "Microsoft.PowerShell")

	// 8193
	easytesting8193 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-PowerShell_Operational_8193.json", eng, "PowerShellUniversal"))
	easytesting8193.CheckDesiredValue("param1", "7f28a475-ed3e-47e1-940d-6ea8610c36bb")

	// 8194
	easytesting8194 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-PowerShell_Operational_8194.json", eng, "PowerShellUniversal"))
	easytesting8194.CheckDesiredValue("InstanceId", "157f7844-6bd7-4dd2-ac6d-fb91496d4688")
	easytesting8194.CheckDesiredValue("MaxRunspaces", "2")
	easytesting8194.CheckDesiredValue("MinRunspaces", "10")

	// 24577
	easytesting24577 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-PowerShell_Operational_24577.json", eng, "PowerShellUniversal"))
	easytesting24577.CheckDesiredValue("FileName", "Noname.ps1")

	// 53504
	easytesting53504 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-PowerShell_Operational_53504.json", eng, "PowerShellUniversal"))
	easytesting53504.CheckDesiredValue("param1", "10368")
	easytesting53504.CheckDesiredValue("param2", "DefaultAppDomain")

}
