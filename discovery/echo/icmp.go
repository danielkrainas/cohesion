package echo

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	fastping "github.com/tatsushid/go-fastping"

	"github.com/danielkrainas/cohesion/context"
	"github.com/danielkrainas/cohesion/discovery"
	"github.com/danielkrainas/cohesion/discovery/factory"
)

type driverFactory struct{}

func (f *driverFactory) Create(parameters map[string]interface{}) (discovery.Strategy, error) {
	cidr, ok := parameters["network"].(string)
	if !ok {
		return nil, errors.New("invalid configuration value for `network`")
	}

	return &driver{
		cidr: cidr,
	}, nil
}

func init() {
	factory.Register("echo", &driverFactory{})
}

type driver struct {
	cidr string
}

func (d *driver) Locate(ctx context.Context) ([]string, error) {
	var hosts []string
	if d.cidr != "" {
		hosts = getHostsFromCIDR(d.cidr)
	} else {
		hosts = getHosts()
	}

	if len(hosts) < 1 {
		context.GetLogger(ctx).Warn("no hosts to ping, skipping")
		return nil, nil
	}

	m := sync.Mutex{}
	p := fastping.NewPinger()
	for _, h := range hosts {
		ra, err := net.ResolveIPAddr("ip4:icmp", h)
		if err != nil {

		}

		p.AddIPAddr(ra)

	}

	results := make([]string, 0)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		m.Lock()
		defer m.Unlock()
		results = append(results, addr.String())
	}

	p.OnIdle = func() {}

	if err := p.Run(); err != nil {
		return nil, err
	}

	return results, nil
}

func getHostsFromCIDR(cidr string) []string {
	results := make([]string, 0)
	lookup := getLocalAddressLookup()
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		if _, found := lookup[ip.String()]; found {
			continue
		}

		results = append(results, fmt.Sprint(ip))
	}

	return results
}

func getLocalAddressLookup() map[string]bool {
	results := make(map[string]bool)
	ifaces, err := net.Interfaces()
	if err != nil {
		return results
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return results
		}

		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPAddr:
				ip = v.IP
			case *net.IPNet:
				ip = v.IP
			}

			if ip == nil || ip.DefaultMask() == nil {
				continue
			}

			results[ip.String()] = true
		}
	}

	return results
}

func getHosts() []string {
	results := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return results
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return results
		}

		for _, a := range addrs {
			var ip net.IP
			var mask net.IPMask
			switch v := a.(type) {
			case *net.IPAddr:
				ip = v.IP
				mask = ip.DefaultMask()
			case *net.IPNet:
				ip = v.IP
				mask = v.Mask
			}

			if ip == nil || ip.DefaultMask() == nil || ip.IsLoopback() {
				continue
			}

			n := net.IPNet{
				IP:   ip.Mask(mask),
				Mask: mask,
			}

			for ip := ip.Mask(mask); n.Contains(ip); inc(ip) {
				results = append(results, fmt.Sprint(ip))
			}
		}
	}

	return results
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
