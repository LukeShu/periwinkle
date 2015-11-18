// Copyright 2015 Luke Shumaker

package negotiate

import "testing"

// RFC 7231 says:
//
//     For example,
//
//       Accept: text/*;q=0.3, text/html;q=0.7, text/html;level=1,
//               text/html;level=2;q=0.4, */*;q=0.5
//
//     would cause the following values to be associated:
//
//     +-------------------+---------------+
//     | Media Type        | Quality Value |
//     +-------------------+---------------+
//     | text/html;level=1 | 1             |
//     | text/html         | 0.7           |
//     | text/plain        | 0.3           |
//     | image/jpeg        | 0.5           |
//     | text/html;level=2 | 0.4           |
//     | text/html;level=3 | 0.7           |
//     +-------------------+---------------+
func TestNegotiateContentType(t *testing.T) {
	header := "text/*;q=0.3, text/html;q=0.7, text/html;level=1, text/html;level=2;q=0.4, */*;q=0.5"
	// qvalues are fixed point * 1000
	correct := map[string]qvalue{
		"text/html;level=1": 1000,
		"text/html":         700,
		"text/plain":        300,
		"image/jpeg":        500,
		"text/html;level=2": 400,
		"text/html;level=3": 700,
	}
	contenttypes := make([]string, len(correct))
	i := uint(0)
	for contenttype := range correct {
		contenttypes[i] = contenttype
		i++
	}

	max, quality, err := NegotiateContentType(&header, contenttypes)
	if err != nil {
		t.Error(err)
	}
	if max != 1000 {
		t.Errorf("max:%d != 1000", max)
	}
	if len(quality) != len(contenttypes) {
		t.Errorf("len(quality):%d != len(contenttypes):%d", len(quality), len(contenttypes))
		for k, v := range quality {
			t.Logf(" - quality[%q]:%d", k, v)
		}
	}
	for k, corrv := range correct {
		testv, _ := quality[k]
		if testv != corrv {
			t.Errorf("quality[%q]:%d != corrv:%d", k, testv, corrv)
		}
	}
}

func TestNegotiateCharset(t *testing.T) {
	// TODO
}

func TestNegotiateEncoding(t *testing.T) {
	// TODO
}

func TestNegotiateLanguage(t *testing.T) {
	header := "da, en-gb;q=0.8, en;q=0.7, es-mx;q=0.5"
	// qvalues are fixed point * 1000
	correct := map[string]qvalue{
		"da":     1000,
		"da-foo": 1000,
		"en":     700,
		"en-us":  700,
		"en-gb":  800,
		"es":     0,
		"es-mx":  500,
	}

	languages := make([]string, len(correct))
	i := uint(0)
	for language := range correct {
		languages[i] = language
		i++
	}

	max, quality, err := NegotiateLanguage(&header, languages)
	if err != nil {
		t.Error(err)
	}
	if max != 1000 {
		t.Errorf("max:%d != 1000", max)
	}
	// subtract 1 because q=0 things aren't included
	if len(quality) != len(languages)-1 {
		t.Errorf("len(quality):%d != len(languages):%d", len(quality), len(languages))
		for k, v := range quality {
			t.Logf(" - quality[%q]:%d", k, v)
		}
	}
	for k, corrv := range correct {
		testv, _ := quality[k]
		if testv != corrv {
			t.Errorf("quality[%q]:%d != corrv:%d", k, testv, corrv)
		}
	}
}
