package main

import (
	"github.com/andfasano/tiny-sushy-tools/internal/redfish"
)

func main() {
	server := redfish.New()
	server.Start("8000")
}
