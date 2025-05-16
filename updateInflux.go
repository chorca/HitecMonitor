package main

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

func updateInflux(inData chargerFrame, tag string, client influxdb2.Client) {
	// Write the datapoint to an Influx DB
	timeStamp := time.Now()
	// Dumb but generate a map we can use cause you can't address struct fields by variable/name
	cellMap := make(map[int]float32)
	cellMap[1] = inData.cell1Volts
	cellMap[2] = inData.cell2Volts
	cellMap[3] = inData.cell3Volts
	cellMap[4] = inData.cell4Volts
	cellMap[5] = inData.cell5Volts
	cellMap[6] = inData.cell6Volts
	cellMap[7] = inData.cell7Volts
	cellMap[8] = inData.cell8Volts
	cellMap[9] = inData.cell9Volts
	cellMap[10] = inData.cell10Volts
	cellMap[11] = inData.cell11Volts
	cellMap[12] = inData.cell12Volts

	// Measurements are made separately and need to be written separately
	chgCurrPoint := fmt.Sprintf("chgcurrent,user=%s,state=%t,setcurrent=%d current=%d %d", tag, inData.chargeDischarge, inData.setChargeCurrent, inData.chargingCurrent, timeStamp.UnixNano())
	chgVoltPoint := fmt.Sprintf("chgvolt,user=%s,state=%t voltage=%.2f %d", tag, inData.chargeDischarge, inData.packVoltage, timeStamp.UnixNano())
	chgMahPoint := fmt.Sprintf("chgmah,user=%s,state=%t,caplimit=%d,caplimiten=%t capacity=%d %d", tag, inData.chargeDischarge, inData.capLimit, inData.capCutoffOn, inData.mahCharged, timeStamp.UnixNano())
	chgTimePoint := fmt.Sprintf("chgtime,user=%s,state=%t,timelimit=%d minutes=%d %d", tag, inData.chargeDischarge, inData.timeLimit, inData.totalTime, timeStamp.UnixNano())
	writeAPI := client.WriteAPI(config.InfluxOrg, config.InfluxBucket)
	writeAPI.WriteRecord(chgCurrPoint)
	writeAPI.WriteRecord(chgVoltPoint)
	writeAPI.WriteRecord(chgMahPoint)
	writeAPI.WriteRecord(chgTimePoint)
	// fmt.Println(chgCurrPoint)
	// fmt.Println(chgVoltPoint)
	// fmt.Println(chgMahPoint)
	// fmt.Println(chgTimePoint)
	if inData.numCells > 1 {
		for i := 1; i <= inData.numCells; i++ {
			chgCellPoint := fmt.Sprintf("chgcell%d,user=%s,state=%t voltage=%.2f %d", i, tag, inData.chargeDischarge, cellMap[i], timeStamp.UnixNano())
			writeAPI.WriteRecord(chgCellPoint)
			// fmt.Println(chgCellPoint)
		}
	}
	writeAPI.Flush()
}
