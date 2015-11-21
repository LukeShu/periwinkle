// Copyright 2015 Luke Shumaker

package postfixpipe

import (
	"os"
	"bufio"
	"io"
)

type Message struct {
	args  []string
	intro *string // a pointer so that it is nullable
	stdin *bufio.Reader
	env   env
}

func Get() *Message {
	return &Message{
		args: os.Args,
		intro: nil,
		stdin: nil,
		env: getEnv(),
	}
}

func (pm *Message) Intro() (string, error) {
	if (pm.stdin == nil) {
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

func (pm *Message) Reader() (io.Reader, error) {
	_, err := pm.Intro()
	if err != nil {
		return nil, err
	}
	return pm.stdin, nil
}

func (pm *Message) Args() []string {
	return pm.args
}
