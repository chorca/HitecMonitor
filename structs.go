package main

type chargerFrame struct {
	// Battery Stats
	chargingCurrent int     // Realtime current being fed to the battery
	packVoltage     float32 // Realtime total pack voltage
	mahCharged      int     // Running total of power fed into the pack during charging in milliamps
	totalTime       int     // Minutes the pack has been charging for
	running         bool    // Is a program currently running
	chargeDischarge bool    // True if charging, false if discharging (Power going into or out of the battery)
	// Battery Settings
	setChargeCurrent int // Charging current specified on the front panel
	numCells         int // How many cells (S) the current pack is set to
	// Charger settings
	trickleMa int // Milliamps used for trickle-charging NiCD, NiMH, and Pb batteries
	nicdSens  int // Sensitivty (mV/Cell) for NiCD batteries (Default 40)
	nimhSens  int // Sensitivty (mV/Cell) for NiMH batteries (Default 40)
	wasteTime int // Amount of time delay in minutes between Discharge to Charge cycles
	// Charger prefs
	keyBeepOn      bool    // Is the key beeper enabled?
	buzzerOn       bool    // Is the buzzer enabled (cycle finish)
	dispBacklight  int     // Backlight percentage
	lowPowerCutoff float32 // Input voltage at which to shut off the charger
	// Safety preferences
	capCutoffOn  bool // Is capacity limit enabled?
	timeCutoffOn bool // Is the charge time limit enabled?
	tempCutoffOn bool // Is the temperature shutoff enabled?
	tempLimit    int  // Cell temperature limit (external port) for shutoff
	timeLimit    int  // Time limit for charging
	capLimit     int  // Capacity limit for pack
	// Misc
	screenMode int     // Current screen mode/function
	inputVolts float32 // Input voltage to the charger
	// Individual Cell voltages
	cell1Volts  float32
	cell2Volts  float32
	cell3Volts  float32
	cell4Volts  float32
	cell5Volts  float32
	cell6Volts  float32
	cell7Volts  float32
	cell8Volts  float32
	cell9Volts  float32
	cell10Volts float32
	cell11Volts float32
	cell12Volts float32
}

// ConfigJSON contains configuration read from the config file
type ConfigJSON struct {
	InfluxToken   string `json:"influxToken"`
	InfluxBucket  string `json:"influxBucket"`
	InfluxOrg     string `json:"influxOrg"`
	ServerAddress string `json:"serverAddress"`
}
