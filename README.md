# SSH Key Scanner

A Go program to scan specified networks for SSH servers and identify reused SSH host keys across different hosts. It accepts CIDR notations as positional arguments to specify the networks to scan.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Options](#options)
- [Notes](#notes)
- [License](#license)

## Features

- Scans multiple CIDR ranges concurrently for SSH servers.
- Retrieves SSH host keys without needing authentication.
- Identifies and reports duplicate SSH host keys used across different IPs.
- Supports rate limiting and concurrency control.
- Provides optional progress bar and verbosity levels.
- Outputs results in table, JSON, or CSV formats.

## Prerequisites

- Go 1.11 or later installed on your system.
- Network access to the IP addresses in the specified CIDR ranges.
- Necessary permissions to perform network scanning.

## Installation

1. **Save the Code**

   Save the `ssh_key_scanner.go` file in a directory of your choice.

2. **Initialize a Go Module**

   Open a terminal in the project directory and run:

   ```bash
   go mod init ssh_key_scanner
   ```

3. **Download Dependencies**

   Retrieve the necessary Go packages:

   ```bash
   go get golang.org/x/crypto/ssh
   go get github.com/schollz/progressbar/v3
   go get github.com/spf13/pflag
   ```

4. **Build the Program**

   Compile the Go program:

   ```bash
   go build -o ssh_key_scanner ssh_key_scanner.go
   ```

## Usage

Run the compiled binary with the desired options and CIDR notations as positional arguments:

```bash
./ssh_key_scanner [options] CIDR1 [CIDR2 ...]
```

**Or**, run the program directly without building:

```bash
go run ssh_key_scanner.go [options] CIDR1 [CIDR2 ...]
```

## Examples

1. **Scanning a Single CIDR Range:**

   ```bash
   ./ssh_key_scanner 192.168.1.0/24
   ```

2. **Scanning Multiple CIDR Ranges:**

   ```bash
   ./ssh_key_scanner 192.168.1.0/24 10.0.0.0/8
   ```

3. **Using Short Options and Enabling Progress Bar:**

   ```bash
   ./ssh_key_scanner -p -v 2 -r 200 -c 100 -o json 192.168.1.0/24 10.0.0.0/8
   ```

   - `-p`: Enable progress bar.
   - `-v 2`: Set verbosity level to 2.
   - `-r 200`: Set rate limit to 200 scans per second.
   - `-c 100`: Set concurrency to 100 goroutines.
   - `-o json`: Output results in JSON format.

## Options

- `--rate-limit, -r int`  
  Number of scan attempts per second (default 100).

- `--concurrency, -c int`  
  Maximum number of concurrent scanning goroutines (default 50).

- `--verbosity, -v int`  
  Verbosity level (0-4) (default 0).

  - **Level 0:** Minimal output (only duplicates).
  - **Level 1:** Basic progress updates.
  - **Level 2:** Scanned IPs and their fingerprints.
  - **Level 3:** Includes errors encountered during connections.
  - **Level 4:** Detailed information, including when no host key is found.

- `--progress, -p`  
  Show progress bar (off by default).

- `--output-format, -o string`  
  Output format: `table` (default), `json`, `csv`.

## Notes

- **Network Permissions:**
  - Ensure you have the necessary permissions to scan the specified networks.
  - Unauthorized scanning may violate network policies and laws.

- **Firewall and Security Software:**
  - Be aware that scanning might trigger alerts or be blocked by firewalls.

- **Performance:**
  - Adjust the `--rate-limit` and `--concurrency` options based on your network environment and system capabilities.
  - High concurrency levels can consume significant system resources.

- **Understanding the Output:**
  - The program reports SSH host key fingerprints (SHA256) that are used by multiple hosts.
  - Reused SSH host keys can be a security concern, indicating misconfiguration.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**Disclaimer:** Use this tool responsibly and only on networks for which you have explicit permission to scan. The authors are not liable for any misuse of this software.
```
