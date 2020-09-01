package redfish

import (
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/beevik/etree"
	libvirt "github.com/libvirt/libvirt-go"
)

type libvirtdriver struct {
}

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

func (l *libvirtdriver) getDomainByID(id string, wr io.Writer) {
	conn, err := libvirt.NewConnect("qemu+ssh://root@192.168.111.1/system?&keyfile=./id_rsa_virt_power&no_verify=1&no_tty=1")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		panic(err)
	}

	tpl, err := template.
		New("system").
		Funcs(template.FuncMap{
			"strContains": func(src string, tgt string) bool {
				res := strings.EqualFold(src, tgt)
				return res
			},
		}).
		Parse(templateSystem)
	if err != nil {
		panic(err)
	}

	for _, d := range doms {

		UUID, _ := d.GetUUIDString()

		if id == UUID {

			//d.SetUserPassword()

			input := struct {
				Identity         string
				Name             string
				UUID             string
				PowerState       string
				BootSourceTarget string
				BootSourceMode   string
				TotalCpus        uint
				TotalMemoryGB    uint64
				IndicatorLed     string
			}{
				UUID:             UUID,
				Identity:         UUID,
				PowerState:       "Off",
				IndicatorLed:     "Lit",
				BootSourceTarget: "None",
				BootSourceMode:   "None",
			}

			input.Name, _ = d.GetName()

			if active, _ := d.IsActive(); active {
				input.PowerState = "On"
			}

			////////// BootSourceTarget
			xmlDesc, _ := d.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
			doc := etree.NewDocument()
			doc.ReadFromString(xmlDesc)

			if boot := doc.FindElement(".//boot"); boot != nil {

				if device := boot.SelectAttr("device"); device != nil {
					input.BootSourceTarget = BOOT_DEVICE_MAP_REV[device.Value]
				}
			}

			minOrder := -1
			if devices := doc.FindElement(".//devices"); devices != nil && input.BootSourceTarget == "None" {

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
						input.BootSourceTarget = bootSourceTarget
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

					input.BootSourceTarget = "Pxe"
					minOrder = order
				}
			}

			////////// BootSourceMode
			loader := doc.FindElement(".//loader")
			if loader != nil {
				input.BootSourceMode = BOOT_MODE_MAP_REV[loader.SelectAttrValue("type", "None")]
			}

			///////
			maxCpus, _ := d.GetMaxVcpus()
			input.TotalCpus = maxCpus

			maxMem, _ := d.GetMaxMemory()
			input.TotalMemoryGB = maxMem / 1024 / 1024

			err = tpl.Execute(wr, input)
			if err != nil {
				panic(err)
			}
		}
	}
}

var templateSystem = `{
    "@odata.type": "#ComputerSystem.v1_1_0.ComputerSystem",
    "Id": "{{ .UUID }}",
    "Name": "{{ .Name }}",
    "UUID": "{{ .UUID }}",
    "Status": {
        "State": "Enabled",
        "Health": "OK",
        "HealthRollUp": "OK"
    },
	{{- if .PowerState }}
    "PowerState": "{{ .PowerState }}",
	{{- end }}
	"Boot": {
		{{- if .BootSourceTarget }}
        "BootSourceOverrideEnabled": "Continuous",
        "BootSourceOverrideTarget": "{{ .BootSourceTarget }}",
        "BootSourceOverrideTarget@Redfish.AllowableValues": [
			"Pxe",
			"Cd",
			"Hdd"
		{{- if .BootSourceMode }}
		],
		{{- if strContains .BootSourceMode "uefi" }}
		"BootSourceOverrideMode": "{{ .BootSourceMode }}",
		"UefiTargetBootSourceOverride": "/0x31/0x33/0x01/0x01"
		{{- else }}
		"BootSourceOverrideMode": {{ .BootSourceMode }},
		{{- end }}
		{{- else }}
		]
		{{- end}}
		{{- else }}
		"BootSourceOverrideEnabled": "Continuous"
		{{- end }}
	},
	"ProcessorSummary": {
        {{- if .TotalCpus }}
        "Count": {{ .TotalCpus }},
        {{- end }}
        "Status": {
            "State": "Enabled",
            "Health": "OK",
            "HealthRollUp": "OK"
        }
	},
	"MemorySummary": {
        {{- if .TotalMemoryGB }}
        "TotalSystemMemoryGiB": {{ .TotalMemoryGB }},
        {{- end }}
        "Status": {
            "State": "Enabled",
            "Health": "OK",
            "HealthRollUp": "OK"
        }
	},
	"Bios": {
        "@odata.id": "/redfish/v1/Systems/{{ .UUID }}/BIOS" 
	},
	"Processors": {
        "@odata.id": "/redfish/v1/Systems/{{ .UUID }}/Processors"
    },
    "Memory": {
        "@odata.id": "/redfish/v1/Systems/{{ .UUID }}/Memory"
    },
    "EthernetInterfaces": {
        "@odata.id": "/redfish/v1/Systems/{{ .UUID }}/EthernetInterfaces"
    },
    "SimpleStorage": {
        "@odata.id": "/redfish/v1/Systems/{{ .UUID }}/SimpleStorage"
    },
    "Storage": {
        "@odata.id": "/redfish/v1/Systems/{{ .UUID }}/Storage"
	},
	{{- if .IndicatorLed }}
    "IndicatorLED": "{{ .IndicatorLed }}",
	{{- end }}
	"Links": {
		"Chassis": [
			{
                "@odata.id": "/redfish/v1/Chassis/15693887-7984-9484-3272-842188918912"
            }
			],
        "ManagedBy": [
            {
                "@odata.id": "/redfish/v1/Managers/{{ .UUID }}"
            }
        ]		
	},
	"Actions": {
        "#ComputerSystem.Reset": {
            "target": "/redfish/v1/Systems/{{ .UUID }}/Actions/ComputerSystem.Reset",
            "ResetType@Redfish.AllowableValues": [
                "On",
                "ForceOff",
                "GracefulShutdown",
                "GracefulRestart",
                "ForceRestart",
                "Nmi",
                "ForceOn"
            ]
        }
	},
	"@odata.context": "/redfish/v1/$metadata#ComputerSystem.ComputerSystem",
	"@odata.id": "/redfish/v1/Systems/{{ .UUID }}",
	"@Redfish.Copyright": "Copyright 2014-2016 Distributed Management Task Force, Inc. (DMTF). For the full DMTF copyright policy, see http://www.dmtf.org/about/policies/copyright."
}`
