package main

import (
	"flag"

	"github.com/andfasano/tiny-sushy-tools/internal/redfish"
)

func main() {
	var tiny_sushy_port string
	var tiny_libvirt_user string
	var tiny_libvirt_ip string
	var tiny_libvirt_key string

	flag.StringVar(&tiny_sushy_port, "port", "8000", "port to listen")
	flag.StringVar(&tiny_libvirt_user, "user", "root", "user for libvirt connection")
	flag.StringVar(&tiny_libvirt_ip, "ip", "127.0.0.1", "ip of libvirt server")
	flag.StringVar(&tiny_libvirt_key, "keyfile", "~/.ssh/id_rsa", "path to ssh key for libvirt server")

	flag.Parse()

	server := redfish.New()
	server.Start(tiny_sushy_port, tiny_libvirt_user, tiny_libvirt_ip, tiny_libvirt_key)
}
