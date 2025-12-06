package mister

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"sni/cmd/sni/config"
	"sni/devices"
	"sni/protos/sni"
	"time"
)

const driverName = "mister"
const misterPort = "23074"

var driver *Driver

type Driver struct {
	container devices.DeviceContainer
	client    *Client
}

func NewDriver(address *net.TCPAddr) *Driver {
	var d *Driver
	d = &Driver{
		container: devices.NewDeviceDriverContainer(d.openDevice),
		client: &Client{
			addr: address,
		},
	}
	return d
}

func (d *Driver) openDevice(uri *url.URL) (q devices.Device, err error) {
	// create a new device with its own connection:
	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", uri.Host)
	if err != nil {
		return
	}

	var c *Client
	c = NewClient(addr, addr.String(), time.Second*5)
	err = c.Connect()
	if err != nil {
		return
	}

	q = c
	return
}
func (d *Driver) DisplayName() string {
	return "MisterFPGA SNES Core"
}

func (d *Driver) DisplayDescription() string {
	return "Connect to a remote Mister FPGA SNES Core over the network"
}

func (d *Driver) Kind() string { return driverName }

var driverCapabilities = []sni.DeviceCapability{
	sni.DeviceCapability_ReadMemory,
	sni.DeviceCapability_WriteMemory,
}

func (d *Driver) HasCapabilities(capabilities ...sni.DeviceCapability) (bool, error) {
	return devices.CheckCapabilities(capabilities, driverCapabilities)
}

func (d *Driver) DeviceKey(uri *url.URL) string {
	return uri.Host
}

func (d *Driver) Device(uri *url.URL) devices.AutoCloseableDevice {
	return devices.NewAutoCloseableDevice(d.container, uri, d.DeviceKey(uri))
}

func (d *Driver) Detect() (devs []devices.DeviceDescriptor, err error) {
	// TODO
	return nil, nil
}

func (d *Driver) DisconnectAll() {
	// TODO
}

func DriverInit() {
	if config.Config.GetBool("mister_disable") {
		log.Printf("Mister: disabling MisterFPGA driver\n")
		return
	}

	hostStr := config.Config.GetString("mister_host")
	portStr := config.Config.GetString("mister_port")
	addrStr := fmt.Sprintf("%s:%s", hostStr, portStr)
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		log.Printf("Mister: resolve('%s'): %v\n", addrStr, err)
	}

	driver = NewDriver(addr)
	devices.Register(driverName, driver)
}
