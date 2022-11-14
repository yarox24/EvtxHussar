package tests

import (
	"testing"
)

func TestAuditPolicyChangedEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine()

	// 4719
	easytesting4719 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4719.json", eng, "AuditPolicyChanged"))
	easytesting4719.CheckDesiredValue("SubjectUserSid", "S-1-5-18")
	easytesting4719.CheckDesiredValue("SubjectUserName", "Yar0$")
	easytesting4719.CheckDesiredValue("SubjectDomainName", "HUSS")
	easytesting4719.CheckDesiredValue("SubjectLogonId", "0x3e7")
	easytesting4719.CheckDesiredValue("CategoryId", "Logon/Logoff")
	easytesting4719.CheckDesiredValue("SubcategoryId", "Logon")
	easytesting4719.CheckDesiredValue("SubcategoryGuid", "0CCE9215-69AE-11D9-BED3-505054503030")
	easytesting4719.CheckDesiredValue("AuditPolicyChanges", "Success removed, Failure removed")

	// 4912
	//easytesting4912 := NewEasyTesting(t, UnmarshallAndParseEvent("Security_4912.json", eng, "AuditPolicyChanged"))
}
