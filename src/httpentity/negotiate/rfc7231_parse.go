// Copyright 2015 Luke Shumaker

package negotiate

import (
	"fmt"
	"strconv"
	"strings"
)

// A fixed-point representation of RFC 7231 `qvalue`
//
// This means that valid values are 0-1000.  Uses a signed integer to
// allow using `-1` to mean "invalid".
type qvalue int16

func parseQvalue(str string) qvalue {
	if len(str) < 1 {
		return -1
	}
	ret := qvalue(-1)
	// check the digit to the left of the point
	switch str[0] {
	case '0':
		ret = 0
	case '1':
		ret = 1000
	default:
		return -1
	}
	if len(str) == 1 {
		return ret
	}
	// check for the point
	if str[1] != '.' {
		return -1
	}
	// strip the runes to the left of the point
	str = str[2:]
	// check the fractional part
	i, err := strconv.Atoi(str)
	if err != nil || i < 0 {
		return -1
	}
	switch len(str) {
	case 1:
		ret += qvalue(i * 100)
	case 2:
		ret += qvalue(i * 10)
	case 3:
		ret += qvalue(i)
	default:
		return -1
	}
	// make sure it wasn't too high
	if ret > 1000 {
		return -1
	}
	return ret
}

// a representation of RFC 7231 `Accept`
type accept struct {
	Type         string
	Subtype      string
	Weight       qvalue
	TypeParams   map[string]string
	AcceptParams map[string]string
}

// this captures both `parameter` and accept-parames/accept-ext, with
// the exception that it converts the accept-ext LHS to lowercase,
// which it shouldn't do.
func parseParameter(lhs, rhs string) (string, string, error) {
	var err error
	if !isToken(lhs) {
		return "", "", fmt.Errorf("%q is not a valid RFC 7230 token", lhs)
	}
	if rhs[0] == '"' {
		rhs, err = parseQuotedString(rhs)
		if err != nil {
			return "", "", fmt.Errorf("%q is not a valid RFC 7230 token or [quoted string]", rhs)
		}
	} else {
		if !isToken(rhs) {
			return "", "", fmt.Errorf("%q is not a valid RFC 7230 [token] or quoted string", rhs)
		}
	}
	lhs = strings.ToLower(lhs)
	return lhs, rhs, nil
}

func parseAccept(header *string) ([]accept, error) {
	if header == nil {
		return []accept{{
			Type:         "*",
			Subtype:      "*",
			Weight:       1000,
			TypeParams:   map[string]string{},
			AcceptParams: map[string]string{},
		}}, nil
	}
	members := strings.Split(*header, ",")
	if strings.TrimSpace(*header) == "" {
		members = []string{}
	}
	ret := make([]accept, len(members))
	for i, member := range members {
		parts := strings.Split(member, ";")
		ret[i] = accept{
			Type:         "*",
			Subtype:      "*",
			TypeParams:   make(map[string]string),
			Weight:       -1,
			AcceptParams: make(map[string]string),
		}
		acc := &ret[i]
		mediatype := strings.Split(strings.ToLower(strings.TrimSpace(parts[0])), "/")
		if len(mediatype) != 2 {
			return nil, perrorf("Accept", "%q is not a valid RFC 7231 media-range", member)
		}
		acc.Type = mediatype[0]
		acc.Subtype = mediatype[1]
		if acc.Type == "*" && acc.Subtype != "*" {
			return nil, perrorf("Accept", "Media type cannot be * if subtype isn't also")
		}
		if acc.Type != "*" && !isToken(acc.Type) {
			return nil, perrorf("Accept", "%q is not a valid RFC 7230 token", acc.Type)
		}
		if acc.Subtype != "*" && !isToken(acc.Subtype) {
			return nil, perrorf("Accept", "%q is not a valid RFC 7230 token", acc.Type)
		}
		for _, param := range parts[1:] {
			asgn := strings.SplitN(strings.TrimSpace(param), "=", 2)
			if len(asgn) != 2 {
				return nil, perrorf("Accept", "%q contains the wrong number of \"=\" to be a valid RFC 7231 parameter", param)
			}
			if acc.Weight < 0 {
				if strings.ToLower(asgn[0]) == "q" {
					acc.Weight = parseQvalue(asgn[1])
					if acc.Weight < 0 {
						return nil, perrorf("Accept", "%q is not a valid RFC 7231 qvalue", asgn[1])
					}
				} else {
					lhs, rhs, err := parseParameter(asgn[0], asgn[1])
					if err != nil {
						return nil, ParseError{"Accept", err}
					}
					acc.TypeParams[lhs] = rhs
				}
			} else {
				_, rhs, err := parseParameter(asgn[0], asgn[1])
				if err != nil {
					return nil, ParseError{"Accept", err}
				}
				// use asgn[0] instead of lhs because
				// parseParameter downcases it.
				acc.AcceptParams[asgn[0]] = rhs
			}
		}
		if acc.Weight < 0 {
			acc.Weight = 1000
		}
	}
	return ret, nil
}

type acceptCharset struct {
	Charset string
	Weight  qvalue
}

