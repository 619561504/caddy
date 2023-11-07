/**
 * Created by wuhanjie on 2023/11/7 10:24
 */

package tproxy

import (
	"errors"
	"fmt"
	"github.com/coreos/go-iptables/iptables"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"net"
	"strconv"
	"syscall"
)

var (
	EmptyIpNetErr = errors.New("ip net is empty")
)

const (
	Mark  = 123
	Table = 100
	Chain = "mangle"
)

type MmProxy struct {
	ipNet   *net.IPNet
	mark    int
	table   int
	rule    *netlink.Rule
	route   *netlink.Route
	iptRule []string
}

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

	rule := netlink.NewRule()
	rule.Table = p.table
	rule.Mark = p.mark
	p.rule = rule

	lo, _ := net.InterfaceByName("lo")
	route := &netlink.Route{
		LinkIndex: lo.Index,
		Type:      unix.RTN_LOCAL,
		Scope:     netlink.SCOPE_HOST,
		Dst:       &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
		Table:     p.table,
	}
	p.route = route
	p.iptRule = []string{"-p", "tcp", "-d", subIpNet, "-j", "MARK", "--set-mark", strconv.Itoa(mark)}

	return p, nil
}

// AddRule 检查添加路由、iptables规则
func (p *MmProxy) AddRule() error {

	if p.ipNet == nil {
		return EmptyIpNetErr
	}

	err := netlink.RuleDel(p.rule)
	if err != nil {
		fmt.Printf("rule del err: %v\n", err)
	}
	err = netlink.RuleAdd(p.rule)

	if err != nil {
		fmt.Printf("rule add err: %v\n", err)
	}

	err = netlink.RouteAdd(p.route)

	if err != nil {
		fmt.Printf("route add err: %v\n", err)
	}
	return err
}

func (p *MmProxy) AddIptables() error {

	ipt, err := iptables.New()
	if err != nil {
		return err
	}
	err = ipt.InsertUnique(Chain, "PREROUTING", 1, p.iptRule...)
	return err
}

func (p *MmProxy) TcpIpTRANSPARENTConn(sIp net.IP, targetAddr string) (net.Conn, error) {
	if !p.ipNet.Contains(sIp) {
		return nil, errors.New(fmt.Sprintf("Source Ip %s is not in ipnet %s", sIp, p.ipNet))
	}
	dialer := net.Dialer{LocalAddr: &net.TCPAddr{IP: sIp}}
	dialer.Control = DialUpstreamControl(0, 0)
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

func DialUpstreamControl(sport int, mark int) func(string, string, syscall.RawConn) error {
	return func(network, address string, c syscall.RawConn) error {
		var syscallErr error
		err := c.Control(func(fd uintptr) {
			syscallErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_SYNCNT, 2)
			if syscallErr != nil {
				syscallErr = fmt.Errorf("setsockopt(IPPROTO_TCP, TCP_SYNCTNT, 2): %w", syscallErr)
				return
			}

			syscallErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TRANSPARENT, 1)
			if syscallErr != nil {
				syscallErr = fmt.Errorf("setsockopt(IPPROTO_IP, IP_TRANSPARENT, 1): %w", syscallErr)
				return
			}

			syscallErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			if syscallErr != nil {
				syscallErr = fmt.Errorf("setsockopt(SOL_SOCKET, SO_REUSEADDR, 1): %w", syscallErr)
				return
			}

			if sport == 0 {
				ipBindAddressNoPort := 24
				syscallErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, ipBindAddressNoPort, 1)
				if syscallErr != nil {
					syscallErr = fmt.Errorf("setsockopt(SOL_SOCKET, IPPROTO_IP, %d): %w", mark, syscallErr)
					return
				}
			}

			if mark != 0 {
				syscallErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, mark)
				if syscallErr != nil {
					syscallErr = fmt.Errorf("setsockopt(SOL_SOCK, SO_MARK, %d): %w", mark, syscallErr)
					return
				}
			}

			if network == "tcp6" || network == "udp6" {
				syscallErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 0)
				if syscallErr != nil {
					syscallErr = fmt.Errorf("setsockopt(IPPROTO_IP, IPV6_ONLY, 0): %w", syscallErr)
					return
				}
			}
		})

		if err != nil {
			return err
		}
		return syscallErr
	}
}

func (p *MmProxy) Clear() error {
	err := netlink.RuleDel(p.rule)

	if err != nil {
		fmt.Printf("rule del err: %v\n", err)
	}

	err = netlink.RouteDel(p.route)
	if err != nil {
		fmt.Printf("route del err: %v\n", err)
	}

	ipt, err := iptables.New()
	if err != nil {
		fmt.Printf("iptables new err: %v\n", err)
		return err
	}

	err = ipt.DeleteIfExists(Chain, "PREROUTING", p.iptRule...)
	if err != nil {
		fmt.Printf("iptables del err: %v\n", err)
	}

	return err
}
