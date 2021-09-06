package redfish

import (
	"io"
	"strings"
	"text/template"
)

type system struct {
	Identity         string
	Name             string
	UUID             string
	PowerState       string
	BootSourceTarget string
	BootSourceMode   string
	TotalCpus        uint
	TotalMemoryGB    uint64
	IndicatorLed     string

	Username string
	Password string
}

func newSystem(UUID string) *system {
	s := &system{
		UUID:             UUID,
		Identity:         UUID,
		PowerState:       "Off",
		IndicatorLed:     "Lit",
		BootSourceTarget: "None",
		BootSourceMode:   "None",
		Username:         "admin",
		Password:         "password",
	}

	lv := newLibvirtDomain(UUID)

	s.UUID = lv.getUUID()
	s.Name = lv.getName()
	s.PowerState = lv.getPowerState()
	s.BootSourceTarget = lv.getBootSourceTarget()
	s.BootSourceMode = lv.getBootSourceMode()
	s.TotalCpus = lv.getTotalCpus()
	s.TotalMemoryGB = lv.getTotalMemory()

	return s
}

//Send renders a System template
func (s *system) Send(wr io.Writer) {
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
	err = tpl.Execute(wr, s)
	if err != nil {
		panic(err)
	}
}
