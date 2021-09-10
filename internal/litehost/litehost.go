package litehost

import (
	"log"

	"github.com/andfasano/tiny-sushy-tools/internal/dummydriver"
	"github.com/andfasano/tiny-sushy-tools/internal/execdriver"
	"github.com/andfasano/tiny-sushy-tools/internal/libvirtdriver"
)

type LiteHost interface {
	BootFromDevice(device string) error
	MountISO(isoUrl string) error
	EjectISO() error
	PowerOn() error
	PowerOff() error
	Reboot() error
	NewConnection(systemID string, tinyOobUser string, tinyOobIP string, tinyOobKey string) (map[string]string, error)
	CloseConnection() error
}

type DriverLoader struct {
	liteHostAttributes map[string]string
}

//make sure host is off and no ISO is mounted to it
func (dl *DriverLoader) StartLiteHost(driver string, systemID string, tinyOobUser string, tinyOobIP string, tinyOobKey string) error {
	var lh LiteHost

	if driver == "libvirt" {
		lh = &libvirtdriver.LibvirtDomainFacade{}
	} else if driver == "execdriver" {
		lh = &execdriver.ExecDriverFacade{}
	} else if driver == "dummydriver" {
		lh = &dummydriver.DummyDriverFacade{}
	} else {
		log.Println("Invalid driver ", driver)
		log.Fatal("Invalid driver")
	}

	dl.liteHostAttributes, _ = lh.NewConnection(systemID, tinyOobUser, tinyOobIP, tinyOobKey)

	return nil
}

func (dl *DriverLoader) LiteHostAttribute(key string, defaultValue string) string {
	if val, ok := dl.liteHostAttributes[key]; ok {
		return val
	}
	return defaultValue
}
