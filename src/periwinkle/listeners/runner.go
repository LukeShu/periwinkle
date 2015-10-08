// Copyright 2015 Luke Shumaker

package main

import (
	//"periwinkle/listeners/twilio"
	"fmt"
	"periwinkle/listeners/maildir"
	"periwinkle/listeners/web"
	"sync"
)

func main() {
	fmt.Println("listener starting")
	var wg sync.WaitGroup
	wg.Add(3)
	//go func() { twilio.Main(); wg.Done() }()
	go func() { maildir.Main(); wg.Done() }()
	go func() { web.Main(); wg.Done() }()
	wg.Wait()
}
