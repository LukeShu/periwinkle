// Copyright 2015 Luke Shumaker

package negotiate

import (
	"regexp"
	"strings"
)

type languageRange string

var languageRangeRegexp = regexp.MustCompile(`^([a-zA-Z]{1,8}(-[a-zA-Z0-9]{1,8})*|\*)$`)

func isLanguageRange(str string) bool {
	return languageRangeRegexp.Match([]byte(str))
}

func languageFilterTag(languageRange, languageTag string) bool {
	if languageRange == "*" {
		return true
	}
	languageRange = strings.ToLower(languageRange)
	languageTag = strings.ToLower(languageTag)
	return languageTag == languageRange || strings.HasPrefix(languageTag, languageRange+"-")
}
