// Copyright 2015 Luke Shumaker

package negotiate

import "fmt"

func isToken(tok string) bool {
	// 7230 is specified using the US-ASCII core rules; if there
	// is UTF-8 data here, it's illegal in a token, so we can just
	// do byte handling.
	for _, b := range []byte(tok) {
		switch b {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			// do nothing
		default:
			if isDIGIT(b) || isALPHA(b) {
				// do nothing
			} else {
				return false
			}
		}
	}
	return true
}

func isQdtext(char byte) bool {
	return isHTAB(char) || isSP(char) || char == 0x21 || (0x23 <= char && char <= 0x5B) || (0x5D <= char && char <= 0x7E) || isObsText(char)
}

func isObsText(char byte) bool {
	return 0x80 <= char && char <= 0xFF
}

func parseQuotedString(str string) (string, error) {
	// strip the quotes
	if str[0] != '"' || str[len(str)-1] != '"' {
		return "", fmt.Errorf("%q does not look like an RFC 7230 quoted-string", str)
	}
	in := str[1 : len(str)-1]
	// main algo
	out := make([]byte, len(in))
	o := uint(0)
	for i := uint(0); i < uint(len(in)); i++ {
		switch {
		case isQdtext(in[i]):
			out[o] = in[i]
			o++
		case in[i] == '\\':
			i++
			switch {
			case isHTAB(in[i]) || isSP(in[i]) || isVCHAR(in[i]) || isObsText(in[i]):
				out[o] = in[i]
				o++
			default:
				return "", fmt.Errorf("%c is not a legal character to follow a backslash in an RFC 7230 quoted-string", in[i])
			}
		default:
			return "", fmt.Errorf("%c is not a legal character in an RFC 7230 quoted-string", in[i])
		}
	}
	// return
	return string(out[:o]), nil
}
