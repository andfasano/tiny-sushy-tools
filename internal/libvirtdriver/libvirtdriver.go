package libvirtdriver

import (
	"strconv"

	"github.com/beevik/etree"
	uuid "github.com/google/uuid"
	libvirt "github.com/libvirt/libvirt-go"
)

const (
	DEVICE_TYPE_PXE    = "Pxe"
	DEVICE_TYPE_HDD    = "Hdd"
	DEVICE_TYPE_CD     = "Cd"
	DEVICE_TYPE_FLOPPY = "Floppy"
)

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

type LibvirtDomainFacade struct {
	conn   *libvirt.Connect
	domain *libvirt.Domain
}

func uri(user string, ip string, keyfilePath string) (string, error) {
	// default u: "qemu+ssh://root@192.168.111.1/system?&keyfile=~/.ssh/id_rsa_virt_power&no_verify=1&no_tty=1"
	var lvUrl string

	if user == "" {
		user = "root"
	}

	if ip == "" {
		ip = "192.168.111.1"
	}

	if keyfilePath == "" {
		keyfilePath = "~/.ssh/id_rsa_virt_power"
	}

	lvUrl = "qemu+ssh://" + user + "@" + ip + "/system?&keyfile=" + keyfilePath + "&no_verify=1&no_tty=1"

	return lvUrl, nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

//Open a new libvirt connection
func (l *LibvirtDomainFacade) NewConnection(systemID string, tinyOobUser string, tinyOobIP string, tinyOobKey string) (map[string]string, error) {
	var libvirtUri string
	var err error

	libvirtUri, _ = uri(tinyOobUser, tinyOobIP, tinyOobKey)

	l.conn, err = libvirt.NewConnect(libvirtUri)
	if err != nil {
		panic(err)
	}

	if isValidUUID(systemID) {
		l.domain, err = l.conn.LookupDomainByUUIDString(systemID)
		if err != nil {
			panic(err)
		}
	} else {
		l.domain, err = l.conn.LookupDomainByName(systemID)
		if err != nil {
			panic(err)
		}
	}

	systemMap := map[string]string{
		"Identity":         l.getUUID(),
		"Name":             l.getName(),
		"UUID":             l.getUUID(),
		"PowerState":       l.getPowerState(),
		"BootSourceTarget": l.getBootSourceTarget(),
		"BootSourceMode":   l.getBootSourceMode(),
		"TotalCpus":        strconv.FormatUint(uint64(l.getTotalCpus()), 10),   //uint
		"TotalMemoryGB":    strconv.FormatUint(uint64(l.getTotalMemory()), 10), //uint64
		"IndicatorLed":     "Lit",
		"Username":         "admin",
		"Password":         "password",
	}

	return systemMap, nil
}

//Closes libvirt connection
func (l *LibvirtDomainFacade) CloseConnection() error {
	l.conn.Close()
	return nil
}

func (l *LibvirtDomainFacade) BootFromDevice(device string) error { return nil }
func (l *LibvirtDomainFacade) MountISO(isoUrl string) error       { return nil }
func (l *LibvirtDomainFacade) EjectISO() error {
	//l.domain.AttachDevice()
	return nil
}

func (l *LibvirtDomainFacade) PowerOn() error {
	//DomainLifecycle type
	// VIR_DOMAIN_LIFECYCLE_POWEROFF	=	0 (0x0)
	// VIR_DOMAIN_LIFECYCLE_REBOOT	=	1 (0x1)
	// VIR_DOMAIN_LIFECYCLE_CRASH	=	2 (0x2)
	// VIR_DOMAIN_LIFECYCLE_LAST	=	3 (0x3)
	//DomainLifecycleAction
	// VIR_DOMAIN_LIFECYCLE_ACTION_DESTROY	=	0 (0x0)
	// VIR_DOMAIN_LIFECYCLE_ACTION_RESTART	=	1 (0x1)
	// VIR_DOMAIN_LIFECYCLE_ACTION_RESTART_RENAME	=	2 (0x2)
	// VIR_DOMAIN_LIFECYCLE_ACTION_PRESERVE	=	3 (0x3)
	// VIR_DOMAIN_LIFECYCLE_ACTION_COREDUMP_DESTROY	=	4 (0x4)
	// VIR_DOMAIN_LIFECYCLE_ACTION_COREDUMP_RESTART	=	5 (0x5)
	// VIR_DOMAIN_LIFECYCLE_ACTION_LAST	=	6 (0x6)
	//Flags - DomainModificationImpact
	// VIR_DOMAIN_AFFECT_CURRENT	=	0 (0x0)
	// Affect current domain state.
	// VIR_DOMAIN_AFFECT_LIVE	=	1 (0x1; 1 << 0)
	// Affect running domain state.
	// VIR_DOMAIN_AFFECT_CONFIG	=	2 (0x2; 1 << 1)
	// Affect persistent domain state. 1 << 2 is reserved for virTypedParameterFlags

	if l.domain.SetLifecycleAction(1, 1, 1) != nil {
		return nil
	}
	// fixme
	return nil
}

func (l *LibvirtDomainFacade) PowerOff() error {
	if l.domain.Shutdown() != nil {
		return nil
	}
	//fixme
	return nil
}
func (l *LibvirtDomainFacade) Reboot() error { return nil }

func (l *LibvirtDomainFacade) getName() string {
	name, err := l.domain.GetName()
	if err != nil {
		panic(err)
	}
	return name
}

func (l *LibvirtDomainFacade) getUUID() string {
	uuidBytes, err := l.domain.GetUUID()
	if err != nil {
		panic(err)
	}
	uuidstr, _ := uuid.FromBytes(uuidBytes)
	return uuidstr.String()
}

func (l *LibvirtDomainFacade) getPowerState() string {
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

func (l *LibvirtDomainFacade) getBootSourceTarget() string {

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

func (l *LibvirtDomainFacade) getBootSourceMode() string {
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

func (l *LibvirtDomainFacade) getTotalCpus() uint {
	maxCpus, err := l.domain.GetMaxVcpus()
	if err != nil {
		panic(err)
	}
	return maxCpus
}

func (l *LibvirtDomainFacade) getTotalMemory() uint64 {
	maxMem, err := l.domain.GetMaxMemory()
	if err != nil {
		panic(err)
	}
	return maxMem / 1024 / 1024
}
