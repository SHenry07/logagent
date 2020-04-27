package utils

import (
	"log"
	"net"
)

func Init() (LocalAddr string) {
	Addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Panicf("Can't resolve local IP, %v", err)
	}

	for _, address := range Addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// 是ipv4   是全局单播地址    不是链路本地单播地址
			if ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
				LocalAddr = ipnet.IP.String()
				log.Printf("%v\n", LocalAddr)
			}
		}
	}

	return LocalAddr
}
