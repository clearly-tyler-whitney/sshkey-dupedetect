# SSH Key Scanner

A Go program to scan specified networks for SSH servers and identify reused SSH host keys across different hosts. It accepts CIDR notations as positional arguments to specify the networks to scan.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Options](#options)
- [Docker Usage](#docker-usage)
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

- **Go 1.20 or later** installed on your system (if building without Docker).
- Network access to the IP addresses in the specified CIDR ranges.
- Necessary permissions to perform network scanning.
- **Docker** installed on your system (if using Docker).

## Installation

### **Cloning the Repository**

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/ssh_key_scanner.git
cd ssh_key_scanner
```

### **Building from Source**

1. **Initialize the Go Module**

   Ensure all dependencies are fetched:

   ```bash
   go mod tidy
   ```

2. **Build the Program**

   Compile the Go program:

   ```bash
   go build -o ssh_key_scanner
   ```

### **Building with Docker**

Build the Docker image using the provided `Dockerfile`:

```bash
docker build -t ssh_key_scanner .
```

## Usage

### **Running the Binary**

Run the compiled binary with the desired options and CIDR notations as positional arguments:

```bash
./ssh_key_scanner [options] CIDR1 [CIDR2 ...]
```

### **Running with Docker**

Run the Docker container with the necessary options:

```bash
docker run --rm ssh_key_scanner [options] CIDR1 [CIDR2 ...]
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

4. **Running with Docker and Saving Output:**

   ```bash
   docker run --rm ssh_key_scanner -o json 192.168.1.0/24 > results.json
   ```

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

## Docker Usage

### **Building the Docker Image**

Build the Docker image using the provided `Dockerfile`:

```bash
docker build -t ssh_key_scanner .
```

### **Running the Docker Container**

1. **Basic Usage**

   ```bash
   docker run --rm ssh_key_scanner [options] CIDR1 [CIDR2 ...]
   ```

2. **Example with Options**

   ```bash
   docker run --rm ssh_key_scanner -p -o json 192.168.1.0/24
   ```

3. **Using Host Network (Linux Only)**

   ```bash
   docker run --rm --network host ssh_key_scanner 192.168.1.0/24
   ```

   **Note:** `--network host` allows the container to access the host's network interfaces, which may be necessary to scan local networks.

4. **Saving Output to a File**

   ```bash
   docker run --rm ssh_key_scanner -o json 192.168.1.0/24 > results.json
   ```

## Notes

- **Network Permissions:**
  - Ensure you have the necessary permissions to scan the specified networks.
  - Unauthorized scanning may violate network policies and laws.

- **Firewall and Security Software:**
  - Be aware that scanning might trigger alerts or be blocked by firewalls.

- **Performance:**
  - Adjust the `--rate-limit` and `--concurrency` options based on your network environment and system capabilities.
  - High concurrency levels can consume significant system resources.

- **Docker on macOS and Windows:**
  - The `--network host` option is not supported on Docker for macOS and Windows due to limitations in the Docker network stack.
  - You may need to run the scanner directly on your host system to scan local networks on these platforms.

- **Understanding the Output:**
  - The program reports SSH host key fingerprints (SHA256) that are used by multiple hosts.
  - Reused SSH host keys can be a security concern, indicating misconfiguration.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**Disclaimer:** Use this tool responsibly and only on networks for which you have explicit permission to scan. The authors are not liable for any misuse of this software.
