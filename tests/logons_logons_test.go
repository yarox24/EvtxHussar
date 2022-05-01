package tests

import (
	"testing"
)

// Regex search: "(\w+)": "?(.*?)"?,
// Regex replace: easytesting4624.CheckDesiredValue\("\1", "\2"\)

func TestLogonEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 4624 Version 2
	easytesting4624 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4624.json", eng, "LogonsUniversal"))
	easytesting4624.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4624.CheckDesiredValue("SubjectUserName", "DESKTOP-EvtxHussar$")
	easytesting4624.CheckDesiredValue("SubjectDomainName", "WORKGROUP")
	easytesting4624.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4624.CheckDesiredValue("TargetUserSid", "S-1-5-18")
	easytesting4624.CheckDesiredValue("TargetUserName", "SYSTEM")
	easytesting4624.CheckDesiredValue("TargetDomainName", "NT AUTHORITY")
	easytesting4624.CheckDesiredValue("TargetLogonId", "0x3e7")
	easytesting4624.CheckDesiredValue("LogonType", "Service [Original value: 5]")
	easytesting4624.CheckDesiredValue("LogonType (Use cases)", "Windows services")
	easytesting4624.CheckDesiredValue("LogonProcessName", "Advapi  ")
	easytesting4624.CheckDesiredValue("AuthenticationPackageName", "Negotiate")
	easytesting4624.CheckDesiredValue("WorkstationName", "-")
	easytesting4624.CheckDesiredValue("LogonGuid", "00000000-0000-0000-0000-000000000000")
	easytesting4624.CheckDesiredValue("TransmittedServices/TransitedServices", "-")
	easytesting4624.CheckDesiredValue("LmPackageName", "-")
	easytesting4624.CheckDesiredValue("KeyLength", "0")
	easytesting4624.CheckDesiredValue("ProcessId", "0x25c")
	easytesting4624.CheckDesiredValue("ProcessName", "C:\\Windows\\System32\\services.exe")
	easytesting4624.CheckDesiredValue("IpAddress", "-")
	easytesting4624.CheckDesiredValue("IpPort", "-")
	easytesting4624.CheckDesiredValue("ImpersonationLevel", "Impersonation")
	easytesting4624.CheckDesiredValue("RestrictedAdminMode", "-")
	easytesting4624.CheckDesiredValue("TargetOutboundUserName", "-")
	easytesting4624.CheckDesiredValue("TargetOutboundDomainName", "-")
	easytesting4624.CheckDesiredValue("VirtualAccount", "No")
	easytesting4624.CheckDesiredValue("TargetLinkedLogonId", "0x0")
	easytesting4624.CheckDesiredValue("ElevatedToken", "Yes")

	// 4625
	easytesting4625 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4625.json", eng, "LogonsUniversal"))
	easytesting4625.CheckDesiredValue("SubjectUserSid", "S-1-0-0")
	easytesting4625.CheckDesiredValue("SubjectUserName", "-")
	easytesting4625.CheckDesiredValue("SubjectDomainName", "-")
	easytesting4625.CheckDesiredValue("SubjectLogonId", "0x0")
	easytesting4625.CheckDesiredValue("TargetUserSid", "S-1-0-0")
	easytesting4625.CheckDesiredValue("TargetUserName", "hussar")
	easytesting4625.CheckDesiredValue("TargetDomainName", "HUSS")
	easytesting4625.CheckDesiredValue("Status", "Logon failure: unknown user name or bad password. [Original value: 0xC000006D]")
	easytesting4625.CheckDesiredValue("FailureReason/FailureCode", "Unknown user name or bad password.")
	easytesting4625.CheckDesiredValue("SubStatus", "The specified network password is not correct. [Original value: 0xC000006A]")
	easytesting4625.CheckDesiredValue("LogonType", "Network [Original value: 3]")
	easytesting4625.CheckDesiredValue("LogonProcessName", "NtLmSsp ")
	easytesting4625.CheckDesiredValue("AuthenticationPackageName", "NTLM")
	easytesting4625.CheckDesiredValue("WorkstationName", "HUSStext")
	easytesting4625.CheckDesiredValue("TransmittedServices/TransitedServices", "-")
	easytesting4625.CheckDesiredValue("LmPackageName", "-")
	easytesting4625.CheckDesiredValue("KeyLength", "0")
	easytesting4625.CheckDesiredValue("ProcessId", "0x0")
	easytesting4625.CheckDesiredValue("ProcessName", "-")
	easytesting4625.CheckDesiredValue("IpAddress", "192.235.140.1")
	easytesting4625.CheckDesiredValue("IpPort", "0")

	// 4648
	easytesting4648 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4648.json", eng, "LogonsUniversal"))
	easytesting4648.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4648.CheckDesiredValue("SubjectUserName", "EvtxHussar$")
	easytesting4648.CheckDesiredValue("SubjectDomainName", "WORKGROUP")
	easytesting4648.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4648.CheckDesiredValue("LogonGuid", "00000000-0000-0000-0000-000000000000")
	easytesting4648.CheckDesiredValue("TargetUserName", "UMFD-5")
	easytesting4648.CheckDesiredValue("TargetDomainName", "Font Driver Host")
	easytesting4648.CheckDesiredValue("TargetLogonGuid", "00000000-0000-0000-0000-000000000000")
	easytesting4648.CheckDesiredValue("TargetServerName", "localhost")
	easytesting4648.CheckDesiredValue("TargetInfo", "localhost")
	easytesting4648.CheckDesiredValue("ProcessId", "0x7d4")
	easytesting4648.CheckDesiredValue("ProcessName", "C:\\Windows\\System32\\winlogon.exe")
	easytesting4648.CheckDesiredValue("IpAddress", "-")
	easytesting4648.CheckDesiredValue("IpPort", "-")

	// 4768
	easytesting4768 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4768.json", eng, "LogonsUniversal"))
	easytesting4768.CheckDesiredValue("TargetUserName", "EvtxHussar$")
	easytesting4768.CheckDesiredValue("TargetDomainName", "HUSS")
	easytesting4768.CheckDesiredValue("TargetUserSid", "S-1-5-18")
	easytesting4768.CheckDesiredValue("ServiceName", "krbtgt")
	easytesting4768.CheckDesiredValue("ServiceSid", "S-1-5-18")
	easytesting4768.CheckDesiredValue("TicketOptions", "Name-canonicalize | Forwardable [Original value: 0x40010000]")
	easytesting4768.CheckDesiredValue("Status", "Status OK. [Original value: 0x0]")
	easytesting4768.CheckDesiredValue("TicketEncryptionType", "AES256-CTS-HMAC-SHA1-96")
	easytesting4768.CheckDesiredValue("PreAuthType", "PA-ENC-TIMESTAMP (This is a normal type for standard password authentication.) [Original value: 2]")
	easytesting4768.CheckDesiredValue("IpAddress", "::ffff:192.168.0.1")
	easytesting4768.CheckDesiredValue("IpPort", "41384")
	easytesting4768.CheckDesiredValue("CertIssuerName", "")
	easytesting4768.CheckDesiredValue("CertSerialNumber", "")
	easytesting4768.CheckDesiredValue("CertThumbprint", "")

	// 4769
	easytesting4769 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4769.json", eng, "LogonsUniversal"))
	easytesting4769.CheckDesiredValue("TargetUserName", "huus@huss")
	easytesting4769.CheckDesiredValue("TargetDomainName", "HUSS")
	easytesting4769.CheckDesiredValue("ServiceName", "WINGED$")
	easytesting4769.CheckDesiredValue("ServiceSid", "S-1-5-19")
	easytesting4769.CheckDesiredValue("TicketOptions", "Name-canonicalize [Original value: 0x10000]")
	easytesting4769.CheckDesiredValue("TicketEncryptionType", "AES256-CTS-HMAC-SHA1-96")
	easytesting4769.CheckDesiredValue("IpAddress", "::ffff:192.1.1.1")
	easytesting4769.CheckDesiredValue("IpPort", "41386")
	easytesting4769.CheckDesiredValue("Status", "Status OK. [Original value: 0x0]")
	easytesting4769.CheckDesiredValue("LogonGuid", "B2ABA7AB-4FDE-A309-598D-A5DCC7553372")
	easytesting4769.CheckDesiredValue("TransmittedServices/TransitedServices", "-")

	// 4776
	easytesting4776 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4776.json", eng, "LogonsUniversal"))
	easytesting4776.CheckDesiredValue("AuthenticationPackageName", "MICROSOFT_AUTHENTICATION_PACKAGE_V1_0")
	easytesting4776.CheckDesiredValue("TargetUserName", "huss")
	easytesting4776.CheckDesiredValue("WorkstationName", "DESKTOP-EvtxHussar")
	easytesting4776.CheckDesiredValue("Status", "Status OK. [Original value: 0x0]")

	// 4778
	easytesting4778 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4778.json", eng, "LogonsUniversal"))
	easytesting4778.CheckDesiredValue("SubjectUserName", "acco")
	easytesting4778.CheckDesiredValue("SubjectDomainName", "HUSS")
	easytesting4778.CheckDesiredValue("SubjectLogonId", "0x5fcc54cd")
	easytesting4778.CheckDesiredValue("SessionName", "RDP-Tcp#26")
	easytesting4778.CheckDesiredValue("WorkstationName", "SourceHuss")
	easytesting4778.CheckDesiredValue("IpAddress", "192.168.0.1")

	// 4964
	easytesting4964 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4964.json", eng, "LogonsUniversal"))
	easytesting4964.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4964.CheckDesiredValue("SubjectUserName", "WIN10-1703$")
	easytesting4964.CheckDesiredValue("SubjectDomainName", "HQCORP")
	easytesting4964.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4964.CheckDesiredValue("LogonGuid", "00000000-0000-0000-0000-000000000000")
	easytesting4964.CheckDesiredValue("TargetUserSid", "S-1-5-21-1913345275-1711810662-261465553-500")
	easytesting4964.CheckDesiredValue("TargetUserName", "Administrator")
	easytesting4964.CheckDesiredValue("TargetDomainName", "HQCORP")
	easytesting4964.CheckDesiredValue("TargetLogonId", "0x2020668")
	easytesting4964.CheckDesiredValue("TargetLogonGuid", "00000000-0000-0000-0000-000000000000")
	easytesting4964.CheckDesiredValue("SidList", "\r\n\t\t%{S-1-5-21-1913345275-1711810662-261465553-512}")

}
