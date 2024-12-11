package utils

import (
	"net"
)

// 将ip地址转换为长整型
func IP2Long(str string) int32 {
	if ip := net.ParseIP(str); ip != nil {
		var n uint32
		ipBytes := ip.To4()
		for i := uint8(0); i <= 3; i++ {
			n |= uint32(ipBytes[i]) << ((3 - i) * 8)
		}
		return int32(n)
	}
	return 0
}

// 将长整型转换为ip地址
func Long2IP(in int32) string {
	ipBytes := net.IP{}
	for i := uint(0); i <= 3; i++ {
		ipBytes = append(ipBytes, byte(in>>((3-i)*8)))
	}
	return ipBytes.String()
}
