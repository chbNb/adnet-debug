package main

import (
	"flag"
	"fmt"
	"testing"
)

var systemTest *bool

func init() {
	systemTest = flag.Bool("systemTest", false, "Set to true when running system tests")
}

// Test started when the test binary is started. Only calls main.
func TestSystem(t *testing.T) {
	_ = t
	if *systemTest {
		fmt.Println("starting coverage test...")
		main()
		fmt.Println("stop coverage test...")
	}
}
