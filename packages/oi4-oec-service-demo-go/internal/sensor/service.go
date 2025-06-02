package sensor

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// =============================
// Adapt these values for debugging with your IOLink Master
// =============================
const IOLinkMasterAddress = "http://192.168.1.11"

// =============================

const timeout = 10 * time.Second

type Service struct {
	IOLinkMasterAddress string
	logger              *zap.SugaredLogger
}

func NewSensorService(IOLinkMasterAddress string, logger *zap.SugaredLogger) *Service {
	return &Service{
		IOLinkMasterAddress: IOLinkMasterAddress,
		logger:              logger,
	}
}

// CreateAcyclicRequestsForParameters creates an acyclic request for each IOLinkParameter.
func CreateAcyclicRequestsForParameters(params []IOLinkParameter) []IOLinkAcyclicRequest {
	var requests []IOLinkAcyclicRequest
	for _, p := range params {
		req := IOLinkAcyclicRequest{
			Code: "request",
			Cid:  4711,
			Adr:  "/iolinkmaster/port[1]/iolinkdevice/iolreadacyclic",
		}
		req.Data.Index = p.Index
		req.Data.Subindex = 0
		requests = append(requests, req)
	}
	return requests
}

// SendAcyclicRequests sends each request as a POST REST API call to IP adress of IOLinkMaster and returns the responses.
func SendAcyclicRequests(requests []IOLinkAcyclicRequest) ([]IOLinkAcyclicResponse, error) {
	var responses []IOLinkAcyclicResponse
	var iolResp IOLinkAcyclicResponse
	client := &http.Client{Timeout: timeout}
	for _, req := range requests {
		body, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		httpReq, err := http.NewRequest("POST", IOLinkMasterAddress, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, err
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(respBody, &iolResp); err != nil {
			return nil, err
		}
		responses = append(responses, iolResp)
	}
	return responses, nil
}

func (s *Service) GetSensorParameterData(params []IOLinkParameter, group string) ([]ResponseValue, error) {
	requests := CreateAcyclicRequestsForParameters(params)
	responses, err := SendAcyclicRequests(requests)
	if err != nil {
		return nil, err
	}

	var values []ResponseValue

	for i, resp := range responses {
		if resp.Data.Value != "" {
			decodedValue, decodeErr := DecodeIOLinkHexValue(resp.Data.Value, params[i].DataType)
			if decodeErr != nil {
				return nil, decodeErr
			}
			if strVal, ok := decodedValue.(string); ok {
				var data = DataDict{
					Timestamp: time.Now().Unix(),
					Value:     strVal,
				}
				values = append(values, ResponseValue{
					Key:   params[i].ID,
					Group: &group,
					Data:  []DataDict{data},
				})
			} else {
				resp.Data.Value = ""
			}
		}
	}
	if len(values) == 0 {
		s.logger.Warn("No valid parameter values for parameter group " + group + " found in the response")
		return nil, errors.New("No valid parameter values for parameter group " + group + " found in the response")
	}

	return values, nil
}

func (s *Service) GetSensorProcessData(unit *string) ([]ResponseValue, error) {
	var iolResp IOLinkAcyclicResponse
	var processValue ResponseValue
	client := &http.Client{Timeout: timeout}
	req := IOLinkAcyclicRequest{
		Code: "request",
		Cid:  4711,
		Adr:  "/iolinkmaster/port[1]/iolinkdevice/getdata",
	}
	body, err := json.Marshal(req)
	httpReq, err := http.NewRequest("POST", IOLinkMasterAddress, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(respBody, &iolResp); err != nil {
		return nil, err
	}
	if iolResp.Data.Value != "" {
		decodedValue, decodedErr := DecodeProcessData(iolResp.Data.Value)
		if decodedErr != nil {
			return nil, decodedErr
		}
		processValue.Key = "Temperature"
		processValue.Unit.Code = unit
		processValue.Group = unit
		processValue.Data = []DataDict{
			{
				Timestamp: time.Now().Unix(),
				Value:     fmt.Sprintf("%v", decodedValue),
			},
		}
	}
	return []ResponseValue{processValue}, nil
}

// DecodeIOLinkValue decodes a hex string from the response according to the DataType into the appropriate Go type.
// Supports Integer, Float32, and ASCII string.
func DecodeIOLinkHexValue(hexStr string, dataType string) (interface{}, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	switch {
	case dataType == "Float32T":
		if len(data) < 4 {
			return nil, errors.New("not enough bytes for float32")
		}
		bits := uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
		floatVal := math.Float32frombits(bits)
		return float64(floatVal), nil
	case dataType == "UIntegerT(8)":
		if len(data) < 1 {
			return nil, errors.New("not enough bytes for uint8")
		}
		return int(data[0]), nil
	case dataType == "UIntegerT(16)":
		if len(data) < 2 {
			return nil, errors.New("not enough bytes for uint16")
		}
		return int(data[0]) | int(data[1])<<8, nil
	case dataType == "UIntegerT(32)":
		if len(data) < 4 {
			return nil, errors.New("not enough bytes for uint32")
		}
		return int(data[0]) | int(data[1])<<8 | int(data[2])<<16 | int(data[3])<<24, nil
	case dataType == "StringT(16)", dataType == "StringT(20)", dataType == "StringT(32)", dataType == "StringT(60)", dataType == "StringT(8)":
		// ASCII-String, Nullbytes entfernen
		return strings.ReplaceAll(string(data), "\x00", ""), nil
	default:
		return hexStr, nil // Fallback: return hex string
	}
}

// Decode ProcessData from hex string to float64 temperature value according to the iTherm process data format.
func DecodeProcessData(hexStr string) (float64, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return 0, err
	}
	if len(data) < 4 {
		return 0, errors.New("not enough bytes for process data")
	}

	// Temperature: Bytes 0 und 1 (Big Endian, sint16)
	tempRaw := int16(data[0])<<8 | int16(data[1])
	temperature := float64(tempRaw) / 10.0
	// Scale: Byte 2 (int8)
	//scale := int8(data[2]) is always 10^-1, so we can ignore it
	// Status and Switch State: Byte 3 (uint8)
	// status := (data[3] >> 1) & 0x07
	// switchState := data[3] & 0x01

	return temperature, nil
}

func (u SensorUnit) UnitToString() string {
	switch u {
	case Celsius:
		return "°C"
	case Fahrenheit:
		return "°F"
	case Kelvin:
		return "K"
	default:
		return "°C"
	}
}
