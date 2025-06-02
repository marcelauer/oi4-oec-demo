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
	} `json:"data,omitempty"`
}

// IOLinkParameter represents a parameter defined in the IODD (IOLink Device Description) for the iTHERM CompactLine TM311.
type IOLinkParameter struct {
	ID           string // e.g. "V_MS_Unit"
	Index        int    // e.g. 5121
	AccessRights string // e.g. "rw", "ro"
	DataType     string // e.g. "UIntegerT", "Float32T"
	DefaultValue string // e.g. "32"
	Name         string // e.g. "Unit"
	Description  string // e.g. "Selection of the unit for all measured values."
}

type SensorUnit int

const (
	Celsius    SensorUnit = 32
	Fahrenheit SensorUnit = 33
	Kelvin     SensorUnit = 35
)

type ResponseValue struct {
	Key  string `json:"key"`
	Unit struct {
		Code *string `json:"code"`
	} `json:"unit,omitempty"`
	Group *string    `json:"group,omitempty"`
	Data  []DataDict `json:"data"`
}

type DataDict struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
}
