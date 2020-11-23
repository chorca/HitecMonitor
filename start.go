package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bugst/go-serial"
	"github.com/bugst/go-serial/enumerator"
)

func start() (serial.Port, error) {
	// Start serial port
	serOptions := serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	portList, err := enumerator.GetDetailedPortsList()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	batVID := "10c4"
	batPID := "ea60"
	var serPort serial.Port

	for i := 0; i < len(portList); i++ {
		if portList[i].VID == batVID && portList[i].PID == batPID {
			serPort, err = serial.Open(portList[i].Name, &serOptions)
			if err != nil {
				fmt.Println(err)
				return serPort, err
			}
			return serPort, nil
		}
	}
	return serPort, errors.New("Serial port not found")
}
