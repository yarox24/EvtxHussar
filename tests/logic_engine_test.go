package tests

import (
	"github.com/Velocidex/ordereddict"
	"github.com/yarox24/EvtxHussar/engine"
	"testing"
)

func PrepareEvent4625(t *testing.T, eng *engine.Engine, global_logic string, Container_Or [][]string, Container_And [][]string) []engine.ChanFullInfo {

	le := engine.NewLogicEngine()
	l2_name := "LogonsUniversal"

	// Prepare supported EID's
	supported_eids := engine.ChannelTOEID{
		Eid: make(map[string]engine.EIDToChan, 0),
	}

	e4625existing := eng.EventsCache[l2_name]["Security"]["4625"]

	var cfi = engine.ChanFullInfo{
		Chan:              nil,
		Provider_guid:     e4625existing.Provider_guid,
		Provider_name:     e4625existing.Provider_name,
		Attrib_extraction: e4625existing.Attrib_extraction,
		Matching_Rules: engine.MatchingRulesT{
			Global_Logic:          global_logic,
			Container_Or:          Container_Or,
			Container_OrEnhanced:  nil,
			Container_And:         Container_And,
			Container_AndEnhanced: nil,
		},
	}
	cfi.Matching_Rules.EnhanceRulesInPlace()

	eid_4625 := supported_eids.Eid["4625"]
	eid_4625.Chans = make([]engine.ChanFullInfo, 0)
	eid_4625.Chans = append(eid_4625.Chans, cfi)
	supported_eids.Eid["4625"] = eid_4625
	le.SetSupportedEIDs(supported_eids)

	//// 4625
	ev := UnmarshallEvent("Security_4625.json", eng)
	ev_map, _ := ordereddict.GetMap(ev, "Event")
	return le.ReturnMatchingChanFullInfo(ev_map)
}

func ExpectedOutcome(cfi []engine.ChanFullInfo, t *testing.T, error_string string, outcome bool) {
	if outcome == true {
		if len(cfi) == 0 {
			t.Errorf(error_string)
		}
	} else {
		if len(cfi) > 0 {
			t.Errorf(error_string)
		}
	}
}

func TestLogicEngine(t *testing.T) {

	// Load Engine
	eng := LoadEngine()

	// A) Global logic = OR

	// A1) Status == 3221225581 (Single rule) | Match
	ExpectedOutcome(PrepareEvent4625(t, eng, "or", [][]string{{"single_match:Function=DecimalEqual,Field=Status,Value=3221225581"}}, [][]string{}),
		t, "Logic - A1) Status == 3221225581 (Single rule) - Failed!\n", true)

	// A2) Status == 555 (Single rule) | No Match
	ExpectedOutcome(PrepareEvent4625(t, eng, "or", [][]string{{"single_match:Function=DecimalEqual,Field=Status,Value=555"}}, [][]string{}),
		t, "Logic - A2) Status == 555 (Single rule) - Failed!\n", false)

	// A3) Status == 3221225581 (Single rule) | Match
	ExpectedOutcome(PrepareEvent4625(t, eng, "or", [][]string{
		{"single_match:Function=DecimalEqual,Field=Status,Value=1", "single_match:Function=DecimalEqual,Field=Status,Value=2", "single_match:Function=DecimalEqual,Field=Status,Value=3",
			"single_match:Function=DecimalEqual,Field=Status,Value=4", "single_match:Function=DecimalEqual,Field=Status,Value=555555", "single_match:Function=DecimalEqual,Field=Status,Value=3221225581"},
	}, [][]string{}), t, "Logic - A2) Status == 3221225581 (Multiple rules) - Failed!\n", true)

	// A4) Status == 2222 (Single rule) | Match
	ExpectedOutcome(PrepareEvent4625(t, eng, "or", [][]string{
		{"single_match:Function=DecimalEqual,Field=Status,Value=1", "single_match:Function=DecimalEqual,Field=Status,Value=2", "single_match:Function=DecimalEqual,Field=Status,Value=3",
			"single_match:Function=DecimalEqual,Field=Status,Value=4", "single_match:Function=DecimalEqual,Field=Status,Value=555555", "single_match:Function=DecimalEqual,Field=Status,Value=2222"},
	}, [][]string{}), t, "Logic - A2) Status == 3221225581 (Multiple rules) - Failed!\n", false)

	// B) Global logic = AND
	// B1) AuthenticationPackageName=NTLM AND WorkstationName=HUSStext | Match
	ExpectedOutcome(PrepareEvent4625(t, eng, "and",
		[][]string{},
		[][]string{{"single_match:Function=Substring,Field=AuthenticationPackageName,Value=NTLM", "single_match:Function=Substring,Field=WorkstationName,Value=HUSStext"}},
	), t, "Logic - B1 - Failed!\n", true)

	// B2) AuthenticationPackageName=NTLM AND WorkstationName=Failedd | Match
	ExpectedOutcome(PrepareEvent4625(t, eng, "and",
		[][]string{},
		[][]string{{"single_match:Function=Substring,Field=AuthenticationPackageName,Value=NTLM", "single_match:Function=Substring,Field=WorkstationName,Value=Failedd"}},
	), t, "Logic - B2 - Failed!\n", false)

	// Global OR check
	// OR1)
	ExpectedOutcome(PrepareEvent4625(t, eng, "or",
		[][]string{{"single_match:Function=DecimalEqual,Field=Status,Value=555"}},
		[][]string{{"single_match:Function=Substring,Field=AuthenticationPackageName,Value=NTLM", "single_match:Function=Substring,Field=WorkstationName,Value=HUSStext"}},
	), t, "Logic - OR1 - Failed!\n", true)

	// Global AND check
	// AND1)
	ExpectedOutcome(PrepareEvent4625(t, eng, "and",
		[][]string{{"single_match:Function=DecimalEqual,Field=Status,Value=555"}},
		[][]string{{"single_match:Function=Substring,Field=AuthenticationPackageName,Value=NTLM", "single_match:Function=Substring,Field=WorkstationName,Value=HUSStext"}},
	), t, "Logic - OR1 - Failed!\n", false)

}
