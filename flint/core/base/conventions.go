package base

import "os"

func BrokeConention(input string, conv string, reason string) {
	println("cannot use string " + input + " as it breaks convention " + conv + " because " + reason)
	os.Exit(1)
}

func RFC1123(input string) string {
	if len(input) > 253 {
		BrokeConention(input, "RFC1123", "length > 253")
	}

	return input
}
