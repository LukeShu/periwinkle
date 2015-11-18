// Copyright 2015 Luke Shumaker

package negotiate

func isALPHA(char byte) bool {
	return (0x41 <= char && char <= 0x5A) || (0x61 <= char && char <= 0x7A)
}

func isDIGIT(char byte) bool {
	return 0x30 <= char && char <= 0x39
}

func isHTAB(char byte) bool {
	return char == 0x09
}

func isSP(char byte) bool {
	return char == 0x20
}

func isVCHAR(char byte) bool {
	return 0x21 <= char && char <= 0x7E
}
