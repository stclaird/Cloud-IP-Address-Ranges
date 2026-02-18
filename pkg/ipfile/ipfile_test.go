package ipfile

import (
	"testing"
)

func TestMatchIp(t *testing.T) {

	ipAddressstring := "IP: 82.12.162.1/32"
	have := MatchIp(ipAddressstring)
	t.Logf("Have: %v", have[0])

	var require [1]string
	require[0] = "82.12.162.1/32"
	t.Log("Require", require[0])

	if require[0] != have[0] {
		t.Errorf("Expected '%s', but got '%s'", require, have)
	}
}

func TestStrInSlice(t *testing.T) {

	var inSlice []string
	inSlice = append(inSlice, "cat")
	have := StrInSlice("cat", inSlice)
	t.Log("Have:", have)

	require := true
	t.Log("Require", require)

	if require != have {
		t.Errorf("Expected '%t', but got '%t'", require, have)
	}
}

func TestAsText(t *testing.T) {

}

func TestMatchIpWithTimestamps(t *testing.T) {
	// Test that timestamps are not matched as IPv6
	jsonWithTimestamp := `{"creationTime": "2026-02-16T12:07:07.518526", "ipv4Prefix": "8.8.8.0/24"}`
	matches := MatchIp(jsonWithTimestamp)
	
	// Should only match the IPv4 address, not the timestamp
	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d: %v", len(matches), matches)
	}
	
	if matches[0] != "8.8.8.0/24" {
		t.Errorf("Expected '8.8.8.0/24', got '%s'", matches[0])
	}
	
	// Verify no timestamp was matched
	for _, match := range matches {
		if match == "12:07:07" || match == "12:07:07/128" {
			t.Errorf("Timestamp was incorrectly matched as IPv6: %s", match)
		}
	}
}

func TestMatchIpIPv6Valid(t *testing.T) {
	// Test valid IPv6 addresses are matched
	testCases := []struct {
		input    string
		expected []string
	}{
		{"2001:4860:4860::8888", []string{"2001:4860:4860::8888"}},
		{"2001:db8::1/64", []string{"2001:db8::1/64"}},
		{"fe80::1", []string{"fe80::1"}},
		{"::1", []string{"::1"}},
		{"2001:0db8:0000:0000:0000:ff00:0042:8329", []string{"2001:0db8:0000:0000:0000:ff00:0042:8329"}},
	}
	
	for _, tc := range testCases {
		matches := MatchIp(tc.input)
		if len(matches) != len(tc.expected) {
			t.Errorf("For input '%s': expected %d matches, got %d: %v", tc.input, len(tc.expected), len(matches), matches)
			continue
		}
		if len(matches) > 0 && matches[0] != tc.expected[0] {
			t.Errorf("For input '%s': expected '%s', got '%s'", tc.input, tc.expected[0], matches[0])
		}
	}
}
