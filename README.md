
# SNI Finder

This app scans a range of IP addresses for domains with TLS 1.3 and HTTP/2 (h2) enabled. It is designed to identify useful SNI (Server Name Indication) domains for various configurations and tests.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Run on Linux/Mac OS](#run-on-linuxmac-os)
  - [Run on Windows](#run-on-windows)
- [Parameters](#parameters)
- [Building from Source](#building-from-source)
  - [Prerequisites](#prerequisites)
  - [Build Steps](#build-steps)
- [Example Output](#example-output)
- [Ethical Usage](#ethical-usage)
- [Contribution](#contribution)
- [License](#license)

## Features

| Feature             | Description                                              |
|---------------------|----------------------------------------------------------|
| **TLS 1.3 and HTTP/2** | Scans for domains supporting TLS 1.3 and HTTP/2 to identify modern, secure configurations. |
| **TLS Ping Display** | Shows TLS Ping results for the scanned domains, helping assess response times. |
| **Top Servers**     | Outputs the top servers based on the lowest ping values, useful for prioritizing fast servers. |

### Output Files
- `results.txt`: Contains the detailed scan log.
- `domains.txt`: Contains only the discovered domain names.

> **Note:** It is recommended to run this scanner locally _(with your residential internet)_. Running a scanner in the cloud may cause VPS providers to flag your account.

## Installation

Download the appropriate binary for your system architecture (`amd64` or `arm64`).

## Usage

### Run on Linux/Mac OS:

1. Download and set up the binary:
    ```bash
    wget "https://github.com/v-kamerdinerov/SNI-Finder/releases/latest/download/SNI-Finder-$(uname -s | tr A-Z a-z)-amd64" -O SNI-Finder && chmod +x SNI-Finder
    ```
2. Run the scanner:
    ```bash
    ./SNI-Finder-linux-amd64 -addr <ip-address>
    ```

### Run on Windows:

1. Download the binary from [Releases](https://github.com/v-kamerdinerov/SNI-Finder/releases/latest).
2. Open `CMD` or `Powershell` in the download directory.
3. Run the scanner:
    ```cmd
    .\SNI-Finder-windows-amd64.exe -addr <ip-address>
    ```

> Replace `<ip-address>` with the IP address range you want to scan.

## Parameters

| Argument       | Type    | Default Value | Description                                         |
|----------------|---------|---------------|-----------------------------------------------------|
| `-addr`        | string  | `0.0.0.0`     | The starting address for the scan.                  |
| `-port`        | string  | `443`         | The port to scan.                                   |
| `-thread`      | int     | `128`         | The number of threads to run in parallel for scanning. |
| `-o`           | bool    | `true`        | Enable or disable output to the `results.txt` file. |
| `-timeOut`     | int     | `4`           | The scan timeout in seconds.                        |
| `-showFail`    | bool    | `false`       | Show logs for failed scans.                         |

## Building from Source

### Prerequisites

Install `wget`:
```bash
sudo apt install -y wget
```

Run this script to install Go and other dependencies _(Debian & Ubuntu & RedOS $ RHEL)_:
```bash
wget "https://raw.githubusercontent.com/v-kamerdinerov/SNI-Finder/main/install-go.sh" -O install-go.sh && chmod +x install-go.sh && bash install-go.sh
```
> **Tip:** A system reboot is recommended after installation.

### Build Steps

1. Clone the repository:
    ```bash
    git clone https://github.com/v-kamerdinerov/SNI-Finder.git 
    ```
2. Navigate into the repository directory:
    ```bash
    cd SNI-Finder 
    ```
3. Download and install the all required package:
    ```bash
    go mod tidy
    ```
4. Build the project:
    ```bash
    CGO_ENABLED=0 go build
    ```

## Example Output

Here’s a sample output for your reference:

```
2024-09-06 22:15:13 95.179.221.159:443    TLS v1.3    ALPN: h2  ilink-app.com         Ping: 100.254ms                     

2024-09-06 22:15:16 95.179.222.66:443     TLS v1.3    ALPN: h2  cretathemes.com       Ping: 90.27ms                       

2024-09-06 22:15:17 95.179.222.98:443     TLS v1.3    ALPN: h2  cityy.net             Ping: 98.162ms                     

2024-09-06 22:15:48 Scan completed.

Top servers by TLS Ping:
1: 95.179.190.252:443    TLS v1.3    ALPN: h2  cloudflare.com        Ping: 29.271ms                       (Ping: 29.271ms)
2: 95.179.211.147:443    TLS v1.3    ALPN: h2  cloudflare-dns.com    Ping: 29.84ms                        (Ping: 29.84ms)
3: 95.179.201.194:443    TLS v1.3    ALPN: h2  belugabahis.com       Ping: 69.107ms                       (Ping: 69.107ms)
4: 95.179.207.32:443     TLS v1.3    ALPN: h2  actlocalmedia.com     Ping: 69.762ms                       (Ping: 69.762ms)
5: 95.179.210.147:443    TLS v1.3    ALPN: h2  webnativ.fr           Ping: 70.294ms                       (Ping: 70.294ms)
```

## Ethical Usage

This tool is intended for educational and legitimate testing purposes only. Unauthorized scanning of networks without permission may violate laws and result in severe penalties. Use responsibly.

## Contribution

Contributions are welcome! Please fork the repository and submit a pull request. Make sure your code adheres to the existing style and is thoroughly tested.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
