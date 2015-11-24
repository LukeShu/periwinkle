// Copyright 2015 Luke Shumaker

package negotiate

import (
	"mime"
	"strings"
)

func (accept accept) Acceptable(contenttype string) (bool, error) {
	mediafulltype, mediaparams, err := mime.ParseMediaType(contenttype)
	if err != nil {
		return false, err
	}
	mediatypeParts := strings.SplitN(mediafulltype, "/", 2)
	mediatype := mediatypeParts[0]
	mediasubtype := mediatypeParts[1]

	if accept.Type != mediatype && accept.Type != "*" {
		return false, nil
	}
	if accept.Subtype != mediasubtype && accept.Subtype != "*" {
		return false, nil
	}
	for key, acceptval := range accept.TypeParams {
		mediaval, ok := mediaparams[key]
		if !ok || mediaval != acceptval {
			return false, nil
		}
	}
	return true, nil
}

type precedence struct {
	major, minor uint
}

func comparePrecedence(a precedence, b precedence) int8 {
	var cmpA, cmpB uint
	if a.major == b.major {
		cmpA = a.minor
		cmpB = b.minor
	} else {
		cmpA = a.major
		cmpB = b.major
	}
	switch {
	case cmpA == cmpB:
		return 0
	case cmpA < cmpB:
		return -1
	default:
		return 1
	}
}

// Higher numbers are higher precedence
//
// NB: it is undefined which of these has higher pecedence:
//  - "text/*;foo=bar"
//  - "text/plain"
// This interprets the latter to have higher precedence.
func (accept accept) Precedence() (ret precedence) {
	switch {
	case accept.Type == "*":
		ret.major = 0
	case accept.Subtype == "*":
		ret.major = 1
	default:
		ret.major = 2
	}
	ret.minor = uint(len(accept.TypeParams))
	return
}

func considerContentTypes(header *string, contenttypes []string) (max qvalue, quality map[string]qvalue, err error) {
	accepts, err := parseAccept(header)
	if err != nil {
		return -1, nil, err
	}
	quality = map[string]qvalue{}
	precedence := map[string]precedence{}
	for _, contenttype := range contenttypes {
		for _, accept := range accepts {
			acceptable, err := accept.Acceptable(contenttype)
			if err != nil {
				return -1, nil, err
			}
			if acceptable {
				old, alreadyset := precedence[contenttype]
				new := accept.Precedence()
				// NB: what do do in a precedence tie is undefined
				if !alreadyset || comparePrecedence(old, new) <= 0 {
					quality[contenttype] = accept.Weight
					precedence[contenttype] = new
					if max < accept.Weight {
						max = accept.Weight
					}
				}
			}
		}
	}
	return
}

func (accept acceptCharset) Acceptable(charset string) bool {
	return strings.EqualFold(accept.Charset, charset) || accept.Charset == "*"
}

func considerCharsets(header *string, charsets []string) (max qvalue, quality map[string]qvalue, err error) {
	accepts, err := parseAcceptCharset(header)
	if err != nil {
		return -1, nil, err
	}
	quality = map[string]qvalue{}
	for _, charset := range charsets {
		for _, accept := range accepts {
			if accept.Acceptable(charset) {
				quality[charset] = accept.Weight
				if max < accept.Weight {
					max = accept.Weight
				}
			}
		}
	}
	return
}

func (accept acceptEncoding) Acceptable(encoding string) bool {
	return strings.EqualFold(accept.Coding, encoding) || accept.Coding == "*"
}

func considerEncodings(header *string, encodings []string) (max qvalue, quality map[string]qvalue, err error) {
	accepts, err := parseAcceptEncoding(header)
	if err != nil {
		return -1, nil, err
	}
	quality = map[string]qvalue{}
	for _, encoding := range encodings {
		for _, accept := range accepts {
			if accept.Acceptable(encoding) {
				quality[encoding] = accept.Weight
				if max < accept.Weight {
					max = accept.Weight
				}
			}
		}
	}
	return
}

func considerLanguages(header *string, languageTags []string) (max qvalue, quality map[string]qvalue, err error) {
	accepts, err := parseAcceptLanguage(header)
	if err != nil {
		return -1, nil, err
	}
	quality = map[string]qvalue{}
	precedence := map[string]uint{}
	var def *acceptLanguage = nil
	for _, language := range languageTags {
		for _, accept := range accepts {
			if accept.LanguageRange == "*" {
				def = &accept
			} else {
				if languageFilterTag(accept.LanguageRange, language) {
					old, alreadyset := precedence[language]
					new := uint(strings.Count(string(accept.LanguageRange), "-")) + 1
					if !alreadyset || old <= new {
						quality[language] = accept.Weight
						precedence[language] = new
						if max < accept.Weight {
							max = accept.Weight
						}
					}
				}
			}
		}
		if _, matched := precedence[language]; def != nil && !matched {
			if languageFilterTag(def.LanguageRange, language) {
				quality[language] = def.Weight
				precedence[language] = 0
				if max < def.Weight {
					max = def.Weight
				}
			}
		}
	}
	return
}
