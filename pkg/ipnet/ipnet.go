package ipnet

import (
	"fmt"
	"log"
	"net"

	"gopkg.in/netaddr.v1"
)

type cidrObject struct {
	CidrDecimal    int64
	NetIP          net.IP
	BcastIP        net.IP
	NetIPDecimal   int64
	BcastIPDecimal int64
	Iptype         string
}

func PrepareCidrforDB(cidrIn string) (cidrOut cidrObject, err error) {
	//Process a Cidr and return first address (net), and last address (bcast)
	_, ipnet, err := net.ParseCIDR(cidrIn)

	if err != nil {
		log.Println("PrepareCidrforDB Error: ", cidrIn, err)
		return cidrOut, err
	}

	cidrOut.BcastIP = netaddr.BroadcastAddr(ipnet)
	cidrOut.NetIP = netaddr.NetworkAddr(ipnet)
	cidrOut.Iptype = IpType(fmt.Sprintf("%q", cidrOut.NetIP))
	
	// Only support IPv4 for now
	if cidrOut.Iptype != "IPv4" {
		return cidrOut, fmt.Errorf("only IPv4 addresses are supported")
	}
	
	cidrOut.NetIPDecimal = IPv4toDecimal(cidrOut.NetIP)
	cidrOut.BcastIPDecimal = IPv4toDecimal(cidrOut.BcastIP)

	return cidrOut, nil
}

func IpType(s string) string {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return "IPv4"
		case ':':
			return "IPv6"
		}
	}
	return "UNKNOWN"

}

func isZeros(p net.IP) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return false
		}
	}
	return true
}

func IPv4toDecimal(ipIn net.IP) (decimalOut int64) {
	//Convert an IP4 Address to a decimal
	ipOct := net.IP.To4(ipIn)
	if ipOct == nil {
		return 0
	}
	octInts := [4]int64{int64(ipOct[0]) * 16777216, int64(ipOct[1]) * 65536, int64(ipOct[2]) * 256, int64(ipOct[3])}

	for _, value := range octInts {
		decimalOut = decimalOut + value
	}
	return decimalOut
}
