// Copyright 2015 Guntas Grewal

package main

import (
	"periwinkle/listeners/twilio"
	"fmt"
)

func main() {
	testNumber, err := twilio.NewPhoneNum()

	if err!=nil{
		fmt.Printf("ERROR CHECK!\n")
	}

	fmt.Printf("%s\n", testNumber)
}