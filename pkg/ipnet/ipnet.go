package ipnet

import (
	"errors"
	"fmt"
	"log"
	"net"

	"gopkg.in/netaddr.v1"
)

type cidrObject struct {
	CidrDecimal    int
	NetIP          net.IP
	BcastIP        net.IP
	NetIPDecimal   int
	BcastIPDecimal int
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
	cidrOut.NetIPDecimal = IPv4toDecimal(cidrOut.NetIP)
	cidrOut.BcastIPDecimal = IPv4toDecimal(cidrOut.BcastIP)

	cidrOut.Iptype = IpType(fmt.Sprintf("%q", cidrOut.NetIP))

	//deal with IPv6 until it is supported
	if cidrOut.Iptype == "IPv6" {
		err := errors.New("IPv6 Not supported Yet")
		log.Println("Error : ", cidrIn, err)
		return cidrOut, err
	}

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

func IPv4toDecimal(ipIn net.IP) (decimalOut int) {
	//Convert an IP4 Address to a decimal
	ipOct := net.IP.To4(ipIn)
	if ipOct == nil {
		return 0
	}
	octInts := [4]int{int(ipOct[0]) * 16777216, int(ipOct[1]) * 65536, int(ipOct[2]) * 256, int(ipOct[3])}

	for _, value := range octInts {
		decimalOut = decimalOut + value
	}
	return decimalOut
}
