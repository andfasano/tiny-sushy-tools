package redfish

import (
	"io"
	"strings"
	"text/template"

	"github.com/andfasano/tiny-sushy-tools/internal/libvirtdriver"
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

func newSystem(systemID string, rf *Server) *system {
	s := &system{
		Identity:         systemID,
		Name:             "",
		UUID:             "",
		PowerState:       "Off",
		BootSourceTarget: "None",
		BootSourceMode:   "None",
		TotalCpus:        0,
		TotalMemoryGB:    0,
		IndicatorLed:     "Lit",
		Username:         "admin",
		Password:         "password",
	}

	lv := libvirtdriver.NewLibvirtDomain(systemID, rf.TinyOobUser, rf.TinyOobIP, rf.TinyOobKey)

	s.UUID = lv.GetUUID()
	s.Name = lv.GetName()
	s.PowerState = lv.GetPowerState()
	s.BootSourceTarget = lv.GetBootSourceTarget()
	s.BootSourceMode = lv.GetBootSourceMode()
	s.TotalCpus = lv.GetTotalCpus()
	s.TotalMemoryGB = lv.GetTotalMemory()

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
