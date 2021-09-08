package litehost

// NOTE: INITIAL SALFOLDING FOR INTERFACE

type liteHost interface {
	bootFromDevice(device string) error
	mountISO(isoUrl string) error
	ejectISO() error
	powerOn() error
	powerOff() error
	reboot() error
}

func StartLiteHost(lh liteHost) error {
	//make sure host is off and no ISO is mounted to it
	lh.powerOff()
	lh.ejectISO()

	return nil
}
