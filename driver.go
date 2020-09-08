package gm1356

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Fatih-Cetinkaya-Bose/hid"
)

// Driver is GM1356 controller
type Driver struct {
	device      *hid.Device
	eventBuffer chan Event
	importer    *importer
	ctx         context.Context
	cancel      context.CancelFunc
	log         func(...interface{})
}

// Open opens GM1356 device
func Open(eventBufferSize uint64) (*Driver, error) {
	devices := hid.Enumerate(vendorID, productID)
	if len(devices) == 0 {
		return nil, errors.New("GM1356 device not found")
	}
	device, err := devices[0].Open()
	if err != nil {
		return nil, err
	}
	eventBuffer := make(chan Event, eventBufferSize)
	logFunc := func(a ...interface{}) {
		fmt.Println(a...)
	}
	driver := &Driver{
		device:      device,
		eventBuffer: eventBuffer,
		importer:    newImporter(eventBuffer),
		log:         logFunc,
	}
	driver.ctx, driver.cancel = context.WithCancel(context.Background())
	go driver.handleInput()
	return driver, nil
}

// Close closes GM1356 device
func (d *Driver) Close() {
	d.cancel()
	d.device.Close()
}

// Configure requests configuration change
func (d *Driver) Configure(config Config) error {
	if d.IsImporting() {
		return errors.New("importing")
	}
	// write configration data
	return d.write(newConfigureRequest(config))
}

// Measure requests current sound level
func (d *Driver) Measure() error {
	if d.IsImporting() {
		return errors.New("importing")
	}
	// write measure request
	return d.write(newMeasureRequest())
}

// Import requests recorded data
func (d *Driver) Import() error {
	if err := d.importer.Start(); err != nil {
		return err
	}
	// write import request
	if err := d.write(newImportRequest()); err != nil {
		return err
	}
	for {
		if d.importer.IsImporting() {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	return nil
}

// IsImporting returns true when importing recorded data
func (d *Driver) IsImporting() bool {
	return d.importer.IsImporting()
}

// EventChannel retuns event channel
func (d *Driver) EventChannel() <-chan Event {
	return d.eventBuffer
}

// handleInput reads date and enqueue event
func (d *Driver) handleInput() {
	isImportingFirstResponse := true
	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			// read
			buf, err := d.read()
			if err != nil {
				if err == hid.ErrDeviceClosed {
					d.log("device closed")
					return
				}
				d.log("failed to read: ", err)
			}
			if d.IsImporting() {
				// process import data
				if err := d.importer.Write(buf); err != nil {
				if err == hid.ErrDeviceClosed {
					d.log("device closed")
					return
				}
				d.log("failed to write: ", err)
			}
				// request next import data
				if err := d.write(newNextImportDataRequest(isImportingFirstResponse)); err != nil {
					if err == hid.ErrDeviceClosed {
						d.log("device closed")
						return
					}
					d.log("failed to write: ", err)
				}
				isImportingFirstResponse = false
			} else {
				if event, ok := parseData(buf); ok {
				d.eventBuffer <- event
			}
		}
	}
}
}

// read reads bytes from GM1356 device
// this will block until any data is recieved
func (d *Driver) read() ([]byte, error) {
	buf := make([]byte, 8)
	n, err := d.device.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != 8 {
		return nil, errors.New("unexpected read size")
	}
	// d.log("read", time.Now(), buf)
	return buf, nil
}

// write writes bytes to GM1356 device
func (d *Driver) write(data []byte) error {
	// d.log("write", time.Now(), data)
	n, err := d.device.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return errors.New("unexpected write size")
	}
	return nil
}
