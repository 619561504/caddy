/**
 * Created by wuhanjie on 2023/11/20 16:36
 */

package tproxy

import (
	"errors"
	"github.com/vishvananda/netlink"
	"net"
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
