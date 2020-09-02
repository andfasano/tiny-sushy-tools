package redfish

import (
	"strconv"

	"github.com/beevik/etree"
	libvirt "github.com/libvirt/libvirt-go"
)

const DEVICE_TYPE_PXE = "Pxe"
const DEVICE_TYPE_HDD = "Hdd"
const DEVICE_TYPE_CD = "Cd"
const DEVICE_TYPE_FLOPPY = "Floppy"

var BOOT_DEVICE_MAP map[string]string = map[string]string{
	DEVICE_TYPE_PXE:    "network",
	DEVICE_TYPE_HDD:    "hd",
	DEVICE_TYPE_CD:     "cdrom",
	DEVICE_TYPE_FLOPPY: "floppy",
}

var BOOT_DEVICE_MAP_REV map[string]string = reverseMap(BOOT_DEVICE_MAP)

var DISK_DEVICE_MAP map[string]string = map[string]string{
	DEVICE_TYPE_HDD:    "disk",
	DEVICE_TYPE_CD:     "cdrom",
	DEVICE_TYPE_FLOPPY: "floppy",
}

var DISK_DEVICE_MAP_REV map[string]string = reverseMap(DISK_DEVICE_MAP)

var BOOT_MODE_MAP map[string]string = map[string]string{
	"Legacy": "rom",
	"UEFI":   "pflash",
}

var BOOT_MODE_MAP_REV map[string]string = reverseMap(BOOT_MODE_MAP)

func reverseMap(src map[string]string) map[string]string {
	dest := make(map[string]string)
	for k, v := range src {
		dest[v] = k
	}
	return dest
}

type libvirtDomainFacade struct {
	conn   *libvirt.Connect
	domain *libvirt.Domain
}

func newLibvirtDomain(UUID string) *libvirtDomainFacade {
	conn, err := libvirt.NewConnect("qemu+ssh://root@192.168.111.1/system?&keyfile=./id_rsa_virt_power&no_verify=1&no_tty=1")
	if err != nil {
		panic(err)
	}

	dom, err := conn.LookupDomainByUUIDString(UUID)
	if err != nil {
		panic(err)
	}

	facade := &libvirtDomainFacade{
		conn:   conn,
		domain: dom,
	}

	return facade
}

func (l *libvirtDomainFacade) close() {
	l.conn.Close()
}

func (l *libvirtDomainFacade) getName() string {
	name, err := l.domain.GetName()
	if err != nil {
		panic(err)
	}
	return name
}

func (l *libvirtDomainFacade) getPowerState() string {
	state := "Off"
	active, err := l.domain.IsActive()
	if err != nil {
		panic(err)
	}

	if active {
		state = "On"
	}

	return state
}

func (l *libvirtDomainFacade) getBootSourceTarget() string {

	xmlDesc, err := l.domain.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
	if err != nil {
		panic(err)
	}
	doc := etree.NewDocument()
	doc.ReadFromString(xmlDesc)

	if boot := doc.FindElement(".//boot"); boot != nil {

		if device := boot.SelectAttr("device"); device != nil {
			return BOOT_DEVICE_MAP_REV[device.Value]
		}
	}

	minOrder := -1
	target := ""
	if devices := doc.FindElement(".//devices"); devices != nil {

		for _, disk := range devices.SelectElements("disk") {
			boot := disk.SelectElement("boot")
			if boot == nil {
				continue
			}
			order, _ := strconv.Atoi(boot.SelectAttrValue("order", "1000"))
			if minOrder != -1 && order >= minOrder {
				continue
			}
			device := disk.SelectAttrValue("device", "")
			if device == "" {
				continue
			}

			if bootSourceTarget, ok := DISK_DEVICE_MAP_REV[device]; ok {
				target = bootSourceTarget
				minOrder = order
			}
		}

		for _, iface := range devices.SelectElements("interface") {
			boot := iface.SelectElement("boot")
			if boot == nil {
				continue
			}
			order, _ := strconv.Atoi(boot.SelectAttrValue("order", "1000"))
			if minOrder != -1 && order >= minOrder {
				continue
			}

			target = "Pxe"
			minOrder = order
		}

	}

	return target
}

func (l *libvirtDomainFacade) getBootSourceMode() string {
	mode := "None"

	xmlDesc, err := l.domain.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
	if err != nil {
		panic(err)
	}
	doc := etree.NewDocument()
	doc.ReadFromString(xmlDesc)

	loader := doc.FindElement(".//loader")
	if loader != nil {
		mode = BOOT_MODE_MAP_REV[loader.SelectAttrValue("type", "None")]
	}

	return mode
}

func (l *libvirtDomainFacade) getTotalCpus() uint {
	maxCpus, err := l.domain.GetMaxVcpus()
	if err != nil {
		panic(err)
	}
	return maxCpus
}

func (l *libvirtDomainFacade) getTotalMemory() uint64 {
	maxMem, err := l.domain.GetMaxMemory()
	if err != nil {
		panic(err)
	}
	return maxMem / 1024 / 1024
}
