package main

import (
	"flag"

	"github.com/andfasano/tiny-sushy-tools/internal/redfish"
)

func main() {
	server := redfish.New()

	flag.StringVar(&server.TinySushyPort, "port", "8000", "port to listen")
	flag.StringVar(&server.TinySushyDriver, "driver", "libvirt", "backend driver to use (libvirt, execdriver, dummydriver)")
	flag.StringVar(&server.TinyOobUser, "user", "root", "user for libvirt connection")
	flag.StringVar(&server.TinyOobIP, "ip", "127.0.0.1", "ip of libvirt server")
	flag.StringVar(&server.TinyOobKey, "keyfile", "~/.ssh/id_rsa", "path to ssh key for libvirt server")

	flag.Parse()

	server.Start()
}
