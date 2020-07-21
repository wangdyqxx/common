package util

import (
	"errors"
	"fmt"
	"net"
)

func GetListenAddr(a string) (string, error) {
	addrTcp, err := net.ResolveTCPAddr("tcp", a)
	if err != nil {
		return "", err
	}
	addr := addrTcp.String()
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	if len(host) == 0 {
		return GetServAddr(addrTcp)
	}
	return addr, nil
}


func GetServAddr(a net.Addr) (string, error) {
	addr := a.String()
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	if len(host) == 0 {
		host = "0.0.0.0"
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return "", fmt.Errorf("ParseIP error:%s", host)
	}
	raddr := addr
	if ip.IsUnspecified() {
		// 没有指定ip的情况下，使用内网地址
		inerip, err := GetInterIp()
		if err != nil {
			return "", err
		}
		raddr = net.JoinHostPort(inerip, port)
	}
	return raddr, nil
}

func GetInterIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println(ipnet.IP.String())
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("no inter ip")
}


// 获取首个外网ip v4
func GetExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
