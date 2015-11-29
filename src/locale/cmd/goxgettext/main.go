// Copyright 2015 Luke Shumaker

package main

import (
	"fmt"
	"os"

	"go/parser"
	"go/token"

	"path/filepath"

	"locale/xgettext"
)

func parseFile(filename string) (*xgettext.File, error) {
	fset := token.NewFileSet()
	AST, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	return &xgettext.File{
		AST:  AST,
		FSet: fset,
	}, err
}

func getFile(searchpath []string, filename string) (*xgettext.File, error) {
	if filepath.IsAbs(filename) {
		return parseFile(filename)
	}
	for _, dir := range searchpath {
		file, err := parseFile(filepath.Join(dir, filename))
		if file != nil {
			return file, err
		}
	}
	return parseFile(filename)
}

// If KEYWORDSPEC is a C identifier ID, ‘xgettext’ looks for strings
// in the first argument of each call to the function ID.
//
// If KEYWORDSPEC is of the form ‘ID:ARGNUM’, ‘xgettext’ looks for
// strings in the ARGNUMth argument of the call.
//
// If KEYWORDSPEC is of the form ‘ID:ARGNUM1,ARGNUM2’, ‘xgettext’
// looks for strings in the ARGNUM1st argument and in the ARGNUM2nd
// argument of the call, and treats them as singular/plural variants
// for a message with plural handling.
//
// If KEYWORDSPEC is of the form ‘ID:CONTEXTARGNUMc,ARGNUM’ or
// ‘ID:ARGNUM,CONTEXTARGNUMc’, ‘xgettext’ treats strings in the
// CONTEXTARGNUMth argument as a context specifier.
//
// If KEYWORDSPEC is of the form ‘ID:…,TOTALNUMARGSt’,
// ‘xgettext’ recognizes this argument specification only if the
// number of actual arguments is equal to TOTALNUMARGS.
//
// If KEYWORDSPEC is of the form ‘ID:ARGNUM...,"XCOMMENT"’,
// ‘xgettext’, when extracting a message from the specified argument
// strings, adds an extracted comment XCOMMENT to the message.  Note
// that when used through a normal shell command line, the
// double-quotes around the XCOMMENT need to be escaped.
func parseKeywordSpec(str string) xgettext.KeywordSpec {
	panic("TODO")
}

func main() {
	//options, err := docopt.Parse(fmt.Sprintf(usage, os.Args[0]), os.Args[1:], true, "0.1", false, false)
	//fmt.Printf("options: %T(%#v)\n", options, options)
	//fmt.Printf("err: %T(%#v)\n", err, err)
	searchpath := []string{""}
	for _, filename := range os.Args[1:] {
		file, err := getFile(searchpath, filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %v\n", os.Args[0], filename, err)
		}
		if file != nil {
			msgs := file.ExtractStrings(true, []xgettext.KeywordSpec{})
			for msg := range msgs {
				fmt.Println(msg)
			}
		}
	}
}
