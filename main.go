package main

import (
	"bufio"
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type ScanReport struct {
	Timestamp    int64             `json:"timestamp,omitempty"`
	SSID         string            `json:"ssid,omitempty"`
	BSSID        string            `json:"bssid"`
	Capabilities string            `json:"capabilities,omitempty"`
	Channel      string            `json:"channel,omitempty"`
	FreqKHz      int               `json:"freqkhz,omitempty"`
	Signal       int               `json:"signal,omitempty"`
	Lat          float64           `json:"lat,omitempty"`
	Lon          float64           `json:"lon,omitempty"`
	Alt          float64           `json:"alt,omitempty"`
	Spd          float64           `json:"spd,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
}

type ScanPayload struct {
	SourceName string       `json:"source_name"`
	SourceUUID string       `json:"source_uuid"`
	Reports    []ScanReport `json:"reports"`
}

const ESP32_BAUD_RATE = 115200

var serialChannel = make(chan string)

func main() {
	portInfo, err := enumerator.GetDetailedPortsList()

	if err != nil {
		log.Fatal(err)
	}

	ports := make([]*enumerator.PortDetails, 0)

	for _, p := range portInfo {
		if !p.IsUSB {
			continue
		}
		ports = append(ports, p)
	}

	for _, p := range ports {
		go handlePort(p)
	}

	for {
		line := <-serialChannel
		fmt.Println(line)
	}
}

func handlePort(port *enumerator.PortDetails) {
	mode := &serial.Mode{
		BaudRate: ESP32_BAUD_RATE,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	p, err := serial.Open(port.Name, mode)
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(p)
	buf.Size()
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			fmt.Println(err)
			break
		}
		serialChannel <- string(line)
	}
	fmt.Println("Port handle closed. Trying to open again in 5 seconds.")
	time.Sleep(time.Second * 5)
	go handlePort(port)
}
