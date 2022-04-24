package tests

import (
	"testing"
)

func TestAccountsUserRelatedOperationsSecurityEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	// 4720
	easytesting4720 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4720.json", eng, "AccountsUserRelatedOperations"))
	easytesting4720.CheckDesiredValue("TargetUserName", "someuser")
	easytesting4720.CheckDesiredValue("TargetDomainName", "DESKTOP-EvtxHussar")
	easytesting4720.CheckDesiredValue("TargetSid", "S-1-5-18")
	easytesting4720.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4720.CheckDesiredValue("SubjectUserName", "DESKTOP-EvtxHussar$")
	easytesting4720.CheckDesiredValue("SubjectDomainName", "HUSS")
	easytesting4720.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4720.CheckDesiredValue("PrivilegeList", "-")
	easytesting4720.CheckDesiredValue("SamAccountName", "someuser")
	easytesting4720.CheckDesiredValue("DisplayName", "-")
	easytesting4720.CheckDesiredValue("UserPrincipalName", "-")
	easytesting4720.CheckDesiredValue("HomeDirectory", "-")
	easytesting4720.CheckDesiredValue("HomePath", "-")
	easytesting4720.CheckDesiredValue("ScriptPath", "-")
	easytesting4720.CheckDesiredValue("ProfilePath", "-")
	easytesting4720.CheckDesiredValue("UserWorkstations", "-")
	easytesting4720.CheckDesiredValue("PasswordLastSet", "<never>")
	easytesting4720.CheckDesiredValue("AccountExpires", "<never>")
	easytesting4720.CheckDesiredValue("PrimaryGroupId", "513")
	easytesting4720.CheckDesiredValue("AllowedToDelegateTo", "-")
	easytesting4720.CheckDesiredValue("OldUacValue", "0x0")
	easytesting4720.CheckDesiredValue("NewUacValue", "SCRIPT (The logon script will be run) | LOCKOUT [Original value: 0x15]")
	easytesting4720.CheckDesiredValue("UserAccountControl", "\r\n\t\tAccount Disabled\r\n\t\t'Password Not Required' - Enabled\r\n\t\t'Normal Account' - Enabled")
	easytesting4720.CheckDesiredValue("UserParameters", "-")
	easytesting4720.CheckDesiredValue("SidHistory/SidList", "-")
	easytesting4720.CheckDesiredValue("LogonHours", "All")

	// 4724
	easytesting4724 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4724.json", eng, "AccountsUserRelatedOperations"))
	easytesting4724.CheckDesiredValue("TargetUserName", "anotheruser")
	easytesting4724.CheckDesiredValue("TargetDomainName", "HUSS")
	easytesting4724.CheckDesiredValue("TargetSid", "S-1-5-18")
	easytesting4724.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4724.CheckDesiredValue("SubjectUserName", "anotheruser")
	easytesting4724.CheckDesiredValue("SubjectDomainName", "HUSS")
	easytesting4724.CheckDesiredValue("SubjectLogonId", "0x104bd867a")

	// 4728
	easytesting4728 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4728.json", eng, "AccountsUserRelatedOperations"))
	easytesting4728.CheckDesiredValue("MemberName", "-")
	easytesting4728.CheckDesiredValue("MemberSid", "S-1-5-18")
	easytesting4728.CheckDesiredValue("TargetUserName", "None")
	easytesting4728.CheckDesiredValue("TargetDomainName", "DESKTOP-EvtxHussar")
	easytesting4728.CheckDesiredValue("TargetSid", "S-1-5-18")
	easytesting4728.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4728.CheckDesiredValue("SubjectUserName", "DESKTOP-EvtxHussar$")
	easytesting4728.CheckDesiredValue("SubjectDomainName", "HUSS")
	easytesting4728.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4728.CheckDesiredValue("PrivilegeList", "-")

	// 4738
	easytesting4738 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4738.json", eng, "AccountsUserRelatedOperations"))
	easytesting4738.CheckDesiredValue("TargetUserName", "Guest")
	easytesting4738.CheckDesiredValue("TargetDomainName", "HUSS")
	easytesting4738.CheckDesiredValue("TargetSid", "S-1-5-19")
	easytesting4738.CheckDesiredValue("SubjectUserSid", "S-1-5-19")
	easytesting4738.CheckDesiredValue("SubjectUserName", "Administrator")
	easytesting4738.CheckDesiredValue("SubjectDomainName", "HUSS")
	easytesting4738.CheckDesiredValue("SubjectLogonId", "0x70adad")
	easytesting4738.CheckDesiredValue("PrivilegeList", "-")
	easytesting4738.CheckDesiredValue("SamAccountName", "Guest")
	easytesting4738.CheckDesiredValue("DisplayName", "-")
	easytesting4738.CheckDesiredValue("UserPrincipalName", "-")
	easytesting4738.CheckDesiredValue("HomeDirectory", "-")
	easytesting4738.CheckDesiredValue("HomePath", "-")
	easytesting4738.CheckDesiredValue("ScriptPath", "-")
	easytesting4738.CheckDesiredValue("ProfilePath", "-")
	easytesting4738.CheckDesiredValue("UserWorkstations", "-")
	easytesting4738.CheckDesiredValue("PasswordLastSet", "12/24/2016 11:11:11 PM")
	easytesting4738.CheckDesiredValue("AccountExpires", "<never>")
	easytesting4738.CheckDesiredValue("PrimaryGroupId", "513")
	easytesting4738.CheckDesiredValue("AllowedToDelegateTo", "-")
	easytesting4738.CheckDesiredValue("OldUacValue", "SCRIPT (The logon script will be run) | LOCKOUT | NORMAL_ACCOUNT (It's a default account type that represents a typical user) [Original value: 0x215]")
	easytesting4738.CheckDesiredValue("NewUacValue", "LOCKOUT | NORMAL_ACCOUNT (It's a default account type that represents a typical user) [Original value: 0x214]")
	easytesting4738.CheckDesiredValue("UserAccountControl", "\r\n\t\tAccount Enabled")
	easytesting4738.CheckDesiredValue("UserParameters", "-")
	easytesting4738.CheckDesiredValue("SidHistory/SidList", "-")
	easytesting4738.CheckDesiredValue("LogonHours", "All")

	// 4781
	easytesting4781 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4781.json", eng, "AccountsUserRelatedOperations"))
	easytesting4781.CheckDesiredValue("OldTargetUserName", "Administrators")
	easytesting4781.CheckDesiredValue("NewTargetUserName", "Administrators")
	easytesting4781.CheckDesiredValue("TargetDomainName", "Builtin")
	easytesting4781.CheckDesiredValue("TargetSid", "S-1-5-32-544")
	easytesting4781.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4781.CheckDesiredValue("SubjectUserName", "DESKTOP-EvtxHussar$")
	easytesting4781.CheckDesiredValue("SubjectDomainName", "WORKGROUP")
	easytesting4781.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4781.CheckDesiredValue("PrivilegeList", "-")

}
