package redfish

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
