package tests

import (
	"testing"
)

func TestComputerAccountsRelatedOperationsSecurityEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 4741
	easytesting4741 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4741.json", eng, "ActiveDirectoryComputerAccounts"))
	easytesting4741.CheckDesiredValue("TargetUserName", "DESKTOP-EvtxHussar$")
	easytesting4741.CheckDesiredValue("TargetDomainName", "HUSS")
	easytesting4741.CheckDesiredValue("TargetSid", "S-1-5-19")
	easytesting4741.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4741.CheckDesiredValue("SubjectUserName", "Administrator")
	easytesting4741.CheckDesiredValue("SubjectDomainName", "HUSS")
	easytesting4741.CheckDesiredValue("SubjectLogonId", "0x241d46")
	easytesting4741.CheckDesiredValue("PrivilegeList", "-")
	easytesting4741.CheckDesiredValue("SamAccountName", "DESKTOP-EvtxHussar$")
	easytesting4741.CheckDesiredValue("DisplayName", "-")
	easytesting4741.CheckDesiredValue("UserPrincipalName", "-")
	easytesting4741.CheckDesiredValue("HomeDirectory", "-")
	easytesting4741.CheckDesiredValue("HomePath", "-")
	easytesting4741.CheckDesiredValue("ScriptPath", "-")
	easytesting4741.CheckDesiredValue("ProfilePath", "-")
	easytesting4741.CheckDesiredValue("UserWorkstations", "-")
	easytesting4741.CheckDesiredValue("PasswordLastSet", "<never>")
	easytesting4741.CheckDesiredValue("AccountExpires", "<never>")
	easytesting4741.CheckDesiredValue("PrimaryGroupId", "515")
	easytesting4741.CheckDesiredValue("AllowedToDelegateTo", "-")
	easytesting4741.CheckDesiredValue("OldUacValue", "0x0")
	easytesting4741.CheckDesiredValue("NewUacValue", "SCRIPT (The logon script will be run) | ENCRYPTED_TEXT_PWD_ALLOWED (The user can send an encrypted password) [Original value: 0x85]")
	easytesting4741.CheckDesiredValue("UserAccountControl", "\r\n\t\tAccount Disabled\r\n\t\t'Password Not Required' - Enabled\r\n\t\t'Workstation Trust Account' - Enabled")
	easytesting4741.CheckDesiredValue("UserParameters", "-")
	easytesting4741.CheckDesiredValue("SidHistory", "-")
	easytesting4741.CheckDesiredValue("LogonHours", "-")
	easytesting4741.CheckDesiredValue("DnsHostName", "-")
	easytesting4741.CheckDesiredValue("ServicePrincipalNames", "-")

	// 4743
	easytesting4743 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4743.json", eng, "ActiveDirectoryComputerAccounts"))
	easytesting4743.CheckDesiredValue("TargetUserName", "DESKTOP-EvtxHussar$")
	easytesting4743.CheckDesiredValue("TargetDomainName", "HUSS")
	easytesting4743.CheckDesiredValue("TargetSid", "S-1-5-19")
	easytesting4743.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4743.CheckDesiredValue("SubjectUserName", "Administrator")
	easytesting4743.CheckDesiredValue("SubjectDomainName", "HUSS")
	easytesting4743.CheckDesiredValue("SubjectLogonId", "0x241d46")
	easytesting4743.CheckDesiredValue("PrivilegeList", "-")

}
