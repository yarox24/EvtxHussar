package tests

import (
	"testing"
)

func TestAuditSMBEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine()

	// Microsoft-Windows-SmbClient_Audit_32002
	easytesting32002 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SmbClient_Audit_32002.json", eng, "SMB_ClientDestinations"))
	easytesting32002.CheckDesiredValue("Reason", "Smb2DiagReasonNotSpecified [Original value: 0x0]")
	easytesting32002.CheckDesiredValue("Dialect", "5")
	easytesting32002.CheckDesiredValue("SecurityMode", "12803")
	easytesting32002.CheckDesiredValue("ServerName", "\\172.16.40.22")

	// Microsoft-Windows-SmbClient_Connectivity_30803
	easytesting30803 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SmbClient_Connectivity_30803.json", eng, "SMB_ClientDestinations"))
	easytesting30803.CheckDesiredValue("Reason", "Smb2DiagReasonNetworkConnect [Original value: 0x4]")
	easytesting30803.CheckDesiredValue("Status", "The I/O request was canceled. [Original value: 0xC0000120]")
	easytesting30803.CheckDesiredValue("ServerName", "evtxhussar.net")
	easytesting30803.CheckDesiredValue("InstanceName", "\\Device\\LanmanRedirector")
	easytesting30803.CheckDesiredValue("ConnectionType", "Wsk [Original value: 0x1]")

	// Microsoft-Windows-SmbClient_Connectivity_30806
	easytesting30806 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SmbClient_Connectivity_30806.json", eng, "SMB_ClientDestinations"))
	easytesting30806.CheckDesiredValue("Status", "The operation completed successfully. [Original value: 0x0]")
	easytesting30806.CheckDesiredValue("SessionId", "0xa80094000055")
	easytesting30806.CheckDesiredValue("TreeId", "0x0")
	easytesting30806.CheckDesiredValue("ServerName", "\\Hussar.NET")

	// Microsoft-Windows-SmbClient_Security_31000
	easytesting31000 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SmbClient_Security_31000.json", eng, "SMB_ClientDestinations"))
	easytesting31000.CheckDesiredValue("Reason", "Smb2DiagReasonAcquireCredHandle [Original value: 0x9]")
	easytesting31000.CheckDesiredValue("Status", "A specified authentication package is unknown. [Original value: 0xC00000FE]")
	easytesting31000.CheckDesiredValue("SecurityStatus", "The requested security package was not found. [Original value: 0x80090305]")
	easytesting31000.CheckDesiredValue("LogonId", "0x3e7")
	easytesting31000.CheckDesiredValue("ServerName", "\\HUSS")
	easytesting31000.CheckDesiredValue("PrincipalName", "")
	easytesting31000.CheckDesiredValue("UserName", "")

	// Microsoft-Windows-SmbClient_Security_31013
	easytesting31013 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SmbClient_Security_31013.json", eng, "SMB_ClientDestinations"))
	easytesting31013.CheckDesiredValue("Smb2Command", "11")
	easytesting31013.CheckDesiredValue("MessageId", "5")
	easytesting31013.CheckDesiredValue("SessionId", "0xd884b400004d")
	easytesting31013.CheckDesiredValue("TreeId", "0x0")
	easytesting31013.CheckDesiredValue("ServerName", "\\huss.net")
	easytesting31013.CheckDesiredValue("Status", "The cryptographic signature is invalid. [Original value: 0xC000A000]")

	// Microsoft-Windows-SmbClient_Security_31019
	easytesting31019 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SmbClient_Security_31019.json", eng, "SMB_ClientDestinations"))
	easytesting31019.CheckDesiredValue("Reason", "Smb2DiagReasonMADowngrade [Original value: 0xC9]")
	easytesting31019.CheckDesiredValue("Status", "The account is not authorized to log on from this station. [Original value: 0xC0000248]")
	easytesting31019.CheckDesiredValue("SecurityStatus", "The account is not authorized to log on from this station. [Original value: 0xC0000248]")
	easytesting31019.CheckDesiredValue("LogonId", "0xd20d5")
	easytesting31019.CheckDesiredValue("ServerName", "\\huss")
	easytesting31019.CheckDesiredValue("UserName", "")

	// Microsoft-Windows-SMBServer_Audit_3000
	easytesting3000 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SMBServer_Audit_3000.json", eng, "SMB_ServerAccessAudit"))
	easytesting3000.CheckDesiredValue("ClientName", "1.1.1.1")

	// Microsoft-Windows-SMBServer_Operational_1016
	easytesting1016 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SMBServer_Operational_1016.json", eng, "SMB_ServerAccessAudit"))
	easytesting1016.CheckDesiredValue("Status", "The object name is not found. [Original value: 0xC0000034]")
	easytesting1016.CheckDesiredValue("TranslatedStatus", "The object name is not found. [Original value: 0xC0000034]")
	easytesting1016.CheckDesiredValue("RKFStatus", "0x0")
	easytesting1016.CheckDesiredValue("TranslatedRKFStatus", "0x0")
	easytesting1016.CheckDesiredValue("ConnectionGUID", "FF2D70A7-31C7-0302-6A0F-43EFC791D301")
	easytesting1016.CheckDesiredValue("ClientName", "\\\\1.1.1.1")
	easytesting1016.CheckDesiredValue("ClientAddress", "")
	easytesting1016.CheckDesiredValue("ShareName", "C$")
	easytesting1016.CheckDesiredValue("SubjectUserName", "Hu\\SS")
	easytesting1016.CheckDesiredValue("SubjectLogonId", "0xa403c000006d")
	easytesting1016.CheckDesiredValue("FileName", "Program Files (x86)\\install\\install.exe")
	easytesting1016.CheckDesiredValue("DurableHandle", "false")
	easytesting1016.CheckDesiredValue("ResilientHandle", "false")
	easytesting1016.CheckDesiredValue("PersistentHandle", "false")
	easytesting1016.CheckDesiredValue("ResumeKey", "2F23241B-1FD5-11A8-81BD-821844E9A9F9")
	easytesting1016.CheckDesiredValue("Reason", "1")

	// Microsoft-Windows-SMBServer_Security_551
	easytesting551 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SMBServer_Security_551.json", eng, "SMB_ServerAccessAudit"))
	easytesting551.CheckDesiredValue("SessionGUID", "84224323-FC1E-0003-BD6D-68841EFCD601")
	easytesting551.CheckDesiredValue("ConnectionGUID", "84334323-FC1E-0001-858B-68841EFCD601")
	easytesting551.CheckDesiredValue("Status", "The attempted logon is invalid. This is either due to a bad username or authentication information. [Original value: 0xC000006D]")
	easytesting551.CheckDesiredValue("TranslatedStatus", "The attempted logon is invalid. This is either due to a bad username or authentication information. [Original value: 0xC000006D]")
	easytesting551.CheckDesiredValue("SubjectLogonId", "0x80000000001")
	easytesting551.CheckDesiredValue("SubjectUserName", "")
	easytesting551.CheckDesiredValue("ClientName", "\\\\1.1.1.1")
	easytesting551.CheckDesiredValue("SPN", "session setup failed before the SPN could be queried")
	easytesting551.CheckDesiredValue("SPNValidationPolicy", "0")

	// Microsoft-Windows-SMBServer_Security_1006
	easytesting1006 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SMBServer_Security_1006.json", eng, "SMB_ServerAccessAudit"))
	easytesting1006.CheckDesiredValue("ShareName", "\\\\*\\ss$")
	easytesting1006.CheckDesiredValue("ShareLocalPath", "\\??\\D:\\ss")
	easytesting1006.CheckDesiredValue("SubjectUserName", "HU\\SS")
	easytesting1006.CheckDesiredValue("ClientName", "\\\\1.1.1.1")
	easytesting1006.CheckDesiredValue("MappedAccess", "0x100081")
	easytesting1006.CheckDesiredValue("GrantedAccess", "0x0")
	easytesting1006.CheckDesiredValue("ShareSecurityDescriptor", "")
	easytesting1006.CheckDesiredValue("Status", "{Access Denied} A process has requested access to an object but has not been granted those access rights. [Original value: 0xC0000022]")
	easytesting1006.CheckDesiredValue("TranslatedStatus", "{Access Denied} A process has requested access to an object but has not been granted those access rights. [Original value: 0xC0000022]")
	easytesting1006.CheckDesiredValue("SubjectLogonId", "0x14c0004000011")

	// Microsoft-Windows-SMBServer_Security_1007
	easytesting1007 := NewEasyTesting(t, UnmarshallAndParseEvent("Microsoft-Windows-SMBServer_Security_1007.json", eng, "SMB_ServerAccessAudit"))
	easytesting1007.CheckDesiredValue("ShareName", "\\\\*\\SYSVOL")
	easytesting1007.CheckDesiredValue("ShareLocalPath", "\\??\\H:\\Windows\\SYSVOL_DFSR\\sysvol")
	easytesting1007.CheckDesiredValue("ClientName", "\\\\1.1.1.1")

	// Security_5140
	easytesting5140 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_5140.json", eng, "SMB_ServerAccessAudit"))
	easytesting5140.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting5140.CheckDesiredValue("SubjectUserName", "user02")
	easytesting5140.CheckDesiredValue("SubjectDomainName", "EXAMPLE")
	easytesting5140.CheckDesiredValue("SubjectLogonId", "0x15e1a7")
	easytesting5140.CheckDesiredValue("ObjectType", "File")
	easytesting5140.CheckDesiredValue("ClientAddress", "1.1.1.1")
	easytesting5140.CheckDesiredValue("IpPort", "49222")
	easytesting5140.CheckDesiredValue("ShareName", "\\\\*\\IPC$")
	easytesting5140.CheckDesiredValue("ShareLocalPath", "")
	easytesting5140.CheckDesiredValue("AccessMask", "0x1")
	easytesting5140.CheckDesiredValue("AccessList", "ReadData (or ListDirectory)\r\n\t\t\t\t")

	// Security_5143
	easytesting5143 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_5143.json", eng, "SMB_ServerModifications"))
	easytesting5143.CheckDesiredValue("SubjectUserSid", "S-1-5-21-3212943211-794299840-588279583-500")
	easytesting5143.CheckDesiredValue("SubjectUserName", "Administrator")
	easytesting5143.CheckDesiredValue("SubjectDomainName", "2016SRV")
	easytesting5143.CheckDesiredValue("SubjectLogonId", "0x2f90a")
	easytesting5143.CheckDesiredValue("ObjectType", "Directory")
	easytesting5143.CheckDesiredValue("ShareName", "\\\\*\\Documents")
	easytesting5143.CheckDesiredValue("ShareLocalPath", "C:\\Documents")
	easytesting5143.CheckDesiredValue("OldRemark", "My Share!")
	easytesting5143.CheckDesiredValue("NewRemark", "My Modified Share!")
	easytesting5143.CheckDesiredValue("OldMaxUsers", "0xa")
	easytesting5143.CheckDesiredValue("NewMaxUsers", "0xf")
	easytesting5143.CheckDesiredValue("OldShareFlags", "0x30")
	easytesting5143.CheckDesiredValue("NewShareFlags", "0x30")
	easytesting5143.CheckDesiredValue("OldSD", "O:BAG:S-1-5-21-3212943211-794299840-588279583-513D:(A;;0x1301bf;;;WD)")
	easytesting5143.CheckDesiredValue("NewSD", "O:BAG:S-1-5-21-3212943211-794299840-588279583-513D:(A;;0x1200a9;;;WD)")

	// Security_5145
	easytesting5145 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_5145.json", eng, "SMB_ServerAccessAudit"))
	easytesting5145.CheckDesiredValue("SubjectUserSid", "S-1-5-21-308926384-506822093-3341789130-107103")
	easytesting5145.CheckDesiredValue("SubjectUserName", "backdoor")
	easytesting5145.CheckDesiredValue("SubjectDomainName", "3B")
	easytesting5145.CheckDesiredValue("SubjectLogonId", "0x4e5197")
	easytesting5145.CheckDesiredValue("ObjectType", "File")
	easytesting5145.CheckDesiredValue("ClientAddress", "172.16.66.1")
	easytesting5145.CheckDesiredValue("IpPort", "52415")
	easytesting5145.CheckDesiredValue("ShareName", "\\\\*\\IPC$")
	easytesting5145.CheckDesiredValue("ShareLocalPath", "")
	easytesting5145.CheckDesiredValue("RelativeTargetName", "protected_storage")
	easytesting5145.CheckDesiredValue("AccessMask", "0x12019f")
	//easytesting5145.CheckDesiredValue("AccessList", "")
	easytesting5145.CheckDesiredValue("AccessReason", "-")

}
