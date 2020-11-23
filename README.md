# HitecMonitor

## About
Hitec is a company who makes various products for RC/hobby use. One of their lines includes battery chargers for various chemistries, some of them equipped with USB ports to allow for monitoring/control of the charging process. Some of these are older and lack modern software to gather the data, the one I have in particular is an X2 Ultima (P/N 44164), which has two commonly-used ATMega-controlled chargers in one package, along with a USB port for monitoring.
I wanted to be able to monitor this data and graph it over time to see how my packs were doing, and so reversed the protocol used by the device to talk to the PC software, and wrote my own in Go which will dump data to the screen, a file, or database (currently only InfluxDB is supported)


## Hardware
 The Hitec X2 Ultima battery charger is 2 commonly-available chargers combined into one. The basic hardware within has been included in hundreds of chargers, and is fairly common. However, Hitec has included some custom logic to switch the display and control keys between the two MCUs onboard, and also included a SILabs CP210x USB-Serial interface to pull data from the charger for logging on the PC. This model only does unidirectional comms (to my knowledge) and can only return data.
Data is returned at 9600 baud, 8 N 1.
Data return is paused when the charger completes a cycle and is waiting on user intervention; there's no way to tell a cycle has completed other than the absence of reports for a period of time.

## Protocol
The information here is based upon what I've observed, and I'm no expert at reverse-engineering, so please let me know if I've missed something.
Decoded.txt has some examples and a breakdown of what I'm seeing.
* The protocol is interesting in that it seems to contain not only battery charge statistics, but also (most of) the configuration of the unit. This is a bit heavy at 9600 baud, but it is plenty for monitoring long-term.
* The serial portion can also detect if the charger is actively charging; when idle it will transmit once every 2 seconds or so, but when charging it will increase reporting speed to once every second.
* The entire message is XOR'd with `0x80`
* Check bits are something I'm not familiar with yet, so I'm not able to easily reverse what's being done here.

|Byte|Purpose|Size|Units|Description|
|----|-------|----|-----|-------|
|0-1|Start Bytes (`FD FB`)|2||These seem to be at the beginning of every transmission; I'm assuming since this comes after the check bytes|
|2|Configuration|1|Bitfield|This byte is a bitfield containing a few settings regarding safety and general preferences|
|3|NiCD Sensitvity|1|mv/Cell|Sensitivity during charging NiCD cells|
|4|NiMH Sensitvity|1|mv/Cell|Sensitivity during charging NiMH cells|
|5|Temp Limit|1|deg F|Safety temperature limit|
|6|Waste Time|1|mins|Time between a discharge and recharge cycle (NiCD/NiMH/Pb)|
|7|Display Backlight|1|*5%|Backlight value 0-100%, in count of 5% increments|
|8|Low Voltage Cutoff|1|Volts*10|Input voltage (from power supply) at which the unit will shut off for safety|
|9|Charge/Discharge|1|Boolean|This is 1 when current is flowing into the battery, and 0 when current is being drained|
|10-17|Unknown|8||Have not seen these change during my time with the charger|
|18|Set Charge Current|1|*100mA|Current in mA the charger is currently set to|
|19|Number of cells|1||The number of cells selected on the charger|
|20-23|Unknown|4||Have not seen these change during my time with the charger|
|24|Screen Selection|1||Currently selected screen. Each screen on the main page seems to have it's own number.|
|25|Running|1|Boolean|1 if the charger is currently running a program, 0 if it is idle.|
|26-30|Unknown|5||Have not seen these change during my time with the charger|
|31|Safety Timer|1|mins*10|Total charging time limit for safety|
|32-33|Capacity Limit|2|mAh|Charging capacity limit for safety|
|34-35|Charging Current|2|mAh*10|Real time charging/discharging current|
|36-37|Pack Voltage|2|Volts/100|Real time pack voltage. Byte 36 is whole number, Byte 37 is remainder|
|38-41|Unknown|4||Have not seen these change during my time with the charger|
|42-43|Input Voltage|2|Volts/100|Real time input voltage (power supply), encoded same as Pack Voltage|
|44-45|Total Charge|2|mAh|Real time transfered energy to or from the pack|
|46-47|Cell 1 Voltage|2|Volts/100|Real time cell 1 voltage. Encoded same as Pack Voltage.|
|48-49|Cell 2 Voltage|2|Volts/100|Real time cell 2 voltage. Encoded same as Pack Voltage.|
|50-51|Cell 3 Voltage|2|Volts/100|Real time cell 3 voltage. Encoded same as Pack Voltage.|
|52-53|Cell 4 Voltage|2|Volts/100|Real time cell 4 voltage. Encoded same as Pack Voltage.|
|54-55|Cell 5 Voltage|2|Volts/100|Real time cell 5 voltage. Encoded same as Pack Voltage.|
|56-57|Cell 6 Voltage|2|Volts/100|Real time cell 6 voltage. Encoded same as Pack Voltage.|
|58-59|Cell 7 Voltage|2|Volts/100|Real time cell 7 voltage. Encoded same as Pack Voltage.|
|60-61|Cell 8 Voltage|2|Volts/100|Real time cell 8 voltage. Encoded same as Pack Voltage.|
|62-63|Cell 9 Voltage|2|Volts/100|Real time cell 9 voltage. Encoded same as Pack Voltage.|
|64-65|Cell 10 Voltage|2|Volts/100|Real time cell 10 voltage. Encoded same as Pack Voltage.|
|66-67|Cell 11 Voltage|2|Volts/100|Real time cell 11 voltage. Encoded same as Pack Voltage.|
|68-69|Cell 12 Voltage|2|Volts/100|Real time cell 12 voltage. Encoded same as Pack Voltage.|
|70-71|Running Time|2|min|Total running time. Undefined behavior when not in running mode.|
|72|Unknown|1||Have not seen these change during my time with the charger|
|73|Trickle Current|1|lookup table|Amount of current for trickle-charging NiCD,NiMH,Pb chemestries.|
|74-75|Check Digits|2||Not sure how these are calculated, but they only change when other values within the packet change.|

### Trickle Current
The trickle current is measured in 10mA steps from 50mA to 200mA.
```
53 = off
54 = 50mA
55 = 60mA
.....
64 = 200mA
```

### Config Bits
```
0 = off
1 = on

Key Beep
| 
| |
| | | Buzzer
| | | | 
| | | | | Capacity Cutoff
| | | | | | Safety Timer 
| | | | | | | Temp Cut off
| | | | | | | |
7 6 5 4 3 2 1 0 
```

### Screen Selection
```
00: Program Menu Item
01: Li*
02: NiMH
03: NiCD
04: Pb
05: Save Menu Item
06: Load Menu Item
```
