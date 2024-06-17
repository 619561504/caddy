//go:build !linux
// +build !linux

/**
 * Created by wuhanjie on 2023/11/7 10:24
 */

package tproxy

import (
	"errors"
	"fmt"
	"net"
)

func NewMmProxy(subIpNet string, mark int, table int) (*MmProxy, error) {

	_, ipNet, err := net.ParseCIDR(subIpNet)
	if err != nil {
		return nil, err
	}

	if mark == 0 {
		mark = Mark
	}
	if table == 0 {
		table = Table
	}

	p := &MmProxy{ipNet: ipNet, mark: mark, table: table}

	return p, nil
}

// AddRule 检查添加路由、iptables规则
func (p *MmProxy) AddRule() error {

	return nil
}

func (p *MmProxy) AddIptables() error {

	return nil
}

func (p *MmProxy) TcpIpTRANSPARENTConn(sIp net.IP, targetAddr string) (net.Conn, error) {
	if !p.ipNet.Contains(sIp) {
		return nil, errors.New(fmt.Sprintf("Source Ip %s is not in ipnet %s", sIp, p.ipNet))
	}
	dialer := net.Dialer{LocalAddr: &net.TCPAddr{IP: sIp}}
	conn, err := dialer.Dial("tcp", targetAddr)
	if err != nil {
		return nil, err
	}
	if err := conn.(*net.TCPConn).SetNoDelay(true); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func (p *MmProxy) Clear() error {
	return nil
}
