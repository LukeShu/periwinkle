// Copyright 2015 Luke Shumaker

package negotiate

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
