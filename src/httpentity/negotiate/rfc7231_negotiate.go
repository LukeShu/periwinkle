// Copyright 2015 Luke Shumaker

package negotiate

// NegotiateContentType negotiates the content type (use the `Accept:` header).
func NegotiateContentType(header *string, contenttypes []string) (options []string, err error) {
	max, qualities, err := considerContentTypes(header, contenttypes)
	if err != nil {
		return nil, err
	}
	options = []string{}
	for val, quality := range qualities {
		if quality == max {
			options = append(options, val)
		}
	}
	return
}

// NegotiateCharset negotiates the character set (use the `Accept-Charset:` header).
func NegotiateCharset(header *string, charsets []string) (options []string, err error) {
	max, qualities, err := considerCharsets(header, charsets)
	if err != nil {
		return nil, err
	}
	options = []string{}
	for val, quality := range qualities {
		if quality == max {
			options = append(options, val)
		}
	}
	return
}

// NegotiateEncoding negotiates the encoding (use the `Accept-Encoding:` header).
func NegotiateEncoding(header *string, encodings []string) (options []string, err error) {
	max, qualities, err := considerEncodings(header, encodings)
	if err != nil {
		return nil, err
	}
	options = []string{}
	for val, quality := range qualities {
		if quality == max {
			options = append(options, val)
		}
	}
	return
}

// NegotiateLanguage negotiates the language (use the `Accept-Language:` header).
func NegotiateLanguage(header *string, languageTags []string) (options []string, err error) {
	max, qualities, err := considerContentTypes(header, languageTags)
	if err != nil {
		return nil, err
	}
	options = []string{}
	for val, quality := range qualities {
		if quality == max {
			options = append(options, val)
		}
	}
	return
}
