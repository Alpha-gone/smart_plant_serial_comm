package serial

import (
	"errors"
	"fmt"

	"go.bug.st/serial"
)

type serialCommPort struct {
	port serial.Port
	mode serial.Mode
	buf  []byte
}

func NewSerialComm(mode serial.Mode, bufSize int) *serialCommPort {
	return &serialCommPort{mode: mode, buf: make([]byte, bufSize)}
}

func (comm *serialCommPort) CheckPort(portName string) error {
	ports, err := serial.GetPortsList()
	if err != nil {
		return fmt.Errorf("failed get serial ports list: %w", err)
	}

	if len(ports) == 0 {
		return errors.New("serial ports not found")
	}

	for _, port := range ports {
		if port == portName {
			return nil
		}
	}

	return fmt.Errorf("port %s not found", portName)
}

func (comm *serialCommPort) OpenPort(portName string) error {
	if err := comm.CheckPort(portName); err != nil {
		return err
	}

	port, err := serial.Open(portName, &comm.mode)

	if err != nil {
		return err
	}

	comm.port = port
	return nil
}

func (comm *serialCommPort) isPortOpen() error {
	if comm.port == nil {
		return errors.New("serial port is not opened")
	}

	return nil
}

func (comm *serialCommPort) ChangePortConfiguration(mode *serial.Mode) error {
	if err := comm.isPortOpen(); err != nil {
		return err
	}

	if err := comm.port.SetMode(mode); err != nil {
		return fmt.Errorf("failed change port configuration: %w", err)
	}

	return nil
}

func (comm *serialCommPort) ClosePort() error {
	if err := comm.isPortOpen(); err != nil {
		return nil
	}

	if err := comm.port.Close(); err != nil {
		return fmt.Errorf("failed close serial port: %w", err)
	}

	comm.port = nil
	return nil
}

func (comm *serialCommPort) Write(message string) (int, error) {
	if err := comm.isPortOpen(); err != nil {
		return 0, err
	}

	writeByteSize, err := comm.port.Write([]byte(message))
	if err != nil {
		return 0, fmt.Errorf("writing error to serial port: %w", err)
	}

	return writeByteSize, nil
}

func (comm *serialCommPort) Read() (string, error) {
	if err := comm.isPortOpen(); err != nil {
		return "", err
	}

	readByteSize, err := comm.port.Read(comm.buf)
	if err != nil {
		return "", fmt.Errorf("reading error from serial port: %w", err)
	}

	if readByteSize == 0 {
		return "", errors.New("EOF")
	}

	return string(comm.buf[:readByteSize]), nil
}
