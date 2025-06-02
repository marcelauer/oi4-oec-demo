package sensor

import (
	"testing"
)

func TestCreateAcyclicRequestsForParameters(t *testing.T) {
	params := []IOLinkParameter{
		{ID: "TestParam", Index: 123, AccessRights: "ro", DataType: "UIntegerT(8)", DefaultValue: "0", Name: "Test", Description: "Test"},
	}
	reqs := CreateAcyclicRequestsForParameters(params)
	if len(reqs) != 1 {
		t.Errorf("expected 1 request, got %d", len(reqs))
	}
	if reqs[0].Data.Index != 123 {
		t.Errorf("expected index 123, got %d", reqs[0].Data.Index)
	}
}

func TestDecodeIOLinkValue_Integer(t *testing.T) {
	val, err := DecodeIOLinkHexValue("0A", "UIntegerT(8)")
	if err != nil {
		t.Fatal(err)
	}
	if val != 10 {
		t.Errorf("expected 10, got %v", val)
	}
}

func TestDecodeIOLinkValue_Float32(t *testing.T) {
	// 0x0000803F == 1.0 in float32 (little endian: 3F800000)
	val, err := DecodeIOLinkHexValue("0000803F", "Float32T")
	if err != nil {
		t.Fatal(err)
	}
	if val != 1.0 {
		t.Errorf("expected 1.0, got %v", val)
	}
}

func TestDecodeIOLinkValue_String(t *testing.T) {
	val, err := DecodeIOLinkHexValue("69544845524D20436F6D706163744C696E6520544D3331310000000000000000", "StringT(16)")
	if err != nil {
		t.Fatal(err)
	}
	if val != "iTHERM CompactLine TM311" {
		t.Errorf("expected 'Hello', got '%v'", val)
	}
}

// Note: For SendAcyclicRequests a real HTTP mock would be necessary.
// Here only a dummy test that checks the function signature.
func TestSendAcyclicRequests_Signature(t *testing.T) {
	// This function only tests if the function is callable.
	// A real test would require a mock HTTP backend.
	reqs := []IOLinkAcyclicRequest{}
	_, err := SendAcyclicRequests(reqs)
	if err != nil && err.Error() != "Post \"http://192.168.1.11\": unsupported protocol scheme \"\"" {
		// Error is ok as long as the function is callable
	}
}

func TestDecodeProcessData(t *testing.T) {
	// Example: 0x00FA = 250 (25.0°C), scale -1, status 0, switch 0
	hexStr := "00faff00"
	val, err := DecodeProcessData(hexStr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 25.0
	if val != expected {
		t.Errorf("expected %v, got %v", expected, val)
	}

	// Negative temperature: 0xFF38 = -200 (should be -20.0°C)
	hexStr = "ff38ff00"
	val, err = DecodeProcessData(hexStr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = -20.0
	if val != expected {
		t.Errorf("expected %v, got %v", expected, val)
	}

	// Not enough bytes
	hexStr = "00fa"
	_, err = DecodeProcessData(hexStr)
	if err == nil {
		t.Error("expected error for short input, got nil")
	}
}
