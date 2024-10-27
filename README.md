
# SNI Finder

![GitHub CI Status](https://github.com/v-kamerdinerov/SNI-Finder/actions/workflows/lint.yml/badge.svg)
![GitHub CI Status](https://github.com/v-kamerdinerov/SNI-Finder/actions/workflows/release.yml/badge.svg)
![GitHub](https://img.shields.io/github/license/v-kamerdinerov/SNI-Finder)
[![GitHub tag](https://img.shields.io/github/tag/v-kamerdinerov/SNI-Finder.svg)](https://github.com/v-kamerdinerov/SNI-Finder/tags)

----

<p align="center">
 <a href="./README.md">
 English
 </a>
 /
 <a href="./README-RU.md">
 Русский
 </a>
</p>

This app scans a range of IP addresses for domains with TLS 1.3 and HTTP/2 (h2) enabled. It is designed to identify useful SNI (Server Name Indication) domains for various configurations and tests.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Run on Linux/Mac OS](#run-on-linuxmac-os)
  - [Run on Windows](#run-on-windows)
- [Parameters](#parameters)
- [Example Output](#example-output)
- [Building from Source](#building-from-source)
  - [Prerequisites](#prerequisites)
  - [Build Steps](#build-steps)
- [Ethical Usage](#ethical-usage)
- [Contribution](#contribution)
- [License](#license)

## Features

| Feature                | Description                                                                                    |
|------------------------|------------------------------------------------------------------------------------------------|
| **TLS 1.3 and HTTP/2** | Scans for domains supporting TLS 1.3 and HTTP/2 to identify modern, secure configurations.     |
| **TLS Ping Display**   | Shows TLS Ping results for the scanned domains, helping assess response times.                 |
| **Top Servers**        | Outputs the top servers based on the lowest ping values, useful for prioritizing fast servers. |

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
    ./SNI-Finder -addr <ip-address>
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

| Argument    | Type   | Default Value | Description                                            |
|-------------|--------|---------------|--------------------------------------------------------|
| `-addr`     | string | `0.0.0.0`     | The starting address for the scan.                     |
| `-num`      | int    | `10000`       | The number of IPs to scan.                             |
| `-o`        | bool   | `true`        | Enable or disable output to the `results.txt` file.    |
| `-port`     | string | `443`         | The port to scan.                                      |
| `-showFail` | bool   | `false`       | Show logs for failed scans.                            |
| `-thread`   | int    | `128`         | The number of threads to run in parallel for scanning. |
| `-timeOut`  | int    | `4`           | The scan timeout in seconds.                           |
| `-top`      | int    | `10`          | The number of top servers to display.                  |

## Example Output

Here’s a sample output for your reference:

```
2024-09-08 22:25:49 37.139.62.38:443      TLS v1.3    ALPN: h2  pixelproof.ru         Ping: 22.315ms                      

2024-09-08 22:25:51 37.139.62.151:443     TLS v1.3    ALPN: h2  maelia.rs             Ping: 34.011ms                      

2024-09-08 22:25:51 37.139.62.192:443     TLS v1.3    ALPN: h2  doterra-challenge.ru  Ping: 22.68ms                       

2024-09-08 22:28:06 Scan completed.

Top servers by TLS Ping:
1: 37.139.41.28:443      TLS v1.3    ALPN: h2  jsdaddy.tech          Ping: 3.121ms                       
2: 37.139.42.255:443     TLS v1.3    ALPN: h2  malyi-biznes.ru       Ping: 3.496ms                       
3: 37.139.42.35:443      TLS v1.3    ALPN: h2  88date.co             Ping: 3.797ms                       
4: 37.139.41.247:443     TLS v1.3    ALPN: h2  you-note.ru           Ping: 3.804ms                       
5: 37.139.43.113:443     TLS v1.3    ALPN: h2  medvuza.ru            Ping: 4.029ms                       
6: 37.139.40.192:443     TLS v1.3    ALPN: h2  xn--8f9ac.xn--p1ai    Ping: 9.772ms                       
7: 37.139.42.68:443      TLS v1.3    ALPN: h2  mega74.ru             Ping: 11.621ms                      
8: 37.139.62.38:443      TLS v1.3    ALPN: h2  pixelproof.ru         Ping: 22.315ms                      
9: 37.139.62.192:443     TLS v1.3    ALPN: h2  doterra-challenge.ru  Ping: 22.68ms                       
10: 37.139.62.151:443     TLS v1.3    ALPN: h2  maelia.rs             Ping: 34.011ms
```

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

## Ethical Usage

This tool is intended for educational and legitimate testing purposes only. Unauthorized scanning of networks without permission may violate laws and result in severe penalties. Use responsibly.

## Contribution

Contributions are welcome! Please fork the repository and submit a pull request. Make sure your code adheres to the existing style and is thoroughly tested.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
