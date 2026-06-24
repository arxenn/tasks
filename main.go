package main

import (
	"github.com/arxenn/tasks/cmd"
)

// TODO LIST:
// - Complete windows shell integration

var version = "dev"

func main() {
	cmd.Execute(version)
}
