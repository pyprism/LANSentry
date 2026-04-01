package resolver

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"lansentry/device"
)

// Resolver enriches scanned devices with hostnames by querying the gateway's
// DNS, the system resolver (which includes mDNS on macOS), and NetBIOS.
type Resolver struct {
	gatewayIP   string
	concurrency int
	verbose     bool
}

// New creates a new Resolver. It auto-detects the default gateway so it can
// query the router's DNS directly (routers know DHCP client hostnames).
func New(verbose bool) *Resolver {
	gw := detectGateway()
	if verbose {
		if gw != "" {
			log.Printf("Gateway detected: %s (will query for DHCP hostnames)", gw)
		} else {
			log.Printf("Gateway not detected; will use system DNS only")
		}
	}

	return &Resolver{
		gatewayIP:   gw,
		concurrency: 20,
		verbose:     verbose,
	}
}

// EnrichDevices resolves hostnames for devices that are still missing them.
// Lookups run concurrently with per-host timeouts.
func (r *Resolver) EnrichDevices(devices []device.Device) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, r.concurrency)

	for i := range devices {
		if devices[i].Hostname != "" {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(d *device.Device) {
			defer wg.Done()
			defer func() { <-sem }()
			d.Hostname = r.resolveHostname(d.IP)
		}(&devices[i])
	}

	wg.Wait()
}

// ---------------------------------------------------------------------------
// Hostname resolution (tried in order)
//  1. Gateway DNS  – router knows DHCP hostnames
//  2. System DNS   – includes mDNS on macOS via mDNSResponder / avahi on Linux
//  3. NetBIOS      – Windows/Samba devices (if nmblookup is installed)
// ---------------------------------------------------------------------------

const lookupTimeout = 2 * time.Second

func (r *Resolver) resolveHostname(ip string) string {
	// 1. Gateway DNS (most home routers serve DHCP names via dnsmasq/etc.)
	if r.gatewayIP != "" {
		if name := r.lookupViaDNS(ip, r.gatewayIP+":53"); name != "" {
			if r.verbose {
				log.Printf("  resolved %s → %s (gateway DNS)", ip, name)
			}
			return name
		}
	}

	// 2. System resolver (mDNS on macOS, avahi on Linux, regular DNS)
	if name := r.lookupSystemDNS(ip); name != "" {
		if r.verbose {
			log.Printf("  resolved %s → %s (system DNS)", ip, name)
		}
		return name
	}

	// 3. NetBIOS
	if name := r.lookupNetBIOS(ip); name != "" {
		if r.verbose {
			log.Printf("  resolved %s → %s (NetBIOS)", ip, name)
		}
		return name
	}

	return ""
}

// lookupViaDNS queries a specific DNS server (e.g. the gateway) for PTR records.
func (r *Resolver) lookupViaDNS(ip, server string) string {
	res := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return (&net.Dialer{Timeout: lookupTimeout}).DialContext(ctx, "udp", server)
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), lookupTimeout)
	defer cancel()

	names, err := res.LookupAddr(ctx, ip)
	if err != nil || len(names) == 0 {
		return ""
	}
	return cleanHostname(names[0], ip)
}

// lookupSystemDNS uses the OS default resolver (includes mDNS on macOS).
func (r *Resolver) lookupSystemDNS(ip string) string {
	ctx, cancel := context.WithTimeout(context.Background(), lookupTimeout)
	defer cancel()

	names, err := (&net.Resolver{}).LookupAddr(ctx, ip)
	if err != nil || len(names) == 0 {
		return ""
	}
	return cleanHostname(names[0], ip)
}

// ---------------------------------------------------------------------------
// Hostname cleaning
// ---------------------------------------------------------------------------

// cleanHostname strips the trailing DNS dot and common LAN domain suffixes
// so names match what a router admin UI typically shows.
func cleanHostname(raw, ip string) string {
	name := strings.TrimSuffix(raw, ".")
	if name == "" || name == ip {
		return ""
	}

	lower := strings.ToLower(name)
	for _, suffix := range lanSuffixes {
		if strings.HasSuffix(lower, suffix) {
			if stripped := name[:len(name)-len(suffix)]; stripped != "" {
				return stripped
			}
		}
	}
	return name
}

var lanSuffixes = []string{
	".lan",
	".local",
	".home",
	".localdomain",
	".internal",
	".private",
	".fritz.box",
	".home.arpa",
}

// ---------------------------------------------------------------------------
// NetBIOS (Windows / Samba devices)
// ---------------------------------------------------------------------------

var nbRe = regexp.MustCompile(`(?m)^\s+(\S+)\s+<00>\s+-\s+[^<]*<ACTIVE>`)

func (r *Resolver) lookupNetBIOS(ip string) string {
	if _, err := exec.LookPath("nmblookup"); err != nil {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), lookupTimeout)
	defer cancel()

	out, err := exec.CommandContext(ctx, "nmblookup", "-A", ip).Output()
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "<GROUP>") {
			continue
		}
		if m := nbRe.FindStringSubmatch(line); len(m) >= 2 {
			if name := strings.TrimSpace(m[1]); name != "" {
				return name
			}
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// Gateway detection
//   Tries in order:
//   1. ip route show default   (modern Linux, iproute2)
//   2. netstat -rn             (macOS, legacy Linux with net-tools)
//   3. /proc/net/route         (Linux fallback, no external tools needed)
// ---------------------------------------------------------------------------

var (
	// "default via 192.168.1.1 dev eth0"
	ipRouteRe = regexp.MustCompile(`default\s+via\s+(\d+\.\d+\.\d+\.\d+)`)
	// "default  192.168.1.1  UGScg  en0" or "0.0.0.0  192.168.1.1  ..."
	netstatRe = regexp.MustCompile(`(?:default|0\.0\.0\.0)\s+(\d+\.\d+\.\d+\.\d+)`)
)

func detectGateway() string {
	// 1. ip route (iproute2, standard on modern Linux)
	if gw := detectGatewayIPRoute(); gw != "" {
		return gw
	}
	// 2. netstat (macOS, older Linux)
	if gw := detectGatewayNetstat(); gw != "" {
		return gw
	}
	// 3. /proc/net/route (Linux without external tools)
	return detectGatewayProc()
}

func detectGatewayIPRoute() string {
	out, err := exec.Command("ip", "route", "show", "default").Output()
	if err != nil {
		return ""
	}
	if m := ipRouteRe.FindStringSubmatch(string(out)); len(m) >= 2 {
		return m[1]
	}
	return ""
}

func detectGatewayNetstat() string {
	out, err := exec.Command("netstat", "-rn").Output()
	if err != nil {
		return ""
	}
	if m := netstatRe.FindStringSubmatch(string(out)); len(m) >= 2 {
		return m[1]
	}
	return ""
}

// detectGatewayProc parses /proc/net/route (Linux only).
// Fields: Iface Destination Gateway Flags RefCnt Use Metric Mask ...
// Gateway for the default route (Destination == "00000000") is in little-endian hex.
func detectGatewayProc() string {
	data, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		if fields[1] != "00000000" { // default route
			continue
		}
		gw := fields[2]
		if len(gw) != 8 {
			continue
		}
		// Parse little-endian hex IP: e.g. "0101A8C0" → 192.168.1.1
		var octets [4]uint64
		for i := 0; i < 4; i++ {
			v, err := strconv.ParseUint(gw[i*2:i*2+2], 16, 8)
			if err != nil {
				return ""
			}
			octets[i] = v
		}
		return fmt.Sprintf("%d.%d.%d.%d", octets[3], octets[2], octets[1], octets[0])
	}
	return ""
}
