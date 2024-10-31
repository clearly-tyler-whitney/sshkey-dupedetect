package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	// Define the flags
	cidrStr := flag.String("cidr", "", "Comma-separated list of CIDRs to scan")
	rateLimit := flag.Int("rate-limit", 100, "Number of scan attempts per second")
	concurrency := flag.Int("concurrency", 50, "Maximum number of concurrent scanning goroutines")
	flag.Parse()

	if *cidrStr == "" {
		fmt.Println("Please provide CIDRs using --cidr flag")
		return
	}

	// Split the CIDR string into a list
	cidrList := strings.Split(*cidrStr, ",")
	var ipList []string

	// Expand each CIDR into individual IP addresses
	for _, cidr := range cidrList {
		cidr = strings.TrimSpace(cidr)
		ips, err := hosts(cidr)
		if err != nil {
			fmt.Printf("Invalid CIDR %s: %v\n", cidr, err)
			continue
		}
		ipList = append(ipList, ips...)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	hostKeyMap := make(map[string][]string) // Map of host key fingerprints to IPs

	// Create rate limiter and concurrency limiter
	rateLimiter := time.Tick(time.Second / time.Duration(*rateLimit))
	concurrencyLimiter := make(chan struct{}, *concurrency)

	// Concurrently scan each IP address
	for _, ip := range ipList {
		<-rateLimiter // Rate limiting

		wg.Add(1)
		concurrencyLimiter <- struct{}{} // Acquire a slot
		go func(ip string) {
			defer wg.Done()
			defer func() { <-concurrencyLimiter }() // Release the slot

			hostKey, err := getHostKey(ip + ":22")
			if err != nil {
				return
			}
			if hostKey == nil {
				return
			}
			fingerprint := ssh.FingerprintSHA256(hostKey)
			mu.Lock()
			hostKeyMap[fingerprint] = append(hostKeyMap[fingerprint], ip)
			mu.Unlock()
		}(ip)
	}

	wg.Wait()

	// Identify and print duplicate host keys
	for fingerprint, ips := range hostKeyMap {
		if len(ips) > 1 {
			fmt.Printf("Duplicate host key %s used by hosts: %v\n", fingerprint, ips)
		}
	}
}

// getHostKey connects to an SSH server and retrieves its host key
func getHostKey(addr string) (ssh.PublicKey, error) {
	var hostKey ssh.PublicKey

	config := &ssh.ClientConfig{
		User: "ignored",
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			hostKey = key
			return nil
		},
		Auth:    []ssh.AuthMethod{},
		Timeout: 5 * time.Second,
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		if strings.Contains(err.Error(), "unable to authenticate") || strings.Contains(err.Error(), "no common algorithm") {
			return hostKey, nil // Host key is retrieved even if authentication fails
		}
		return nil, err
	}
	defer conn.Close()

	return hostKey, nil
}

// hosts generates a list of IP addresses from a CIDR notation
func hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// Remove network and broadcast addresses
	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}
	return ips, nil
}

// inc increments an IP address
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
