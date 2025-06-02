package sensor

// IOLinkAcyclicResponse represents the response structure for an IOLink acyclic read request.
type IOLinkAcyclicResponse struct {
	Cid  int `json:"cid"`
	Data struct {
		Value string `json:"value"`
	} `json:"data"`
	Code int `json:"code"`
}

// IOLinkAcyclicRequest represents the request structure for an IOLink acyclic write request.
type IOLinkAcyclicRequest struct {
	Code string `json:"code"`
	Cid  int    `json:"cid"`
	Adr  string `json:"adr"`
	Data struct {
		Index    int `json:"index"`
		Subindex int `json:"subindex"`
	} `json:"data"`
}

// IOLinkParameter represents a parameter defined in the IODD (IOLink Device Description) for the iTHERM CompactLine TM311.
type IOLinkParameter struct {
	ID           string // z.B. "V_MS_Unit"
	Index        int    // z.B. 5121
	AccessRights string // z.B. "rw", "ro"
	DataType     string // z.B. "UIntegerT", "Float32T"
	DefaultValue string // z.B. "32"
	Name         string // z.B. "Unit"
	Description  string // z.B. "Selection of the unit for all measured values."
}

// ParameterList contains all parameters defined in the IODD for the iTHERM CompactLine TM311
var parameterList = []IOLinkParameter{
	// --- Switch Output ---
	{ID: "V_SWO_OperatingMode", Index: 2050, AccessRights: "rw", DataType: "UIntegerT(16)", DefaultValue: "3228", Name: "Operating mode", Description: "Indicates the operating mode of the switch output."},
	{ID: "V_SWO_SwitchPointValue", Index: 2051, AccessRights: "rw", DataType: "Float32T", DefaultValue: "100", Name: "Switch point value", Description: "Switch point value for hysteresis / Upper value for window function."},
	{ID: "V_SWO_SwitchbackPointValue", Index: 2052, AccessRights: "rw", DataType: "Float32T", DefaultValue: "90", Name: "Switchback point value", Description: "Switchback point value for hysteresis / Lower value for window function."},
	{ID: "V_SWO_SwitchDelay", Index: 2053, AccessRights: "rw", DataType: "UIntegerT(8)", DefaultValue: "0", Name: "Switch delay", Description: "Setting of the switch delay."},
	{ID: "V_SWO_SwitchbackDelay", Index: 2054, AccessRights: "rw", DataType: "UIntegerT(8)", DefaultValue: "0", Name: "Switchback delay", Description: "Setting of the switchback delay."},
	{ID: "V_SWO_Simulation", Index: 2056, AccessRights: "rw", DataType: "UIntegerT(16)", DefaultValue: "0", Name: "Switch output simulation", Description: "Use this function to enter a switch simulation."},

	// --- Sensor ---
	{ID: "V_ST1_SensorMaxIndicator", Index: 3080, AccessRights: "ro", DataType: "Float32T", DefaultValue: "", Name: "Sensor max value", Description: "Indication of the maximum measured temperature in the past."},
	{ID: "V_ST1_SensorMinIndicator", Index: 3081, AccessRights: "ro", DataType: "Float32T", DefaultValue: "", Name: "Sensor min value", Description: "Indication of the minimum measured temperature in the past."},
	{ID: "V_ST1_SensorOffset", Index: 3082, AccessRights: "rw", DataType: "Float32T", DefaultValue: "0", Name: "Sensor offset", Description: "Setting of the zero adjustment of the measured sensor value."},
	{ID: "V_ST1_SimSensorValue", Index: 3104, AccessRights: "rw", DataType: "Float32T", DefaultValue: "0", Name: "Sensor simulation value", Description: "Use this function to enter a simulation value of the process variable."},
	{ID: "V_ST1_SimEnable", Index: 3109, AccessRights: "rw", DataType: "UIntegerT(8)", DefaultValue: "0", Name: "Sensor simulation", Description: "Use this function to activate the simulation of the process variable."},
	{ID: "V_ST1_OperatingTime1", Index: 3132, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Lower boundary operating time sensor", Description: "Indicates how long the sensor has been in operation in the lower boundary temperature range."},
	{ID: "V_ST1_OperatingTime2", Index: 3133, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Lower extended operating time sensor", Description: "Indicates how long the sensor has been in operation in the lower extended temperature range."},
	{ID: "V_ST1_OperatingTime3", Index: 3134, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Standard operating time sensor", Description: "Indicates how long the sensor has been in operation in the standard temperature range."},
	{ID: "V_ST1_OperatingTime4", Index: 3135, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Upper extended operating time sensor", Description: "Indicates how long the sensor has been in operation in the upper extended temperature range."},
	{ID: "V_ST1_OperatingTime5", Index: 3136, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Upper boundary operating time sensor", Description: "Indicates how long the sensor has been in operation in the upper boundary temperature range."},

	// --- Device Temperature ---
	{ID: "V_DevTemp_Value", Index: 4096, AccessRights: "ro", DataType: "Float32T", DefaultValue: "", Name: "Device temperature", Description: "Indication of the measured device temperature."},
	{ID: "V_DevTemp_MaxIndicator", Index: 4106, AccessRights: "ro", DataType: "Float32T", DefaultValue: "", Name: "Device temperature max", Description: "Indication of the maximum measured device temperature in the past."},
	{ID: "V_DevTemp_MinIndicator", Index: 4107, AccessRights: "ro", DataType: "Float32T", DefaultValue: "", Name: "Device temperature min", Description: "Indication of the minimum measured device temperature in the past."},
	{ID: "V_DevTemp_OperatingTime1", Index: 4109, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Lower boundary operating time device", Description: "Indicates how long the device temperature has been in operation in the lower boundary temperature range."},
	{ID: "V_DevTemp_OperatingTime2", Index: 4110, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Lower extended operating time device", Description: "Indicates how long the device temperature has been in operation in the lower extended temperature range."},
	{ID: "V_DevTemp_OperatingTime3", Index: 4111, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Standard operating time device", Description: "Indicates how long the device temperature has been in operation in the standard temperature range."},
	{ID: "V_DevTemp_OperatingTime4", Index: 4112, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Upper extended operating time device", Description: "Indicates how long the device temperature has been in operation in the upper extended temperature range."},
	{ID: "V_DevTemp_OperatingTime5", Index: 4113, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Upper boundary operating time device", Description: "Indicates how long the device temperature has been in operation in the upper boundary temperature range."},

	// --- Measuring Data Channel Descriptor ---
	{ID: "V_MDC_Descriptor", Index: 16512, AccessRights: "ro", DataType: "RecordT(56)", DefaultValue: "", Name: "MDC Descriptor", Description: ""},

	// --- Unit & Damping ---
	{ID: "V_MS_Unit", Index: 5121, AccessRights: "rw", DataType: "UIntegerT(8)", DefaultValue: "32", Name: "Unit", Description: "Selection of the unit for all measured values."},
	{ID: "V_DV_PVDamping", Index: 7271, AccessRights: "rw", DataType: "UIntegerT(8)", DefaultValue: "0", Name: "Damping", Description: "Setting of the time constant for the damping of the measured value."},

	// --- Diagnostic System ---
	{ID: "V_DS_AlarmDelay", Index: 6147, AccessRights: "rw", DataType: "UIntegerT(8)", DefaultValue: "2", Name: "Alarm delay", Description: "Setting of the duration how long diagnostic messages are suppressed"},
	{ID: "V_DS_OperatingTime", Index: 6148, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Operating time", Description: "Indicates how long the device has been in operation"},
	{ID: "V_DS_ActualDiagnostics1", Index: 6184, AccessRights: "ro", DataType: "D_DiagnosticsType", DefaultValue: "", Name: "Actual diagnostics 1", Description: "Displays the currently active diagnostic message with the highest priority."},
	{ID: "V_DS_ActualDiagnostics2", Index: 6186, AccessRights: "ro", DataType: "D_DiagnosticsType", DefaultValue: "", Name: "Actual diagnostics 2", Description: "Displays the currently active diagnostic message with the second highest priority."},
	{ID: "V_DS_ActualDiagnostics3", Index: 6188, AccessRights: "ro", DataType: "D_DiagnosticsType", DefaultValue: "", Name: "Actual diagnostics 3", Description: "Displays the currently active diagnostic message with the third highest priority."},
	{ID: "V_DS_LastDiagTimestamp01", Index: 6204, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Timestamp 1", Description: "Shows the timestamp of the previous diagnostic message."},
	{ID: "V_DS_LastDiagTimestamp02", Index: 6205, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Timestamp 2", Description: "Shows the timestamp of the previous diagnostic message."},
	{ID: "V_DS_LastDiagTimestamp03", Index: 6206, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Timestamp 3", Description: "Shows the timestamp of the previous diagnostic message."},
	{ID: "V_DS_LastDiagTimestamp04", Index: 6207, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Timestamp 4", Description: "Shows the timestamp of the previous diagnostic message."},
	{ID: "V_DS_LastDiagTimestamp05", Index: 6208, AccessRights: "ro", DataType: "UIntegerT(32)", DefaultValue: "", Name: "Timestamp 5", Description: "Shows the timestamp of the previous diagnostic message."},
	{ID: "V_DS_LastDiagnostics1", Index: 6214, AccessRights: "ro", DataType: "D_DiagnosticsType", DefaultValue: "", Name: "Previous diagnostics 1", Description: "Display of the diagnostic messages that have occurred in the past."},
	{ID: "V_DS_LastDiagnostics2", Index: 6216, AccessRights: "ro", DataType: "D_DiagnosticsType", DefaultValue: "", Name: "Previous diagnostics 2", Description: "Display of the diagnostic messages that have occurred in the past."},
	{ID: "V_DS_LastDiagnostics3", Index: 6218, AccessRights: "ro", DataType: "D_DiagnosticsType", DefaultValue: "", Name: "Previous diagnostics 3", Description: "Display of the diagnostic messages that have occurred in the past."},
	{ID: "V_DS_LastDiagnostics4", Index: 6220, AccessRights: "ro", DataType: "D_DiagnosticsType", DefaultValue: "", Name: "Previous diagnostics 4", Description: "Display of the diagnostic messages that have occurred in the past."},
	{ID: "V_DS_LastDiagnostics5", Index: 6222, AccessRights: "ro", DataType: "D_DiagnosticsType", DefaultValue: "", Name: "Previous diagnostics 5", Description: "Display of the diagnostic messages that have occurred in the past."},

	// --- Current Output ---
	{ID: "V_CO_SimulationMode", Index: 8210, AccessRights: "rw", DataType: "UIntegerT(16)", DefaultValue: "33004", Name: "Current output simulation", Description: "Switch the simulation of the current output on and off"},
	{ID: "V_CO_SimulationValue", Index: 8211, AccessRights: "rw", DataType: "Float32T", DefaultValue: "3.58", Name: "Value current output", Description: "Enter the current value for simulation"},
	{ID: "V_CO_TrimValueHi", Index: 8212, AccessRights: "rw", DataType: "Float32T", DefaultValue: "20.00", Name: "Current trimming 20 mA", Description: "Setting of the correction value for the current output at the upper range value 20 mA."},
	{ID: "V_CO_TrimValueLo", Index: 8213, AccessRights: "rw", DataType: "Float32T", DefaultValue: "4.00", Name: "Current trimming 4 mA", Description: "Setting of the correction value for the current output at the lower range value 4 mA."},
	{ID: "V_CO_Input04mA", Index: 8218, AccessRights: "rw", DataType: "Float32T", DefaultValue: "0", Name: "4 mA value", Description: "Enter 4 mA value"},
	{ID: "V_CO_Input20mA", Index: 8219, AccessRights: "rw", DataType: "Float32T", DefaultValue: "150", Name: "20 mA value", Description: "Enter 20 mA value"},
	{ID: "V_CO_HiAlarmCurrent", Index: 8232, AccessRights: "rw", DataType: "Float32T", DefaultValue: "22.5", Name: "Failure current", Description: "Setting of the value the current output adopts in an alarm condition."},
	{ID: "V_CO_FailureMode", Index: 8234, AccessRights: "rw", DataType: "UIntegerT(8)", DefaultValue: "0", Name: "Failure mode", Description: "Selection of the failure signal level of the current output in case of an error."},

	// --- System ---
	{ID: "V_IOLINK_DeviceType", Index: 256, AccessRights: "ro", DataType: "UIntegerT(16)", DefaultValue: "", Name: "Device type", Description: "Shows the device type"},
	{ID: "V_STD_ENPDeviceOrderIdent", Index: 1054, AccessRights: "ro", DataType: "StringT(20)", DefaultValue: "", Name: "Order code", Description: "Shows the device order code"},
	{ID: "V_STD_ENPDeviceOrderCode", Index: 259, AccessRights: "ro", DataType: "StringT(60)", DefaultValue: "", Name: "Extended order code", Description: "The extended order code indicates the version of all the features of the product structure for the device."},
}

type Value struct {
	Key  string `json:"key"`
	Unit struct {
		Code *string `json:"code"`
	} `json:"unit"`
	Data []struct {
		Timestamp *string `json:"timestamp,omitempty"`
		Value     float64 `json:"value"`
		Status    *string `json:"status,omitempty"`
		Simulated *bool   `json:"simulated,omitempty"`
		Hold      *bool   `json:"hold,omitempty"`
	} `json:"data,omitempty"`
}
