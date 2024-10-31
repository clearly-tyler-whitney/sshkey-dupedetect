package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
)

type DuplicateHostKey struct {
	Fingerprint string   `json:"fingerprint"`
	Hosts       []string `json:"hosts"`
}

func main() {
	// Define flags with short versions using pflag
	var (
		rateLimit    int
		concurrency  int
		verbosity    int
		progress     bool
		outputFormat string
	)

	pflag.IntVarP(&rateLimit, "rate-limit", "r", 100, "Number of scan attempts per second")
	pflag.IntVarP(&concurrency, "concurrency", "c", 50, "Maximum number of concurrent scanning goroutines")
	pflag.IntVarP(&verbosity, "verbosity", "v", 0, "Verbosity level (0-4)")
	pflag.BoolVarP(&progress, "progress", "p", false, "Show progress bar")
	pflag.StringVarP(&outputFormat, "output-format", "o", "table", "Output format: table, json, csv")

	pflag.Parse()

	// Get CIDR arguments from positional arguments
	cidrArgs := pflag.Args()
	if len(cidrArgs) == 0 {
		fmt.Println("Please provide at least one CIDR as a positional argument")
		fmt.Println("Usage: ssh_key_scanner [options] CIDR1 CIDR2 ...")
		pflag.PrintDefaults()
		return
	}

	var cidrList []string
	for _, arg := range cidrArgs {
		cidrs := strings.Split(arg, ",")
		for _, cidr := range cidrs {
			cidrList = append(cidrList, strings.TrimSpace(cidr))
		}
	}

	var ipList []string

	// Expand each CIDR into individual IP addresses
	for _, cidr := range cidrList {
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
	rateLimiter := time.Tick(time.Second / time.Duration(rateLimit))
	concurrencyLimiter := make(chan struct{}, concurrency)

	// Initialize progress bar if enabled
	var bar *progressbar.ProgressBar
	if progress {
		bar = progressbar.Default(int64(len(ipList)))
	}

	// Concurrently scan each IP address
	for _, ip := range ipList {
		<-rateLimiter // Rate limiting

		wg.Add(1)
		concurrencyLimiter <- struct{}{} // Acquire a slot
		go func(ip string) {
			defer wg.Done()
			defer func() { <-concurrencyLimiter }() // Release the slot

			hostKey, err := getHostKey(ip+":22", verbosity)
			if err != nil {
				if verbosity >= 3 {
					fmt.Printf("Error connecting to %s: %v\n", ip, err)
				}
			} else if hostKey == nil {
				if verbosity >= 4 {
					fmt.Printf("No host key found for %s\n", ip)
				}
			} else {
				fingerprint := ssh.FingerprintSHA256(hostKey)
				mu.Lock()
				hostKeyMap[fingerprint] = append(hostKeyMap[fingerprint], ip)
				mu.Unlock()
				if verbosity >= 2 {
					fmt.Printf("Scanned %s: %s\n", ip, fingerprint)
				}
			}

			if bar != nil {
				bar.Add(1)
			}
		}(ip)
	}

	wg.Wait()

	// Identify duplicate host keys
	var duplicates []DuplicateHostKey
	for fingerprint, ips := range hostKeyMap {
		if len(ips) > 1 {
			// Sort the IPs for consistent output
			sort.Strings(ips)
			duplicates = append(duplicates, DuplicateHostKey{
				Fingerprint: fingerprint,
				Hosts:       ips,
			})
		}
	}

	// Sort duplicates for consistent output
	sort.Slice(duplicates, func(i, j int) bool {
		return duplicates[i].Fingerprint < duplicates[j].Fingerprint
	})

	// Output the duplicates in the specified format
	switch strings.ToLower(outputFormat) {
	case "table", "pretty":
		outputTable(duplicates)
	case "json":
		outputJSON(duplicates)
	case "csv":
		outputCSV(duplicates)
	default:
		fmt.Printf("Unknown output format: %s\n", outputFormat)
	}
}

// getHostKey connects to an SSH server and retrieves its host key
func getHostKey(addr string, verbosity int) (ssh.PublicKey, error) {
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

func outputTable(duplicates []DuplicateHostKey) {
	for _, dup := range duplicates {
		fmt.Printf("Duplicate host key %s used by hosts:\n", dup.Fingerprint)
		for _, host := range dup.Hosts {
			fmt.Printf("  - %s\n", host)
		}
	}
}

func outputJSON(duplicates []DuplicateHostKey) {
	data, err := json.MarshalIndent(duplicates, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func outputCSV(duplicates []DuplicateHostKey) {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	// Write header
	w.Write([]string{"Fingerprint", "Hosts"})

	for _, dup := range duplicates {
		w.Write([]string{dup.Fingerprint, strings.Join(dup.Hosts, ";")})
	}
}
