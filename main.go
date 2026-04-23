// Package main is the entry point for the powertracker application.
// powertracker monitors and tracks power consumption data from smart meters.
//
// This is a personal fork of github.com/poolski/powertracker.
// Primary use case: monitoring home solar + grid usage on a Raspberry Pi.
package main

import "github.com/poolski/powertracker/cmd"

func main() {
	cmd.Execute()
}
