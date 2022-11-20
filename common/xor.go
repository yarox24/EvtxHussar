package common

func EncryptDecrypt(input string, key []byte) []byte {
	output := []byte(input)

	for i, val := range output {
		output[i] = val ^ key[i%len(key)]
	}

	return output
}
