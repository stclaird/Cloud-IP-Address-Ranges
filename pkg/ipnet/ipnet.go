package ipnet

import (
	"fmt"
	"log"
	"net"
	"strings"

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
	// Normalize some malformed IPv6 forms like "2602:ff03:4d5/128"
	// to a valid compressed form "2602:ff03:4d5::/128" so parsing succeeds.
	if strings.Contains(cidrIn, ":") && !strings.Contains(cidrIn, "::") {
		parts := strings.SplitN(cidrIn, "/", 2)
		addr := parts[0]
		if strings.Count(addr, ":") < 7 {
			addr = addr + "::"
			if len(parts) == 2 {
				cidrIn = addr + "/" + parts[1]
			} else {
				cidrIn = addr
			}
		}
	}

	_, ipnet, err := net.ParseCIDR(cidrIn)

	if err != nil {
		log.Println("PrepareCidrforDB Error: ", cidrIn, err)
		return cidrOut, err
	}

	cidrOut.BcastIP = netaddr.BroadcastAddr(ipnet)
	cidrOut.NetIP = netaddr.NetworkAddr(ipnet)
	cidrOut.Iptype = IpType(fmt.Sprintf("%q", cidrOut.NetIP))

	// If IPv4 populate decimal representations; for IPv6 leave decimals zero and use textual fields
	if cidrOut.Iptype == "IPv4" {
		cidrOut.NetIPDecimal = IPv4toDecimal(cidrOut.NetIP)
		cidrOut.BcastIPDecimal = IPv4toDecimal(cidrOut.BcastIP)
	} else {
		cidrOut.NetIPDecimal = 0
		cidrOut.BcastIPDecimal = 0
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
