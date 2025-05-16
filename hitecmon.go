package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

func main() {
	var userTag string
	var useInflux bool
	var configFile string
	flag.StringVar(&userTag, "u", "", "Data to be added to the user tag")
	flag.BoolVar(&useInflux, "i", false, "Enable writing to influxDB (must populate config file)")
	flag.StringVar(&configFile, "c", "", "Configuration file (Default is $HOME/.hitec.json)")
	flag.Parse()

	// Defined cause go doesn't like when you define inside ifs.
	var influxClient influxdb2.Client

	// Start up influx client
	if useInflux {
		config, err := parseConfig(configFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Start Influx client
		influxClient = influxdb2.NewClient(config.ServerAddress, config.InfluxToken)
		// Always close client at the end
		defer influxClient.Close()
	}

	serPort, err := start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer serPort.Close()

	serByte := make([]byte, 1)
	var lineData []byte
	var lastTime time.Time
	// Pre-populate to avoid throwing out the first result
	lastTime = time.Now()
	printBanner()
	for {
		_, err := serPort.Read(serByte)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Append the xored data
		lineData = append(lineData, serByte[0]^0x80)
		// Watch for the start bytes
		if serByte[0]^0x80 == 0xFB && len(lineData) > 1 && lineData[len(lineData)-2] == 0xFD {
			/* Found the start byte
			This means it's time to dump what we have and clear out for
			the next bit of data.
			However, sometimes it seems that the unit stores a reading in it's buffer
			after a cycle finishes and spits it out after the 'stop' button is pressed.
			To help eliminiate this, a simple timer.*/

			// Unit updates at 1 or 2s intervals, this should be safe.
			if time.Now().Sub(lastTime) > (time.Second * 5) {
				// Throw out this reading
				fmt.Println("Throwing out old reading")
				lineData = nil
				fmt.Printf("Current time %s\n", time.Now().String())
				fmt.Printf("last time: %s\n", lastTime.String())
				fmt.Printf("Difference: %s\n", time.Now().Sub(lastTime).String())
				lastTime = time.Now()
				continue
			}

			// Log the current time

			lineData = append([]byte{0xFD, 0xFB}, lineData[:len(lineData)-2]...)
			if len(lineData) == 76 {
				parsedData := processData(lineData)
				printData(parsedData)
				if useInflux {
					// Only send updates to influx if in state: running
					if parsedData.running {
						updateInflux(parsedData, userTag, influxClient)
					}
				}
				// Unknown areas
				// fmt.Printf("%X ", lineData[10:18])
				// fmt.Printf("%X ", lineData[20:24])
				// fmt.Printf("%X ", lineData[26:31])
				// fmt.Printf("%X ", lineData[38:42])
				// fmt.Printf("%X \n", lineData[72])
			}
			// Reset the line in prep for data
			lineData = nil
			// Reset time
			lastTime = time.Now()
		}
	}
}

func printBanner() {
	fmt.Printf("Running     Time    mAh   Pack Voltage  Current\n")
}
func printData(inData chargerFrame) {
	fmt.Printf("  %t     %dm   %dmAh     %.2fv       %dmA\n", inData.running, inData.totalTime, inData.mahCharged, inData.packVoltage, inData.chargingCurrent)
}

// getBitValue Gets a bitfield value
func getBitValue(bitHolder uint64, bitPos int) int {
	setBit := bitHolder & (1 << bitPos)
	setBit = setBit >> bitPos
	return int(setBit)
}
