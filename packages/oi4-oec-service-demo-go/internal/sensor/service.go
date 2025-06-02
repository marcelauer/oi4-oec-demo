package sensor

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

const timeout = 10 * time.Second

type Service struct {
	IOLinkMasterAddress string
}

func NewSensorService(IOLinkMasterAddress string) *Service {
	return &Service{
		IOLinkMasterAddress: IOLinkMasterAddress,
	}
}

func (s *Service) SendAcyclicRequests(params []IOLinkParameter) ([]IOLinkAcyclicResponse, error) {
	requests := CreateAcyclicRequestsForParameters(params)
	return s.sendAcyclicRequests(requests)
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

func (s *Service) GetSensorData(params []IOLinkParameter) ([]IOLinkAcyclicResponse, error) {
	requests := CreateAcyclicRequestsForParameters(params)
	responses, err := s.sendAcyclicRequests(requests)
	if err != nil {
		return nil, err
	}

	for i, resp := range responses {
		if resp.Data.Value != "" {
			decodedValue, decodeErr := DecodeIOLinkValue(resp.Data.Value, params[i].DataType)
			if decodeErr != nil {
				return nil, decodeErr
			}
			if strVal, ok := decodedValue.(string); ok {
				resp.Data.Value = strVal
			} else {
				resp.Data.Value = ""
			}
		}
	}

	return responses, nil
}

// sendAcyclicRequests sends each request as a POST REST API call to the IOLinkMasterAddress and returns the responses.
func (s *Service) sendAcyclicRequests(requests []IOLinkAcyclicRequest) ([]IOLinkAcyclicResponse, error) {
	var responses []IOLinkAcyclicResponse
	var iolResp IOLinkAcyclicResponse
	client := &http.Client{}
	for _, req := range requests {
		body, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		httpReq, err := http.NewRequest("POST", s.IOLinkMasterAddress, bytes.NewBuffer(body))
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

// DecodeIOLinkValue decodes a hex string from the response according to the DataType into the appropriate Go type.
// Supports Integer, Float32, and ASCII string.
func DecodeIOLinkValue(hexStr string, dataType string) (interface{}, error) {
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
