package main

import (
	"fmt"
	"time"

	"github.com/paypal/gatt"
)

func main() {
	for {
		fmt.Println("Started scan")
		scan()
		fmt.Println("Finished scan")

		time.Sleep(5 * time.Second)
	}
}

func scan() (err error) {
	d, err := gatt.NewDevice()
	if err != nil {
		return
	}

	d.Handle(gatt.PeripheralDiscovered(onPeriphDiscovered))
	d.Init(onStateChanged)

	select {
	case <-time.After(time.Second * 10):
		d.StopScanning() //Stops scanning but only works once
		//d.Stop() <-- This is defined in device_linux.go, why can't I access it?
		//d.Stop() changes the state which would trigger the onStateChanged handler
		//below to stop scanning and also call d.hci.Close() to free the connection
	}

	return
}

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		d.Scan([]gatt.UUID{}, false)
	default:
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	fmt.Println("  Local Name        =", a.LocalName)
	fmt.Println("  TX Power Level    =", a.TxPowerLevel)
	fmt.Println("  Manufacturer Data =", a.ManufacturerData)
	fmt.Println("  Service Data      =", a.ServiceData)
}
