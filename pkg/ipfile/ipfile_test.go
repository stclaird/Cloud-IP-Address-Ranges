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
