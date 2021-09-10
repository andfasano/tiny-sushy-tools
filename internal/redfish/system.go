package redfish

import (
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/andfasano/tiny-sushy-tools/internal/litehost"
)

type system struct {
	Identity         string
	Name             string
	UUID             string
	PowerState       string
	BootSourceTarget string
	BootSourceMode   string
	TotalCpus        uint64
	TotalMemoryGB    uint64
	IndicatorLed     string

	Username string
	Password string
}

func (s *system) newSystem(systemID string, rf *Server) error {

	dl := litehost.DriverLoader{}
	dl.StartLiteHost(rf.TinySushyDriver, systemID, rf.TinyOobUser, rf.TinyOobIP, rf.TinyOobKey)

	s.Identity = systemID
	s.Name = dl.LiteHostAttribute("Name", "unknown.litehost.local")
	s.UUID = dl.LiteHostAttribute("UUID", "010294aa-a89c-40ea-982a-a1c64d6f4509")
	s.PowerState = dl.LiteHostAttribute("PowerState", "Off")
	s.BootSourceTarget = dl.LiteHostAttribute("BootSourceTarget", "None")
	s.BootSourceMode = dl.LiteHostAttribute("BootSourceMode", "None")
	s.TotalCpus, _ = strconv.ParseUint(dl.LiteHostAttribute("TotalCpus", "0"), 10, 0)
	s.TotalMemoryGB, _ = strconv.ParseUint(dl.LiteHostAttribute("TotalMemoryGB", "0"), 10, 0)
	s.IndicatorLed = dl.LiteHostAttribute("IndicatorLed", "Lit")
	s.Username = dl.LiteHostAttribute("Username", "admin")
	s.Password = dl.LiteHostAttribute("Password", "password")

	return nil
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
