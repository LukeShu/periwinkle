// Copyright 2015 Luke Shumaker

package negotiate

import "testing"

// Test parsing several example headers from RFC 7231
func TestParseAccept(t *testing.T) {
	correct := map[string]int{
		"": 0,
		"audio/*; q=0.2, audio/basic":                                                          2,
		"text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c":                            4,
		"text/*, text/plain, text/plain;format=flowed, */*":                                    4,
		"text/*;q=0.3, text/html;q=0.7, text/html;level=1, text/html;level=2;q=0.4, */*;q=0.5": 5,
	}
	for header, count := range correct {
		accepts, err := parseAccept(&header)
		if count < 0 {
			if accepts != nil || err == nil {
				t.Errorf("%q", header)
			}
		} else {
			if len(accepts) != count || err != nil {
				t.Errorf("%q", header)
			}
		}
	}
}

// Test parsing several example headers from RFC 7231
func TestParseAcceptCharset(t *testing.T) {
	correct := map[string]int{
		"": -1,
		"iso-8859-5, unicode-1-1;q=0.8": 2,
	}
	for header, count := range correct {
		accepts, err := parseAcceptCharset(&header)
		if count < 0 {
			if accepts != nil || err == nil {
				t.Errorf("%q", header)
			}
		} else {
			if len(accepts) != count || err != nil {
				t.Errorf("%q", header)
			}
		}
	}
}

// Test parsing several example headers from RFC 7231
func TestParseAcceptEncoding(t *testing.T) {
	correct := map[string]int{
		"compress, gzip": 2,
		"":               0,
		"*":              1,
		"compress;q=0.5, gzip;q=1.0":         2,
		"gzip;q=1.0, identity; q=0.5, *;q=0": 3,
	}
	for header, count := range correct {
		accepts, err := parseAcceptEncoding(&header)
		if count < 0 {
			if accepts != nil || err == nil {
				t.Errorf("%q", header)
			}
		} else {
			if len(accepts) != count || err != nil {
				t.Errorf("%q", header)
			}
		}
	}
}

// Test parsing several example headers from RFC 7231
func TestParseAcceptLanguage(t *testing.T) {
	correct := map[string]int{
		"": -1,
		"da, en-gb;q=0.8, en;q=0.7": 3,
	}
	for header, count := range correct {
		accepts, err := parseAcceptLanguage(&header)
		if count < 0 {
			if accepts != nil || err == nil {
				t.Errorf("%q", header)
			}
		} else {
			if len(accepts) != count || err != nil {
				t.Errorf("%q", header)
			}
		}
	}
}
