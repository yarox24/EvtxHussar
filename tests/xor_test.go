package tests

import (
	"bytes"
	"github.com/yarox24/EvtxHussar/common"
	"testing"
)

func TestXOR_forPowerShellScriptBlocks(t *testing.T) {

	// "# Copyright © 2" xored with key: "EH"
	xored_array := []byte{0x66, 0x68, 0x06, 0x27, 0x35, 0x31, 0x37, 0x21, 0x22, 0x20, 0x31, 0x68, 0x87, 0xe1, 0x65, 0x7a}
	xor_key := []byte{'E', 'H'}

	// Valid expected result
	valid_array := []byte{0x23, 0x20, 0x43, 0x6F, 0x70, 0x79, 0x72, 0x69, 0x67, 0x68, 0x74, 0x20, 0xC2, 0xA9, 0x20, 0x32}
	output_array := common.EncryptDecrypt(string(xored_array), xor_key)

	if bytes.Compare(output_array, valid_array) != 0 {
		t.Errorf("common.EncryptDecrypt() - XOR error. Output mismatch")
	}

}
