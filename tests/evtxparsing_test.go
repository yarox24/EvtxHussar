package tests

import (
	"testing"
)

func TestEvtxSecurity4624(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	efi_after_parsing1 := FakeInspectEvtx(eng.Maps_path, "Security4624.evtx")

	UniversalCheckDesiredValue(t, efi_after_parsing1.GetChannel(), "Security")
	UniversalCheckDesiredValue(t, efi_after_parsing1.GetLatestComputer(), "win10")
	UniversalCheckDesiredValueBoolean(t, efi_after_parsing1.IsEmpty(), false)
	UniversalCheckDesiredValueBoolean(t, efi_after_parsing1.IsValid(), true)

	// Alternative  (Fixed) header checks
	aheader := efi_after_parsing1.GetAlternativeHeader()

	UniversalCheckDesiredValueUint64(t, aheader.OldestChunk, 0)
	UniversalCheckDesiredValueUint64(t, aheader.CurrentChunkNum, 81)
	UniversalCheckDesiredValueUint64(t, aheader.NextRecordNum, 7210)
	UniversalCheckDesiredValueUint64(t, uint64(aheader.ChunkCount), 82)
	UniversalCheckDesiredValueUint64(t, uint64(aheader.FileFlags), 0)
	UniversalCheckDesiredValueUint64(t, uint64(aheader.CheckSum), 34572473)

	// Test events count
	records_counter := TestL1Worker(eng.Maps_path, "Security4624.evtx")
	UniversalCheckDesiredValueint64(t, records_counter, 7209)

}

func TestEvtxApplication15(t *testing.T) {

	// Load Engine
	eng := LoadEngine() //

	efi_after_parsing1 := FakeInspectEvtx(eng.Maps_path, "Application15.evtx")

	UniversalCheckDesiredValue(t, efi_after_parsing1.GetChannel(), "Application")
	UniversalCheckDesiredValue(t, efi_after_parsing1.GetLatestComputer(), "win10")
	UniversalCheckDesiredValueBoolean(t, efi_after_parsing1.IsEmpty(), false)
	UniversalCheckDesiredValueBoolean(t, efi_after_parsing1.IsValid(), true)

	// Alternative  (Fixed) header checks
	aheader := efi_after_parsing1.GetAlternativeHeader()

	UniversalCheckDesiredValueUint64(t, aheader.OldestChunk, 0)
	UniversalCheckDesiredValueUint64(t, aheader.CurrentChunkNum, 0)
	UniversalCheckDesiredValueUint64(t, aheader.NextRecordNum, 2)
	UniversalCheckDesiredValueUint64(t, uint64(aheader.ChunkCount), 1)
	UniversalCheckDesiredValueUint64(t, uint64(aheader.FileFlags), 0)
	UniversalCheckDesiredValueUint64(t, uint64(aheader.CheckSum), 1056700737)

	// Test events count
	records_counter := TestL1Worker(eng.Maps_path, "Application15.evtx")
	UniversalCheckDesiredValueint64(t, records_counter, 1)

}
