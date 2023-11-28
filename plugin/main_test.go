package main

import (
	"os"
	"testing"
)

func TestInstalling(t *testing.T) {
	testMain(t, []string{"cmd", "{\"key\":\"install\"}"}, "Expected output for arg1")
}

// Add more test functions as needed...

// testMain is a helper function to test the main function with given arguments
func testMain(t *testing.T, args []string, expected string) {

	// Set os.Args to simulate command-line arguments
	os.Args = args

	// Call main function
	main()
}
