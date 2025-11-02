package main

import (
	"fmt"
	rs485 "smart_plant_serial_comm/serial"
	"time"

	"go.bug.st/serial"
)

func main() {
	commPort := rs485.NewSerialComm(
		serial.Mode{BaudRate: 9600},
		128)
	done := make(chan bool)

	err := commPort.OpenPort("/dev/ttyAMA5")
	if err != nil {
		fmt.Println("failed serial port open:", err)
	}
	defer commPort.ClosePort()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				recv, err := commPort.Read()
				if err != nil {
					continue
				}

				fmt.Println("received: ", recv)
			}
		}
	}()

	for ch := 'a'; ch <= 'z'; ch++ {
		if _, err := commPort.Write(string(ch)); err != nil {
			fmt.Println("failed serial message write:", err)
		} else {
			fmt.Println("send message: ", ch)
		}

		time.Sleep(2 * time.Second)
	}

	time.Sleep(10 * time.Second)
	close(done)
	time.Sleep(100 * time.Millisecond)
}
