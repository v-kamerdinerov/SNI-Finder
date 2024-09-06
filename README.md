# SNI Finder


This app scans a range of IP addresses for domains with TLS 1.3 and HTTP/2 (h2) enabled. It is designed to identify useful SNI (Server Name Indication) domains for various configurations and tests.

### Features

| Feature                | Description                                              |
|------------------------|----------------------------------------------------------|
| TLS 1.3 and HTTP/2     | Scans for domains supporting TLS 1.3 and HTTP/2.        |
| TLS Ping Display       | Displays TLS Ping results for the scanned domains.      |
| Top 5 Servers          | Outputs the top 5 servers based on the lowest ping values. |

When you begin the scan, two files are created: `results.txt` contains the output log, while `domains.txt` only contains the domain names.

It is recommended to run this scanner locally _(with your residential internet)_. It may cause VPS to be flagged if you run a scanner in the cloud.


## Run

### Linux/Mac OS:

Be careful when choosing an architecture, all binaries are available in two versions - `amd64` and `arm64`.

1.
    ```
    wget "https://github.com/hawshemi/SNI-Finder/releases/latest/download/SNI-Finder-$(uname -s | tr A-Z a-z)-amd64" -O SNI-Finder && chmod +x SNI-Finder
    ```
2. 
    ```
    ./SNI-Finder -addr ip
    ```

### Windows:

1. Download from [Releases](https://github.com/hawshemi/SNI-Finder/releases/latest).
2. Open `CMD` or `Powershell` in the directory.
3.
    ```
    .\SNI-Finder-windows-amd64.exe -addr ip
    ```

#### Replace `ip` with your VPS IP Address.

## Params

| Argument       | Type    | Default Value | Description                                         |
|----------------|---------|---------------|-----------------------------------------------------|
| `-addr`        | string  | `0.0.0.0`     | The starting address for the scan.                  |
| `-port`        | string  | `443`         | The port to scan.                                   |
| `-thread`      | int     | `128`         | The number of threads to run in parallel for scanning. |
| `-o`           | bool    | `true`        | Enable or disable output to the `results.txt` file. |
| `-timeOut`     | int     | `4`           | The scan timeout in seconds.                        |
| `-showFail`    | bool    | `false`       | Show logs for failed scans.                         |



## Build

### Prerequisites

#### Install `wget`:
```
sudo apt install -y wget
```

#### First run this script to install `Go` & other dependencies _(Debian & Ubuntu)_:
```
wget "https://raw.githubusercontent.com/hawshemi/SNI-Finder/main/install-go.sh" -O install-go.sh && chmod +x install-go.sh && bash install-go.sh
```
- Reboot is recommended.


#### Then:

#### 1. Clone the repository
```
git clone https://github.com/hawshemi/SNI-Finder.git 
```

#### 2. Navigate into the repository directory
```
cd SNI-Finder 
```

#### 3. Download and install `logrus` package
```
go get github.com/sirupsen/logrus
```

#### 4. Build
```
CGO_ENABLED=0 go build
```
