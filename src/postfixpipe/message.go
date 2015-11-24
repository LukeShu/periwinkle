// Copyright 2015 Luke Shumaker

// Package postfixpipe provides helpers for writing programs for
// Postfix dot-forwards and aliases.
package postfixpipe

import (
	"bufio"
	"io"
	"os"
)

// A Message is a context for a Postfix-pipe helper process; it is
// what Postix passes us.
type Message struct {
	args  []string
	intro *string // a pointer so that it is nullable
	stdin *bufio.Reader
	env   env
}

// Get instantiates a new Postfix-pipe context for this process.
func Get() *Message {
	return &Message{
		args:  os.Args,
		intro: nil,
		stdin: nil,
		env:   getEnv(),
	}
}

// Intro returns the envelope line insertet by Postfix before the RFC
// 822-style message.
func (pm *Message) Intro() (string, error) {
	if pm.stdin == nil {
		pm.stdin = bufio.NewReader(os.Stdin)
	}
	if pm.intro == nil {
		intro, err := pm.stdin.ReadString('\n')
		if err != nil {
			return "", err
		}
		pm.intro = &intro
	}
	return *pm.intro, nil
}

// Reader returns an io.Reader for the RFC 822-style message.
func (pm *Message) Reader() (io.Reader, error) {
	_, err := pm.Intro()
	if err != nil {
		return nil, err
	}
	return pm.stdin, nil
}

// Args returns the command line arguments passed in by Postfix.
func (pm *Message) Args() []string {
	return pm.args
}
