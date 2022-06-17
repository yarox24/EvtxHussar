package tests

import (
	"testing"
)

func TestAVSymantecEvents(t *testing.T) {

	// Load Engine
	eng := LoadEngine()

	// 400
	easytesting400 := NewEasyTesting(t, UnmarshallAndParseEvent("Application_400.json", eng, "AV_SymantecNetwork"))
	easytesting400.CheckDesiredValue("Description", "[SID: 30369] Audit: Nessus Vulnerability Scanner Activity 3 attack detected but not blocked. Application path: C:\\PROGRAM FILES\\JAVA\\JAVA.EXE")
	easytesting400.CheckDesiredValue("Description (Path)", "C:\\PROGRAM FILES\\JAVA\\JAVA.EXE")

}
