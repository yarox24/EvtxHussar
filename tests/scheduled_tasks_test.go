package tests

import (
	"testing"
)

func TestScheduledTasksSecurityEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 4698 Version 1
	easytesting4698 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4698.json", eng, "ScheduledTasks_CreationModification"))
	easytesting4698.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4698.CheckDesiredValue("SubjectUserName", "DESKTOP-EvtxHussar$")
	easytesting4698.CheckDesiredValue("SubjectDomainName", "HUSS")
	easytesting4698.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4698.CheckDesiredValue("TaskName", "\\Microsoft\\Windows\\UpdateMe\\Musake")
	easytesting4698.CheckDesiredValue("ClientProcessStartKey", "10123099161583875")
	easytesting4698.CheckDesiredValue("ClientProcessId", "10012")
	easytesting4698.CheckDesiredValue("ParentProcessId", "1016")
	easytesting4698.CheckDesiredValue("RpcCallClientLocality", "0")
	easytesting4698.CheckDesiredValue("FQDN", "DESKTOP-EvtxHussar")

	// XML Parser
	easytesting4698.CheckDesiredValue("Actions summary (TaskContent XML)", "Exec: 1")
	easytesting4698.CheckDesiredValue("XML Task Version (TaskContent XML)", "1.2")
	easytesting4698.CheckDesiredValue("RegistrationInfo URI (TaskContent XML)", "\\Microsoft\\Windows\\UpdateMe\\Musake")
	easytesting4698.CheckDesiredValue("Triggers summary (TaskContent XML)", "LogonTrigger: 1")
	easytesting4698.CheckDesiredValue("Task Hidden (TaskContent XML)", "false")
	easytesting4698.CheckDesiredValue("Exec Command with Arguments (TaskContent XML)", "%systemroot%\\system32\\Musake.exe UpdateMe")
	easytesting4698.CheckDesiredValue("Principal UserId (TaskContent XML)", "S-1-5-18")
	easytesting4698.CheckDesiredValue("Principal RunLevel (TaskContent XML)", "LeastPrivilege")
	easytesting4698.CheckDesiredValue("Task enabled (TaskContent XML)", "true")
}

func TestScheduledTasksMicrosoftWindowsTaskSchedulerOperationalEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 100
	easytesting100 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-TaskScheduler_Operational_100.json", eng, "ScheduledTasks_Execution"))
	easytesting100.CheckDesiredValue("TaskName", "\\Microsoft\\Windows\\DirectX\\AdapterCache")
	easytesting100.CheckDesiredValue("UserContext", "NT AUTHORITY\\SYSTEM")
	easytesting100.CheckDesiredValue("InstanceId", "4B0F2D5B-2182-291A-BF46-1017A6B18D05")

	// 129
	easytesting129 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-TaskScheduler_Operational_129.json", eng, "ScheduledTasks_Execution"))
	easytesting129.CheckDesiredValue("TaskName", "\\Microsoft\\Windows\\DirectX\\AdapterCache")
	easytesting129.CheckDesiredValue("Path", "%windir%\\system32\\adaptercache.exe")
	easytesting129.CheckDesiredValue("ProcessID", "0x1aac")
	easytesting129.CheckDesiredValue("Priority", "16384")

	// 201
	easytesting201 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-TaskScheduler_Operational_201.json", eng, "ScheduledTasks_Execution"))
	easytesting201.CheckDesiredValue("TaskName", "\\Microsoft\\Windows\\Subscription\\Renew")
	easytesting201.CheckDesiredValue("TaskInstanceId", "0B7D84B2-B42A-18E7-8127-B56B858191B1")
	easytesting201.CheckDesiredValue("ActionName", "%SystemRoot%\\system32\\Renew.exe")
	easytesting201.CheckDesiredValue("ResultCode", "0")
	easytesting201.CheckDesiredValue("EnginePID", "6848")

}
