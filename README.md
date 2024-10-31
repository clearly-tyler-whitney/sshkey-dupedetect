# SSH Key Scanner

A Go program to scan specified networks for SSH servers and identify reused SSH host keys across different hosts. It accepts a `--cidr` flag with comma-separated CIDR notations to specify the networks to scan.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Example](#example)
- [Notes](#notes)
- [License](#license)

## Features

- Scans multiple CIDR ranges concurrently for SSH servers.
- Retrieves SSH host keys without needing authentication.
- Identifies and reports duplicate SSH host keys used across different IPs.
- Utilizes concurrency to speed up the scanning process.

## Prerequisites

- Go 1.11 or later installed on your system.
- Network access to the IP addresses in the specified CIDR ranges.
- Necessary permissions to perform network scanning.

## Installation

1. **Clone the Repository or Save the Code**

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
   ```

4. **Build the Program**

   Compile the Go program:

   ```bash
   go build -o ssh_key_scanner ssh_key_scanner.go
   ```

## Usage

Run the compiled binary with the `--cidr` flag, providing a comma-separated list of CIDR notations:

```bash
./ssh_key_scanner --cidr="CIDR1,CIDR2,..."
```

**Or**, run the program directly without building:

```bash
go run ssh_key_scanner.go --cidr="CIDR1,CIDR2,..."
```

## Example

Scanning the networks `192.168.1.0/24` and `10.0.0.0/8`:

```bash
./ssh_key_scanner --cidr="192.168.1.0/24,10.0.0.0/8"
```

**Sample Output:**

```
Duplicate host key SHA256:abcd1234... used by hosts: [192.168.1.10 192.168.1.20]
Duplicate host key SHA256:efgh5678... used by hosts: [10.0.0.5 10.0.0.25]
```

## Notes

- **Network Permissions:**
  - Ensure you have the necessary permissions to scan the specified networks.
  - Unauthorized scanning may violate network policies and laws.

- **Firewall and Security Software:**
  - Be aware that scanning might trigger alerts or be blocked by firewalls.

- **Performance:**
  - The program uses concurrency (`goroutines`) to scan multiple IPs simultaneously.
  - Adjust system limits if scanning very large networks to prevent resource exhaustion.

- **Understanding the Output:**
  - The program reports SSH host key fingerprints (SHA256) that are used by multiple hosts.
  - Reused SSH host keys can be a security concern, indicating misconfiguration.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**Disclaimer:** Use this tool responsibly and only on networks for which you have explicit permission to scan. The authors are not liable for any misuse of this software.