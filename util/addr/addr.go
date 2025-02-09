// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/util/addr/addr.go

package addr

import (
	"errors"
	"net"
)

var (
	privateBlocks []*net.IPNet
)

func init() {
	blocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"100.64.0.0/10",
		"fd00::/8",
	}
	AppendPrivateBlocks(blocks...)
}

// AppendPrivateBlocks append private network blocks
func AppendPrivateBlocks(bs ...string) {
	for _, b := range bs {
		if _, block, err := net.ParseCIDR(b); err == nil {
			privateBlocks = append(privateBlocks, block)
		}
	}
}

func isPrivateIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return false
	}

	for _, blocks := range privateBlocks {
		if blocks.Contains(ip) {
			return true
		}
	}
	return false
}

func addrToIP(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPAddr:
		return v.IP
	case *net.IPNet:
		return v.IP
	default:
		return nil
	}
}

func localIPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ipAddrs []string

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue // ignore error
		}

		for _, addr := range addrs {
			if ip := addrToIP(addr); ip != nil {
				ipAddrs = append(ipAddrs, ip.String())
			}
		}
	}

	return ipAddrs
}

// IsLocal tells us whether an ip is local
func IsLocal(addr string) bool {
	// extract the host
	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		addr = host
	}

	// check if its localhost
	if addr == "localhost" {
		return true
	}

	// check against all local ips
	for _, ip := range localIPs() {
		if addr == ip {
			return true
		}
	}

	return false
}

// Extract returns a real ip
func Extract(addr string) (string, error) {
	// if addr specified then its returned
	if len(addr) > 0 {
		if addr != "0.0.0.0" && addr != "[::]" && addr != "::" {
			return addr, nil
		}
	}

	var privateAddrs []string
	var publicAddrs []string
	var loopbackAddrs []string

	for _, ipAddr := range localIPs() {
		ip := net.ParseIP(ipAddr)
		if ip == nil {
			continue
		}

		if ip.IsUnspecified() {
			continue
		}

		if ip.IsLoopback() {
			loopbackAddrs = append(loopbackAddrs, ipAddr)
		} else if isPrivateIP(ipAddr) {
			privateAddrs = append(privateAddrs, ipAddr)
		} else {
			publicAddrs = append(publicAddrs, ipAddr)
		}
	}

	if len(privateAddrs) > 0 {
		return privateAddrs[0], nil
	} else if len(publicAddrs) > 0 {
		return publicAddrs[0], nil
	} else if len(loopbackAddrs) > 0 {
		return loopbackAddrs[0], nil
	}

	return "", errors.New("No IP address found, and explicit IP not provided")
}

// IPs returns all known ips
func IPs() []string {
	return localIPs()
}