func parseAcceptCharset(header *string) ([]acceptCharset, error) {
	if header == nil {
		return []acceptCharset{{
			Charset: "*",
			Weight:  1000,
		}}, nil
	}
	if strings.TrimSpace(*header) == "" {
		return nil, perrorf("Accept-Charset", "May not be empty")
	}
	members := strings.Split(strings.ToLower(*header), ",")
	ret := make([]acceptCharset, len(members))
	for i, member := range members {
		parts := strings.Split(member, ";")
		switch len(parts) {
		case 1:
			charset := strings.TrimSpace(parts[0])
			if charset != "*" && !isToken(charset) {
				return nil, perrorf("Accept-Charset", "%q is not a valid RFC 7231 charset", charset)
			}
			ret[i].Charset = charset
			ret[i].Weight = 1000
		case 2:
			charset := strings.TrimSpace(parts[0])
			if charset != "*" && !isToken(charset) {
				return nil, perrorf("Accept-Charset", "%q is not a valid RFC 7231 charset", charset)
			}
			ret[i].Charset = charset
			weight := strings.TrimSpace(parts[1])
			if !strings.HasPrefix(weight, "q=") {
				return nil, perrorf("Accept-Charset", "%q is not a valid RFC 7231 weight", ";"+weight)
			}
			ret[i].Weight = parseQvalue(weight[2:])
			if ret[i].Weight < 0 {
				return nil, perrorf("Accept-Charset", "%q is not a valid RFC 7231 qvalue", weight[2:])
			}
		default:
			return nil, perrorf("Accept-Charset", "%q is not a valid RFC 7231 charset/weight pair", member)
		}
	}
	return ret, nil
}

type acceptEncoding struct {
	Coding string
	Weight qvalue
}

func parseAcceptEncoding(header *string) ([]acceptEncoding, error) {
	if header == nil {
		return []acceptEncoding{{
			Coding: "*",
			Weight: 1000,
		}}, nil
	}
	members := strings.Split(strings.ToLower(*header), ",")
	if strings.TrimSpace(*header) == "" {
		members = []string{}
	}
	// don't require >= 1 entries
	ret := make([]acceptEncoding, len(members))
	for i, member := range members {
		parts := strings.Split(member, ";")
		switch len(parts) {
		case 1:
			coding := strings.TrimSpace(parts[0])
			if coding != "*" && !isToken(coding) {
				return nil, perrorf("Accept-Encoding", "%q is not a valid RFC 7231 codings", coding)
			}
			ret[i].Coding = coding
			ret[i].Weight = 1000
		case 2:
			coding := strings.TrimSpace(parts[0])
			if coding != "*" && !isToken(coding) {
				return nil, perrorf("Accept-Encoding", "%q is not a valid RFC 7231 codings", coding)
			}
			ret[i].Coding = coding
			weight := strings.TrimSpace(parts[1])
			if !strings.HasPrefix(weight, "q=") {
				return nil, perrorf("Accept-Encoding", "%q is not a valid RFC 7231 weight", ";"+weight)
			}
			ret[i].Weight = parseQvalue(weight[2:])
			if ret[i].Weight < 0 {
				return nil, perrorf("Accept-Encoding", "%q is not a valid RFC 7231 qvalue", weight[2:])
			}
		default:
			return nil, perrorf("Accept-Encoding", "%q is not a valid RFC 7231 codings/weight pair", member)
		}
	}
	return ret, nil
}

type acceptLanguage struct {
	LanguageRange string
	Weight        qvalue
}

func parseAcceptLanguage(header *string) ([]acceptLanguage, error) {
	if header == nil {
		return []acceptLanguage{{
			LanguageRange: "*",
			Weight:        1000,
		}}, nil
	}
	if strings.TrimSpace(*header) == "" {
		return nil, perrorf("Accept-Language", "May not be empty")
	}
	members := strings.Split(*header, ",")
	ret := make([]acceptLanguage, len(members))
	for i, member := range members {
		parts := strings.Split(member, ";")
		switch len(parts) {
		case 1:
			languageRange := strings.TrimSpace(parts[0])
			if languageRange != "*" && !isLanguageRange(languageRange) {
				return nil, perrorf("Accept-Language", "%q is not a valid RFC 7231 language-range", languageRange)
			}
			ret[i].LanguageRange = languageRange
			ret[i].Weight = 1000
		case 2:
			languageRange := strings.TrimSpace(parts[0])
			if languageRange != "*" && !isLanguageRange(languageRange) {
				return nil, perrorf("Accept-Language", "%q is not a valid RFC 7231 language-range", languageRange)
			}
			ret[i].LanguageRange = languageRange
			weight := strings.TrimSpace(parts[1])
			if !strings.HasPrefix(weight, "q=") {
				return nil, perrorf("Accept-Language", "%q is not a valid RFC 7231 weight", ";"+weight)
			}
			ret[i].Weight = parseQvalue(weight[2:])
			if ret[i].Weight < 0 {
				return nil, perrorf("Accept-Language", "%q is not a valid RFC 7231 qvalue", weight[2:])
			}
		default:
			return nil, perrorf("Accept-Language", "%q is not a valid RFC 7231 language-range/weight pair", member)
		}
	}
	return ret, nil
}
