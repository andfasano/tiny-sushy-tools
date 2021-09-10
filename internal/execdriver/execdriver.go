package execdriver

import "strconv"

type ExecDriverFacade struct {
	dummy string
}

func (exec *ExecDriverFacade) NewConnection(systemID string, tinyOobUser string, tinyOobIP string, tinyOobKey string) (map[string]string, error) {
	systemMap := map[string]string{
		"Identity":         "id-6e162f250552",
		"Name":             "dummy",
		"UUID":             "4ef0b933-5439-42bd-bc97-6e162f250552",
		"PowerState":       "On",
		"BootSourceTarget": "cdrom",
		"BootSourceMode":   exec.getBootSourceMode(),
		"TotalCpus":        strconv.FormatUint(uint64(exec.getTotalCpus()), 10),   //uint
		"TotalMemoryGB":    strconv.FormatUint(uint64(exec.getTotalMemory()), 10), //uint64
		"IndicatorLed":     "Lit",
		"Username":         "admin",
		"Password":         "password",
	}

	return systemMap, nil
}

//Not implemented
func (exec *ExecDriverFacade) CloseConnection() error             { return nil }
func (exec *ExecDriverFacade) BootFromDevice(device string) error { return nil }
func (exec *ExecDriverFacade) MountISO(isoUrl string) error       { return nil }
func (exec *ExecDriverFacade) EjectISO() error                    { return nil }
func (exec *ExecDriverFacade) PowerOn() error                     { return nil }
func (exec *ExecDriverFacade) PowerOff() error                    { return nil }
func (exec *ExecDriverFacade) Reboot() error                      { return nil }

func (exec *ExecDriverFacade) getBootSourceMode() string { return "None" }
func (exec *ExecDriverFacade) getTotalCpus() uint        { return 0 }
func (exec *ExecDriverFacade) getTotalMemory() uint64    { return 0 }
