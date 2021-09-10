package dummydriver

//Sample scalfolding for a driver

type DummyDriverFacade struct {
	dummy string
}

func (l *DummyDriverFacade) NewConnection(systemID string, tinyOobUser string, tinyOobIP string, tinyOobKey string) (map[string]string, error) {
	return map[string]string{}, nil
}
func (l *DummyDriverFacade) CloseConnection() error             { return nil }
func (l *DummyDriverFacade) BootFromDevice(device string) error { return nil }
func (l *DummyDriverFacade) MountISO(isoUrl string) error       { return nil }
func (l *DummyDriverFacade) EjectISO() error                    { return nil }
func (l *DummyDriverFacade) PowerOn() error                     { return nil }
func (l *DummyDriverFacade) PowerOff() error                    { return nil }
func (l *DummyDriverFacade) Reboot() error                      { return nil }
