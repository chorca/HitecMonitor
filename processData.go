package main

func processData(inData []byte) chargerFrame {
	// Process the data into a struct
	var outData = chargerFrame{}
	// Configs
	outData.tempCutoffOn = getBitValue(uint64(inData[2]), 0) != 0
	outData.timeCutoffOn = getBitValue(uint64(inData[2]), 1) != 0
	outData.capCutoffOn = getBitValue(uint64(inData[2]), 2) != 0
	outData.buzzerOn = getBitValue(uint64(inData[2]), 4) != 0
	outData.keyBeepOn = getBitValue(uint64(inData[2]), 7) != 0
	// Battery stuff
	outData.nicdSens = int(inData[3])
	outData.nimhSens = int(inData[4])
	outData.tempLimit = int(inData[5])
	outData.wasteTime = int(inData[6])
	outData.dispBacklight = int(inData[7])
	outData.lowPowerCutoff = float32(inData[8]) / 10
	outData.chargeDischarge = int(inData[9]) != 0
	outData.setChargeCurrent = int(inData[18]) * 100
	outData.numCells = int(inData[19])
	outData.running = int(inData[25]) != 0
	outData.screenMode = int(inData[24])
	outData.timeLimit = int(inData[31]) * 10
	outData.capLimit = int(inData[32])*100 + int(inData[33])
	outData.chargingCurrent = int(inData[34])*1000 + int(inData[35])*10
	outData.packVoltage = float32(inData[36]) + (float32(inData[37]) / 100)
	outData.inputVolts = float32(inData[42]) + (float32(inData[43]) / 100)
	outData.mahCharged = int(inData[44])*100 + int(inData[45])
	outData.cell1Volts = float32(inData[46]) + (float32(inData[47]) / 100)
	outData.cell2Volts = float32(inData[48]) + (float32(inData[49]) / 100)
	outData.cell3Volts = float32(inData[50]) + (float32(inData[51]) / 100)
	outData.cell4Volts = float32(inData[52]) + (float32(inData[53]) / 100)
	outData.cell5Volts = float32(inData[54]) + (float32(inData[55]) / 100)
	outData.cell6Volts = float32(inData[56]) + (float32(inData[57]) / 100)
	outData.cell7Volts = float32(inData[58]) + (float32(inData[59]) / 100)
	outData.cell8Volts = float32(inData[60]) + (float32(inData[61]) / 100)
	outData.cell9Volts = float32(inData[62]) + (float32(inData[63]) / 100)
	outData.cell10Volts = float32(inData[64]) + (float32(inData[65]) / 100)
	outData.cell11Volts = float32(inData[66]) + (float32(inData[67]) / 100)
	outData.cell12Volts = float32(inData[68]) + (float32(inData[69]) / 100)
	outData.totalTime = int(inData[70])*100 + int(inData[71])
	outData.trickleMa = int(inData[73])

	return outData
}
